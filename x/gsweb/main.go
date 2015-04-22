// Package main implement gsweb cli program
// usage : gsweb {options} [golang package name]
// options
//          -run            : run target gsweb app, and watching app's file changed event.
//          -create         : create new gsweb app by golang package full name.
package main

import (
	"flag"

	"github.com/gsdocker/gslogger"
)

var log = gslogger.Get("gsweb")
var runflag = flag.Bool("run", false, "start a gsweb app indicate by golang package name")
var createflag = flag.Bool("create", false, "start a gsweb app indicate by golang package name")

func main() {

	log.I("starting gsweb damon ...")

	flag.Parse()

	if flag.NFlag() != 1 {
		log.E("invalid command line options")
		flag.PrintDefaults()
		return
	}

	// start gsweb app
	if *runflag {

	}

}
