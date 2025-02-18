package download

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/aria2"
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"movieSpider/internal/log"
	"movieSpider/internal/model"
	"movieSpider/internal/spider/searchspider"
	"movieSpider/internal/types"
	"os"
	"sync"
)

type Download struct {
	scheduling string
	types.Resolution
}

func NewDownloader(scheduling string) *Download {
	//nolint:exhaustruct
	return &Download{scheduling: scheduling}
}

func (d *Download) downloadTask() {
	err := d.downloadMovieTask()
	if err != nil {
		log.Error(err)
	}
	err = d.downloadTvTask()
	if err != nil {
		log.Error(err)
	}
}

// downloadTvTask
//
//	@Description: 下载电视剧
//	@receiver d
//	@return err
//
//nolint:nakedret
func (d *Download) downloadTvTask() (err error) {
	log.Info("Downloader tv working...")
	videos, err := model.NewMovieDB().FetchDouBanVideoByType(types.VideoTypeTV)
	if err != nil {
		return errors.WithMessage(err, "FetchDouBanVideoByType")
	}

	//  FilterMap 暂存 电视剧名相同的视频
	var FilterMap = make(map[string][]*types.FeedVideo)
	log.Info("查找需要下载的tv.")

	// 归类同一个电视剧名的 feedVideo
	for douBanVideo, name := range videos {
		// 获取 tv
		videoList, err := model.NewMovieDB().GetFeedVideoTVByName(douBanVideo.DoubanID, name...)
		if err != nil {
			log.Warn(err)
		}
		if len(videoList) == 0 {
			continue
		}

		// 归类同一个电视剧名的视频
		for _, video := range videoList {
			//  将此次所有feedVideo的下载状态更新为3
			video.Download = 3

			err = model.NewMovieDB().UpdateFeedVideo(video)
			if err != nil {
				log.Error(err)
			}

			// 如果 feedVideo 不能转化为 downloadHistory 则跳过
			downloadHistory := video.Convert2DownloadHistory()
			if downloadHistory == nil {
				log.Debugf("TorrentName: %#v 不能转化为 downloadHistory ", video.TorrentName)
				continue
			}

			FilterMap[douBanVideo.Names] = append(FilterMap[douBanVideo.Names], video)
		}
	}

	// 根据 清晰度 季数和集数过滤
	needDownloadFeedVideo := make([]*types.FeedVideo, 0)
	for _, v := range FilterMap {
		list := FilterByResolution(types.VideoTypeTV, v...)
		needDownloadFeedVideo = append(needDownloadFeedVideo, list...)
	}

	//  如果没有需要下载的视频 则返回
	if len(needDownloadFeedVideo) == 0 {
		log.Warn("此次没有要下载的tv.")
		return
	}

	// 推送 磁力连接至 aria2
	err = d.aria2Download(needDownloadFeedVideo...)
	if err != nil {
		log.Error(err)
	}

	// 更新feedVideo的下载状态，记录这一次下载的视频
	for _, video := range needDownloadFeedVideo {
		UpdateFeedVideoAndDownloadHistory(video)
	}
	return
}

