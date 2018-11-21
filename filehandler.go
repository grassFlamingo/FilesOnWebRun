package main

import (
	"log"
	"net/http"
	"path"
	"strings"
)

func newTemplateFileServer(fsys http.FileSystem) http.Handler {
	return &fileServer{root: fsys}
}

type fileServer struct {
	root http.FileSystem
	w    http.ResponseWriter
	r    *http.Request
}

func (self *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	self.w = w
	self.r = r
	self.ServeFile(upath)
}

func (self *fileServer) ServeFile(upath string) {
	const indexPage = "/index.html"

	// redirect .../index.html to .../
	// can't use Redirect() because that would make the path absolute,
	// which would be a problem running under StripPrefix
	if strings.HasSuffix(upath, indexPage) {
		localRedirect(self.w, self.r, "./")
		return
	}
	file, err := self.root.Open(path.Clean(upath))
	if err != nil {
		http.Error(self.w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()
	log.Println("Open File", upath)

	info, err := file.Stat()
	if err != nil {
		http.Error(self.w, err.Error(), http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		vp := newViewPacker(self.w, self.r)
		vp.basepath = info.Name()
		vp.ShowDir(info.Name())
	} else {
		http.ServeContent(self.w, self.r, info.Name(), info.ModTime(), file)
	}
}

func simpleErrorHandler(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// localRedirect gives a Moved Permanently response.
// It does not convert relative paths to absolute paths like Redirect does.
func localRedirect(w http.ResponseWriter, r *http.Request, newPath string) {
	if q := r.URL.RawQuery; q != "" {
		newPath += "?" + q
	}
	w.Header().Set("Location", newPath)
	w.WriteHeader(http.StatusMovedPermanently)
}
