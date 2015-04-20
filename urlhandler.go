package gsweb

import "github.com/gsdocker/gslogger"

// URLHandler .
type URLHandler struct {
	gslogger.Log   // Mixin log APIs
	MethodHandlers // Mixin MethodHandlers
}

// NewURLHandler create new URLHandler
func NewURLHandler() *URLHandler {
	return &URLHandler{
		Log:            gslogger.Get("URI"),
		MethodHandlers: make(MethodHandlers),
	}
}

// URLHandle bind url's handle.
func (handler *URLHandler) URLHandle(url string, methods MethodHandlers) {

}
