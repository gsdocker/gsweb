package main

import (
	"path/filepath"

	"github.com/gsdocker/gsos/fs"
	"github.com/gsdocker/gsweb"
)

func main() {

	website := gsweb.NewWebSite()

	filehandle := gsweb.NewFileHandler()

	filehandle.RegisterPath("/", filepath.Join(fs.Current(), "static"))

	website.ChainHandle("static", filehandle)

	website.RunHTTP(":8080")
}
