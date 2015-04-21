package test

import (
	"net/http"
	"testing"

	"github.com/gsdocker/gsweb"
)

type testHandler struct {
}

func (handler *testHandler) HandleGet(context *gsweb.Context) error {
	context.Redirect("http://baidu.com", http.StatusFound)
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

func TestRedirect(t *testing.T) {
	webSite := gsweb.NewWebSite()

	uri := gsweb.NewURIHandler()
	uri.Handle("/", &testHandler{})
	webSite.ChainHandle("uri", uri)

	webSite.RunHTTP(":8080")
}
