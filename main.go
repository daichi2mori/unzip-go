package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/bodgit/sevenzip"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var version string = "v1.0.1"
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
			fmt.Println("\nすべてのファイルが解凍されました")
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

func extractArchive(archive string) error {
	r, err := sevenzip.OpenReader(archive)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if err = extractFile(f); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(f *sevenzip.File) error {
	// 解凍するファイルの読み取りストリームを開く
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// ファイル名のパスを作成
	outPath := filepath.Join(".", f.Name)

	// ディレクトリの場合、ディレクトリを作成して終了
	if f.FileInfo().IsDir() {
		return os.MkdirAll(outPath, f.Mode())
	}

	// ファイルを作成して中身を書き出す
	outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	// ストリームからファイルへコピー
	_, err = io.Copy(outFile, rc)
	if err != nil {
		return err
	}

	return nil
}
