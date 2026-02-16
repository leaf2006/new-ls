package render

import (
	"os"

	"golang.org/x/term"
)

func GetTerminalWidth() (width int, isTTY bool) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return 0, false // 未检测到终端
	}
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, true
	}
	return width, true
}
