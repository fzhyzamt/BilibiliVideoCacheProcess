package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getFFMPEG() (string, error) {
	execPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", err
	}
	return execPath, nil
}

func processM4sFile(sourcePath string) (string, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer sourceFile.Close()

	sourceFileStat, err := sourceFile.Stat()
	if err != nil {
		return "", err
	}

	var BiliM4sFileDataPrefix = []byte{
		0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30,
		0x30, 0x00, 0x00, 0x00, 0x24, 0x66, 0x74, 0x79,
		0x70, 0x69, 0x73, 0x6f, 0x35, 0x00, 0x00, 0x00,
		0x01, 0x61, 0x76, 0x63, 0x31, 0x69, 0x73, 0x6f,
	}
	var m4sFileDataPrefix = []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x69, 0x73, 0x6f, 0x35, 0x00, 0x00, 0x00, 0x01,
		0x69, 0x73, 0x6f,
	}
	headerBuffer := make([]byte, len(BiliM4sFileDataPrefix))
	_, err = io.ReadFull(sourceFile, headerBuffer)
	if err != nil {
		return "", err
	}

	if !bytes.Equal(headerBuffer, BiliM4sFileDataPrefix) {
		return sourcePath, nil
	}

	tempPath := sourcePath + ".temp"
	tempStat, err := os.Stat(tempPath)
	if err == nil {
		if tempStat.Size()+int64(len(BiliM4sFileDataPrefix)-len(m4sFileDataPrefix)) == sourceFileStat.Size() {
			return tempPath, nil
		}
		err = os.Remove(tempPath)
		if err != nil {
			return "", err
		}
	}

	targetFile, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	_, err = targetFile.Write(m4sFileDataPrefix)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, 40*1024*1024) // 40MB
	for {
		n, readErr := sourceFile.Read(buffer)
		if readErr != nil && readErr != io.EOF {
			return "", readErr
		}

		_, err = targetFile.Write(buffer[:n])
		if err != nil {
			return "", err
		}

		if readErr == io.EOF {
			break
		}
	}
	return tempPath, nil
}

func parseVideoInfo(path string) (VideoInfo, error) {
	videoInfoData, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("读取.videoInfo失败", err)
		return VideoInfo{}, err
	}
	var videoInfo VideoInfo
	err = json.Unmarshal(videoInfoData, &videoInfo)
	if err != nil {
		fmt.Println("解析.videoInfo失败", err)
		return VideoInfo{}, err
	}
	coverImageState, err := os.Stat(videoInfo.CoverPath)
	if err != nil || coverImageState.Size() == 0 {
		videoInfo.CoverPath = ""
	}
	videoInfo.Title = removeInvalidCharFromPathname(videoInfo.Title)
	return videoInfo, nil
}

func removeInvalidCharFromPathname(pathname string) string {
	var builder strings.Builder
	builder.Grow(len(pathname))

	for _, char := range pathname {
		if isInvalidChar(char) {
			builder.WriteRune('_')
		} else {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

// 检查字符是否在字符数组中
func isInvalidChar(char rune) bool {
	invalidChars := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	for _, c := range invalidChars {
		if char == c {
			return true
		}
	}
	return false
}

func mergeToMp4(video VideoDir, targetPath string, coverImagePath string) error {

	inputStreamCount := 0
	// Example: ffmpeg -i cover.jpg -i audio.m4s -i video.m4s -map 0 -map 1 -map 2 -codec copy -disposition:v:0 attached_pic output.mp4
	var params []string
	if coverImagePath != "" {
		params = append(params, "-i", coverImagePath)
		inputStreamCount++
	}
	for _, m4s := range video.m4sPath {
		params = append(params, "-i", m4s)
		inputStreamCount++
	}
	for i := 0; i < inputStreamCount; i++ {
		params = append(params, "-map", strconv.Itoa(i))
	}
	params = append(params, "-codec", "copy")
	if coverImagePath != "" {
		params = append(params, "-disposition:v:0", "attached_pic")
	}
	params = append(params, targetPath)

	fmt.Println("进行ffmpeg转换", params)
	cmd := exec.Command("ffmpeg", params...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("ffmpeg转换失败", err)
		return err
	}
	fmt.Println(string(output))
	return nil
}
