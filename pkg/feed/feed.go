package feed

import (
	"movieSpider/pkg/config"
	"movieSpider/pkg/feed/btbt"
	"movieSpider/pkg/feed/eztv"
	"movieSpider/pkg/feed/glodls"
	"movieSpider/pkg/feed/rarbg"
	"movieSpider/pkg/log"
	"movieSpider/pkg/types"
)

type Feeder interface {
	Crawler() ([]*types.FeedVideo, error)
	Run()
}

func feederTask(feeds ...Feeder) {
	for _, f := range feeds {
		go f.Run()
	}
}

func RunFeed() {
	if config.EZTV != nil {
		feedEztv := eztv.NewFeedEztv(config.EZTV.Url, config.EZTV.Scheduling)
		feederTask(feedEztv)
	}

	for _, r := range config.RARBG {
		if r != nil {
			log.Debug(r)
			feedRarbg := rarbg.NewFeedRarbg(r.Url, r.Scheduling, r.Typ)
			feederTask(feedRarbg)
		}
	}
	if config.GLODLS != nil {
		feedGlodls := glodls.NewFeedGlodls(config.GLODLS.Url, config.GLODLS.Scheduling)
		feederTask(feedGlodls)
	}
	if config.BTBT != nil {
		feedBTBT := btbt.NewFeedBTBT(config.BTBT.Url, config.BTBT.Scheduling)
		feederTask(feedBTBT)
	}

}
