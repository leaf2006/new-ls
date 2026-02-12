package output

import (
	"fmt"

	"github.com/leaf2006/new-ls/internal/core"
)

func SimpleOutput() {
	// TODO:未完工，仅供测试
	rows := core.Rows
	for count, row := range rows {
		if count == len(rows)-1 {
			fmt.Printf(" %s %s\n", row.Icon, row.Name)
		} else {
			fmt.Printf(" %s %s", row.Icon, row.Name)
		}
	}
}
