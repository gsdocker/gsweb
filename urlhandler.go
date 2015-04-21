package gsweb

import "github.com/gsdocker/gslogger"

// URIHandler .
type URIHandler struct {
	gslogger.Log                                     // Mixin log APIs
	handlers     map[string]map[string]MethodHandler // uri handlers
}

// NewURIHandler create new URIHandler
func NewURIHandler() *URIHandler {
	return &URIHandler{
		Log:      gslogger.Get("URI"),
		handlers: make(map[string]map[string]MethodHandler),
	}
}

// HandleUnknown implement Get interface
func (uri *URIHandler) HandleUnknown(context *Context) error {

	requestMethod := context.RequestMethod()
	requestURI := context.RequestURI()

	uri.V("%s %s forward processing", requestMethod, requestURI)

	if handler, ok := uri.handlers[requestURI]; ok {
		if method, ok := handler[requestMethod]; ok {

			uri.D("%s %s handler -- found", requestMethod, requestURI)

			if err := method(context); err != nil {
				uri.E("%s %s handler execute error : %s ", requestMethod, requestURI, err)
				return context.Failed(err, "%s %s handler error : %s", requestMethod, requestURI, err)
			}

			uri.D("%s %s handler execute -- success ", requestMethod, requestURI)
		}
	}

	err := context.Forward()

	uri.V("%s %s backward processing", requestMethod, requestURI)

	return err
}

// Handle register uri handler
func (uri *URIHandler) Handle(requestURI string, handler interface{}) {
	uri.handlers[requestURI] = ExtractMethods(handler)
}
