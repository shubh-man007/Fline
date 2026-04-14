package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const PORT = "1324"

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

type WorkFlowDisable struct {
	ID uuid.UUID `json:"id"`
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
	TimeOut  int               `json:"timeout"`
	Retries  int               `json:"retries"`
	Client   *http.Client
}

func (h *HTTPRequest) Execute(ctx context.Context, wfCtx WorkFlowCtx) (WorkFlowCtx, error) {
	log.Printf("Executing HTTP step: %s %s\n", h.Method, h.Endpoint)

	req, err := http.NewRequestWithContext(ctx, h.Method, h.Endpoint, nil)
	if err != nil {
		return wfCtx, fmt.Errorf("error building request: %w", err)
	}

	for k, v := range h.Headers {
		req.Header.Add(k, v)
	}

	res, err := h.Client.Do(req)
	if err != nil {
		return wfCtx, fmt.Errorf("error executing request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return wfCtx, fmt.Errorf("error reading response body: %w", err)
	}

	var parsed WorkFlowCtx
	if err = json.Unmarshal(body, &parsed); err != nil {
		wfCtx["result"] = string(body)
	} else {
		wfCtx["result"] = parsed
	}

	return wfCtx, nil
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
		// case StepTypeFilter:
		// 	var s FilterStep
		// 	err = json.Unmarshal(rawStep, &s)
		// 	step = &s
		// case StepTypeTransform:
		// 	var s TransformStep
		// 	err = json.Unmarshal(rawStep, &s)
		// 	step = &s
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

func (ws *WorkFlowStore) DisableWorkflow(w http.ResponseWriter, r *http.Request) {
	wfID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		log.Printf("Could not parse workflow ID: %v", err)
		errJSON(w, http.StatusBadRequest, "invalid workflow ID")
		return
	}

	wf, err := ws.GetWorkFlow(wfID)
	if err != nil {
		log.Printf("Workflow: %s not found", wfID)
		errJSON(w, http.StatusNotFound, "workflow not found")
		return
	}

	wf.Enabled = false
	log.Printf("Workflow %s is disabled", wf.ID)
	respondJSON(w, http.StatusOK, "workflow disabled")
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
		fmt.Printf("Current Context: %v\n", wf.CurrentCtx)
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
	r.Patch("/workflows/{workflowID}/disable", store.DisableWorkflow)

	log.Printf("Server running at port:%s\n", PORT)
	err := http.ListenAndServe(":"+PORT, r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed\n")
	} else if err != nil {
		log.Printf("[ERROR] error starting server: %s\n", err)
		os.Exit(1)
	}
}

// {
//   "name": "Random Anime Quote",
//   "enabled": true,
//   "trigger": {
//     "type": "http",
//     "path": "/t/83064af3-bb81-4514-a6d4-afba340825cd"
//   },
//   "steps": [
//     {
//       "type": "http_request",
//       "method": "GET",
//       "url": "https://api.animechan.io/v1/quotes/random",
//       "headers": {
//         "Accept": "application/json"
//       },
//       "body": {
//         "mode": "ctx"
//       },
//       "timeoutMs": 5000,
//       "retries": 3
//     }
//   ]
// }

// https://animechan.io/docs/quote/random-via-anime
// BASE_URL: https://api.animechan.io/v1   {animechan}
// curl -v https://api.animechan.io/v1/quotes/random

// {
//   "status": "success",
//   "data": {
//     "content": "I told you before, Komamura. The only paths that I see with these eyes are the ones not dyed with blood. Those paths are the paths to justice. So whichever path I choose...Is justice.",
//     "anime": {
//       "id": 222,
//       "name": "Bleach",
//       "altName": "Bleach"
//     },
//     "character": {
//       "id": 2143,
//       "name": "Tousen Kaname"
//     }
//   }
// }

// curl -X POST "http://localhost:1324/workflows" -H "Content-Type: application/json" -d '{"name":"Random Anime Quote","steps":[{"type":"http_request","method":"GET","URL":"https://api.animechan.io/v1/quotes/random","headers":{"Accept":"application/json"},"timeout":5000,"retries":3}]}'
// curl -X POST http://localhost:1324/t/fa35033d-2843-4b69-a6e1-9b0df5326a48 -H "Content-Type: application/json" -d '{}'
