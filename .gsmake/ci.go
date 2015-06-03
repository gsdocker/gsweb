package plugins

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gsos/fs"
	"github.com/gsdocker/gsos/fsnotify"
	"github.com/gsmake/gsmake"
	"github.com/gsmake/gsmake/property"
	"github.com/gsmake/gsmake/vfs"
)

// GswebConfig .
type GswebConfig struct {
	FSWatch []string // watch directories
}

// CompileSession .
type buildBrief struct {
	StartTime time.Time // Current build session start time
	EndTime   time.Time // Current build session end time
	Md5Check  []byte    // Current build session product md5 check value
	Binary    string    // Current build session product file path
}

type buildServe struct {
	runner    *gsmake.Runner      // gstask runner
	fswatcher *fsnotify.FSWatcher // The source file watcher
	target    string              // build target
	binary    string              // binary path
	md5Check  []byte              // binary md5
	Notify    chan buildBrief     // compile session notify
	cmd       *exec.Cmd           // running command

}

func newBuildServe(runner *gsmake.Runner, target string) (*buildServe, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	var config GswebConfig

	err = runner.Property("golang", runner.Name(), "gsweb", &config)

	if err != nil && !vfs.NotFound(err) && !property.NotFound(err) {
		return nil, err
	}

	for _, dir := range config.FSWatch {
		watcher.Add(filepath.Join(runner.StartDir(), dir), true)
	}

	watcher.Add(filepath.Join(runner.StartDir(), "src", target), true)

	serve := &buildServe{
		runner:    runner,
		fswatcher: watcher,
		target:    target,
		binary:    filepath.Join(runner.StartDir(), "bin", target+fs.ExeSuffix),
		Notify:    make(chan buildBrief, 100),
	}

	return serve, nil
}

func (serve *buildServe) Start() error {

	if err := serve.build(); err != nil {
		return err
	}

	serve.start()

	for _ = range serve.fswatcher.Events {

		serve.kill()

		if err := serve.build(); err != nil {
			return err
		}

		serve.start()
	}

	return nil
}

func (serve *buildServe) kill() {
	if serve.cmd != nil {

		for {
			serve.runner.I("kill process %d ...", serve.cmd.Process.Pid)

			err := serve.cmd.Process.Kill()

			if err != nil {
				serve.runner.W("kill process %d error\n%s", serve.cmd.Process.Pid, err)
				<-time.After(time.Second * 5)
				continue
			}

			break
		}

	}

}

func (serve *buildServe) build() error {

	startTime := time.Now()

	err := serve.runner.Run("compile", serve.runner.StartDir())

	if err != nil {
		return err
	}

	serve.runner.I("build times %v", time.Now().Sub(startTime))

	return nil
}

func (serve *buildServe) start() {

	serve.cmd = exec.Command(serve.binary, serve.runner.StartDir())

	serve.cmd.Stdout = os.Stdout
	serve.cmd.Stderr = os.Stderr

	err := serve.cmd.Start()

	if err != nil {
		serve.runner.E("exec :%s\n\terr:%s", serve.binary, err)
	}
}

// TaskGsweb implement task gsweb
func TaskGsweb(runner *gsmake.Runner, args ...string) error {

	if len(args) == 0 {
		return gserrors.Newf(nil, "expect serve name")
	}

	build, err := newBuildServe(runner, args[0])

	if err != nil {
		return err
	}

	return build.Start()
}
