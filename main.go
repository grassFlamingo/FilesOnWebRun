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

var (
	fileOnlyServer = http.FileServer(FileOnlyDir("."))
)

const (
	WORKING_PATH = "/files"
	PATH_INDEX   = "index.html"
)

type MainConfig struct {
	WorkingDir string
	Port       string
}

func unPackConfigFile(dir string) *MainConfig {
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
	workdir, err := os.Open(conf.WorkingDir)
	if err != nil {
		home := os.Getenv("HOME")
		if home == "" {
			conf.WorkingDir = "."
		} else {
			conf.WorkingDir = home
		}

	}
	workdir.Close()
	return conf
}

func main() {
	log.Println("Starting App")
	conf := unPackConfigFile("filesonweb.json")
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler(WORKING_PATH, http.StatusFound))
	mux.Handle(WORKING_PATH+"/", newFileServer(http.Dir(conf.WorkingDir)))
	mux.Handle("/img/", fileOnlyServer)
	mux.Handle("/js/", fileOnlyServer)
	mux.Handle("/css/", fileOnlyServer)

	err := http.ListenAndServe(conf.Port, mux)
	fmt.Println(err)
}
