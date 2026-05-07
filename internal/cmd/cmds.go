package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/leaf2006/new-ls/internal/core"
	"github.com/leaf2006/new-ls/internal/output"
	"github.com/urfave/cli/v3"
)

const version = "v0.0.1"

var (
	Args                 cli.Args
	enableOutputAllFiles bool
	enableByteOutput     bool
)

func Commands() {
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			outputVersion := cmd.Bool("version")
			isSimple := cmd.Bool("simple")
			isOutputAllFiles := cmd.Bool("All")
			isByteOutput := cmd.Bool("byte")

			if outputVersion {
				fmt.Printf("New-ls \nA beautiful,powershell-style ls tool built from Golang \nVersion: %s \nCopyright Leafdeveloper(C) 2026 \n", version)
			} else {
				Args = cmd.Args() // 在internal/core/entry.go中引用
				enableOutputAllFiles = false
				enableByteOutput = false

				if isOutputAllFiles {
					enableOutputAllFiles = true
				} // 输出所有文件/文件夹（输出隐藏文件）

				if isByteOutput {
					enableByteOutput = true
				} // 以字节形式输出文件大小，默认会以更符合人类日常习惯的方式输出

				var filepath string
				var enableEntrySimple bool
				if Args.Len() > 0 {
					filepath = Args.First()
				}
				// _, err := core.Entry(filepath, enableOutputAllFiles) //传递至internal/core/entry.go
				// if err != nil {
				// 	return err
				// }

				if isSimple {
					enableEntrySimple = true
					_, err := core.Entry(filepath, enableOutputAllFiles, enableEntrySimple, enableByteOutput) //传递至internal/core/entry.go
					if err != nil {
						return err
					}
					output.SimpleOutput()
					return nil
				}

				enableEntrySimple = false                                                                 // normal output
				_, err := core.Entry(filepath, enableOutputAllFiles, enableEntrySimple, enableByteOutput) //传递至internal/core/entry.go
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
