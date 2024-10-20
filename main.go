package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gen2brain/go-unarr"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var version string = "v1.1.0"
var versionFlag bool

var all bool

var extractableExtensions = []string{".7z", ".zip", ".rar", ".tar", ".gz", ".bz2", ".xz"}
var selectedFilesToExtract []string

var rootCmd = &cobra.Command{
	Use:   "unzip",
	Short: "解凍フォルダを指定して実行できます",
	Long:  `選択した圧縮ファイルを解凍できます`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Println("バージョン: ", version)
			os.Exit(0)
		}

		files, err := filepath.Glob("*")
		if err != nil {
			fmt.Println("ファイルを取得できませんでした:", err)
			return
		}

		for _, file := range files {
			if isExtractableFile(file) {
				if all || askUser(file) {
					selectedFilesToExtract = append(selectedFilesToExtract, file)
				}
			}
		}

		if len(selectedFilesToExtract) > 0 {
			bar := progressbar.Default(int64(len(selectedFilesToExtract)), "解凍中")
			for _, file := range selectedFilesToExtract {
				if err := extractArchive(file); err != nil {
					fmt.Printf("%sの解凍に失敗しました: %v\n", file, err)
				}
				bar.Add(1)
			}
			fmt.Println("\n終了しました")
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "バージョン情報を表示します")
	rootCmd.Flags().BoolVarP(&all, "all", "a", false, "すべてのファイルを解凍しますか？")
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func isExtractableFile(fileName string) bool {
	ext := strings.ToLower(filepath.Ext(fileName))
	for _, extractableExt := range extractableExtensions {
		if ext == extractableExt {
			return true
		}
	}
	return false
}

func askUser(file string) bool {
	var answer string
	prompt := &survey.Select{
		Message: fmt.Sprintf("%sを解凍しますか？", file),
		Options: []string{"Yes", "No"},
	}
	survey.AskOne(prompt, &answer)
	return answer == "Yes"
}

func extractArchive(file string) error {
	archive, err := unarr.NewArchive(file)
	if err != nil {
		return err
	}
	defer archive.Close()

	dirName := strings.TrimSuffix(file, filepath.Ext(file)) // フォルダ名取得
	if err = os.Mkdir(dirName, os.ModePerm); err != nil {
		return err
	}

	_, err = archive.Extract(dirName)
	if err != nil {
		return err
	}

	return nil
}
