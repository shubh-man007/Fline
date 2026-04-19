package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const PORT = "1324"

var re = regexp.MustCompile(`{{([^{}]*)}}`)

func getByPath(data any, path string) any {
	parts := strings.SplitN(path, ".", 2)
	key := parts[0]

	switch v := data.(type) {
	case map[string]any:
		val, ok := v[key]
		if !ok {
			return nil
		}
		if len(parts) == 1 {
			return val
		}
		return getByPath(val, parts[1])

	case []any:
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(v) {
			return nil
		}
		if len(parts) == 1 {
			return v[idx]
		}
		return getByPath(v[idx], parts[1])
	}

	return nil
}

func resolveTemplate(template string, wfCtx WorkFlowCtx) string {
	return re.ReplaceAllStringFunc(template, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		val := getByPath(map[string]any(wfCtx), key)
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%v", val)
	})
}

type WorkFlowCtx map[string]any

type Fliner interface {
	Execute(ctx context.Context, wfCtx WorkFlowCtx) (WorkFlowCtx, error)
	// GetContext()
}

type WorkFlow struct {
	ID         uuid.UUID
	Name       string
	Enabled    bool
	Steps      []Fliner
	CurrentCtx WorkFlowCtx
}

func NewWorkFlow() *WorkFlow {
	return &WorkFlow{}
}

type WorkFlowInit struct {
	Name  string          `json:"name"`
	Steps json.RawMessage `json:"steps"`
}

type WorkFlowResponse struct {
	ID      uuid.UUID         `json:"id"`
	Name    string            `json:"name"`
	Enabled bool              `json:"enabled"`
	Steps   []json.RawMessage `json:"steps"`
}

type WorkFlowStore struct {
	mu    sync.RWMutex
	Store map[uuid.UUID]*WorkFlow
}

func NewWorkFlowStore() *WorkFlowStore {
	return &WorkFlowStore{
		Store: make(map[uuid.UUID]*WorkFlow),
	}
}

func (ws *WorkFlowStore) AddWorkFlow(wf *WorkFlow) {
	ws.mu.Lock()
	ws.Store[wf.ID] = wf
	ws.mu.Unlock()

	log.Printf("[%s] workflow added to store", wf.ID.String())
}

func (ws *WorkFlowStore) DeleteWorkFlow(wfID uuid.UUID) {
	ws.mu.Lock()
	delete(ws.Store, wfID)
	ws.mu.Unlock()

	log.Printf("[%s] workflow deleted from store", wfID.String())
}

func (ws *WorkFlowStore) FlipWorkFlow(wfID uuid.UUID) (bool, error) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	wf, prs := ws.Store[wfID]

	if prs != true {
		return false, fmt.Errorf("WorkFlow: %s does not exists", wfID.String())
	}

	wf.Enabled = !wf.Enabled
	return wf.Enabled, nil
}

func (ws *WorkFlowStore) GetWorkFlow(wfID uuid.UUID) (*WorkFlow, error) {
	ws.mu.RLock()
	wf, prs := ws.Store[wfID]
	ws.mu.RUnlock()

	if prs != true {
		return nil, fmt.Errorf("WorkFlow: %s does not exists", wfID.String())
	}

	return wf, nil
}

