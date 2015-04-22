package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gslogger"
	"github.com/howeyc/fsnotify"
)

// The app runner errors
var (
	ErrAppRunner = errors.New("gsweb app runner error")
)

type eventType int

const (
	filesChanged eventType = iota
)

// AppConfig The AppRunner's app config object
type AppConfig struct {
	Pid int // running app pid
}

// AppRunner the gsweb app runner
type AppRunner struct {
	gslogger.Log                   // mixin log APIs
	packageName  string            // the gsweb app's golang full package name
	fswatcher    *fsnotify.Watcher // the app files watcher
	tempDir      string            // the template directory
	appPath      string            // the built app's path
	appSrcPath   string            // the app source path
	events       chan eventType    // event channel
	appConfig    AppConfig         // app config
}

func newAppRunner(packageName string) (*AppRunner, error) {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	runner := &AppRunner{
		Log:         gslogger.Get("gsweb"),
		packageName: packageName,
		fswatcher:   watcher,
		tempDir:     filepath.Join(os.TempDir(), packageName),
		events:      make(chan eventType),
	}

	runner.appPath = filepath.Join(runner.tempDir, "app")

	// try get packagePath
	GOPATH := os.Getenv("GOPATH")

	if GOPATH == "" {
		return nil, gserrors.Newf(ErrAppRunner, "must set GOPATH first")
	}

	var packagePaths []string

	goPath := strings.Split(GOPATH, string(os.PathListSeparator))

	for _, path := range goPath {
		fullpath := filepath.Join(path, "src", packageName)
		fi, err := os.Stat(path)

		if err == nil && fi.IsDir() {
			packagePaths = append(packagePaths, fullpath)
		}
	}

	if len(packagePaths) > 1 {
		var stream bytes.Buffer

		stream.WriteString(fmt.Sprintf("found more than one package named :%s", packageName))

		for i, path := range packagePaths {
			stream.WriteString(fmt.Sprintf("\n\t%d) %s", i, path))
		}

		return nil, gserrors.Newf(ErrAppRunner, stream.String())
	}

	runner.appSrcPath = packagePaths[0]

	return runner, nil
}

// Run start app runner
func (runner *AppRunner) Run() {
	runner.events <- filesChanged

	runner.eventLoop()
}

func (runner *AppRunner) eventLoop() {
	for event := range runner.events {
		switch event {
		case filesChanged:
			runner.I("compile gsweb app : %s", runner.packageName)
			// first compile app
			if err := runner.compileApp(); err != nil {
				runner.E("compile gsweb app -- failed :\n\tapp: %s\n\terror: %s", runner.packageName, err)
				continue
			}

			runner.I("compile gsweb app -- success")

			runner.I("start gsweb app : %s", runner.packageName)

			if err := runner.startApp(); err != nil {
				runner.E("start gsweb app -- failed :\n\tapp: %s\n\terror: %s", runner.packageName, err)
				continue
			}

			runner.I("start gsweb app -- success")
		}
	}
}

func (runner *AppRunner) compileApp() error {

	cmd := exec.Command("go build -o %s %s", runner.appPath, runner.packageName)

	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (runner *AppRunner) stopApp() error {
	file, err := ioutil.ReadFile(filepath.Join(runner.tempDir, "runner.json"))

	if err != nil {
		return err
	}

	json.Unmarshal(file, &runner.appConfig)

	return nil
}

func (runner *AppRunner) startApp() error {

	runner.stopApp()

	cmd := exec.Command("%s %s", runner.appPath, runner.appSrcPath)

	cmd.Stdout = os.Stdout

	err := cmd.Start()

	if err != nil {
		return err
	}

	_, err = json.Marshal(&runner.appConfig)

	return err
}
