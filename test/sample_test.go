package test

import (
	"testing"

	"github.com/gsdocker/gsweb"
)

type testHandler struct {
}

func (handler *testHandler) HandleGet(context *gsweb.Context) error {
	//context.Response().WriteHeader(200)
	context.Response().Write([]byte("hello"))
	return context.Success()
}

func TestHandler(t *testing.T) {
	webSite := gsweb.NewWebSite()

	// uri := gsweb.NewURIHandler()
	// uri.Handle("/", &testHandler{})
	// webSite.ChainHandle("uri", uri)

	file := gsweb.NewFileHandler()

	path, _ := file.RegisterPath("/", "./")

	path.EnableGetDir(true)

	webSite.ChainHandle("staticfiles", file)

	go webSite.RunHTTP(":8080")

	webSite.RunHTTPS(":4343", "./cert.pem", "./key.pem")
}
