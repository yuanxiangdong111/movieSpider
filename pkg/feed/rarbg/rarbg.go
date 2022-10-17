package rarbg

import (
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
	"movieSpider/pkg"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"movieSpider/pkg/utils"
	"os"
	"regexp"
	"strings"
)

type rarbg struct {
	typ        types.Resource
	url        string
	web        string
	scheduling string
}

func NewFeedRarbg(url, scheduling string, resourceType types.Resource) *rarbg {

	if resourceType == types.ResourceMovie {
		return &rarbg{
			resourceType,
			url,
			"rarbg",
			scheduling,
		}
	} else {
		return &rarbg{
			resourceType,
			url,
			"rarbg",
			scheduling,
		}
	}

}
func (r *rarbg) Crawler() (Videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()

	if r.typ == types.ResourceMovie {
		fd, _ := fp.ParseURL(r.url)
		if fd == nil {
			return nil, pkg.ErrRARBFeedNull
		}
		log.Debugf("RARBG Config: %#v", fd)
		log.Debugf("RARBG Data: %#v", fd.String())
		compileRegex := regexp.MustCompile("(.*)\\.([0-9][0-9][0-9][0-9])\\.")
		for _, v := range fd.Items {
			// 片名
			name := strings.ReplaceAll(v.Title, " ", ".")
			ok := utils.ExcludeVideo(name)
			if ok {
				continue
			}

			var fVideo types.FeedVideo
			fVideo.Web = r.web
			fVideo.TorrentName = name
			fVideo.Magnet = v.Link
			fVideo.Type = "movie"

			// 原始数据
			bytes, _ := json.Marshal(v)
			fVideo.RowData = string(bytes)

			// 片名
			matchArr := compileRegex.FindStringSubmatch(name)
			if len(matchArr) > 0 {
				fVideo.Name = matchArr[1]
			} else {
				fVideo.Name = name
			}
			// 年份
			if len(matchArr) > 0 {
				fVideo.Year = matchArr[2]
			}
			Videos = append(Videos, &fVideo)
		}
	}
	if r.typ == types.ResourceTV {
		fd, _ := fp.ParseURL(r.url)
		if fd == nil {
			return nil, pkg.ErrRARBFeedNull
		}
		log.Debugf("RARBG Config: %#v", fd)
		log.Debugf("RARBG Data: %#v", fd.String())
		compileRegex := regexp.MustCompile("(.*)\\.[sS][0-9][0-9]|[Ee][0-9][0-9]?\\.")
		for _, v := range fd.Items {
			// 片名
			name := strings.ReplaceAll(v.Title, " ", ".")
			ok := utils.ExcludeVideo(name)
			if ok {
				continue
			}

			matchArr := compileRegex.FindStringSubmatch(name)

			var fVideo types.FeedVideo
			fVideo.TorrentName = name
			fVideo.Magnet = v.Link
			fVideo.Type = "tv"
			// 原始数据
			bytes, _ := json.Marshal(v)
			fVideo.RowData = string(bytes)

			// 片名
			if len(matchArr) > 0 {
				fVideo.Name = matchArr[1]
			} else {
				fVideo.Name = name
			}

			Videos = append(Videos, &fVideo)
		}
	}

	return
}
func (r *rarbg) Run() {
	if r.scheduling == "" {
		log.Error("RARBG Scheduling is null")
		os.Exit(1)
	}
	log.Infof("RARBG Scheduling is: [%s]", r.scheduling)
	c := cron.New()
	_, err := c.AddFunc(r.scheduling, func() {
		videos, err := r.Crawler()
		pkg.CheckError("RARBG", err)
		for _, v := range videos {
			go func(video *types.FeedVideo) {
				err = model.MovieDB.CreatFeedVideo(video)
				if err != nil {
					pkg.CheckError("RARBG", err)
				} else {
					log.Infof("RARBG: %s 保存完毕", video.Name)
				}
			}(v)
		}
	})
	if err != nil {
		log.Error("RARBG: AddFunc is null")
		os.Exit(1)
	}
	c.Start()

}
