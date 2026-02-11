package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	// "text/tabwriter"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
)

var (
	version    = "v0.0.1"
	TitleColor = color.New(color.FgHiGreen, color.Bold)
	// ApplicationColor = color.New(color.FgHiGreen)
	// FolderColor      = color.New(color.FgBlue)
	// FileColor        = color.New(color.FgWhite)
)

type FileRow struct {
	RawFile os.DirEntry
	Mode    string
	Time    string
	Size    string
	Name    string
	Icon    string
	IsDir   bool
}

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
				var isSimpleColor bool
				enableOutputAllFiles = false

				if isOutputAllFiles {
					enableOutputAllFiles = true
				}

				if Args.Len() != 0 {
					FilePath = Args.First()

				} else {
					FilePath = "."
				}

				files, err := os.ReadDir(FilePath)
				if err != nil {
					// fmt.Printf("Error:%s", err)
					return err
				}

				var rows []FileRow
				maxSizeLen := 4 //default size length
				for _, file := range files {
					if enableOutputAllFiles == false && strings.HasPrefix(file.Name(), ".") {
						continue
					}

					modeStr := FileMode(file)
					sizeStr, timeStr := FileInfo(file, FilePath)
					isDirBool := file.IsDir()
					// 更新最大宽度：如果当前文件大小字符串长度超过了目前的记录，就更新
					if len(sizeStr) > maxSizeLen {
						maxSizeLen = len(sizeStr)
					}

					rows = append(rows, FileRow{
						RawFile: file,
						Mode:    modeStr,
						Time:    timeStr,
						Size:    sizeStr,
						// Name: ColorFormatter(file, isDirBool),
						Name:  file.Name(),
						Icon:  IconMap(file),
						IsDir: isDirBool,
					})
				}
				if isSimple {
					isSimpleColor = true
					// TODO 暂时先不搞
				}
				// normal output
				isSimpleColor = false
				TitleColor.Printf("%-10s  %-16s  %-*s  %s\n", "Mode", "LastWriteTime", maxSizeLen, "Size", "Name")
				dashSize := strings.Repeat("-", maxSizeLen)
				TitleColor.Printf("%-10s  %-16s  %-*s  %s\n", "--------", "-------------", maxSizeLen, dashSize, "----")

				for _, row := range rows {
					var namePrinter *color.Color
					fmt.Printf("%-10s  ", row.Mode)
					fmt.Printf("%-16s  ", row.Time)
					fmt.Printf("%*s  ", maxSizeLen, row.Size)

					namePrinter = ColorFormatter(row.RawFile, row.IsDir, isSimpleColor)
					if namePrinter != nil {
						namePrinter.Printf("%s %s", row.Icon, row.Name)
						fmt.Printf("\n") //消除蓝色底色色块溢出的问题
					} else {
						fmt.Printf("%s %s\n", row.Icon, row.Name)
					}

				}
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
