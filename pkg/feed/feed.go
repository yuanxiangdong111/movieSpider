package feed

import (
	"movieSpider/pkg/config"
	"movieSpider/pkg/feed/btbt"
	"movieSpider/pkg/feed/eztv"
	"movieSpider/pkg/feed/glodls"
	"movieSpider/pkg/feed/rarbg"
	"movieSpider/pkg/feed/tgx"
	"movieSpider/pkg/log"
	"movieSpider/pkg/report"
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
		feedEztv := eztv.NewFeedEztv(config.EZTV.Scheduling)
		feederTask(feedEztv)
	}

	for _, r := range config.RARBG {
		if r != nil {
			log.Debug(r)
			feedRarbg := rarbg.NewFeedRarbg(r.Scheduling, r.Typ)
			feederTask(feedRarbg)
		}
	}
	if config.GLODLS != nil {
		feedGlodls := glodls.NewFeedGlodls(config.GLODLS.Scheduling)
		feederTask(feedGlodls)
	}
	if config.BTBT != nil {
		feedBTBT := btbt.NewFeedBTBT(config.BTBT.Scheduling)
		feederTask(feedBTBT)
	}

	if config.TGX != nil {
		feedBTBT := tgx.NewTGx(config.TGX.Scheduling)
		feederTask(feedBTBT)
	}

	if config.Global.Report {
		report.NewReport("*/1 * * * *").Run()
	}
}
