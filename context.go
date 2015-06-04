package gsweb

import (
	"net/http"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslogger"
)

// Context the request handler context
type Context struct {
	gslogger.Log
	Router         *Router             // router belongs
	responseWriter http.ResponseWriter // response writer
	request        *http.Request       // request
	forwardCursor  int                 // The forward chain cursor
}

func newContext(
	router *Router,
	request *http.Request,
	response http.ResponseWriter) *Context {

	return &Context{
		Log:            router.Log,
		Router:         router,
		request:        request,
		responseWriter: response,
		forwardCursor:  0,
	}
}

// Forward forward request to next handler
func (context *Context) Forward() error {

	for len(context.Router.handleChain) > context.forwardCursor {
		cursor := context.forwardCursor
		context.forwardCursor++

		handler := context.Router.handleChain[cursor]

		method, ok := handler.methods[context.RequestMethod()]

		if !ok {

			method, ok = handler.methods["UNKNOWN"]
		}

		if ok {
			err := method(context)

			gserrors.Assert(
				context.forwardCursor == len(context.Router.handleChain),
				"handler method %s#%s must call context.Success or context.Failed before return",
				handler.name,
				context.RequestMethod(),
			)

			return err
		}
	}

	return nil
}

// Failed break the request handler chain processing and return error
func (context *Context) Failed(err error, fmt string, args ...interface{}) error {
	context.forwardCursor = len(context.Router.handleChain) // request handler cotract check codes
	return gserrors.Newf(err, fmt, args...)
}

// Success break the request handler chain processing and return success
func (context *Context) Success() error {
	context.forwardCursor = len(context.Router.handleChain) // request handler cotract check codes
	return nil
}

// Request  get response writer
func (context *Context) Request() *http.Request {
	return context.request
}

// Response  get response writer
func (context *Context) Response() http.ResponseWriter {
	return context.responseWriter
}

// RequestMethod get http request's method
func (context *Context) RequestMethod() string {
	return context.request.Method
}

// RequestURI get http request's URI
func (context *Context) RequestURI() string {
	return context.request.URL.Path
}

// Redirect redirect url
func (context *Context) Redirect(urlStr string, code int) {
	http.Redirect(context.responseWriter, context.request, urlStr, code)
}
