package rarbg

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/pkg/ipProxy"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"movieSpider/pkg/utils"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type rarbg struct {
	typ        types.Resource
	url        string
	web        string
	scheduling string
	httpClient *http.Client
}

func NewFeedRarbg(url, scheduling string, resourceType types.Resource) *rarbg {

	if resourceType == types.ResourceMovie {
		return &rarbg{
			resourceType,
			url,
			"rarbg",
			scheduling,
			&http.Client{},
		}
	} else {
		return &rarbg{
			resourceType,
			url,
			"rarbg",
			scheduling,
			&http.Client{},
		}
	}

}
func (r *rarbg) Crawler() (Videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fp.Client = r.httpClient
	fp.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	if r.typ == types.ResourceMovie {
		fd, err := fp.ParseURL(r.url)
		if err == nil {
			log.Error(err)
		}
		if fd == nil {
			return nil, errors.New(fmt.Sprintf("RARBG: %s feed is nil.", r.typ.Typ()))
		}
		log.Debugf("RARBG movie Data: %#v", fd.String())
		if len(fd.Items) == 0 {
			return nil, errors.New(fmt.Sprintf("RARBG: 没有movie feed数据."))
		}
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
			return nil, errors.New("RARBG: 没有tv feed数据")
		}
		log.Debugf("RARBG tv Data: %#v", fd.String())
		if len(fd.Items) == 0 {
			return nil, errors.New("RARBG: 没有tv feed数据")
		}
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
			fVideo.Web = r.web
			// 片名
			if len(matchArr) > 0 {
				fVideo.Name = matchArr[1]
			} else {
				fVideo.Name = name
			}
			Videos = append(Videos, &fVideo)
		}
	}

	//request, err := http.NewRequest(http.MethodGet, r.url, nil)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//request.Header.Set("Cookie", "aby=2; gaDts48g=q8h5pp9t; gaDts48g=q8h5pp9t; skt=fp5kiozvk8; skt=fp5kiozvk8")
	//request.Header.Set("sec-ch-ua", `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`)
	//request.Header.Set("sec-ch-ua-mobile", "?0")
	//request.Header.Set("sec-ch-ua-platform", "Linux")
	//request.Header.Set("Sec-Fetch-Dest", "document")
	//request.Header.Set("Sec-Fetch-Mode", "navigate")
	//request.Header.Set("Sec-Fetch-Site", "none")
	//request.Header.Set("Sec-Fetch-User", "?1")
	//request.Header.Set("Upgrade-Insecure-Requests", "1")
	//do, err := r.httpClient.Do(request)
	//if err != nil {
	//	return nil, err
	//}
	//
	//defer func() {
	//	do.Body.Close()
	//}()
	//
	//all, err := io.ReadAll(do.Body)
	//if err != nil {
	//	return nil, err
	//}
	//log.Debugf("RARBG %s Data: %#v", r.typ.Typ(), string(all))
	//doc := etree.NewDocument()
	//err = doc.ReadFromString(string(all))
	//if err != nil {
	//	return nil, err
	//}
	//
	//elements := doc.SelectElement("rss")
	//element := elements.SelectElement("channel")
	//
	//if r.typ == types.ResourceMovie {
	//	compileRegex := regexp.MustCompile("(.*)\\.([0-9][0-9][0-9][0-9])\\.")
	//
	//	for _, e := range element.SelectElements("item") {
	//		// 种子名
	//		title := e.SelectElement("title").Text()
	//		name := strings.ReplaceAll(title, " ", ".")
	//		ok := utils.ExcludeVideo(name)
	//		if ok {
	//			continue
	//		}
	//		var fVideo types.FeedVideo
	//		fVideo.Web = r.web
	//		fVideo.TorrentName = name
	//		// 磁力连接
	//		fVideo.Magnet = e.SelectElement("link").Text()
	//		fVideo.Type = "movie"
	//
	//		// 片名
	//		matchArr := compileRegex.FindStringSubmatch(name)
	//		if len(matchArr) > 0 {
	//			fVideo.Name = matchArr[1]
	//		} else {
	//			fVideo.Name = name
	//		}
	//		// 年份
	//		if len(matchArr) > 0 {
	//			fVideo.Year = matchArr[2]
	//		}
	//		Videos = append(Videos, &fVideo)
	//	}
	//}
	//
	//if r.typ == types.ResourceTV {
	//	compileRegex := regexp.MustCompile("(.*)\\.([0-9][0-9][0-9][0-9])\\.")
	//	for _, e := range element.SelectElements("item") {
	//		// 种子名
	//		title := e.SelectElement("title").Text()
	//		name := strings.ReplaceAll(title, " ", ".")
	//		ok := utils.ExcludeVideo(name)
	//		if ok {
	//			continue
	//		}
	//		var fVideo types.FeedVideo
	//		fVideo.Web = r.web
	//		fVideo.TorrentName = name
	//		fVideo.Magnet = e.SelectElement("link").Text()
	//		fVideo.Type = "movie"
	//		// 片名
	//		matchArr := compileRegex.FindStringSubmatch(name)
	//		if len(matchArr) > 0 {
	//			fVideo.Name = matchArr[1]
	//		} else {
	//			fVideo.Name = name
	//		}
	//		// 年份
	//		if len(matchArr) > 0 {
	//			fVideo.Year = matchArr[2]
	//		}
	//		Videos = append(Videos, &fVideo)
	//	}
	//
	//}

	return
}
func (r *rarbg) Run() {
	if r.scheduling == "" {
		log.Errorf("RARBG %s: Scheduling is null", r.typ.Typ())
		os.Exit(1)
	}
	log.Infof("RARBG %s: Scheduling is: [%s]", r.typ.Typ(), r.scheduling)
	c := cron.New()
	_, err := c.AddFunc(r.scheduling, func() {
		videos, err := r.Crawler()
		if err != nil {
			//for{
			r.switchClient()
			videos, err = r.Crawler()
			if err != nil {
				log.Error(err)
				return
			}
			r.save2DB(videos)
			//}
		}
		r.save2DB(videos)
	})
	if err != nil {
		log.Error("RARBG: AddFunc is null")
		os.Exit(1)
	}
	c.Start()

}

