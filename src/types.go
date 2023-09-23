package main

import "path/filepath"

type RunArgs struct {
	path            string
	mergeOutputPath string
	ffmpegPath      string
}
type VideoDir struct {
	m4sPath       []string
	videoInfoPath string
}
type VideoInfo struct {
	Title string `json:"title"`
}

func (r *VideoDir) appendM4s(path string) {
	if r.m4sPath == nil {
		r.m4sPath = []string{path}
	} else {
		r.m4sPath = append(r.m4sPath, path)
	}
}

func (r *VideoDir) getDirPath() string {
	return filepath.Dir(r.m4sPath[0])
}
