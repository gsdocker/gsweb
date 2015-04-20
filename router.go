package gsweb

import (
	"net/http"

	"github.com/gsdocker/gslogger"
)

// Router resource router
type Router struct {
	gslogger.Log                   // Mixin log APIs
	handleChain  []*requestHandler // request handle chain
	started      bool              // the gsweb state flag
}

func newRouter() *Router {
	return &Router{
		Log: gslogger.Get("router"),
	}
}

// ServeHTTP implement http handler
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// empty handleChain optimize
	if len(router.handleChain) == 0 {
		router.W("gsweb empty handle chain warning !!!!!! ")
		return
	}

	context := newContext(router, r, w)
	err := context.Forward()

	if err != nil {
		router.E(
			"handle request err :\n\tfrom:%s\n\trequest-uri:%s\n\terr:%s",
			r.RemoteAddr,
			r.RequestURI,
			err,
		)
	}
}

// Handle register request handle chain node named by parameter name
// the method is not thread safe,so call it before calling WebSite#Run method
func (router *Router) Handle(name string, methods MethodHandlers) {

	router.handleChain = append(router.handleChain, &requestHandler{name: name, methods: methods})
}
