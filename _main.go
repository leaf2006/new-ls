package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	// "github.com/fatih/color"
)

var (
	version = "v0.0.1"
)

// func OutputReadDir (filePath string) string {
// 	files ,err := os.ReadDir(filePath)
// 	if err != nil {
// 		var error string = fmt.Sprint(err)
// 		return error
// 	}
// 	for fileCount := 0; fileCount < len(files); fileCount++ {
// 		fileResult := files[fileCount].Name()
// 		return fileResult
// 	}
// 	return ""
// }

func main() {
	cmd := &cli.Command{
		Name:      "new-ls",
		Usage:     "list directory contents",
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

			printVersion := cmd.Bool("version")
			isSimple := cmd.Bool("simple")
			isOutputAllFiles := cmd.Bool("All")

			if printVersion {
				fmt.Printf("New-ls \n Version: %s \n Copyright Leafdeveloper(C) 2026 \n", version)
			} else {
				Args := cmd.Args()
				var FilePath string
				var enableOutputAllFiles bool

				if isOutputAllFiles {
					enableOutputAllFiles = true
				} else {
					enableOutputAllFiles = false
				}

				if Args.Len() != 0 {
					FilePath = Args.First() // 获取路径
				} else {
					FilePath = "."
				}

				files, err := os.ReadDir(FilePath)
				if err != nil {
					// fmt.Printf("Error:%s", err)
					return err
				}

				if isSimple { // 简化输出模式
					// simpleFormat := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0) //TODO:很有问题，需要获取终端宽度并计算
					for fileCount, file := range files {
						if enableOutputAllFiles == false && strings.HasPrefix(file.Name(), ".") { //TODO:我个人认为带"."的文件不需要隐藏
							continue
						}

						fileIcon := IconMap(file)
						output := fmt.Sprintf(" %s %s ", fileIcon, file.Name())
						if fileCount == len(files)-1 {
							fmt.Printf("%s\n", output)
							// fmt.Fprintf(simpleFormat, "%s\t\n", output)
						} else {
							fmt.Printf("%s", output)
							// fmt.Fprintf(simpleFormat, "%s\t", output)
						}

					}
					// simpleFormat.Flush()
				} else { //普通输出模式
					normalFormat := tabwriter.
				}
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
