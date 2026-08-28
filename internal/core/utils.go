package core

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func FileInfo(file os.DirEntry, path string, enableByteOutput bool) (string, string, time.Time) { // 文件大小 ， 文件上次修改时间 ， 修改时间原始值(用于按时间排序)
	var fileName string

	if path == "." {
		fileName = file.Name()
	} else {
		fileName = filepath.Join(path, file.Name())
	}

	info, err := os.Stat(fileName)

	if err != nil {
		errorMessage := fmt.Sprintf("ERR:%s", err)
		return errorMessage, errorMessage, time.Time{}
	}

	timeFormatter := "2006/01/02 15:04"

	lastWriteTimeRaw := info.ModTime()
	lastWriteTime := lastWriteTimeRaw.Format(timeFormatter)

	if file.IsDir() {
		return "", lastWriteTime, lastWriteTimeRaw
	} else {
		fileSize := strconv.FormatInt(info.Size(), 10) //10为转为十进制
		if enableByteOutput == false {
			fileSize = FormatFileSize(fileSize)
		}
		return fileSize, lastWriteTime, lastWriteTimeRaw
	}
}

func FileMode(file os.DirEntry) string {
	info, _ := file.Info()
	fileMode := info.Mode().String()
	return fileMode
}

func FormatFileSize(bytesStr string) string {
	bytes, _ := strconv.ParseInt(bytesStr, 10, 64)
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	e := math.Floor(math.Log(float64(bytes)) / math.Log(1024)) // 计算指数 e = log1024(b)
	suffix := units[int(e)]                                    // 计算对应的单位后缀
	val := float64(bytes) / math.Pow(1024, e)                  // 计算在该单位下的数值

	if val >= 10 {
		return fmt.Sprintf("%.0f %s", val, suffix)
	}
	return fmt.Sprintf("%.1f %s", val, suffix)

}