type HTTPRequest struct {
	Endpoint string            `json:"URL"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Body     map[string]any    `json:"body"`
	TimeOut  int               `json:"timeout"`
	Retries  int               `json:"retries"`
	Client   *http.Client
}

func (h *HTTPRequest) Execute(ctx context.Context, wfCtx WorkFlowCtx) (WorkFlowCtx, error) {
	endpoint := resolveTemplate(h.Endpoint, wfCtx)
	log.Printf("Executing HTTP step: %s %s\n", h.Method, endpoint)

	var bodyReader io.Reader
	if h.Body != nil {
		b, err := json.Marshal(h.Body)
		if err != nil {
			return wfCtx, fmt.Errorf("error marshalling step body: %w", err)
		}
		resolved := resolveTemplate(string(b), wfCtx)
		bodyReader = strings.NewReader(resolved)
	} else if len(wfCtx) > 0 && h.Method != "GET" {
		b, err := json.Marshal(wfCtx)
		if err != nil {
			return wfCtx, fmt.Errorf("error marshalling context body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, h.Method, endpoint, bodyReader)
	if err != nil {
		return wfCtx, fmt.Errorf("error building request: %w", err)
	}

	for k, v := range h.Headers {
		req.Header.Add(k, resolveTemplate(v, wfCtx))
	}

	res, err := h.Client.Do(req)
	if err != nil {
		return wfCtx, fmt.Errorf("error executing request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, _ := io.ReadAll(res.Body)
		return wfCtx, fmt.Errorf("request failed with status %d: %s", res.StatusCode, string(body))
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return wfCtx, fmt.Errorf("error reading response body: %w", err)
	}

	var parsed WorkFlowCtx
	if err = json.Unmarshal(body, &parsed); err == nil {
		wfCtx["result"] = parsed
	} else {
		var parsedArray []any
		if err = json.Unmarshal(body, &parsedArray); err == nil {
			wfCtx["result"] = parsedArray
		} else {
			wfCtx["result"] = string(body)
		}
	}

	return wfCtx, nil
}

type TransformOp string

const (
	DefaultOp  TransformOp = "default"
	TemplateOp TransformOp = "template"
	PickOp     TransformOp = "pick"
)

type Transformer interface {
	TransformOperate(WorkFlowCtx) WorkFlowCtx
}

type TransformDefault struct {
	Op    TransformOp `json:"op"`
	Field string      `json:"path"`
	Value string      `json:"value"`
}

func (def *TransformDefault) TransformOperate(wfCtx WorkFlowCtx) WorkFlowCtx {
	_, ok := wfCtx[def.Field]
	if !ok {
		wfCtx[def.Field] = def.Value
		return wfCtx
	}
	return wfCtx
}

type TransformTemplate struct {
	Op       TransformOp `json:"op"`
	To       string      `json:"to"`
	Template string      `json:"template"`
}

func (temp *TransformTemplate) TransformOperate(wfCtx WorkFlowCtx) WorkFlowCtx {
	result := temp.Template
	matches := re.FindAllStringSubmatch(temp.Template, -1)

	for _, match := range matches {
		placeholder := match[0]
		key := strings.TrimSpace(match[1])

		val := getByPath(map[string]any(wfCtx), key)
		if val == nil {
			result = strings.ReplaceAll(result, placeholder, "")
			continue
		}
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", val))
	}

	wfCtx[temp.To] = result
	return wfCtx
}

type TransformPick struct {
	Op     TransformOp `json:"op"`
	Fields []string    `json:"paths"`
}

func (pick *TransformPick) TransformOperate(wfCtx WorkFlowCtx) WorkFlowCtx {
	tempCtx := make(WorkFlowCtx)
	for _, key := range pick.Fields {
		if val, ok := wfCtx[key]; ok {
			tempCtx[key] = val
		}
	}
	return tempCtx
}

type Transform struct {
	Ops []Transformer `json:"-"`
}

func (t *Transform) Execute(ctx context.Context, wfCtx WorkFlowCtx) (WorkFlowCtx, error) {
	for _, tstep := range t.Ops {
		wfCtx = tstep.TransformOperate(wfCtx)
	}

	return wfCtx, nil
}

type OpDiscriminator struct {
	Op TransformOp `json:"op"`
}

func ParseOps(raw json.RawMessage) ([]Transformer, error) {
	var rawOps []json.RawMessage
	if err := json.Unmarshal(raw, &rawOps); err != nil {
		return nil, fmt.Errorf("ops must be a JSON array: %w", err)
	}

	ops := make([]Transformer, 0, len(rawOps))

	for i, rawOp := range rawOps {
		var disc OpDiscriminator
		if err := json.Unmarshal(rawOp, &disc); err != nil {
			return nil, fmt.Errorf("op[%d]: missing or invalid op: %w", i, err)
		}

		var op Transformer
		var err error

		switch disc.Op {
		case DefaultOp:
			var o TransformDefault
			err = json.Unmarshal(rawOp, &o)
			op = &o
		case TemplateOp:
			var o TransformTemplate
			err = json.Unmarshal(rawOp, &o)
			op = &o
		case PickOp:
			var o TransformPick
			err = json.Unmarshal(rawOp, &o)
			op = &o
		default:
			return nil, fmt.Errorf("op[%d]: unknown op %q", i, disc.Op)
		}

		if err != nil {
			return nil, fmt.Errorf("op[%d] (%s): parse error: %w", i, disc.Op, err)
		}

		ops = append(ops, op)
	}

	return ops, nil
}

type StepType string

const (
	StepTypeHTTPRequest StepType = "http_request"
	StepTypeFilter      StepType = "filter"
	StepTypeTransform   StepType = "transform"
)

type stepDiscriminator struct {
	Type StepType `json:"type"`
}

func ParseSteps(raw json.RawMessage) ([]Fliner, error) {
	var rawSteps []json.RawMessage
	if err := json.Unmarshal(raw, &rawSteps); err != nil {
		return nil, fmt.Errorf("steps must be a JSON array: %w", err)
	}

	steps := make([]Fliner, 0, len(rawSteps))

	for i, rawStep := range rawSteps {
		var disc stepDiscriminator
		if err := json.Unmarshal(rawStep, &disc); err != nil {
			return nil, fmt.Errorf("step[%d]: missing or invalid type: %w", i, err)
		}

		var step Fliner
		var err error

		switch disc.Type {
		case StepTypeHTTPRequest:
			var s HTTPRequest
			err = json.Unmarshal(rawStep, &s)
			if err != nil {
				log.Printf("Errors unmarshalling step")
				return steps, err
			}
			s.Method = strings.ToUpper(s.Method)
			if s.TimeOut == 0 {
				s.TimeOut = 500
			}
			s.Client = &http.Client{Timeout: time.Duration(s.TimeOut) * time.Millisecond}
			step = &s

		case StepTypeTransform:
			var rawTransform struct {
				Ops json.RawMessage `json:"ops"`
			}
			if err = json.Unmarshal(rawStep, &rawTransform); err != nil {
				return nil, fmt.Errorf("step[%d]: invalid transform: %w", i, err)
			}

			var ops []Transformer
			ops, err = ParseOps(rawTransform.Ops)
			if err != nil {
				return nil, fmt.Errorf("step[%d]: %w", i, err)
			}
			step = &Transform{Ops: ops}

		default:
			return nil, fmt.Errorf("step[%d]: unknown type %q", i, disc.Type)
		}

		if err != nil {
			return nil, fmt.Errorf("step[%d] (%s): parse error: %w", i, disc.Type, err)
		}

		steps = append(steps, step)
	}

	return steps, nil
}

// func (h *HTTPRequest) GetContext()

func errJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %s", err)
	}
}

func respondJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("Error encoding response: %s", err)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	type Health struct {
		Status string `json:"status"`
	}
	respondJSON(w, http.StatusOK, Health{Status: "Running"})
}

func (ws *WorkFlowStore) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var wfReq WorkFlowInit
	if err := json.NewDecoder(r.Body).Decode(&wfReq); err != nil {
		log.Printf("Error decoding request: %v\n", err)
		errJSON(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if wfReq.Name == "" {
		errJSON(w, http.StatusBadRequest, "name is required")
		return
	}

	if len(wfReq.Steps) == 0 {
		errJSON(w, http.StatusBadRequest, "at least one step is required")
		return
	}

	steps, err := ParseSteps(wfReq.Steps)
	if err != nil {
		log.Printf("Error parsing steps: %v\n", err)
		errJSON(w, http.StatusBadRequest, fmt.Sprintf("invalid steps: %v", err))
		return
	}

	wf := NewWorkFlow()
	wf.ID = uuid.New()
	wf.Name = wfReq.Name
	wf.Enabled = true
	wf.Steps = steps

	ws.AddWorkFlow(wf)

	var rawSteps []json.RawMessage
	if err = json.Unmarshal(wfReq.Steps, &rawSteps); err != nil {
		errJSON(w, http.StatusInternalServerError, "error processing steps")
		return
	}

	respondJSON(w, http.StatusCreated, WorkFlowResponse{
		ID:      wf.ID,
		Name:    wf.Name,
		Enabled: wf.Enabled,
		Steps:   rawSteps,
	})
}

func (ws *WorkFlowStore) ToggleWorkflow(w http.ResponseWriter, r *http.Request) {
	wfID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		log.Printf("Could not parse workflow ID: %v", err)
		errJSON(w, http.StatusBadRequest, "invalid workflow ID")
		return
	}

	enabled, err := ws.FlipWorkFlow(wfID)
	if err != nil {
		errJSON(w, http.StatusNotFound, "workflow not found")
		return
	}

	log.Printf("Workflow %s toggled: enabled=%v", wfID, enabled)
	respondJSON(w, http.StatusOK, map[string]any{"enabled": enabled})
}

func (ws *WorkFlowStore) TriggerWorkflow(w http.ResponseWriter, r *http.Request) {
	wfID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		log.Printf("Could not parse workflow %s: %v", wfID.String(), err)
		errJSON(w, http.StatusBadRequest, "invalid workflow ID")
		return
	}

	wf, err := ws.GetWorkFlow(wfID)
	if err != nil {
		log.Printf("Workflow %s not found", wfID)
		errJSON(w, http.StatusNotFound, "workflow not found")
		return
	}

	if !wf.Enabled {
		log.Printf("Workflow %s is disabled", wf.ID)
		errJSON(w, http.StatusForbidden, "workflow is disabled")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var wfCtx WorkFlowCtx
	if err = json.NewDecoder(r.Body).Decode(&wfCtx); err != nil {
		wfCtx = WorkFlowCtx{}
	}

	for i, step := range wf.Steps {
		select {
		case <-ctx.Done():
			log.Printf("Workflow %s timed out before step %d", wf.ID, i)
			errJSON(w, http.StatusGatewayTimeout, "workflow timed out")
			return
		default:
		}

		wfCtx, err = step.Execute(ctx, wfCtx)
		if err != nil {
			log.Printf("Step %d failed: %v", i, err)
			errJSON(w, http.StatusInternalServerError, fmt.Sprintf("step %d failed: %v", i, err))
			return
		}
		wf.CurrentCtx = wfCtx
		b, _ := json.MarshalIndent(wf.CurrentCtx, "", "  ")
		log.Printf("[CONTEXT]:\n%s\n", b)
	}

	respondJSON(w, http.StatusOK, wfCtx)
}

func main() {
	store := NewWorkFlowStore()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", HealthCheck)
	r.Post("/workflows", store.CreateWorkflow)
	r.Post("/t/{workflowID}", store.TriggerWorkflow)
	r.Post("/workflows/{workflowID}/toggle", store.ToggleWorkflow)

	log.Printf("Server running at port:%s\n", PORT)
	err := http.ListenAndServe(":"+PORT, r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("[ERROR] error starting server: %s\n", err)
		os.Exit(1)
	}
}
