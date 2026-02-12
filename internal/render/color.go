package render

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fatih/color"
)

func ColorFormatter(file os.DirEntry, isDir bool, isSimpleColor bool) *color.Color {
	var (
		isExecutable bool
		// TitleColor       = color.New(color.FgHiGreen, color.Bold).SprintFunc()
		ExecutableColor   = color.New(color.FgHiGreen)
		FolderColorSimple = color.New(color.FgBlue)             // 普通蓝色
		FolderColor       = color.New(color.BgBlue, color.Bold) // 背景蓝色（powershell-style）
		// FileColor        = color.New(color.FgWhite).SprintFunc()
	)
	// fileName := file.Name()
	info, _ := file.Info()
	isExecutable = info.Mode().Perm()&0111 != 0 // 判断文件是否为可执行文件，如果是则值为true，反则为false
	if isDir == true {
		if isSimpleColor == true {
			return FolderColorSimple
		} else {
			return FolderColor
		}
	} else {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if runtime.GOOS == "windows" {
			if ext == ".exe" || ext == ".msi" || ext == ".bat" || ext == ".cmd" || ext == ".ps1" {
				return ExecutableColor
			}
		} else {
			if isExecutable == true {
				return ExecutableColor
			}
		}
	}
	return nil // 普通颜色输出
}
