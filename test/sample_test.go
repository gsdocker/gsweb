package test

import (
	"testing"

	"github.com/gsdocker/gsweb"
)

func TestRun(t *testing.T) {
	webSite := gsweb.NewWebSite(":8080")
	webSite.Run()
}
