package main

import (
	"os"
	"path"
	"strings"
)

type fileEchoItem struct {
	Name   string // file name
	Path   string
	Icon   string // icon path
	IsDir  bool
	IsImg  bool
	Width  int
	Height int
}

func newfileEchoItem(basepath string, info os.FileInfo) *fileEchoItem {
	vf := &fileEchoItem{
		Name:   info.Name(),
		Path:   path.Join(WORKING_PATH, basepath, info.Name()),
		IsDir:  info.IsDir(),
		IsImg:  false,
		Width:  160,
		Height: 120,
	}
	vf.GetIconFile()
	return vf
}

// return icon path
func (self *fileEchoItem) GetIconFile() string {
	if self.IsDir {
		self.Icon = "/img/folder.svg"
		return self.Icon
	}
	suffix := "unknwon"
	for i := len(self.Name) - 1; i >= 0; i-- {
		if self.Name[i] == '.' {
			suffix = strings.ToLower(self.Name[i:])
			break
		}
	}
	var icon string
	switch suffix {
	case ".pdf":
		icon = "pdf.svg"
	case ".txt":
		icon = "txt.svg"
	case ".rar", ".zip", ".gz":
		icon = "achive.svg"
	case ".doc", ".docx":
		icon = "word.svg"
	case ".png", ".svg", ".jpeg", ".jpg", ".bmp", ".gif":
		icon = "img.svg"
		self.IsImg = true
	default:
		icon = "unknow.svg"
	}
	self.Icon = "/img/" + icon
	return self.Icon
}

/*
https://www.jianshu.com/p/05671bab2357
http://www.admpub.com/blog/post-221.html
https://www.quantamagazine.org/new-theory-cracks-open-the-black-box-of-deep-learning-20170921/?utm_source=Quanta+Magazine&utm_campaign=49cbae757d-EMAIL_CAMPAIGN_2017_09_21&utm_medium=email&utm_term=0_f0cb61321c-49cbae757d-389736021
*/
