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
	"time"
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
		if h > 180 || h < 100 {
			w, h = int(180.0*f), 180
		}
	} else {
		// fat picture
		// if the image too fat (w > 300) or not fat enough (w < 120)
		if w > 180 || w < 120 {
			w, h = 180, int(180.0/f)
		}
	}
	return w, h
}

var (
	BUF_IMG_SIZE_MAP = make(map[string]int)
)

func init() {
	go func() {
		for {
			time.Sleep(80 * time.Hour)
			BUF_IMG_SIZE_MAP = nil
			BUF_IMG_SIZE_MAP = make(map[string]int)
		}
	}()
}

func GetBufImageSize(key string, r io.Reader) (int, int) {
	var w, h int
	if v, ok := BUF_IMG_SIZE_MAP[key]; ok {
		w = (v & 0xffff0000) >> 16
		h = v & 0x0000ffff
	} else {
		w, h = GetImageSize(r)
		wh := ((w & 0xffff) << 16) | (h & 0xffff)
		BUF_IMG_SIZE_MAP[key] = wh
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
