package main

import (
	"errors"
	"os"
	"path/filepath"
)

func eachCollectVideoInfoByArgs(args RunArgs) ([]VideoDir, error) {
	var videos []VideoDir
	stat, err := os.Stat(args.path)
	if err != nil {
		return videos, err
	}
	if !stat.IsDir() {
		if filepath.Ext(args.path) == M4sSuffix {
			videos = []VideoDir{{m4sPath: []string{args.path}}}
		} else {
			return videos, errors.New("指定路径既不是文件夹，也不是.m4s文件")
		}
	} else {
		videos, err = eachCollectVideoInfo(args.path, 0)
		if err != nil {
			return videos, err
		}
	}
	return videos, nil
}

func eachCollectVideoInfo(rootPath string, depth int) ([]VideoDir, error) {
	currentVideo := VideoDir{m4sPath: []string{}}
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, err
	}
	var videos []VideoDir
	for _, entry := range entries {
		entryAbsPath := filepath.Join(rootPath, entry.Name())
		if entry.IsDir() {
			childVideos, err := eachCollectVideoInfo(entryAbsPath, depth+1)
			if err != nil {
				continue
			}
			videos = append(videos, childVideos...)
			continue
		}
		if filepath.Ext(entry.Name()) == M4sSuffix {
			currentVideo.appendM4s(entryAbsPath)
		} else if entry.Name() == ".videoInfo" {
			currentVideo.videoInfoPath = entryAbsPath
		}
	}
	if len(currentVideo.m4sPath) > 0 {
		videos = append(videos, currentVideo)
	}

	return videos, nil
}
