/*
 * Created by Aliy At December 28, 2018
 *
 * Simple Filter
 * File Name Filter
 */

package main

import (
	"log"
	"regexp"
)

type StringFilter interface {
	DoFilter(fileName string) bool
}

type FileNameFilter struct {
	blockHiddenFiles bool
	pattern          *regexp.Regexp
}

// New File Name Filter
// - blockHidden: if true then the file start with '.' won't pass
// - pattern: regexp pattern in string. It will return nil if pattern is not right.
// 			if pattern is "", this filter does not use pattern
func NewFileNameFilter(blockHidden bool, pattern string) *FileNameFilter {
	var pat *regexp.Regexp
	var err error
	if len(pattern) <= 0 {
		pat = nil
	} else {
		pat, err = regexp.Compile(pattern)
		if err == nil {
			log.Printf("Error at Compile Pattern %s. Error: %s", pattern, err.Error())
			return nil
		}
	}
	return &FileNameFilter{
		blockHiddenFiles: blockHidden,
		pattern:          pat,
	}
}

func (self *FileNameFilter) DoFilter(fileName string) bool {
	if len(fileName) <= 0 {
		return false
	}
	if fileName[0] == '.' && self.blockHiddenFiles {
		// Hidden file
		return false
	}
	if self.pattern == nil {
		return true
	}
	if self.pattern.MatchString(fileName) {
		return true
	}
	return false
}

type HiddenFileFilter struct{}

func (self *HiddenFileFilter) DoFilter(fileName string) bool {
	if len(fileName) <= 0 {
		return false
	}
	if fileName[0] == '.' || fileName[0] == '$' {
		return false
	}
	return true
}
