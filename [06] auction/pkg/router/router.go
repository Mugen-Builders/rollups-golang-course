package router

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rollmelette/rollmelette"
)

type AdvanceHandlerFunc func(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error

type InspectHandlerFunc func(env rollmelette.EnvInspector, payload []byte) error

type Router struct {
	advanceHandlers map[string]AdvanceHandlerFunc
	inspectHandlers map[string]InspectHandlerFunc
	middlewares     []Middleware
}

func NewRouter() *Router {
	return &Router{
		advanceHandlers: make(map[string]AdvanceHandlerFunc),
		inspectHandlers: make(map[string]InspectHandlerFunc),
		middlewares:     make([]Middleware, 0),
	}
}

func (r *Router) Use(middleware ...Middleware) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) Group(prefix string) *Group {
	return &Group{
		router: r,
		prefix: prefix,
	}
}

func (r *Router) HandleAdvance(path string, handler AdvanceHandlerFunc) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler).(AdvanceHandlerFunc)
	}

	r.advanceHandlers[path] = handler
}

func (r *Router) HandleInspect(path string, handler InspectHandlerFunc) {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler).(InspectHandlerFunc)
	}

	r.inspectHandlers[path] = handler
}

type Request struct {
	Path string          `json:"path"`
	Data json.RawMessage `json:"data"`
}

func parseRequestRawPayload(payload []byte) (*Request, error) {
	var req Request
	if err := json.Unmarshal(payload, &req); err != nil {
		return nil, fmt.Errorf("invalid request format: %v", err)
	}
	if len(req.Data) == 0 {
		fmt.Printf("[parseRequestRawPayload] Empty payload: %s\n", string(payload))
		return nil, fmt.Errorf("empty payload")
	}
	return &req, nil
}

func (r *Router) Advance(env rollmelette.Env, metadata rollmelette.Metadata, deposit rollmelette.Deposit, payload []byte) error {
	req, err := parseRequestRawPayload(payload)
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	path := strings.Trim(req.Path, "/")
	handler, exists := r.advanceHandlers[path]
	if !exists {
		return fmt.Errorf("no handler found for path: %s", path)
	}

	return handler(env, metadata, deposit, req.Data)
}

func (r *Router) Inspect(env rollmelette.EnvInspector, payload []byte) error {
	req, err := parseRequestRawPayload(payload)
	if err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	path := strings.Trim(req.Path, "/")
	handler, exists := r.inspectHandlers[path]
	if !exists {
		return fmt.Errorf("no handler found for path: %s", path)
	}

	return handler(env, req.Data)
}
