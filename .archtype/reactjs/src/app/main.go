package main

import (
	"bytes"
	"path/filepath"
	"text/template"
	"time"

	"github.com/gsdocker/gsos/fs"
	"github.com/gsdocker/gsweb"
)

var comments = `
	[
		{"author": "Pete Hunt", "text": "This is one comment {{.}}"},
		{"author": "Jordan Walke", "text": "This is *another* comment {{.}}"}
	]
`

type commentQuery struct {
	tpl *template.Template
}

func newCommentQuery() *commentQuery {

	tpl, _ := template.New("comments.json").Parse(comments)

	query := &commentQuery{
		tpl: tpl,
	}

	return query
}

func (comment *commentQuery) HandleGet(context *gsweb.Context) error {

	var buff bytes.Buffer

	comment.tpl.Execute(&buff, time.Now())

	context.Response().Write(buff.Bytes())

	return nil
}

func main() {

	website := gsweb.NewWebSite()

	filehandle := gsweb.NewFileHandler()

	filehandle.RegisterPath("/", filepath.Join(fs.Current(), "static"))

	website.ChainHandle("static", filehandle)

	urihandle := gsweb.NewURIHandler()

	urihandle.Handle("/comments.json", newCommentQuery())

	website.ChainHandle("comment", urihandle)

	website.RunHTTP(":8081")
}
