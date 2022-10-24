package feed

import (
	"movieSpider/pkg/config"
	"movieSpider/pkg/log"
	"movieSpider/pkg/report"
	"movieSpider/pkg/types"
)

type Feeder interface {
	Crawler() ([]*types.FeedVideo, error)
	Run()
}

type FeederAbstractFactory interface {
	CreateFeeder(args ...interface{}) Feeder
}

func RunFeed() {

	// BTBT
	facFeedBTBT := new(FactoryBTBT)
	feedBTBT := facFeedBTBT.CreateFeeder(config.BTBT.Scheduling)

	// EZTV
	facFeedEZTV := new(FactoryEZTV)
	feedEZTV := facFeedEZTV.CreateFeeder(config.EZTV.Scheduling)

	// rarbg TV
	facFeedRarbgTV := new(FactoryRARBG)
	var feedRarbgTV Feeder

	// rarbg Movie
	facFeedRarbgMovie := new(FactoryRARBG)
	var feedRarbgMovie Feeder

	for _, r := range config.RARBG {
		if r != nil {
			if r.Typ == types.ResourceTV {
				feedRarbgTV = facFeedRarbgTV.CreateFeeder(r.Scheduling, r.Typ)
			}
			log.Debug(r)
			if r.Typ == types.ResourceMovie {
				feedRarbgMovie = facFeedRarbgMovie.CreateFeeder(r.Scheduling, r.Typ)
			}
			log.Debug(r)
		}
	}

	// GLODLS
	facFeedGLODLS := new(FactoryGLODLS)
	feedGLODLS := facFeedGLODLS.CreateFeeder(config.GLODLS.Scheduling)

	// TGX
	facFeedTGX := new(FactoryTGX)
	feedTGXS := facFeedTGX.CreateFeeder(config.TGX.Scheduling)

	feederTask(feedBTBT, feedEZTV, feedRarbgTV, feedRarbgMovie, feedGLODLS, feedTGXS)

	if config.Global.Report {
		report.NewReport("*/1 * * * *").Run()
	}
}

func feederTask(feeds ...Feeder) {
	for _, f := range feeds {
		go f.Run()
	}
}
