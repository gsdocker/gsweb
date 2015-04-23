package toolkit

import (
	"bytes"
	"crypto/md5"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gsdocker/gsconfig"
	"github.com/gsdocker/gslogger"
	"github.com/gsdocker/gsos"
	"github.com/howeyc/fsnotify"
)

// CompileSession .
type CompileSession struct {
	StartTime time.Time // Current build session start time
	EndTime   time.Time // Current build session end time
	Md5Check  []byte    // Current build session product md5 check value
	Binary    string    // Current build session product file path
}

// AppCompileS create app compiler service
type AppCompileS struct {
	gslogger.Log                     // Mixin log APIs
	app          *App                // The app config object
	fileWatcher  *fsnotify.Watcher   // The source file watcher
	buildDir     string              // The app's build directory
	binaryPath   string              // The app's build target fullpath
	md5Check     []byte              // The app's md5Check
	Notify       chan CompileSession // compile session notify

}

// NewAppCompileS create new app compile service
func NewAppCompileS(app *App) (*AppCompileS, error) {

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, err
	}

	cs := &AppCompileS{
		Log:         gslogger.Get("CompileS"),
		app:         app,
		fileWatcher: watcher,
		buildDir:    filepath.Join(os.TempDir(), app.String()),
		binaryPath:  filepath.Join(os.TempDir(), app.String(), "app", gsos.ExeSuffix),
		Notify:      make(chan CompileSession, gsconfig.Uint("compile.notify.maxsize", 100)),
	}

	// try create build directory
	if !gsos.IsExist(cs.buildDir) {
		err = os.MkdirAll(cs.buildDir, 0777)

		if err != nil {
			return nil, err
		}
	}

	cs.D("build directory : %s", cs.buildDir)

	return cs, cs.fileWatcher.Watch(app.SrcPath())
}

func (cs *AppCompileS) needCompile(file string) bool {
	extension := filepath.Ext(file)

	for _, target := range cs.app.WatchFiles {
		if target == extension {
			return true
		}
	}

	return false
}

// Start start app compile service
func (cs *AppCompileS) Start() {
	go func() {

		// build app at least once
		cs.processBuild()

		for event := range cs.fileWatcher.Event {

			if cs.needCompile(event.Name) {
				cs.processBuild()
			}
		}
	}()
}

func (cs *AppCompileS) processBuild() {
	cs.I("start compile app ...")

	startTime := time.Now()

	if err := cs.compileApp(); err != nil {
		cs.E("compile app -- failed\n%s", err)
		return
	}

	endTime := time.Now()

	cs.I("compile app -- success")

	md5Check := cs.calcMd5Check()

	cs.D("app binary md5 check is :%v", md5Check)

	if bytes.Compare(md5Check, cs.md5Check) == 0 {
		return
	}

	// change file mod

	err := os.Chmod(cs.binaryPath, 0777)

	if err != nil {
		cs.E("compile app -- failed\n%s", err)
	}

	// notify new version binary is valid

	cs.md5Check = md5Check

	// notify compile event

	session := CompileSession{
		StartTime: startTime,
		EndTime:   endTime,
		Binary:    cs.binaryPath,
		Md5Check:  md5Check,
	}

	select {
	case cs.Notify <- session:
	default:
		cs.W("max notify event queue size reach !!!")
	}
}

func (cs *AppCompileS) calcMd5Check() []byte {

	file, err := os.Open(cs.binaryPath)
	if err != nil {

		cs.W("generate binary md5 check err :%s", err)

		return nil
	}

	md5h := md5.New()
	io.Copy(md5h, file)
	return md5h.Sum([]byte(""))
}

// BuildDir compiles build directory
func (cs *AppCompileS) BuildDir() string {
	return cs.buildDir
}

func (cs *AppCompileS) compileApp() error {
	cmd := exec.Command("go", "build", "-o", cs.binaryPath, cs.app.String())

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
