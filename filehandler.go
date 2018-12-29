package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type fileServer struct {
	root   http.FileSystem
	filter StringFilter
}

func newFileServer(fsys http.FileSystem) http.Handler {
	return &fileServer{root: fsys, filter: &HiddenFileFilter{}}
}

func (self *fileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	if strings.HasPrefix(upath, WORKING_PATH) {
		upath = upath[len(WORKING_PATH):]
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
	if r.FormValue("dironly") == "yes" {
		JsonResponse(w, self.echoDirOnly(upath, file))
	} else {
		JsonResponse(w, self.echoDirAndFile(upath, file))
	}
}

func (self *fileServer) echoDirOnly(upath string, file http.File) []*dirEchoItem {
	files, err := file.Readdir(-1)
	if err != nil {
		log.Println("Show Dir", err)
		return nil
	}
	echo := make([]*dirEchoItem, 0, len(files))
	for _, f := range files {
		if !self.filter.DoFilter(f.Name()) {
			continue
		}
		if f.IsDir() {
			fei := &dirEchoItem{
				Name: f.Name(),
				Path: path.Join(WORKING_PATH, upath, f.Name()),
				Icon: "/img/folder.svg",
			}
			echo = append(echo, fei)
		}
	}
	return echo
}

func (self *fileServer) echoDirAndFile(upath string, file http.File) []*fileEchoItem {
	files, err := file.Readdir(-1)
	if err != nil {
		log.Println("Show Dir", err)
		return nil
	}
	echo := make([]*fileEchoItem, 0, len(files))

	for _, f := range files {
		if !self.filter.DoFilter(f.Name()) {
			continue
		}
		fei := newfileEchoItem(upath, f)
		echo = append(echo, fei)
		if !fei.IsImg {
			continue
		}
		if ok, w, h := IsBuffedImageSize(fei.Path); ok {
			fei.Width, fei.Height = w, h
			continue
		}
		imginfo, err := self.root.Open(path.Join(upath, f.Name()))
		if err != nil {
			log.Println("Get Image size in servehttp", err)
		} else {
			fei.Width, fei.Height = BufImageSize(fei.Path, bufio.NewReader(imginfo))
			imginfo.Close()
		}
	}
	return echo
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
		http.SetCookie(w, &http.Cookie{
			Name:     "pwd",
			Value:    upath,
			Path:     "/",
			HttpOnly: true,
		})
		// TODO: move this to global if not debug
		var home_page_doc = template.Must(template.ParseFiles(PATH_INDEX))
		err := home_page_doc.Execute(w, newDirPwd(upath))
		if err != nil {
			log.Println(err)
		}
	} else {
		http.ServeContent(w, r, info.Name(), info.ModTime(), file)
	}
}

type dirpwd string

func newDirPwd(pwd string) dirpwd {
	lpwd := len(pwd)
	if lpwd <= 1 {
		pwd = "/"
	} else {
		if pwd[lpwd-1] == '/' {
			pwd = pwd[0 : lpwd-1]
		}
	}
	return dirpwd(pwd)
}

type dirpwditem struct {
	PWD  string
	Name string
}

func (self dirpwd) GetPackedPWD() []dirpwditem {
	upath := string(self)
	lenght := len(upath)
	out := make([]dirpwditem, 0, 20)

	if lenght <= 0 {
		return out
	}

	out = append(out, dirpwditem{PWD: "/", Name: "ROOT"})

	cend := 1
	cbeg := 1
	for i, c := range upath[1:] {
		if c == '/' {
			cend = i + 1
			out = append(out, dirpwditem{PWD: upath[0:cend], Name: upath[cbeg:cend]})
			cbeg = cend + 1
		}
	}
	if cend < lenght {
		out = append(out, dirpwditem{PWD: upath, Name: upath[cbeg:]})
	}

	return out
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
