package spider

import (
	"movieSpider/pkg/config"
	"movieSpider/pkg/spider/douban"
	"movieSpider/pkg/spider/wmdb"
)

type Spider interface {
	Run()
}

func spiderTask(spiders ...Spider) {
	for _, f := range spiders {
		go f.Run()
	}
}

func RunSpider() {
	// Spider
	if config.DouBan != nil {
		if config.DouBan.DoubanUrl != "" {
			spiderDouBan := douban.NewSpiderDouBan(config.DouBan.DoubanUrl, config.DouBan.Scheduling)
			spiderTask(spiderDouBan)
		}
		spiderWmdb := wmdb.NewSpiderWmdb(config.DouBan.Scheduling)
		spiderTask(spiderWmdb)
	}

}
