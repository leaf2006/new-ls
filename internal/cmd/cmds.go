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
				Usage:   "print only the version",
			},
			&cli.BoolFlag{
				Name:    "simple",
				Aliases: []string{"s"},
				Usage:   "simple the output",
			},
			&cli.BoolFlag{
				Name:    "All",
				Aliases: []string{"A"},
				Usage:   "output all the files",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			outputVersion := cmd.Bool("version")
			isSimple := cmd.Bool("simple")
			isOutputAllFiles := cmd.Bool("All")

			if outputVersion {
				fmt.Printf("New-ls \nA beautiful,powershell-style ls tool built from Golang \nVersion: %s \nCopyright Leafdeveloper(C) 2026 \n", version)
			} else {
				Args = cmd.Args() // 在internal/core/entry.go中引用
				enableOutputAllFiles = false

				if isOutputAllFiles {
					enableOutputAllFiles = true
				}

				var filepath string
				if Args.Len() > 0 {
					filepath = Args.First()
				}
				_, err := core.Entry(filepath, enableOutputAllFiles) //传递至internal/core/entry.go
				if err != nil {
					return err
				}

				if isSimple {
					output.SimpleOutput()
					return nil
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
