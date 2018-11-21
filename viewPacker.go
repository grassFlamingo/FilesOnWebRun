package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type viewPacker struct {
	w        http.ResponseWriter
	r        *http.Request
	basepath string
}

const (
	PATH_INDEX = "index.html"
)

func newViewPacker(w http.ResponseWriter, r *http.Request) *viewPacker {
	return &viewPacker{w: w, r: r}
}

// name is the full path of dir
func (self *viewPacker) ShowDir(name string) {
	log.Println("show dir")
	doc := template.Must(template.ParseFiles(PATH_INDEX))

	files, err := ioutil.ReadDir(name)
	if err != nil {
		log.Println(err)
		return
	}
	echo := make([]string, 0, len(files))
	for _, f := range files {
		if f.Name()[0] == '.' {
			continue
		}
		echo = append(echo, f.Name())
	}
	err = doc.Execute(self.w, &viewEcho{
		BasePath: self.basepath,
		Files:    echo,
	})
	if err != nil {
		log.Println("error", err.Error())
	}
}

type viewEcho struct {
	BasePath string
	Files    []string
}

/*
https://www.jianshu.com/p/05671bab2357
http://www.admpub.com/blog/post-221.html
https://www.quantamagazine.org/new-theory-cracks-open-the-black-box-of-deep-learning-20170921/?utm_source=Quanta+Magazine&utm_campaign=49cbae757d-EMAIL_CAMPAIGN_2017_09_21&utm_medium=email&utm_term=0_f0cb61321c-49cbae757d-389736021
*/
