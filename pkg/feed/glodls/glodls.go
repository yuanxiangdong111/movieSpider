package glodls

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
	"movieSpider/pkg"
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
		return nil, pkg.ErrGLODLSFeedNull
	}
	log.Debugf("GLODLS Config: %#v", fd)
	log.Debugf("GLODLS Data: %#v", fd.String())
	if len(fd.Items) == 0 {
		return nil, pkg.ErrGLODLSFeedNull
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
		parse, err := url.Parse(v.Link)
		if err != nil {
			log.Error(err)
		}
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
				pkg.CheckError("GLODLS", err)
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
		return "", pkg.ErrGLODLSFeedMagnetFetch
	}
	client := httpClient2.GetHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		return "", pkg.ErrGLODLSFeedMagnetFetch
	}
	if resp == nil {
		return "", pkg.ErrGLODLSFeedMagnetFetch
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", pkg.ErrGLODLSFeedMagnetFetch
	}
	selector := "#downloadbox > table > tbody > tr > td:nth-child(1) > a:nth-child(2)"
	magnet, _ = doc.Find(selector).Attr("href")

	return
}
func (g *glodls) Run() {

	if g.scheduling == "" {
		log.Error("GLODLS Scheduling is null")
		os.Exit(1)
	}
	log.Infof("GLODLS Scheduling is: [%s]", g.scheduling)
	c := cron.New()
	_, err := c.AddFunc(g.scheduling, func() {
		videos, err := g.Crawler()
		pkg.CheckError("GLODLS", err)
		for _, v := range videos {
			go func(video *types.FeedVideo) {
				err = model.MovieDB.CreatFeedVideo(video)
				if err != nil {
					pkg.CheckError("GLODLS", err)
				} else {
					log.Infof("GLODLS: %s 保存完毕", video.Name)
				}
			}(v)
		}
	})
	if err != nil {
		log.Error("GLODLS: AddFunc is null")
		os.Exit(1)
	}
	c.Start()

}
