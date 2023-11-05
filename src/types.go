package main

import "path/filepath"

type RunArgs struct {
	path            string
	mergeOutputPath string
	ffmpegPath      string
	deleteTempFile  bool
}
type VideoDir struct {
	m4sPath       []string
	videoInfoPath string
}
type VideoInfo struct {
	/**
	视频标题
	*/
	Title string `json:"title"`
	/**
	封面图
	*/
	CoverPath string `json:"coverPath"`
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
