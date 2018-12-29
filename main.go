/**
 * This file is created by aliy at November 21, 2018
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	WORKING_PATH = "/root"
	PATH_INDEX   = "index.html"
)

type MainConfig struct {
	RootDir string
	Port    string
}

func upackConfigFile(dir string) *MainConfig {
	file, err := os.OpenFile(dir, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Panicln("Error", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicln("Error", err)
	}
	conf := &MainConfig{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		log.Panicln("Error", err)
	}
	workdir, err := os.Open(conf.RootDir)
	if err != nil {
		home := os.Getenv("HOME")
		if home == "" {
			conf.RootDir = "."
		} else {
			conf.RootDir = home
		}

	}
	workdir.Close()
	return conf
}

func main() {
	log.Println("Starting App")
	conf := upackConfigFile("filesonweb.json")
	dir := http.Dir(conf.RootDir)
	fileOnly := http.FileServer(FileOnlyDir("."))

	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler(WORKING_PATH, http.StatusFound))
	mux.Handle(WORKING_PATH+"/", newFileServer(dir))
	mux.Handle("/exec/", NewExecServer(dir))
	mux.Handle("/img/", fileOnly)
	mux.Handle("/js/", fileOnly)
	mux.Handle("/css/", fileOnly)

	err := http.ListenAndServe(conf.Port, mux)
	fmt.Println(err)
}
