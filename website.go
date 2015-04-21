package gsweb

import (
	"net/http"
	"time"

	"github.com/gsdocker/gsconfig"
	"github.com/gsdocker/gslogger"
)

// WebSite The website object
type WebSite struct {
	gslogger.Log // Mixin log apis
	*Router      // Minx Router
}

// NewWebSite create new gsweb instance
func NewWebSite() *WebSite {
	return &WebSite{
		Log:    gslogger.Get("gsweb"),
		Router: newRouter(),
	}

}

// RunHTTP start listen in connection and run dispatch loop
func (website *WebSite) RunHTTP(laddr string) {

	for {
		website.D("start http server : %s", laddr)

		server := http.Server{
			Addr:           laddr,
			ReadTimeout:    gsconfig.Seconds("read_timeout", 10),
			WriteTimeout:   gsconfig.Seconds("writet_imeout", 10),
			MaxHeaderBytes: gsconfig.Int("max_header_bytes", 1<<20),
		}

		server.Handler = website

		err := server.ListenAndServe()

		if err != nil {
			website.E("start http err :%s", err)

			timeout := gsconfig.Seconds("retry_timeout", 5)

			website.E("retry start http server %v later", timeout)

			<-time.After(timeout)
		}

	}

}

// RunHTTPS start listen in connection and run dispatch loop
func (website *WebSite) RunHTTPS(laddr string, certfile string, keyfile string) {

	for {
		website.D("start https server : %s", laddr)

		server := http.Server{
			Addr:           laddr,
			ReadTimeout:    gsconfig.Seconds("read_timeout", 10),
			WriteTimeout:   gsconfig.Seconds("writet_imeout", 10),
			MaxHeaderBytes: gsconfig.Int("max_header_bytes", 1<<20),
		}

		server.Handler = website

		err := server.ListenAndServeTLS(certfile, keyfile)

		if err != nil {
			website.E("start https err :%s", err)

			timeout := gsconfig.Seconds("retry_timeout", 5)

			website.E("retry start https server %v later", timeout)

			<-time.After(timeout)
		}

	}

}
