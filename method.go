package gsweb

// Get .
type Get interface {
	HandleGet(context *Context) error
}

// Put .
type Put interface {
	HandlePut(context *Context) error
}

// Post .
type Post interface {
	HandlePost(context *Context) error
}

// Delete .
type Delete interface {
	HandleDelete(context *Context) error
}

// Unknown .
type Unknown interface {
	HandleUnknown(context *Context) error
}

// MethodHandler .
type MethodHandler func(context *Context) error

// MethodExtractor .
type MethodExtractor func(interface{}) (func(context *Context) error, bool)

var methodExtractors = map[string]MethodExtractor{
	"GET": func(handler interface{}) (func(context *Context) error, bool) {
		if get, ok := handler.(Get); ok {
			return get.HandleGet, true
		}

		return nil, false
	},

	"PUT": func(handler interface{}) (func(context *Context) error, bool) {
		if get, ok := handler.(Put); ok {
			return get.HandlePut, true
		}

		return nil, false
	},

	"POST": func(handler interface{}) (func(context *Context) error, bool) {
		if get, ok := handler.(Post); ok {
			return get.HandlePost, true
		}

		return nil, false
	},

	"DELETE": func(handler interface{}) (func(context *Context) error, bool) {
		if get, ok := handler.(Delete); ok {
			return get.HandleDelete, true
		}

		return nil, false
	},

	"UNKNOWN": func(handler interface{}) (func(context *Context) error, bool) {
		if get, ok := handler.(Unknown); ok {
			return get.HandleUnknown, true
		}

		return nil, false
	},
}

//ExtractMethods .
func ExtractMethods(handler interface{}) map[string]MethodHandler {
	handlers := make(map[string]MethodHandler)

	for k, v := range methodExtractors {
		if methodHandler, ok := v(handler); ok {
			handlers[k] = methodHandler
		}
	}

	return handlers
}

// HTTPMethod register customer http method
func HTTPMethod(name string, extractor MethodExtractor) {
	methodExtractors[name] = extractor
}
