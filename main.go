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
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

const PORT = "1324"

// type Hook struct {
// 	HookCtx     string `json:"context"`
// 	TriggerPath uuid.UUID `json:"workflow"`
// }

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

type HTTPRequest struct {
	Endpoint string            `json:"URL"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	TimeOut  int               `json:"timeout"`
	Retries  int               `json:"retries"`
	Client   *http.Client
}

func (h *HTTPRequest) Execute(ctx WorkFlowCtx) (WorkFlowCtx, error) {
	log.Printf("Executing Workflow\n")
	if h.Method == "GET" {
		log.Printf("GET on %s", h.Endpoint)
		req, err := http.NewRequest("GET", h.Endpoint, nil)
		if err != nil {
			log.Printf("Error sending request: %v", err.Error())
			return ctx, err
		}

		for k, v := range h.Headers {
			req.Header.Add(k, v)
		}

		res, err := h.Client.Do(req)
		if err != nil {
			log.Printf("Error getting response: %v", err.Error())
			return ctx, err
		}

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		var parsed WorkFlowCtx
		if err = json.Unmarshal(body, &parsed); err != nil {
			ctx["Result"] = string(body)
		}
		ctx["Result"] = parsed

		return ctx, nil
	}

	return ctx, errors.New("Method not defined or incorrect format")
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

func (wf *WorkFlow) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var wfReq WorkFlowInit
	err := json.NewDecoder(r.Body).Decode(&wfReq)
	if err != nil {
		log.Printf("Error creating workflow: %v\n", err.Error())
		errJSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	wf.ID = uuid.New()
	wf.Name = wfReq.Name
	wf.Enabled = true
	wf.Steps, err = ParseSteps(wfReq.Steps)
	if err != nil {
		log.Printf("Error parsing steps: %v\n", err.Error())
		errJSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	respondJSON(w, http.StatusCreated, wf)
}

func (wf *WorkFlow) DisableWorkflow(w http.ResponseWriter, r *http.Request) {
	wf.Enabled = false
	respondJSON(w, http.StatusOK, wf)
}

func (wf *WorkFlow) TriggerWorkflow(w http.ResponseWriter, r *http.Request) {
	wfID, err := uuid.Parse(chi.URLParam(r, "workflowID"))
	if err != nil {
		log.Printf("Could not parse workflow ID: %v", err)
		errJSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if wfID != wf.ID {
		log.Printf("Incorrect trigger path: %s", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		errJSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if !wf.Enabled {
		log.Printf("WorkFlow ID: %s disabled", wf.ID)
		errJSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	if r.Method != "POST" {
		log.Printf("Wrong HTTP trigger method")
		errJSON(w, http.StatusMethodNotAllowed, "Something went wrong")
		return
	}

	var ctx WorkFlowCtx
	if err = json.NewDecoder(r.Body).Decode(&ctx); err != nil {
		log.Printf("Wrong body form")
		errJSON(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	wf.CurrentCtx = ctx
	for _, fline := range wf.Steps {
		ctx, err = fline.Execute(wf.CurrentCtx)
		if err != nil {
			log.Printf("Error executing step: \n%s\n", err.Error())
			errJSON(w, http.StatusInternalServerError, "Something went wrong")
			return
		}
		wf.CurrentCtx = ctx
	}

	log.Printf("Current Context:\n%v\n", ctx)
	respondJSON(w, http.StatusOK, wf.CurrentCtx)
}

func main() {
	wf := NewWorkFlow()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", HealthCheck)
	r.Post("/create", wf.CreateWorkflow)
	r.Post("/t/{workflowID}", wf.TriggerWorkflow)

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