// downloadMovieTask
//
//	@Description: 下载电影
//	@receiver d
//	@return error
//
//nolint:nakedret
func (d *Download) downloadMovieTask() (err error) {
	// 获取 豆瓣 数据
	log.Info("Downloader movie working...")
	videos, err := model.NewMovieDB().FetchDouBanVideoByType(types.VideoTypeMovie)
	if err != nil {
		return errors.WithMessage(err, "FetchDouBanVideoByType")
	}

	//  FilterMap 暂存 电视剧名相同的视频
	var FilterMap = make(map[string][]*types.FeedVideo)
	log.Info("查找需要下载的movie.")

	// 归类同一个电视剧名的 feedVideo
	for douBanVideo, names := range videos {
		// 获取 feedVideo movie
		videoList, err := model.NewMovieDB().GetFeedVideoMovieByNameAndDoubanID(douBanVideo.DoubanID, names...)
		if err != nil {
			return errors.WithMessage(err, "douBanVideo: "+douBanVideo.Names)
		}
		if len(videoList) == 0 {
			log.Debugf("douBanVideo: %s 已全部下载完毕，或该影片没有更新.", douBanVideo.Names)
			continue
		}

		for _, video := range videoList {
			//  将此次所有feedVideo的下载状态更新为3
			video.Download = 3
			err = model.NewMovieDB().UpdateFeedVideo(video)
			if err != nil {
				log.Error(err)
			}
			// 如果 feedVideo 不能转化为 downloadHistory 则跳过
			downloadHistory := video.Convert2DownloadHistory()
			if downloadHistory == nil {
				log.Debugf("TorrentName: %#v 不能转化为 downloadHistory ", video.TorrentName)
				continue
			}
			FilterMap[douBanVideo.Names] = append(FilterMap[douBanVideo.Names], video)
		}
	}

	// 根据 清晰度 季数和集数过滤
	needDownloadFeedVideo := make([]*types.FeedVideo, 0)
	for _, v := range FilterMap {
		list := FilterByResolution(types.VideoTypeMovie, v...)
		needDownloadFeedVideo = append(needDownloadFeedVideo, list...)
	}
	//  如果没有需要下载的视频 则返回
	if len(needDownloadFeedVideo) == 0 {
		log.Warn("此次没有要下载的movie.")
		return
	}

	// 推送 磁力连接至 aria2
	err = d.aria2Download(needDownloadFeedVideo...)
	if err != nil {
		log.Warn(err)
	}

	for _, video := range needDownloadFeedVideo {
		log.Debugf("更新 %s ", video.TorrentName)
		UpdateFeedVideoAndDownloadHistory(video)
	}

	return err
}

// aria2Download
//
//	@Description: 通过aria2下载
//	@receiver d
//	@param videos
//	@return err
func (d *Download) aria2Download(videos ...*types.FeedVideo) error {
	if len(videos) < 1 {
		return errors.New("没有需要下载的视频")
	}
	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		return errors.WithMessage(err, "aria2 初始化失败")
	}
	for _, v := range videos {
		video, err := model.NewMovieDB().FetchOneDouBanVideoByDouBanID(v.DoubanID)
		if err != nil {
			log.Error(err)
		}

		gid, err := newAria2.DownloadByWithVideo(video, v.Magnet)
		if err != nil {
			log.Error(err)
			continue
		}

		// 如果开启了tg推送 则推送
		if config.TG.Enable {
			go func() {
				video, err = model.NewMovieDB().FetchOneDouBanVideoByDouBanID(v.DoubanID)
				if err != nil {
					log.Error(err)
				}
				//nolint:exhaustruct
				bus.DownloadNotifyChan <- &types.DownloadNotifyVideo{
					Video: video,
					File:  v.TorrentName,
					Gid:   gid,
				}
			}()
		}

		log.Infof(" 开始下载: %s. videoType: %s.  GID: %s.", v.TorrentName, v.Type, gid)
	}
	return nil
}

func (d *Download) Run() {
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

//nolint:gochecknoglobals
var wg sync.WaitGroup

//nolint:nakedret
func (d *Download) DownloadByName(name, resolution string) (msg string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		feedKnaben := searchspider.NewFeedKnaben(name, d.ResolutionStr2Int(resolution))
		_, err := feedKnaben.Search()
		if err != nil {
			log.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		feedBt4g := searchspider.NewFeedBt4g(name, d.ResolutionStr2Int(resolution))
		_, err := feedBt4g.Search()
		if err != nil {
			log.Error(err)
		}
	}()
	wg.Wait()

	// 获取 磁力连接
	videos, err := model.NewMovieDB().GetFeedVideoMovieByName([]string{name}...)
	if err != nil {
		log.Error(err)
	}

	if len(videos) == 0 {
		return "所有资源已下载过,或没有可下载资源."
	}

	// 推送 磁力连接至 aria2
	newAria2, err := aria2.NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		log.Error(err)
	}

	for _, v := range videos {
		video, err := model.NewMovieDB().FetchOneDouBanVideoByDouBanID(v.DoubanID)
		if err != nil {
			log.Error(err)
		}

		gid, err := newAria2.DownloadByWithVideo(video, v.Magnet)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Infof("Downloader: %s 开始下载. GID: %s", v.Name, gid)
		err = model.NewMovieDB().UpdateFeedVideoDownloadByID(v.ID, 1)
		if err != nil {
			log.Error(err)
		}
	}

	return fmt.Sprintf("已将 %d 资源加入下载.", len(videos))
}
