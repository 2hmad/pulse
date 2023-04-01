package routing

import (
	"github.com/valyala/fasthttp"
)

type Handler func(ctx *Context) error

type Router struct {
	routes          map[string]*route
	stores          map[string]RouteStore
	notFoundHandler Handler
	maxParams       int
}

func CombineHandlers(handlers ...Handler) Handler {
	return func(ctx *Context) error {
		for _, handler := range handlers {
			if err := handler(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func New() *Router {
	return &Router{
		routes: make(map[string]*route),
		stores: make(map[string]RouteStore),
	}
}

func (r *Router) add(method, path string, handlers []Handler) {
	s := r.stores[method]
	if s == nil {
		s = newStore()
		r.stores[method] = s
	}
	if n := s.Add(path, handlers); n > r.maxParams {
		r.maxParams = n
	}
	r.routes[method+path] = &route{method, path, CombineHandlers(handlers...)}
}

func (r *Router) find(method, path string) []Handler {
	s := r.stores[method]
	if s == nil {
		return []Handler{}
	}
	handlers, ok := s.Get(path).([]Handler)
	if !ok {
		return []Handler{}
	}
	return handlers
}

func RouterHandler(router *Router) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		method := string(ctx.Method())
		handlers := router.find(method, path)
		if handlers == nil {
			ctx.Error("Not found", fasthttp.StatusNotFound)
			return
		}
		c := NewContext(ctx)
		for _, h := range handlers {
			err := h(c)
			if err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
				return
			}
		}
	}
}
