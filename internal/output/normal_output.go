package output

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/leaf2006/new-ls/internal/core"
	"github.com/leaf2006/new-ls/internal/render"
)

var TitleColor = color.New(color.FgHiGreen, color.Bold)

func NormalOutput() {
	rows := core.Rows
	isSimpleColor := false

	TitleColor.Printf("%-10s  %-16s  %-*s  %s\n", "Mode", "LastWriteTime", core.MaxSizeLen, "Size", "Name")
	dashSize := strings.Repeat("-", core.MaxSizeLen)
	TitleColor.Printf("%-10s  %-16s  %-*s  %s\n", "----", "-------------", core.MaxSizeLen, dashSize, "----")

	for _, row := range rows {
		var namePrinter *color.Color
		fmt.Printf("%-10s  ", row.Mode)
		fmt.Printf("%-16s  ", row.Time)
		fmt.Printf("%*s  ", core.MaxSizeLen, row.Size)

		namePrinter = render.ColorFormatter(row.RawFile, row.IsDir, isSimpleColor)
		if namePrinter != nil {
			namePrinter.Printf("%s %s", row.Icon, row.Name)
			fmt.Printf("\n") // 消除蓝色底色色块溢出的问题
		} else {
			fmt.Printf("%s %s\n", row.Icon, row.Name)
		}
	}
}
