package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"movieSpider/pkg/config"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"os"
)

var (
	cfgTmp = `MySQL:
  Host: 127.0.0.1
  Port: 3306
  Database: movie
  User: root
  Password: PASSWORD

Douban:
  DoubanUrl: "https://movie.douban.com/people/251312920/wish"
  Scheduling: "*/5 * * * *"
  WMDBPrefix: "https://api.wmdb.tv/movie/api?id="
Feed:
  BTBT:
    Url: https://www.btbtt12.com/forum-index-fid-951.htm
    #    Url: https://www.btbtt12.com/forum-index-fid-951-page-5.htm
    Scheduling: "*/5 * * * *"
  EZTV:
    Scheduling: "*/5 * * * *"
    Url: "https://eztv.re/ezrss.xml"
  GLODLS:
    Scheduling: "*/3 * * * *"
    Url: "https://glodls.to/rss.php?cat=1,41"
  RARBG:
    - Scheduling: "*/5 * * * *"
      Url: "https://rarbgprx.org/rssdd.php?categories=14;15;16;17;21;22;42;44;45;46;47;48"
      ResourceType: movie
    - Scheduling: "*/5 * * * *"
      Url: "https://rarbg.to/rssdd.php?categories=14;15;16;17;21;22;42;44;45;46;47;48"
      ResourceType: tv
  KNABEN:
    Url: "https://rss.knaben.eu"
  Bt4G:
    Url: "https://bt4g.org"

Global:
  LogLevel: info
#  Proxy: socks5://127.0.0.1:1089
Downloader:
  Scheduling: "*/1 * * * *"
  Aria2Label: "home"

Aria2cList:
  - Url: "http://127.0.0.1:6800"
    Token: TOKEN
    Label: home
  - Url: "http://127.0.0.1:6801"
    Token: TOKEN
    Label: nas

TG:
  BotToken: "TOKEN"
  TgIDs: [ 123456 ]
`
	outFile string
	initDB  bool
)
var configCmd = &cobra.Command{
	Use:   "config",
	Short: fmt.Sprintf("generate %s config file.", Name),
	Run: func(cmd *cobra.Command, args []string) {
		switch initDB {
		case false:
			if outFile == "" {
				fmt.Println(cfgTmp)
			}
			if outFile != "" {
				err := ioutil.WriteFile(outFile, []byte(cfgTmp), 0644)
				if err != nil {
					log.Error(err)
					os.Exit(-1)
				}
			}
		case true:
			config.InitConfig(configFile)
			model.NewMovieDB()
			err := model.MovieDB.InitDBTable()
			if err != nil {
				log.Error(err)
				return
			}
			log.Infof("db: %s 数据库初始化完毕.", config.MySQL.Database)
		}

	},
}

func init() {
	configCmd.Flags().StringVarP(&outFile, "out.file", "o", "", "指定输出的文件")
	configCmd.Flags().BoolVar(&initDB, "init.db", false, "初始化DB")
}
