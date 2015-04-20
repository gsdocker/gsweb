package test

import (
	"testing"

	"github.com/gsdocker/gsweb"
)

func handleGet(context *gsweb.Context) error {
	return context.Success()
}

func TestRun(t *testing.T) {
	webSite := gsweb.NewWebSite(":8080")

	webSite.Handle("chain0", gsweb.MethodHandlers{
		"GET": handleGet,
	})

	webSite.Run()
}
