// Package toolkit The gsweb app's build/test library tools
package toolkit

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gsdocker/gsconfig"
	"github.com/gsdocker/gslogger"
)

// AppRunner the gsweb app's test container object
type AppRunner struct {
	gslogger.Log `json:"-"`   // Mixin Log APIs
	app          *App         // The app config object
	compileS     *AppCompileS // The app compile service object
	binary       string       // The app binary path
	Md5Check     []byte       // The app's md5 check
	BuildTime    time.Time    // The app build time
	StartTime    time.Time    // The app start time
	Pid          int          // The started app process id
}

// NewAppRunner .
func NewAppRunner(app *App) (*AppRunner, error) {

	cs, err := NewAppCompileS(app)

	if err != nil {
		return nil, err
	}

	runner := &AppRunner{
		Log:      gslogger.Get("AppRunner"),
		app:      app,
		compileS: cs,
	}

	runner.compileS.Start()

	return runner, nil
}

// Run Start App runner
func (runner *AppRunner) Run() {

	retryTimeout := gsconfig.Seconds("app.runner.retry_timeout", 5)

	done := make(chan *exec.Cmd)

	var cmd *exec.Cmd

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
		case build := <-runner.compileS.Notify:
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

			runner.BuildTime = build.StartTime
			runner.Md5Check = build.Md5Check
			runner.binary = build.Binary
		}

		if runner.binary == "" {
			continue
		}

		currentCmd := exec.Command(runner.binary, runner.app.SrcPath())

		currentCmd.Stdout = os.Stdout
		currentCmd.Stderr = os.Stderr

		err := currentCmd.Start()

		if err != nil {
			runner.E("start app -- failed\n%s", err)
			continue
		}

		cmd = currentCmd

		runner.Pid = cmd.Process.Pid

		file, err := os.Create(filepath.Join(runner.compileS.BuildDir(), "runner.json"))

		if err == nil {
			err = json.NewEncoder(file).Encode(runner)
			file.Close()
		}

		if err != nil {
			runner.W("marshal runner info err :%s", err)
		}

		go func() {
			if err := currentCmd.Wait(); err != nil {
				runner.W("app exit with error :\n\t%s", err)
			}

			done <- currentCmd
		}()
	}

}
