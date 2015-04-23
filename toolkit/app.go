package toolkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gsdocker/gserrors"
	"github.com/gsdocker/gsos"
)

// App the app's config
type App struct {
	name       string   // The app name
	srcPath    string   // The app's source code directory full path
	WatchFiles []string `json:"autocompile.files"` // The files status event watcher config
}

func searchAppPath(packageName string) (string, error) {
	GOPATH := os.Getenv("GOPATH")

	if GOPATH == "" {
		return "", gserrors.Newf(ErrCompile, "must set GOPATH first")
	}

	goPath := strings.Split(GOPATH, string(os.PathListSeparator))

	var found []string
	for _, path := range goPath {
		fullpath := filepath.Join(path, "src", packageName)
		fi, err := os.Stat(path)

		if err == nil && fi.IsDir() {
			found = append(found, fullpath)
		}
	}

	if len(found) > 1 {
		var stream bytes.Buffer

		stream.WriteString(fmt.Sprintf("found more than one package named :%s", packageName))

		for i, path := range found {
			stream.WriteString(fmt.Sprintf("\n\t%d) %s", i, path))
		}
		return "", gserrors.Newf(ErrCompile, stream.String())
	} else if len(found) == 0 {
		var stream bytes.Buffer

		stream.WriteString(fmt.Sprintf("package %s not found in any of:\n", packageName))

		for i, path := range found {
			stream.WriteString(fmt.Sprintf("\n\t%d) %s", i, path))
		}
		return "", gserrors.Newf(ErrCompile, stream.String())
	}

	return found[0], nil
}

// LoadApp load app config by app's package name
func LoadApp(packageName string) (*App, error) {

	srcPath, err := searchAppPath(packageName)

	if err != nil {
		return nil, err
	}

	// Unmarsh json from app's config file
	configfile := filepath.Join(srcPath, ".gsweb")

	if !gsos.IsExist(configfile) {
		return nil, gserrors.Newf(ErrApp, "not found .gsweb file in app %s directory :\n\t%s", packageName, srcPath)
	}

	config, err := ioutil.ReadFile(configfile)

	if err != nil {
		return nil, gserrors.Newf(ErrApp, "read .gsweb config -- failed\n\tfile: %s\n\terr :%s", configfile, err)
	}

	app := &App{
		name:    packageName,
		srcPath: srcPath,
	}

	err = json.Unmarshal(config, &app)

	return app, err
}

func (app *App) String() string {
	return app.name
}

// SrcPath get app's source path
func (app *App) SrcPath() string {
	return app.srcPath
}
