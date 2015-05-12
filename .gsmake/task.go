package plugins

import "github.com/gsdocker/gsmake"

// TaskGsweb implement task gsweb
func TaskGsweb(context *gsmake.Runner, args ...string) error {
	context.I("hello gsweb")
	return nil
}
