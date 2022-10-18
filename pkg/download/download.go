package download

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/pkg"
	"movieSpider/pkg/aria2"
	"movieSpider/pkg/config"
	"movieSpider/pkg/feed/bt4g"
	"movieSpider/pkg/feed/knaben"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"os"
)

type download struct {
	scheduling string
	types.Resolution
}

func NewDownloader(scheduling string) *download {
	return &download{scheduling: scheduling}
}

func (d *download) downloadTask() {
	err := d.downloadMovieTask()
	pkg.CheckError("Downloader", err)
	err = d.downloadTvTask()
	pkg.CheckError("Downloader", err)
}

func (d *download) downloadTvTask() (err error) {
	log.Info("Downloader tv.")
	tvs, err := model.MovieDB.FetchDouBanVideoByType(types.ResourceTV)
	if err != nil {
		pkg.CheckError("Downloader", err)
	}

	// 获取 磁力连接
	tvVides, err := model.MovieDB.FetchMagnetByName(tvs)
	if err != nil {
		pkg.CheckError("Downloader", err)
	}
	if len(tvVides) == 0 {
		log.Warn("Downloader: 此次没有查询到要下载的资源.")
		return
	}

	// 推送 磁力连接至 aria2
	err = d.aria2Download(tvVides)
	return
}

func (d *download) downloadMovieTask() error {
	// 获取 豆瓣 数据
	log.Info("Downloader movie.")
	names, err := model.MovieDB.FetchDouBanVideoByType(types.ResourceMovie)
	if err != nil {
		pkg.CheckError("Downloader", err)
	}

	// 获取 磁力连接
	MovieVides, err := model.MovieDB.FetchMagnetByName(names)
	if err != nil {
		pkg.CheckError("Downloader", err)
	}
	if len(MovieVides) == 0 {
		log.Warn("Downloader: 此次没有查询到要下载的资源.")
		return nil
	}

	// 推送 磁力连接至 aria2
	err = d.aria2Download(MovieVides)
	pkg.CheckError("Downloader", err)
	return err
}

func (d *download) aria2Download(vides []*types.FeedVideo) (err error) {

	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		pkg.CheckError("Downloader", errors.WithMessage(pkg.ErrDownloaderAria2, err.Error()))
	}
	for _, v := range vides {
		gid, err := newAria2.DownloadByUrl(v.Magnet)
		if err != nil {
			return err
		}
		err = model.MovieDB.UpdateFeedVideoDownloadByID(v.ID)
		pkg.CheckError("Downloader", err)
		log.Infof("Downloader: %s 开始下载. GID: %s", v.Name, gid)
	}
	return nil
}

func (d *download) Run() {
	if d.scheduling == "" {
		log.Error("Downloader: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("Downloader: Scheduling is: [%s]", d.scheduling)
	c := cron.New()
	_, err := c.AddFunc(d.scheduling, func() {
		d.downloadTask()
	})
	if err != nil {
		log.Error("Downloader: AddFunc is null")
		os.Exit(1)
	}
	c.Start()
}

func (d *download) DownloadByName(name, Resolution string) (msg string) {
	// 从 knaben 搜索
	feedKnaben := knaben.NewFeedKnaben(config.KNABEN.Url, name, d.ResolutionStr2Int(Resolution))
	_, err := feedKnaben.Crawler()
	pkg.CheckError("Downloader", err)

	// 从 Bt4g 搜索
	feedBt4g := bt4g.NewFeedBt4g(config.Bt4G.Url, name, d.ResolutionStr2Int(Resolution))
	_, err = feedBt4g.Crawler()
	pkg.CheckError("Downloader", err)

	// 获取 磁力连接
	vides, err := model.MovieDB.FetchMagnetByName([]string{name})
	if err != nil {
		pkg.CheckError("Downloader", err)
	}

	if len(vides) == 0 {
		return fmt.Sprint("所有资源已下载过,或没有可下载资源.")
	}

	// 推送 磁力连接至 aria2
	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	pkg.CheckError("Downloader", err)
	for _, v := range vides {
		if v.Magnet == "" {
			continue
		}

		gid, err := newAria2.DownloadByUrl(v.Magnet)
		if err != nil {
			log.Error(err)
			return
		}
		err = model.MovieDB.UpdateFeedVideoDownloadByID(v.ID)
		pkg.CheckError("Downloader", err)
		log.Infof("Downloader: %s 开始下载. GID: %s", v.Name, gid)
	}

	return fmt.Sprintf("已将 %d 资源加入下载.", len(vides))
}
