package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"movieSpider/pkg/bot"
	"movieSpider/pkg/config"
	"movieSpider/pkg/download"
	"movieSpider/pkg/feed"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/spider"
	"os"
)

var (
	Name       = "movieSpider"
	configFile string
	runBotFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   Name,
	Short: fmt.Sprintf("%s 电影助手，自动获取电影种子信息，自动刮取豆瓣电影想看列表，自动下载", Name),

	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig(configFile)
		model.NewMovieDB()
		RunMovieSpider()
		switch {
		case runBotFlag == true:
			if config.TG != nil {
				bot := bot.NewTgBot(config.TG.BotToken, config.TG.TgIDs)
				bot.StartBot()
			} else {
				fmt.Println("请设置TG参数")
				os.Exit(1)
			}

		}

		select {}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config.file", "f", "", "指定配置文件")
	rootCmd.Flags().BoolVar(&runBotFlag, "run.bot", false, "同时运行Telegram bot")

}

func RunMovieSpider() {

	// 	Feed
	feed.RunFeed()
	// Spider
	spider.RunSpider()
	// Downloader
	if config.Downloader != nil {
		download.NewDownloader(config.Downloader.Scheduling).Run()
	}
}
