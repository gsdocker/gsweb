package gsweb

import (
	"net/http"

	"github.com/gsdocker/gslogger"
)

// Handler the http process handler
type Handler struct {
	name    string                   // The handler name
	methods map[string]MethodHandler // methods
}

// Router resource router
type Router struct {
	gslogger.Log            // Mixin log APIs
	handleChain  []*Handler // request handle chain
	started      bool       // the gsweb state flag
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

// ChainHandle register request handle chain node named by parameter name
// the method is not thread safe,so call it before calling WebSite#Run method
func (router *Router) ChainHandle(name string, handler interface{}) {
	router.handleChain = append(router.handleChain, &Handler{name: name, methods: ExtractMethods(handler)})
}
