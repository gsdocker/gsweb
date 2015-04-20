package gsweb

import "reflect"

// MethodGet .
type MethodGet interface {
	HandleGet(context *Context) error
}

// MethodPut .
type MethodPut interface {
	HandlePut(context *Context) error
}

// MethodPost .
type MethodPost interface {
	HandlePost(context *Context) error
}

// MethodDelete .
type MethodDelete interface {
	HandleDelete(context *Context) error
}

var registerMethodTypes = map[string]reflect.Type{}

func init() {
	var get MethodGet
	registerMethodTypes["GET"] = reflect.TypeOf(get)
	var put MethodPut
	registerMethodTypes["PUT"] = reflect.TypeOf(put)
	var post MethodPost
	registerMethodTypes["PUT"] = reflect.TypeOf(post)
	var del MethodDelete
	registerMethodTypes["PUT"] = reflect.TypeOf(del)
}

// HTTPMethod get http method interface's type
func HTTPMethod(name string) (t reflect.Type, ok bool) {

	t, ok = registerMethodTypes[name]

	return
}

// NewHTTPMethod set http method interface's type
func NewHTTPMethod(name string, t interface{}) {
	registerMethodTypes[name] = reflect.TypeOf(t)
}
