/*
 * Created by Aliy At December 28, 2018
 *
 * Sending HTTP/JSON Response
 */

package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type JsonResponseState uint16

const (
	JSRState_OK  = 0x0100
	JSRState_ERR = 0x0200

	JSRState_BAD_CMD  = 0x0201
	JSRState_BAD_DIR  = 0x0202
	JSRState_BAD_OPEN = 0x0203
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

func JsonStateResponse(w io.Writer, state JsonResponseState, data interface{}) {
	res, err := json.Marshal(DataResponse{State: state, Data: data})
	if err != nil {
		log.Fatalln("JsonResponse", err)
	}
	w.Write(res)
}

func httpErrorResponse(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
