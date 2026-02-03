package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
