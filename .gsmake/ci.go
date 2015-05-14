package plugins

import (
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gsdocker/gsconfig"
	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gsmake"
	"github.com/gsdocker/gsos"
	"github.com/gsdocker/gsos/fsnotify"
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

}

func newBuildServe(runner *gsmake.Runner, target string) (*buildServe, error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	var config GswebConfig

	if runner.PackageProperty(runner.Name(), "gsweb", &config) {
		for _, dir := range config.FSWatch {
			watcher.Add(filepath.Join(runner.StartDir(), dir), true)
		}
	}

	watcher.Add(filepath.Join(runner.StartDir(), "src", target), true)

	serve := &buildServe{
		runner:    runner,
		fswatcher: watcher,
		target:    target,
		binary:    filepath.Join(runner.StartDir(), "bin", target+gsos.ExeSuffix),
		Notify:    make(chan buildBrief, 100),
	}

	go serve.Start()

	return serve, nil
}

func (serve *buildServe) Start() {

	serve.build()

	for _ = range serve.fswatcher.Events {
		serve.build()
	}

}

func (serve *buildServe) build() {

	startTime := time.Now()

	err := serve.runner.Run("setup", serve.runner.StartDir())

	if err != nil {
		serve.runner.E("%s", err)
	}

	// calc md5
	file, err := os.Open(serve.binary)

	if err != nil {

		serve.runner.W("generate binary md5 check err :%s", err)

		return
	}

	md5h := md5.New()
	io.Copy(md5h, file)
	md5Check := md5h.Sum([]byte(""))

	if bytes.Compare(md5Check, serve.md5Check) == 0 {
		return
	}

	serve.md5Check = md5Check

	endTime := time.Now()

	brief := buildBrief{
		StartTime: startTime,
		EndTime:   endTime,
		Binary:    serve.binary,
		Md5Check:  md5Check,
	}

	select {
	case serve.Notify <- brief:
	default:
		serve.runner.W("max notify event queue size reach !!!")
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

	retryTimeout := gsconfig.Seconds("app.runner.retry_timeout", 5)

	done := make(chan *exec.Cmd)

	var cmd *exec.Cmd

	var brief buildBrief

	for {
		// wait retry timeout or kill command
		select {
		case <-time.After(retryTimeout):
			if cmd != nil {
				continue
			}

		case currentCmd := <-done:
			{
				runner.I(
					"app process %d exit => systime: %s usertime: %s",
					currentCmd.Process.Pid,
					currentCmd.ProcessState.SystemTime(),
					currentCmd.ProcessState.UserTime(),
				)

				if cmd == currentCmd {
					cmd = nil
				}

				continue
			}
		case build := <-build.Notify:
			if cmd != nil && cmd.Process != nil {

				for {

					err := cmd.Process.Kill()

					if err == nil {
						break
					}

					runner.E("kill app -- failed\n%s", err)

					<-time.After(retryTimeout)
				}

				cmd = nil
			}

			brief = build
		}

		if brief.Binary == "" {
			continue
		}

		currentCmd := exec.Command(brief.Binary, runner.StartDir())

		currentCmd.Stdout = os.Stdout
		currentCmd.Stderr = os.Stderr

		err := currentCmd.Start()

		if err != nil {
			runner.E("start app -- failed\n%s", err)
			continue
		}

		cmd = currentCmd

		runner.I("app process %d started", cmd.Process.Pid)

		go func() {
			if err := currentCmd.Wait(); err != nil {
				runner.W("app exit with error :\n\t%s", err)
			}

			done <- currentCmd
		}()
	}

}
