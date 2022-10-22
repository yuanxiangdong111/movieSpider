package knaben

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"movieSpider/pkg/utils"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

type knaben struct {
	url        string
	resolution types.Resolution
	web        string
}

func NewFeedKnaben(knabenUrl, name string, resolution types.Resolution) *knaben {
	parse, err := url.Parse(knabenUrl)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	strData := url.QueryEscape(name)

	kUrl := fmt.Sprintf("%s://%s/%s", parse.Scheme, parse.Host, strData)

	return &knaben{url: kUrl, resolution: resolution, web: "knaben"}
}

func (k *knaben) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(k.url)
	if fd == nil {
		return nil, errors.New("KNABEN: 没有feed数据")
	}
	log.Debugf("KNABEN: Config %#v", fd)
	log.Debugf("KNABEN: Data %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("KNABEN: 没有feed数据")
	}
	for _, v := range fd.Items {
		// 片名
		name := strings.ReplaceAll(v.Title, " ", ".")
		ok := utils.ExcludeVideo(name)
		if ok {
			continue
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = k.web

		if len(v.Categories) > 0 {
			ty := strings.ToLower(v.Categories[0])
			if ty == "movies" {
				fVideo.Type = "movie"
			} else {
				fVideo.Type = ty
			}
		} else if len(v.Categories) == 0 {
			fVideo.Type = "other"
		}

		for _, m := range strings.Split(v.Description, "\n") {
			if strings.Contains(m, "Magnet") {
				fVideo.Name = fVideo.FormatName(name)
				magnet := k.parseMagnet(m)
				fVideo.Magnet = magnet
				// 种子名
				fVideo.TorrentName = fVideo.Name
				videos = append(videos, fVideo)

			}
		}
	}

	// 异步保存至 数据库
	var wg sync.WaitGroup
	for _, v := range videos {
		wg.Add(1)
		// 异步保存至 数据库
		go func(video *types.FeedVideo) {
			err := model.MovieDB.CreatFeedVideo(video)
			if err != nil {
				if errors.Is(err, model.ErrorDataExist) {
					log.Warn(err)
					return
				}
				log.Error(err)
				return
			}
			log.Infof("KNABEN: %s", video.TorrentName)
		}(v)
		wg.Done()
	}
	wg.Wait()

	// 指定清晰度
	if k.resolution.Res() != "" {
		var resolutionVideos []*types.FeedVideo
		for _, v := range videos {
			if strings.Contains(v.Name, k.resolution.Res()) {
				resolutionVideos = append(resolutionVideos, v)
			}
		}
		return resolutionVideos, nil
	}

	return
}

func (k *knaben) parseMagnet(str string) string {
	compileRegex := regexp.MustCompile(".*(magnet.*)\">Magnet")
	matchArr := compileRegex.FindStringSubmatch(str)
	if len(matchArr) >= 2 {
		return matchArr[1]
	}

	return ""
}
