package core

import (
	"os"
	"strings"

	"github.com/leaf2006/new-ls/internal/render"
)

type FileRow struct {
	RawFile os.DirEntry
	Mode    string
	Time    string
	Size    string
	Name    string
	Icon    string
	IsDir   bool
}

var Rows []FileRow
var actualFilePath string
var MaxSizeLen int
var MaxFileNameLen int

func Entry(filepath string, enableAllFiles bool, enableEntrySimple bool) ([]FileRow, error) {
	// var FilePath string
	if filepath != "" {
		actualFilePath = filepath
	} else {
		actualFilePath = "."
	}

	files, err := os.ReadDir(actualFilePath)
	if err != nil {
		return nil, err
	}

	if enableEntrySimple == false {
		MaxSizeLen = 4
		for _, file := range files {
			if enableAllFiles == false && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			modeStr := FileMode(file)
			sizeStr, timeStr := FileInfo(file, actualFilePath)
			isDirBool := file.IsDir()

			// 更新最大宽度：如果当前文件大小字符串长度超过了目前的记录，就更新
			if len(sizeStr) > MaxSizeLen {
				MaxSizeLen = len(sizeStr)
			}

			Rows = append(Rows, FileRow{
				RawFile: file,
				Mode:    modeStr,
				Time:    timeStr,
				Size:    sizeStr,
				Name:    file.Name(),
				Icon:    render.IconMap(file),
				IsDir:   isDirBool,
			})
		}
		return Rows, nil
	} else {
		MaxFileNameLen = 8 //原来是8
		for _, file := range files {
			if enableAllFiles == false && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			fileName := file.Name()
			isDirBool := file.IsDir()
			if len(fileName) > MaxFileNameLen {
				MaxFileNameLen = len(fileName) + 2 // 为图标预留2个空位
			}

			Rows = append(Rows, FileRow{
				RawFile: file,
				Name:    file.Name(),
				Icon:    render.IconMap(file),
				IsDir:   isDirBool,
			})
		}
		return Rows, nil
	}

}
