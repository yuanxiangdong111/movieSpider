package config

import (
	"bytes"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
	"io/ioutil"
	"movieSpider/pkg/log"
	"movieSpider/pkg/types"
	"os"
)

type btbt struct {
	Url        string `json:"Url"`
	Scheduling string `json:"Scheduling"`
}

type eztv struct {
	Scheduling string `json:"Scheduling"`
	Url        string `json:"Url"`
}
type glodls struct {
	Scheduling string `json:"Scheduling"`
	Url        string `json:"Url"`
}

type rarbg struct {
	Scheduling   string `json:"Scheduling"`
	Url          string `json:"Url"`
	ResourceType string `json:"ResourceType"`
	Typ          types.Resource
}

type knaben struct {
	Url string `json:"Url"`
}
type downloader struct {
	Scheduling string `json:"Scheduling"`
	Aria2Label string `json:"Aria2Label"`
}
type bt4g struct {
	Url string `json:"Url"`
}

var (
	Global     global
	BtSpider   btSpider
	Aria2cList []aria2
	TG         = new(tg)
	MySQL      mysql
	DouBan     = new(douban)
	BTBT       = new(btbt)
	EZTV       = new(eztv)
	GLODLS     = new(glodls)
	KNABEN     = new(knaben)
	Bt4G       = new(bt4g)
	RARBG      []*rarbg
	Downloader *downloader
	ProxyPool  string
)

type global struct {
	LogLevel string
}

type btSpider struct {
	DBfile     string
	Scheduling string
	Port       string
}

type aria2 struct {
	Url   string
	Token string
	Label string
}

type tg struct {
	BotToken string
	TgIDs    []int
	Proxy    struct {
		Url    string
		Enable bool
	}
}

type mysql struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}
type douban struct {
	DoubanUrl  string
	Scheduling string
	WMDBPrefix string
}

func InitConfig(config string) {
	v := viper.New()

	fmt.Printf("config file is %s.\n", config)
	v.SetConfigType("yaml")
	b, err := ioutil.ReadFile(config)
	if err != nil {
		fmt.Printf("配置文件读取错误,err:%s", err.Error())
		os.Exit(1)
	}

	err = v.ReadConfig(bytes.NewReader(b))
	if err != nil {
		fmt.Printf("配置文件错误.")
		os.Exit(1)
	}

	err = v.UnmarshalKey("Global", &Global)
	if err != nil {
		fmt.Println("读取Global配置错误")
		os.Exit(-1)
	}
	if Global.LogLevel == "debug" {
		fmt.Println(string(b))
	}
	log.NewLogger(Global.LogLevel)

	err = v.UnmarshalKey("BtSpider", &BtSpider)
	if err != nil {
		fmt.Println("读取BtSpider配置错误")
		os.Exit(-1)
	}

	BtSpider.Port = ":" + BtSpider.Port
	err = v.UnmarshalKey("Aria2cList", &Aria2cList)
	if err != nil {
		fmt.Println("读取Aria2cList配置错误")
		os.Exit(-1)
	}
	for index, aria := range Aria2cList {
		Aria2cList[index].Url = aria.Url + "/jsonrpc"
	}

	err = v.UnmarshalKey("TG", &TG)
	if err != nil {
		fmt.Println("读取TG配置错误")
		os.Exit(-1)
	}

	// 2022
	err = v.UnmarshalKey("MySQL", &MySQL)
	if err != nil {
		fmt.Println("读取MySQL配置错误")
		os.Exit(-1)
	}

	if err = v.UnmarshalKey("Douban", &DouBan); err != nil {
		fmt.Println("读取Douban配置错误")
		os.Exit(-1)
	}
	if !govalidator.IsURL(DouBan.DoubanUrl) {
		DouBan.DoubanUrl = ""
	}
	if !govalidator.IsURL(DouBan.WMDBPrefix) {
		DouBan.WMDBPrefix = ""
	}

	// btbt
	if err = v.UnmarshalKey("Feed.BTBT", &BTBT); err != nil {
		fmt.Println("读取BTBT配置错误")
		os.Exit(-1)
	}
	if !govalidator.IsURL(BTBT.Url) {
		BTBT = nil
	}

	if err = v.UnmarshalKey("Feed.EZTV", &EZTV); err != nil {
		fmt.Println("读取EZTV配置错误")
		os.Exit(-1)
	}
	if !govalidator.IsURL(EZTV.Url) {
		EZTV = nil
	}

	if err = v.UnmarshalKey("Feed.GLODLS", &GLODLS); err != nil {
		fmt.Println("读取GLODLS配置错误")
		os.Exit(-1)
	}
	if !govalidator.IsURL(GLODLS.Url) {
		GLODLS = nil
	}

	if err = v.UnmarshalKey("Feed.KNABEN", &KNABEN); err != nil {
		fmt.Println("读取KNABEN配置错误")
		os.Exit(-1)
	}
	if !govalidator.IsURL(KNABEN.Url) {
		KNABEN = nil
	}

	if err = v.UnmarshalKey("Feed.Bt4G", &Bt4G); err != nil {
		fmt.Println("读取Bt4G配置错误")
		os.Exit(-1)
	}
	if !govalidator.IsURL(Bt4G.Url) {
		Bt4G = nil
	}

	if err = v.UnmarshalKey("Feed.RARBG", &RARBG); err != nil {
		fmt.Println("读取GLODLS配置错误")
		os.Exit(-1)
	}

	for _, v := range RARBG {
		if !govalidator.IsURL(v.Url) {
			v = nil
			continue
		}
		switch v.ResourceType {
		case "movie":
			v.Typ = types.ResourceMovie
		case "tv":
			v.Typ = types.ResourceTV
		default:
			v.Typ = types.ResourceTV
		}
	}

	if err = v.UnmarshalKey("Feed.ProxyPool", &ProxyPool); err != nil {
		fmt.Println("读取Feed.ProxyPool配置错误")
		os.Exit(-1)
	}

	if err = v.UnmarshalKey("Downloader", &Downloader); err != nil {
		fmt.Println("读取Downloader配置错误")
		os.Exit(-1)
	}

	return

}
