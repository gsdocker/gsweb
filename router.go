package gsweb

import (
	"net/http"

	"github.com/gsdocker/gslogger"
)

// Controller The resource controller object
type Controller interface {
}

// Router resource router
type Router struct {
	gslogger.Log // Mixin log APIs
}

func newRouter() *Router {
	return &Router{
		Log: gslogger.Get("router"),
	}
}

// ServeHTTP implement http handler
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.D("handle request")
}
