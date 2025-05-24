package router

import (
	"strings"
)

type Group struct {
	router     *Router
	prefix     string
	middleware []Middleware
}

func (g *Group) Use(middleware ...Middleware) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *Group) Group(prefix string) *Group {
	if prefix == "" {
		prefix = "/"
	}
	prefix = strings.TrimPrefix(prefix, "/")
	fullPrefix := g.prefix + "/" + prefix
	fullPrefix = strings.TrimPrefix(fullPrefix, "/")

	return &Group{
		router: g.router,
		prefix: fullPrefix,
	}
}

func (g *Group) registerHandler(path string, wrap func(handler interface{}) interface{}, register func(fullPath string, handler interface{}), handler interface{}) {
	var fullPath string
	if path == "" {
		fullPath = g.prefix
	} else {
		fullPath = g.prefix + "/" + strings.Trim(path, "/")
	}
	fullPath = strings.Trim(fullPath, "/")

	h := handler
	for i := len(g.middleware) - 1; i >= 0; i-- {
		h = wrap(g.middleware[i](h))
	}
	register(fullPath, h)
}

func (g *Group) HandleAdvance(path string, handler AdvanceHandlerFunc) {
	g.registerHandler(
		path,
		func(h interface{}) interface{} { return h.(AdvanceHandlerFunc) },
		func(fullPath string, h interface{}) {
			g.router.HandleAdvance(fullPath, h.(AdvanceHandlerFunc))
		},
		handler,
	)
}

func (g *Group) HandleInspect(path string, handler InspectHandlerFunc) {
	g.registerHandler(
		path,
		func(h interface{}) interface{} { return h.(InspectHandlerFunc) },
		func(fullPath string, h interface{}) {
			g.router.HandleInspect(fullPath, h.(InspectHandlerFunc))
		},
		handler,
	)
}
