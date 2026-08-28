package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/leaf2006/new-ls/internal/config" // nls全局设置
	"github.com/leaf2006/new-ls/internal/core"
	"github.com/leaf2006/new-ls/internal/output"
	"github.com/urfave/cli/v3"
)

const version = "v0.0.5"

var (
	Args                 cli.Args
	enableOutputAllFiles bool
	enableByteOutput     bool
)

func Commands() {
	config.Load()

	cmd := &cli.Command{
		Name:      "new-ls",
		Usage:     "A beautiful,powershell-style ls tool built from Golang",
		ArgsUsage: "[path...]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print only the version",
			},
			&cli.BoolFlag{
				Name:    "simple",
				Aliases: []string{"s"},
				Usage:   "Simple the output",
			},
			&cli.BoolFlag{
				Name:    "All",
				Aliases: []string{"A"},
				Usage:   "Output all the files",
			},
			&cli.BoolFlag{
				Name:    "byte",
				Aliases: []string{"b"},
				Usage:   "Print sizes in bytes; not available in -s mode",
			},
			&cli.BoolFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "Sort by modification time, newest first",
			},
			&cli.BoolFlag{
				Name:    "time-reverse",
				Aliases: []string{"tr"},
				Usage:   "Sort by modification time, oldest first",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			outputVersion := cmd.Bool("version")
			isSimple := cmd.Bool("simple")
			isOutputAllFiles := cmd.Bool("All")
			isByteOutput := cmd.Bool("byte")
			isTimeSort := cmd.Bool("time")
			isTimeReverse := cmd.Bool("time-reverse")

			if outputVersion {
				fmt.Printf("New-ls \nA beautiful,powershell-style ls tool built from Golang \nVersion: %s \nCopyright Leafdeveloper(C) 2026 \n", version)
			} else {
				Args = cmd.Args() // 在internal/core/entry.go中引用
				enableOutputAllFiles = false
				enableByteOutput = false

				if isOutputAllFiles || config.Global.AllFile == true {
					enableOutputAllFiles = true
				} // 输出所有文件/文件夹（输出隐藏文件）

				if isByteOutput || config.Global.ByteOutput == true {
					enableByteOutput = true
				} // 以字节形式输出文件大小，默认会以更符合人类日常习惯的方式输出

				// 按修改时间排序：-t 为新到旧，-tr 为旧到新；两者同时给出时以 -tr 为准
				enableTimeSort := isTimeSort || isTimeReverse
				timeNewestFirst := isTimeSort && !isTimeReverse

				var filepath string
				var enableEntrySimple bool
				if Args.Len() > 0 {
					filepath = Args.First()
				}
				// _, err := core.Entry(filepath, enableOutputAllFiles) //传递至internal/core/entry.go
				// if err != nil {
				// 	return err
				// }

				if isSimple || config.Global.Simple == true {
					enableEntrySimple = true
					_, err := core.Entry(filepath, enableOutputAllFiles, enableEntrySimple, enableByteOutput, enableTimeSort, timeNewestFirst) //传递至internal/core/entry.go
					if err != nil {
						return err
					}
					output.SimpleOutput()
					return nil
				}

				enableEntrySimple = false                                                                                                        // normal output
				_, err := core.Entry(filepath, enableOutputAllFiles, enableEntrySimple, enableByteOutput, enableTimeSort, timeNewestFirst) //传递至internal/core/entry.go
				if err != nil {
					return err
				}
				output.NormalOutput()
			}

			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
