package gsweb

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gsdocker/gslogger"
	"github.com/gsdocker/gsos/fs"
)

// RegisterPath the file handler's register path object
type RegisterPath struct {
	enableListChild bool         // Indicate if allow list dir child items
	path            string       // The fileHandler path name
	handler         http.Handler // the fileHandler path's handler
}

// EnableGetDir set flag, true enable list directory's child items,otherwise
// disable it
func (path *RegisterPath) EnableGetDir(flag bool) {
	path.enableListChild = flag
}

// FileHandler The static fileHandler handler
type FileHandler struct {
	gslogger.Log                           // Mixin log APIs
	registerPaths map[string]*RegisterPath // The register fileHandlerpath
}

// NewFileHandler create new fileHandler handler
func NewFileHandler() *FileHandler {
	fileHandler := &FileHandler{
		Log:           gslogger.Get("fileHandler"),
		registerPaths: make(map[string]*RegisterPath),
	}

	return fileHandler
}

// HandleGet implement Get interface
func (fileHandler *FileHandler) HandleGet(context *Context) error {

	uri := context.RequestURI()

	fileHandler.V("GET %s forward processing", uri)

	var matchedPrefix string
	var registerPath *RegisterPath

	for prefix, path := range fileHandler.registerPaths {
		if strings.HasPrefix(uri, prefix) && len(prefix) > len(matchedPrefix) {
			matchedPrefix = prefix
			registerPath = path
			break
		}
	}

	path := strings.Replace(uri, matchedPrefix, registerPath.path+"/", 1)

	// if the target uri is a filesystem's dir try load index.html file
	if fs.IsDir(path) {
		indexfile := filepath.Join(path, "index.html")

		fileHandler.V("try get file : %s", indexfile)

		// if target is not exist or is a directory and disable directory child list
		// break processing and foward this request to next chain handler
		if (!fs.Exists(indexfile) || fs.IsDir(indexfile)) && !registerPath.enableListChild {
			fileHandler.V("not found file : %s", indexfile)
			goto FORWARD
		}
	}

	if fs.Exists(path) {

		fileHandler.D("GET %s handler -- found", uri)

		registerPath.handler.ServeHTTP(context.Response(), context.Request())

		return context.Success()
	}

FORWARD:

	// forward this request to next chain handler
	err := context.Forward()

	fileHandler.V("GET %s backward processing", uri)

	return err
}

// RegisterPath register fileHandler path by uri prefix
func (fileHandler *FileHandler) RegisterPath(uriprefix string, dir string) (*RegisterPath, error) {

	fileHandler.D("register path %s => %s", uriprefix, dir)

	path, err := filepath.Abs(dir)

	if err != nil {
		fileHandler.E("register path error : \n\turi:%s\n\tdir:%s\n\terror:%s", uriprefix, dir, err)
		return nil, err
	}

	registerPath := &RegisterPath{path: dir, handler: http.FileServer(http.Dir(path))}

	fileHandler.registerPaths[uriprefix] = registerPath

	fileHandler.D("register path -- success")

	return registerPath, nil
}
