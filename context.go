package gsweb

import (
	"net/http"

	"github.com/gsdocker/gserrors"
)

// MethodHandlers .
type MethodHandlers map[string]func(context *Context) error

type requestHandler struct {
	name    string
	methods MethodHandlers
}

// Context the request handler context
type Context struct {
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

		if method, ok := handler.methods[context.RequestMethod()]; ok {
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

// RequestMethod get http request's method
func (context *Context) RequestMethod() string {
	return context.request.Method
}
