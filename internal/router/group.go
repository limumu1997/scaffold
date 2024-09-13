package router

import (
	"net/http"
	"path"
)

type RouteGroup struct {
    prefix string
    mux    *http.ServeMux
}

func NewRouteGroup(mux *http.ServeMux, prefix string) *RouteGroup {
    return &RouteGroup{
        prefix: prefix,
        mux:    mux,
    }
}

func (g *RouteGroup) Handle(pattern string, handler http.HandlerFunc) {
    g.mux.HandleFunc(path.Join(g.prefix, pattern), handler)
}