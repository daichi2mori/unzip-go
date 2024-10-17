package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	version string
}

var config Config
var cfgFile string = "config.yml"
var versionFlag bool

var rootCmd = &cobra.Command{
	Use:   "unzip",
	Short: "解凍フォルダを指定して実行できます",
	Long:  `選択した圧縮ファイルを解凍できます`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			version := viper.GetString("version")
			if version == "" {
				fmt.Println("バージョン情報が見つかりません")
			} else {
				fmt.Println("バージョン: ", version)
			}
			os.Exit(0)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "バージョン情報を表示します")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.ConfigFileUsed()
}
