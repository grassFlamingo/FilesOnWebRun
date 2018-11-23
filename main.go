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
}

func unPackConfigFile(dir string) *MainConfig {
	file, err := os.OpenFile(dir, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Panicln("Error", err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicln("Error", err)
	}
	conf := &MainConfig{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		log.Panicln("Error", err)
	}
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

	err := http.ListenAndServe(":9898", mux)
	fmt.Println(err)
}
