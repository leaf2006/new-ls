package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/leaf2006/new-ls/internal/core"
	"github.com/leaf2006/new-ls/internal/render"

	"github.com/leaf2006/new-ls/internal/config"
)

var TitleColor = color.New(color.FgHiGreen, color.Bold)

func NormalOutput() {
	rows := core.Rows
	if len(rows) > 0 {

		isSimpleColor := false // 控制是否为简单输出，当前为正常输出，所以是false
		TitleColor.Printf("%-10s    %-16s    %-*s    %s\n", "Mode", "LastWriteTime", core.MaxSizeLen, "Size", "Name")
		dashSize := strings.Repeat("-", core.MaxSizeLen)
		TitleColor.Printf("%-10s    %-16s    %-*s    %s\n", "----", "-------------", core.MaxSizeLen, dashSize, "----")

		// 计算 Name 列起始位置：Mode(10) + 4 + Time(16) + 4 + Size(MaxSizeLen) + 4
		prefixWidth := 10 + 4 + 16 + 4 + core.MaxSizeLen + 4

		// 获取终端宽度，用于计算文件名可占用的最大宽度
		terminalWidth := 80
		if w, isTTY := render.GetTerminalWidth(); isTTY {
			terminalWidth = w
		}
		maxNameWidth := terminalWidth - prefixWidth
		if maxNameWidth < 20 {
			maxNameWidth = 20
		}

		for _, row := range rows {
			var namePrinter *color.Color
			fmt.Printf("%-10s    ", row.Mode)
			fmt.Printf("%-16s    ", row.Time)
			fmt.Printf("%*s    ", core.MaxSizeLen, row.Size)

			namePrinter = render.ColorFormatter(row.RawFile, row.IsDir, isSimpleColor)
			var EnableIcon string
			if config.Global.Icon == true {
				EnableIcon = row.Icon + " "
			} else {
				EnableIcon = ""
			}

			fullDisplayName := EnableIcon + row.Name
			if len(fullDisplayName) <= maxNameWidth {
				// 文件名未超宽，直接输出
				if namePrinter != nil {
					namePrinter.Printf("%s%s", EnableIcon, row.Name)
					fmt.Printf("\n")
				} else {
					fmt.Printf("%s%s\n", EnableIcon, row.Name)
				}
			} else {
				// 文件名超宽，按 maxNameWidth 分段，换行并对齐到 Name 列
				remaining := fullDisplayName
				firstLine := true
				for len(remaining) > 0 {
					end := maxNameWidth
					if end > len(remaining) {
						end = len(remaining)
					}
					chunk := remaining[:end]
					remaining = remaining[end:]

					if !firstLine {
						fmt.Printf("%*s", prefixWidth, "")
					}
					if namePrinter != nil {
						namePrinter.Printf("%s", chunk)
					} else {
						fmt.Printf("%s", chunk)
					}
					fmt.Printf("\n")
					firstLine = false
				}
			}
		}
	}
}
