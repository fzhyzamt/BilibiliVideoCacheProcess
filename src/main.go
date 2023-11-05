package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const M4sSuffix = ".m4s"

func main() {
	args, err := processArgs()
	if err != nil {
		fmt.Println("初始化配置失败", err)
		return
	}
	fmt.Println(args.String())

	videoDirArray, err := eachCollectVideoInfoByArgs(args)
	if err != nil {
		fmt.Println("遍历目录失败", err)
		return
	}

	for _, videoDir := range videoDirArray {
		fmt.Printf("处理开始 %s\n", videoDir)
		err := processOneVideoDir(args, videoDir)
		if err != nil {
			fmt.Printf("处理文件 %s 时发生错误: %v\n", videoDir, err)
		} else {
			fmt.Printf("处理结束 %s \n", videoDir)
		}
		fmt.Printf("%s\n\n", strings.Repeat("-", 70))
	}
}

func processArgs() (RunArgs, error) {
	path := flag.String("path", "", "视频文件或文件夹路径，默认取当前文件夹")
	mergeOutputPath := flag.String("merge-output", "", "输出视频文件的路径，默认取当前文件夹")
	ffmpeg := flag.String("ffmpeg", "", "ffmpeg执行文件路径，默认从环境变量中获取")
	deleteTempFile := flag.Bool("delete-temp", true, "执行完毕后是否删除临时文件")
	flag.Usage = func() {
		fmt.Println(`注1：windows资源管理器不支持的字符（如冒号、斜杠等）会替换为下划线`)
		//goland:noinspection GoPrintFunctions
		fmt.Println(`注2：哔哩哔哩客户端默认缓存文件夹为%USERPROFILE%\Videos\bilibili\`)
		fmt.Println("命令行参数:")
		flag.PrintDefaults()
	}
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
	args.deleteTempFile = *deleteTempFile

	return args, nil
}

func processOneVideoDir(args RunArgs, oneVideo VideoDir) error {
	if len(oneVideo.m4sPath) == 1 {
		// 如果只有一个文件，不需要合并
		return nil
	}

	targetVideoPathname, coverImagePath := guessTargetVideoPathname(args, oneVideo)
	targetStat, err := os.Stat(targetVideoPathname)
	if err == nil && targetStat.Size() != 0 {
		fmt.Println("目标文件已存在，跳过", targetVideoPathname)
		return nil
	}
	var newM4sFileArray []string
	for _, m4sFile := range oneVideo.m4sPath {
		processedFile, err := processM4sFile(m4sFile)
		if err != nil {
			return err
		}
		if args.deleteTempFile && strings.HasSuffix(processedFile, ".temp") {
			//goland:noinspection GoDeferInLoop
			defer func(name string) {
				err := os.Remove(name)
				if err != nil {
					fmt.Printf("删除临时文件失败：%s\n", processedFile)
				}
			}(processedFile)
		}
		fmt.Printf("处理m4s文件：%s\t-->\t%s\n", m4sFile, processedFile)
		newM4sFileArray = append(newM4sFileArray, processedFile)
	}
	oneVideo.m4sPath = newM4sFileArray
	return mergeToMp4(args, oneVideo, targetVideoPathname, coverImagePath)
}

func guessTargetVideoPathname(args RunArgs, oneVideo VideoDir) (string, string) {
	var targetVideoName string
	var coverImagePath string
	if oneVideo.videoInfoPath != "" {
		videoInfo, err := parseVideoInfo(oneVideo.videoInfoPath)
		if err == nil {
			targetVideoName = videoInfo.Title + ".mp4"
			coverImagePath = videoInfo.CoverPath
		}
	}
	if targetVideoName == "" {
		targetVideoName = filepath.Base(oneVideo.m4sPath[0]) + ".mp4"
	}
	targetVideoPathname := filepath.Join(args.mergeOutputPath, targetVideoName)
	return targetVideoPathname, coverImagePath
}
