package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func FileInfo(file os.DirEntry, path string) (string, string) { // 文件大小 ， 文件上次修改时间
	var fileName string

	if path == "." {
		fileName = file.Name()
	} else {
		fileName = filepath.Join(path, file.Name())
	}

	info, err := os.Stat(fileName)

	if err != nil {
		errorMessage := fmt.Sprintf("ERR:%s", err)
		return errorMessage, errorMessage
	}

	timeFormatter := "2006/01/02 15:04"

	lastWriteTimeRaw := info.ModTime()
	lastWriteTime := lastWriteTimeRaw.Format(timeFormatter)

	if file.IsDir() {
		return "", lastWriteTime
	} else {
		fileSize := strconv.FormatInt(info.Size(), 10) //10为转为十进制
		return fileSize, lastWriteTime
	}
}

func FileMode(file os.DirEntry) string {
	info, _ := file.Info()
	fileMode := info.Mode().String()
	return fileMode
}

func ColorFormatter(file os.DirEntry, isDir bool) string {
	var (
		isExecutable bool
		// TitleColor       = color.New(color.FgHiGreen, color.Bold).SprintFunc()
		ApplicationColor = color.New(color.FgHiGreen).SprintFunc()
		FolderColor      = color.New(color.FgBlue).SprintFunc()
		// FileColor        = color.New(color.FgWhite).SprintFunc()
	)
	fileName := file.Name()
	info, _ := file.Info()
	isExecutable = info.Mode().Perm()&0111 != 0 // 判断文件是否为可执行文件，如果是则值为true，反则为false
	if isDir == true {
		return FolderColor(fileName)
	} else {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if runtime.GOOS == "windows" && ext == ".exe" || ext == ".msi" || ext == ".bat" || ext == ".cmd" || ext == ".ps1" {
			return ApplicationColor(fileName)
		} else {
			if isExecutable == true {
				return ApplicationColor(fileName)
			}
		}
	}
	return fileName
}

// func ColorFormatter(file os.DirEntry, isDir bool) *color.Color {
// 	var (
// 		isExecutable bool
// 		// TitleColor       = color.New(color.FgHiGreen, color.Bold).SprintFunc()
// 		ExecutableColor = color.New(color.FgHiGreen)
// 		FolderColor     = color.New(color.FgBlue)
// 		// FileColor        = color.New(color.FgWhite).SprintFunc()
// 	)
// 	// fileName := file.Name()
// 	info, _ := file.Info()
// 	isExecutable = info.Mode().Perm()&0111 != 0 // 判断文件是否为可执行文件，如果是则值为true，反则为false
// 	if isDir == true {
// 		return FolderColor
// 	} else {
// 		ext := strings.ToLower(filepath.Ext(file.Name()))
// 		if runtime.GOOS == "windows" {
// 			if ext == ".exe" || ext == ".msi" || ext == ".bat" || ext == ".cmd" || ext == ".ps1" {
// 				return ExecutableColor
// 			}
// 		} else {
// 			if isExecutable == true {
// 				return ExecutableColor
// 			}
// 		}
// 	}
// 	return nil
// }
