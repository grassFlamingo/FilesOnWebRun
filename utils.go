/**
 * Created by aliy at November 22, 2018
 *
 */

package main

import (
	"encoding/json"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// this result was scaled limit (300,300)
// return (width, height)
func GetImageSize(r io.Reader) (int, int) {
	conf, format, err := image.DecodeConfig(r)
	if err != nil {
		log.Println("GetImageSize", format, err)
		return 120, 160
	}
	w, h := conf.Width, conf.Height
	f := float64(w) / float64(h)

	if f < 1 {
		// tall picture
		// if the image too high (h > 300) or not tall enough (h < 120)
		if h > 300 || h < 120 {
			w, h = int(300.0*f), 300
		}
	} else {
		// fat picture
		// if the image too fat (w > 300) or not fat enough (w < 120)
		if w > 300 || w < 160 {
			w, h = 300, int(300.0/f)
		}
	}
	return w, h
}

func JsonResponse(w io.Writer, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		log.Fatalln("JsonResponse", err)
	}
	w.Write(res)
	// log.Println("JsonResponse", string(res))
}

func intMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func httpErrorResponse(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

type FileOnlyDir string

func (self FileOnlyDir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	dir := string(self)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	info, err := os.Stat(fullName)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, errors.New("File Only Dir Recive A Dir " + fullName)
	}
	f, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	return f, nil
}
