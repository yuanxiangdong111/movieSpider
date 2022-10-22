package eztv

import (
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"movieSpider/pkg/utils"
	"os"
	"regexp"
	"strings"
)

type eztv struct {
	scheduling string
	url        string
	web        string
}

func NewFeedEztv(url, scheduling string) *eztv {
	return &eztv{
		scheduling,
		url,
		"eztv",
	}
}
func (f *eztv) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(f.url)
	if fd == nil {
		return nil, errors.New("EZTV: 没有feed数据")
	}
	log.Debugf("EZTV Config: %#v", fd)
	log.Debugf("EZTV Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("EZTV: 没有feed数据")
	}
	for _, v := range fd.Items {
		torrentName := strings.ReplaceAll(v.Title, " ", ".")

		ok := utils.ExcludeVideo(torrentName)
		if ok {
			continue
		}

		var name string
		var year string
		compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
		matchArr := compileRegex.FindStringSubmatch(torrentName)
		if len(matchArr) == 0 {
			tvReg := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			TVNameArr := tvReg.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(TVNameArr) == 0 {
				continue
			}
			name = TVNameArr[1]

		} else {
			year = matchArr[2]
			name = matchArr[1]
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = f.web
		fVideo.Year = year

		// 片名
		fVideo.Name = name
		// 种子名
		fVideo.TorrentName = torrentName
		fVideo.TorrentUrl = v.Link
		fVideo.Magnet = v.Extensions["torrent"]["magnetURI"][0].Value
		bytes, _ := json.Marshal(v)
		fVideo.Type = strings.ToLower(v.Categories[0])
		fVideo.RowData = string(bytes)

		videos = append(videos, fVideo)

	}
	return
}
func (f *eztv) Run() {
	if f.scheduling == "" {
		log.Error("EZTV Scheduling is null")
		os.Exit(1)
	}
	log.Infof("EZTV Scheduling is: [%s]", f.scheduling)
	c := cron.New()
	_, err := c.AddFunc(f.scheduling, func() {
		videos, err := f.Crawler()
		if err != nil {
			log.Error(err)
			return
		}
		for _, v := range videos {
			go func(video *types.FeedVideo) {
				err = model.MovieDB.CreatFeedVideo(video)
				if err != nil {
					if errors.Is(err, model.ErrorDataExist) {
						log.Warn(err)
						return
					}
					log.Error(err)
					return
				}
				log.Infof("EZTV: %s 保存完毕", video.Name)
			}(v)
		}

	})
	if err != nil {
		log.Error("EZTV: AddFunc is null")
		os.Exit(1)
	}
	c.Start()
}
