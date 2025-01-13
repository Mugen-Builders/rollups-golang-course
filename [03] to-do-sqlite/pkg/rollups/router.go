package rollups

import (
	"fmt"
)

type AdvanceHandlerFunc func(data *AdvanceResponse) error

type Router struct {
	AdvanceHandlers map[string]AdvanceHandlerFunc
}

func NewRouter() *Router {
	return &Router{
		AdvanceHandlers: make(map[string]AdvanceHandlerFunc),
	}
}

func (r *Router) HandleAdvance(path string, handler AdvanceHandlerFunc) {
	r.AdvanceHandlers[path] = handler
}

func (r *Router) Advance(data *AdvanceResponse) error {
	handler, ok := r.AdvanceHandlers[data.Path]
	if !ok {
		return fmt.Errorf("path '%s' not found", data.Path)
	}
	return handler(data)
}
