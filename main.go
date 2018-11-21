/**
 * This file is created by aliy at November 21, 2018
 */

package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting App")
	mux := http.NewServeMux()
	mux.Handle("/", newTemplateFileServer(http.Dir(".")))
	err := http.ListenAndServe(":9898", mux)
	fmt.Println(err)
}
