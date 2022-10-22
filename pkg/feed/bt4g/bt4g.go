package bt4g

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"movieSpider/pkg/utils"
	"net/url"
	"os"
	"strings"
	"sync"
)

type bt4g struct {
	url string
	//typ        types.Resource
	resolution types.Resolution
	web        string
}

func NewFeedBt4g(urlStr, name string, resolution types.Resolution) *bt4g {
	parse, err := url.Parse(urlStr)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	//strData := url.QueryEscape(name)
	bUrl := fmt.Sprintf("%s://%s/search/%s?page=rss", parse.Scheme, parse.Host, name)
	return &bt4g{url: bUrl, resolution: resolution, web: "bt4g"}
}

func (b *bt4g) Crawler() (videos []*types.FeedVideo, err error) {
	f := gofeed.NewParser()
	fd, err := f.ParseURL(b.url)
	if fd == nil {
		return nil, errors.New("BT4G: 没有feed数据")
	}
	log.Debugf("BT4G Config: %#v", b)
	log.Debugf("BT4G Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("BT4G: 没有feed数据")
	}
	for _, v := range fd.Items {
		// 片名
		name := strings.ReplaceAll(v.Title, " ", ".")
		ok := utils.ExcludeVideo(name)
		if ok {
			continue
		}
		if v.Link == "" {
			continue
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = b.web
		fVideo.Name = name
		fVideo.Magnet = v.Link
		// 种子名
		fVideo.TorrentName = fVideo.Name

		fVideo.TorrentUrl = v.GUID
		fVideo.Type = "other"
		bytes, _ := json.Marshal(v)
		fVideo.RowData = string(bytes)
		videos = append(videos, fVideo)
	}

	var wg sync.WaitGroup
	for _, v := range videos {
		wg.Add(1)
		// 异步保存至 数据库
		go func(video *types.FeedVideo) {
			err := model.MovieDB.CreatFeedVideo(video)
			if err != nil {
				if errors.Is(err, model.ErrorDataExist) {
					log.Debug(err)
					return
				}
				log.Error(err)
				return
			}
			log.Infof("BT4G: %s", video.TorrentName)
		}(v)
		wg.Done()
	}
	wg.Wait()
	// 指定清晰度
	if b.resolution != types.ResolutionOther {
		var resolutionVideos []*types.FeedVideo
		for _, v := range videos {
			if strings.Contains(v.Name, b.resolution.Res()) {
				resolutionVideos = append(resolutionVideos, v)
			}
		}
		return resolutionVideos, nil
	}
	return
}
