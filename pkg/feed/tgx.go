package feed

import (
	"encoding/json"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/pkg/convert"
	httpClient2 "movieSpider/pkg/httpClient"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"movieSpider/pkg/utils"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const urlTgx = "https://tgx.rs/rss"

type tgx struct {
	scheduling string
	url        string
	web        string
}

func (t *tgx) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(t.url)
	if fd == nil {
		return nil, errors.New("TGx: 没有feed数据")
	}
	log.Debugf("TGx Config: %#v", fd)
	log.Debugf("TGx Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("TGx: 没有feed数据")
	}

	var videos1 []*types.FeedVideo
	for _, v := range fd.Items {

		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		ok := utils.ExcludeVideo(torrentName)
		if ok {
			continue
		}

		var name string
		var year string
		//
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

		var typ string
		// 过滤掉 其他类型的种子
		if strings.HasPrefix(strings.ToLower(v.Categories[0]), "tv :") {
			typ = "tv"
		} else if strings.HasPrefix(strings.ToLower(v.Categories[0]), "movies :") {
			typ = "movie"
		} else {
			continue
		}

		fVideo := new(types.FeedVideo)
		fVideo.Web = t.web
		fVideo.Year = year

		// 片名
		fVideo.Name = fVideo.FormatName(name)
		// 种子名
		fVideo.TorrentName = fVideo.FormatName(torrentName)
		fVideo.Type = typ
		fVideo.TorrentUrl = v.Link
		bytes, _ := json.Marshal(v)
		fVideo.RowData = string(bytes)
		videos1 = append(videos1, fVideo)

	}
	var wg sync.WaitGroup
	for _, video := range videos1 {
		wg.Add(1)
		magnet, err := t.fetchMagnet(video.TorrentUrl)
		if err != nil {
			wg.Done()
			return nil, err
		}
		video.Magnet = magnet
		videos = append(videos, video)
		wg.Done()
	}
	wg.Wait()
	return
}

func (t *tgx) fetchMagnet(url string) (magnet string, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.WithMessage(err, "TGx: 磁链获取错误")
	}
	client := httpClient2.NewHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		return "", errors.WithMessage(err, "TGx: 磁链获取错误")
	}
	if resp == nil {
		return "", errors.New("TGx: response is nil")
	}
	defer resp.Body.Close()

	magnet, err = convert.IO2Magnet(resp.Body)
	if err != nil {
		return "", errors.New("TGx: 磁链转换错误")
	}

	return magnet, nil
}
func (t *tgx) Run() {

	if t.scheduling == "" {
		log.Error("TGx Scheduling is null")
		os.Exit(1)
	}
	log.Infof("TGx Scheduling is: [%s]", t.scheduling)
	c := cron.New()
	c.AddFunc(t.scheduling, func() {
		videos, err := t.Crawler()
		if err != nil {
			log.Error(err)
			return
		}

		for _, v := range videos {
			go func(video *types.FeedVideo) {
				err = model.NewMovieDB().CreatFeedVideo(video)
				if err != nil {
					if errors.Is(err, model.ErrorDataExist) {
						log.Warn(err)
						return
					}
					log.Error(err)
					return
				}
				log.Infof("TGx: %s 保存完毕", video.Name)
			}(v)
		}
	})
	c.Start()
}