func (r *rarbg) useProxyClient() {
	proxyStr := ipProxy.FetchProxy()
	if proxyStr == "" {
		log.Error("useProxyClient: proxy is null")
		return
	}

	proxyUrl, err := url.Parse(proxyStr)
	if err != nil {
		log.Error(err)
		return
	}
	if proxyUrl != nil {
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		log.Errorf("useProxyClient: use proxy %#v", proxyUrl.String())
		httpClient := &http.Client{Transport: transport, Timeout: time.Minute * 5}

		r.httpClient = httpClient
	}
	return
}
func (r *rarbg) switchClient() {
	if r.httpClient.Transport == nil {

		proxyStr := ipProxy.FetchProxy()
		if proxyStr == "" {
			log.Infof("RARBG.%s: proxy is null.", r.typ.Typ())
			return
		}
		proxyUrl, err := url.Parse(proxyStr)
		if err != nil {
			log.Error(err)
		}
		if proxyUrl != nil {
			proxy := http.ProxyURL(proxyUrl)
			transport := &http.Transport{Proxy: proxy}
			httpClient := &http.Client{Transport: transport, Timeout: time.Minute * 5}
			r.httpClient = httpClient
			log.Infof("RARBG.%s: 添加代理. proxyStr: %s", r.typ.Typ(), proxyUrl)
		} else {
			log.Warnf("RARBG.%s: 请添加Global.Proxy.Url配置", r.typ.Typ())
		}

	} else {
		r.httpClient = &http.Client{}
		log.Infof("RARBG.%s: 删除代理.", r.typ.Typ())
	}
}

func (r *rarbg) save2DB(videos []*types.FeedVideo) {
	if videos == nil || len(videos) == 0 {
		log.Warnf("RARBG: %s: 没有数据", r.typ.Typ())
		return
	}
	for _, v := range videos {
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
			log.Infof("RARBG: %s: %s 保存完毕", r.typ.Typ(), video.Name)

		}(v)
	}
}
