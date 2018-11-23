package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

func newFileServer(fsys http.FileSystem) http.Handler {
	return &fileServer{root: fsys}
}

type fileServer struct {
	root http.FileSystem
}

func (self *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	if strings.HasPrefix(upath, WORKING_PATH) {
		upath = upath[6:]
	}
	r.URL.Path = path.Clean(upath)
	apireq := r.FormValue("requestfor")
	if apireq != "data" {
		self.ServeView(w, r)
	} else {
		self.ServeJson(w, r)
	}
}

func (self *fileServer) ServeJson(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	file, err := self.root.Open(upath)
	if err != nil {
		log.Println("Open File", upath, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !info.IsDir() {
		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
		return
	}

	// json response
	files, err := file.Readdir(-1)
	if err != nil {
		log.Println("Show Dir", err)
		return
	}
	echo := make([]*fileEchoItem, 0, len(files))
	for _, f := range files {
		if f.Name()[0] == '.' {
			continue
		}
		fei := newfileEchoItem(upath, f)
		if fei.IsImg {
			imginfo, err := self.root.Open(path.Join(upath, f.Name()))
			if err != nil {
				log.Println("Get Image size in servehttp", err)
			} else {
				fei.Width, fei.Height = GetImageSize(bufio.NewReader(imginfo))
				imginfo.Close()
			}
		}
		echo = append(echo, fei)
	}
	JsonResponse(w, echo)
}

func (self *fileServer) ServeView(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path

	file, err := self.root.Open(upath)
	if err != nil {
		log.Println("Open File", upath, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if info.IsDir() {
		doc := template.Must(template.ParseFiles(PATH_INDEX))
		doc.Execute(w, upath)
	} else {
		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	}
}

type dirpwd struct {
	PWD    string
	cindex int
}

func (self dirpwd) HasNext() bool {
	return self.cindex < len(self.PWD)
}
func (self dirpwd) SplidPWD() (string, string) {
	upath := self.PWD
	oind := self.cindex
	lenght := len(upath)
	if oind >= lenght {
		return "", ""
	}
	for i := oind; i < len(upath); i++ {
		if upath[i] == '/' {
			self.cindex = i
			break
		}
	}
	return upath[0:self.cindex], upath[oind:self.cindex]
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
