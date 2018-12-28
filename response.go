/*
 * Created by Aliy At December 28, 2018
 *
 */

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type JsonResponseState uint8

const (
	JSRState_OK  = 0x01
	JSRState_ERR = 0xF0
)

type DataResponse struct {
	State JsonResponseState
	Data  interface{}
}

func JsonResponse(w io.Writer, data interface{}) {
	res, err := json.Marshal(DataResponse{State: JSRState_OK, Data: data})
	if err != nil {
		log.Fatalln("JsonResponse", err)
	}
	w.Write(res)
}

func JsonErrorResponse(w io.Writer, state JsonResponseState) {
	res, err := json.Marshal(DataResponse{State: state, Data: nil})
	if err != nil {
		log.Fatalln("JsonResponse", err)
	}
	w.Write(res)
}

func httpErrorResponse(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
