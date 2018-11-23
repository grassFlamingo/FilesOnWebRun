/**
 * This file is created by aliy at November 21, 2018
 */

package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	fileOnlyServer = http.FileServer(FileOnlyDir("."))
)

const (
	WORKING_PATH = "/files"
	PATH_INDEX   = "index.html"
)

func main() {
	log.Println("Starting App")
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler(WORKING_PATH, http.StatusFound))
	mux.Handle(WORKING_PATH+"/", newFileServer(http.Dir("/media/aliy/DATA/PHOTO/")))
	mux.Handle("/img/", fileOnlyServer)
	mux.Handle("/js/", fileOnlyServer)
	mux.Handle("/css/", fileOnlyServer)

	err := http.ListenAndServe(":9898", mux)
	fmt.Println(err)
}
