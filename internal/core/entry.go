package core

import (
	"os"
	"sort"
	"strings"
	"time"

	"github.com/leaf2006/new-ls/internal/render"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type FileRow struct {
	RawFile os.DirEntry
	Mode    string
	Time    string
	ModTime time.Time // 修改时间原始值，用于 -t/-tr 按时间排序
	Size    string
	Name    string
	Icon    string
	IsDir   bool
}

var Rows []FileRow
var actualFilePath string
var MaxSizeLen int
var MaxFileNameLen int

func Entry(filepath string, enableAllFiles bool, enableEntrySimple bool, enableByteOutput bool, enableTimeSort bool, timeNewestFirst bool) ([]FileRow, error) {
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

	if enableEntrySimple == false { // 详细输出
		MaxSizeLen = 4
		for _, file := range files {
			// 如果文件前带.就是隐藏文件，如果没有开-A就不显示
			if enableAllFiles == false && strings.HasPrefix(file.Name(), ".") {
				continue
			}

			modeStr := FileMode(file)
			sizeStr, timeStr, modTime := FileInfo(file, actualFilePath, enableByteOutput)

			// if enableByteOutput == false {
			// 	sizeStr = FormatFileSize(sizeStr)
			// }

			isDirBool := file.IsDir()

			// 更新最大宽度：如果当前文件大小字符串长度超过了目前的记录，就更新
			if len(sizeStr) > MaxSizeLen {
				MaxSizeLen = len(sizeStr)
			}

			Rows = append(Rows, FileRow{
				RawFile: file,
				Mode:    modeStr,
				Time:    timeStr,
				ModTime: modTime,
				Size:    sizeStr,
				Name:    file.Name(),
				Icon:    render.IconMap(file),
				IsDir:   isDirBool,
			})
		}

		sortRows(enableTimeSort, timeNewestFirst)
		return Rows, nil
	} else { // 简化输出
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

			var modTime time.Time
			if enableTimeSort { // 仅在需要按时间排序时才获取修改时间，避免简单模式下的额外 stat 开销
				if info, err := file.Info(); err == nil {
					modTime = info.ModTime()
				}
			}

			Rows = append(Rows, FileRow{
				RawFile: file,
				ModTime: modTime,
				Name:    file.Name(),
				Icon:    render.IconMap(file),
				IsDir:   isDirBool,
			})
		}

		sortRows(enableTimeSort, timeNewestFirst)
		return Rows, nil
	}

}

// sortRows 对 Rows 进行排序：默认按文件名排序；enableTimeSort 为 true 时按修改时间排序，
// timeNewestFirst 为 true 表示新的在前（-t），false 表示旧的在前（-tr）。
// 修改时间相同时回退为按文件名排序，保证输出顺序稳定。
func sortRows(enableTimeSort bool, timeNewestFirst bool) {
	fileSort := collate.New(language.English) // 暂且用English进行排序
	sort.Slice(Rows, func(i, j int) bool {
		if enableTimeSort {
			ti, tj := Rows[i].ModTime, Rows[j].ModTime
			if !ti.Equal(tj) {
				if timeNewestFirst {
					return ti.After(tj)
				}
				return ti.Before(tj)
			}
		}
		return fileSort.CompareString(Rows[i].Name, Rows[j].Name) < 0
	})
}
