package gsweb

import (
	"net/http"
	"time"

	"github.com/gsdocker/gsconfig"
	"github.com/gsdocker/gslogger"
)

// WebSite The website object
type WebSite struct {
	gslogger.Log             // Mixin log apis
	*Router                  // Minx Router
	server       http.Server // http server
}

// NewWebSite create new gsweb instance
func NewWebSite(laddr string) *WebSite {
	gsweb := &WebSite{
		Log:    gslogger.Get("gsweb"),
		Router: newRouter(),
		server: http.Server{
			Addr:           laddr,
			ReadTimeout:    gsconfig.Seconds("read_timeout", 10),
			WriteTimeout:   gsconfig.Seconds("writet_imeout", 10),
			MaxHeaderBytes: gsconfig.Int("max_header_bytes", 1<<20),
		},
	}

	gsweb.server.Handler = gsweb

	return gsweb
}

// Run start listen in connection and run dispatch loop
func (website *WebSite) Run() {

	website.D("start http server : %s", website.server.Addr)

	err := website.server.ListenAndServe()

	if err != nil {
		website.E("http listen err :%s", err)
		time.AfterFunc(gsconfig.Seconds("retry_timeout", 5), website.Run)
		return
	}
}
