/**
 * Created By Aliy At December 29, 2018
 *
 * This handler execute some basic shell commands and return its results
 * Avaliable Commands:
 * - ls:
 */

package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type ExecServer struct {
	exec *ExecCMD
}

func NewExecServer(root http.FileSystem) *ExecServer {
	return &ExecServer{exec: NewExecCMD(root)}
}

func (self *ExecServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pwd")
	if err != nil {
		log.Println("Error at Cookie", err)
		httpErrorResponse(w, http.StatusBadRequest)
		return
	}
	log.Println("Cookie", cookie.Value)
	cmd := r.FormValue("cmd")
	arg := r.FormValue("args")
	var args []string
	if len(arg) > 0 {
		args = strings.Split(arg, ",")
		for i := 0; i < len(args); i++ {
			args[i] = strings.TrimSpace(args[i])
		}
	}
	log.Println("Cmd is", cmd)
	self.exec.pwd = cookie.Value
	data, jsr := self.exec.Exec(cmd, args)
	JsonStateResponse(w, jsr, data)
}

type ExecCMD struct {
	root       http.FileSystem
	pwd        string
	fileFilter StringFilter
	handmap    map[string]func(args []string) (interface{}, JsonResponseState)
}

func NewExecCMD(root http.FileSystem) *ExecCMD {
	out := &ExecCMD{
		root:       root,
		pwd:        "/",
		fileFilter: &HiddenFileFilter{},
		handmap:    make(map[string]func(args []string) (interface{}, JsonResponseState)),
	}
	out.handmap["ls"] = out.LS
	return out
}

func (self *ExecCMD) Exec(cmd string, args []string) (interface{}, JsonResponseState) {
	cmd = strings.ToLower(cmd)
	hand, exis := self.handmap[cmd]
	if !exis {
		return nil, JSRState_BAD_CMD
	}
	return hand(args)
}

func (self *ExecCMD) LS(args []string) (interface{}, JsonResponseState) {
	file, err := self.root.Open(self.pwd)
	if err != nil {
		log.Println("E Ls", err.Error())
		return nil, JSRState_BAD_OPEN
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil || !stat.IsDir() {
		return nil, JSRState_BAD_DIR
	}
	subfiles, err := file.Readdir(-1)
	if err != nil {
		return nil, JSRState_BAD_DIR
	}
	out := make([]ExecLSItem, 0, len(subfiles))

	for _, f := range subfiles {
		if !self.fileFilter.DoFilter(f.Name()) {
			continue
		}
		out = append(out, ExecLSItem{
			Name:     f.Name(),
			IsDir:    f.IsDir(),
			ModeTime: f.ModTime(),
		})
	}
	return out, JSRState_OK
}

/**
 * This below are definitions of exec response struct
 */
type ExecLSItem struct {
	Name     string
	IsDir    bool
	ModeTime time.Time
}
