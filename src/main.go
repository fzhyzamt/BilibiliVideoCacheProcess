package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const M4sSuffix = ".m4s"

func main() {
	args, err := processArgs()
	if err != nil {
		fmt.Println("初始化配置失败", err)
		return
	}
	fmt.Println("Args", args)

	videoDirArray, err := eachCollectVideoInfoByArgs(args)
	if err != nil {
		fmt.Println("遍历目录失败", err)
		return
	}

	for _, videoDir := range videoDirArray {
		err := processOneVideoDir(args, videoDir)
		if err != nil {
			fmt.Printf("处理文件 %s 时发生错误: %v\n", videoDir, err)
			return
		}
		fmt.Printf("文件 %s 处理完成\n", videoDir)
	}
}

func processArgs() (RunArgs, error) {
	path := flag.String("path", "", "视频文件或文件夹路径")
	mergeOutputPath := flag.String("merge-output", "", "输出视频文件的路径")
	ffmpeg := flag.String("ffmpeg", "", "ffmpeg路径")
	flag.Parse()
	var args = RunArgs{}

	if len(*path) != 0 {
		path, err := filepath.Abs(*path)
		if err != nil {
			return args, err
		}
		args.path = path
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return args, err
		}
		args.path = wd
	}

	if len(*mergeOutputPath) != 0 {
		outputDir, err := filepath.Abs(*path)
		if err != nil {
			return args, err
		}
		args.mergeOutputPath = outputDir
	} else {
		args.mergeOutputPath = args.path
	}

	if len(*ffmpeg) != 0 {
		path, err := filepath.Abs(*path)
		if err != nil {
			return args, err
		}
		args.ffmpegPath = path
	} else {
		ffmpeg, err := getFFMPEG()
		if err != nil {
			return args, err
		}
		args.ffmpegPath = ffmpeg
	}

	return args, nil
}

func processOneVideoDir(args RunArgs, oneVideo VideoDir) error {
	var newM4sFileArray []string
	for _, m4sFile := range oneVideo.m4sPath {
		processedFile, err := processM4sFile(m4sFile)
		if err != nil {
			return err
		}
		fmt.Printf("处理m4s文件：%s\t-->\t%s\n", m4sFile, processedFile)
		newM4sFileArray = append(newM4sFileArray, processedFile)
	}
	oneVideo.m4sPath = newM4sFileArray

	if len(oneVideo.m4sPath) == 1 {
		// 如果只有一个文件，不需要合并
		return nil
	}

	var targetVideoName string
	if len(oneVideo.videoInfoPath) > 0 {
		videoInfo, err := parseVideoInfo(oneVideo.videoInfoPath)
		if err == nil {
			targetVideoName = videoInfo.Title + ".mp4"
		}
	}
	if len(targetVideoName) == 0 {
		targetVideoName = oneVideo.m4sPath[0] + ".mp4"
	}

	return mergeToMp4(oneVideo, filepath.Join(args.mergeOutputPath, targetVideoName))
}
