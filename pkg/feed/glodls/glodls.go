package glodls

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	httpClient2 "movieSpider/pkg/httpClient"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"movieSpider/pkg/utils"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
)

type glodls struct {
	url        string
	scheduling string
	web        string
}

func NewFeedGlodls(url, scheduling string) *glodls {
	return &glodls{
		url,
		scheduling,
		"glodls",
	}
}
func (g *glodls) Crawler() (videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fd, err := fp.ParseURL(g.url)
	if fd == nil {
		return nil, errors.New("GLODLS: 没有feed数据")
	}
	log.Debugf("GLODLS Config: %#v", fd)
	log.Debugf("GLODLS Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, errors.New("GLODLS: 没有feed数据")
	}
	var videosA []*types.FeedVideo
	for _, v := range fd.Items {
		// 片名
		torrentName := strings.ReplaceAll(v.Title, " ", ".")
		ok := utils.ExcludeVideo(torrentName)
		if ok {
			continue
		}

		// 片名处理
		var name string
		var year string

		if strings.ToLower(v.Categories[0]) == "tv" {
			compileRegex := regexp.MustCompile("(.*)(\\.[Ss][0-9][0-9][eE][0-9][0-9])")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			// 如果 正则匹配过后 没有结果直接 过滤掉
			if len(matchArr) == 0 {
				continue
			}
			name = matchArr[1]
		} else if strings.ToLower(v.Categories[0]) == "movies" {
			compileRegex := regexp.MustCompile("(.*)\\.(\\d{4})\\.")
			matchArr := compileRegex.FindStringSubmatch(torrentName)
			if len(matchArr) == 0 {
				name = torrentName
			} else {
				name = matchArr[1]
				year = matchArr[2]
			}

		} else {
			name = torrentName
		}

		fVideo := new(types.FeedVideo)
		fVideo.Name = fVideo.FormatName(name)
		fVideo.Year = year

		fVideo.Web = g.web
		parse, _ := url.Parse(v.Link)
		// 种子名
		fVideo.TorrentName = torrentName

		if len(parse.Query()["id"]) == 0 {
			log.Error("没有ID")
		}
		id := parse.Query()["id"][0]
		all := strings.ReplaceAll(v.Title, " ", "-")

		TorrentUrl := fmt.Sprintf("https://glodls.to/%s-f-%s.html", strings.ToLower(all), id)

		fVideo.TorrentUrl = TorrentUrl

		// 处理 资源类型 是 电影 还是电视剧
		typ := strings.ToLower(v.Categories[0])
		if typ == "movies" {
			fVideo.Type = "movie"
		} else {
			fVideo.Type = typ
		}

		bytes, _ := json.Marshal(v)
		fVideo.RowData = string(bytes)

		videosA = append(videosA, fVideo)

	}
	var wg sync.WaitGroup

	for _, v := range videosA {
		wg.Add(1)
		go func(video *types.FeedVideo) {
			magnet, err := g.fetchMagnet(video.TorrentUrl)
			if err != nil {
				log.Error(err)
			}
			if magnet == "" {
				return
			}
			video.Magnet = magnet
			videos = append(videos, video)
			wg.Done()
		}(v)
	}
	wg.Wait()

	return
}

func (g *glodls) fetchMagnet(url string) (magnet string, err error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.WithMessage(err, "GLODLS: 磁链获取错误")
	}
	client := httpClient2.NewHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		return "", errors.WithMessage(err, "GLODLS: 磁链获取错误")
	}
	if resp == nil {
		return "", errors.New("GLODLS: response is nil")
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", errors.WithMessage(err, "GLODLS: 磁链获取错误")
	}
	selector := "#downloadbox > table > tbody > tr > td:nth-child(1) > a:nth-child(2)"
	magnet, exists := doc.Find(selector).Attr("href")
	if !exists {
		return "", errors.WithMessage(err, "GLODLS: 磁链获取错误")
	}
	return magnet, nil
}
func (g *glodls) Run() {

	if g.scheduling == "" {
		log.Error("GLODLS Scheduling is null")
		os.Exit(1)
	}
	log.Infof("GLODLS Scheduling is: [%s]", g.scheduling)
	c := cron.New()
	c.AddFunc(g.scheduling, func() {
		videos, err := g.Crawler()
		if err != nil {
			log.Error(err)
			return
		}

		for _, v := range videos {
			go func(video *types.FeedVideo) {
				err = model.MovieDB.CreatFeedVideo(video)
				if err != nil {
					if errors.Is(err, model.ErrorDataExist) {
						log.Debug(err)
						return
					}
					log.Error(err)
					return
				}
				log.Infof("GLODLS: %s 保存完毕", video.Name)
			}(v)
		}
	})
	c.Start()

}
