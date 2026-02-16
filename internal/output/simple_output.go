package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/leaf2006/new-ls/internal/core"
	"github.com/leaf2006/new-ls/internal/render"
)

var columnsCount int

func SimpleOutput() {
	// TODO:未完工，仅供测试
	rows := core.Rows
	isSimpleColor := true
	TTYWidth, isTTY := render.GetTerminalWidth()
	if isTTY != true { //对于当前若不是终端环境的异常处理
		fmt.Fprintln(os.Stderr, "ERR:Not running in a terminal environment!")
		os.Exit(1)
	}

	maxSingleFileLength := core.MaxFileNameLen + 2 // 左右两个空格
	columnsCount = TTYWidth / maxSingleFileLength  // 每行可以打印colunmCount个文件

	// 计算每列的实际宽度，用于更好的对齐
	columnWidth := TTYWidth / columnsCount
	if columnsCount <= 0 {
		columnsCount = 1
		columnWidth = TTYWidth
	}

	for count, row := range rows {
		var namePrinter *color.Color
		namePrinter = render.ColorFormatter(row.RawFile, row.IsDir, isSimpleColor)

		// 计算需要的填充空格数
		iconAndNameLen := len(row.Icon) + len(row.Name) + 1 // +1 for space between icon and name
		padding := columnWidth - iconAndNameLen

		// 确保不会出现负数填充
		if padding < 0 {
			padding = 0
		}

		// 格式化输出
		if namePrinter != nil {
			// 彩色输出：图标 + 文件名 + 填充空格
			namePrinter.Printf("%s %s", row.Icon, row.Name)
			// 添加填充空格
			fmt.Printf("%*s", padding, "")
		} else {
			// 非彩色输出
			fmt.Printf("%s %s%*s", row.Icon, row.Name, padding, "")
		}

		// 换行条件：达到列数限制或最后一个元素
		if (count+1)%columnsCount == 0 || count == len(rows)-1 {
			fmt.Printf("\n")
		}
	}
}
