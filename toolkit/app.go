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
	name      string // The app name
	srcPath   string // The app's source code directory full path
	gswebPath string // The gsweb package path
}

func goPaths() ([]string, error) {
	GOPATH := os.Getenv("GOPATH")

	if GOPATH == "" {
		return nil, gserrors.Newf(ErrCompile, "must set GOPATH first")
	}

	return strings.Split(GOPATH, string(os.PathListSeparator)), nil

}

// IsAppDuplicate .
func IsAppDuplicate(err error) bool {
	if gserr, ok := err.(gserrors.GSError); ok {
		return gserr.Origin() == ErrAppDuplicate
	}

	return false
}

// IsAppNotFound .
func IsAppNotFound(err error) bool {
	if gserr, ok := err.(gserrors.GSError); ok {
		return gserr.Origin() == ErrAppNotFound
	}

	return false
}

// SearchAppPath .
func SearchAppPath(packageName string) (string, error) {
	goPath, err := goPaths()

	if err != nil {
		return "", err
	}

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
		return "", gserrors.Newf(ErrAppDuplicate, stream.String())
	} else if len(found) == 0 {
		var stream bytes.Buffer

		stream.WriteString(fmt.Sprintf("package %s not found in any of:\n", packageName))

		for i, path := range found {
			stream.WriteString(fmt.Sprintf("\n\t%d) %s", i, path))
		}
		return "", gserrors.Newf(ErrAppNotFound, stream.String())
	}

	return found[0], nil
}

// CreateApp create new app package base on template
func CreateApp(packageName string, template string, searchpath []string) (*App, error) {

	_, err := SearchAppPath(packageName)

	if IsAppDuplicate(err) {
		return nil, err
	}

	gswebPath, err := SearchAppPath("github.com/gsdocker/gsweb/")

	if err != nil {
		return nil, err
	}

	templateDir := filepath.Join(gswebPath, "template")

	prototypeDir := filepath.Join(templateDir, template)

	if !gsos.IsDir(prototypeDir) {
		return nil, gserrors.Newf(ErrApp, "gsweb app template not found :%s\n\t search path :%s", template, templateDir)
	}

	goPaths, err := goPaths()

	if err != nil {
		return nil, err
	}

	targetDir := filepath.Join(goPaths[0], "src", packageName)

	if gsos.IsDir(targetDir) {
		os.RemoveAll(targetDir)
	}

	err = gsos.CopyDir(prototypeDir, targetDir)

	if err != nil {
		return nil, err
	}

	app, err := LoadApp(packageName)

	return app, err
}

// LoadApp load app config by app's package name
func LoadApp(packageName string) (*App, error) {

	srcPath, err := SearchAppPath(packageName)

	if err != nil {
		return nil, err
	}

	gswebPath, err := SearchAppPath("github.com/gsdocker/gsweb")

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
		name:      packageName,
		srcPath:   srcPath,
		gswebPath: gswebPath,
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

// GSWebPath get gsweb source path
func (app *App) GSWebPath() string {
	return app.gswebPath
}
