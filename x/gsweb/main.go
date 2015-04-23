// Package main implement gsweb cli program
// usage : gsweb {options} [golang package name]
// options
//          -run            : run target gsweb app, and watching app's file changed event.
//          -create         : create new gsweb app by golang package full name.
package main

import (
	"flag"

	"github.com/gsdocker/gslogger"
	"github.com/gsdocker/gsweb/toolkit"
)

var log = gslogger.Get("gsweb")
var runflag = flag.Bool("run", false, "start a gsweb app indicate by golang package name")
var createflag = flag.String("create", "", "create a new gsweb app base on template app")

func main() {

	log.I("starting gsweb damon ...")

	flag.Parse()

	if flag.NFlag() != 1 {
		log.E("invalid command line options")
		flag.PrintDefaults()
		return
	}

	if flag.NArg() != 1 {
		log.E("invalid command line options")
		flag.PrintDefaults()
		return
	}

	var app *toolkit.App

	var err error

	// start gsweb app
	if *runflag {
		app, err = toolkit.LoadApp(flag.Arg(0))

		if err != nil {
			log.E("can't load app\n\tapp:%s\n\terr:%s", flag.Arg(0), err)
			return
		}

	} else if *createflag != "" {
		app, err = toolkit.CreateApp(flag.Arg(0), *createflag, nil)

		if err != nil {
			log.E("can't load app\n\tapp:%s\n\terr:%s", flag.Arg(0), err)
			return
		}
	}

	runner, err := toolkit.NewAppRunner(app)

	if err != nil {
		log.E("create app runner -- failed\n\tapp:%s\n\terr:%s", flag.Arg(0), err)
		return
	}

	runner.Run()

}
