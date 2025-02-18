package feedspider

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"movieSpider/internal/httpclient"
	"movieSpider/internal/log"
	"movieSpider/internal/magnetconvert"
	"movieSpider/internal/types"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	urlTorlockMovie = "https://www.torlock.com/movies/rss.xml"
	urlTorlockTV    = "https://www.torlock.com/television/rss.xml"
)

type torlock struct {
	typ types.VideoType
	// url        string
	web        string
	scheduling string
	// httpClient *http.Client
}

//nolint:gosimple,gocognit,gocritic
func (t *torlock) Crawler() (Videos []*types.FeedVideo, err error) {
	fp := gofeed.NewParser()
	fp.Client = httpclient.NewHTTPClient()
	fp.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36"
	if t.typ == types.VideoTypeMovie {
		fd, err := fp.ParseURL(urlTorlockMovie)
		if err != nil {
			log.Error(err)
		}
		if fd == nil {
			//nolint:revive,goerr113
			return nil, fmt.Errorf("TORLOCK.%s feed is nil.", t.typ)
		}
		log.Debugf("TORLOCK.movie Data: %#v", fd.String())
		if len(fd.Items) == 0 {
			//nolint:revive
			return nil, errors.New("TORLOCK.movie: 没有feed数据.")
		}

		var videos1 []*types.FeedVideo
		nameReg := regexp.MustCompile("(.*)\\.([0-9][0-9][0-9][0-9])\\.")
		yearReg := regexp.MustCompile("(.*)\\.\\(([0-9][0-9][0-9][0-9])\\)\\.")
		for _, v := range fd.Items {
			// 片名
			name := strings.ReplaceAll(v.Title, " ", ".")

			var fVideo types.FeedVideo
			fVideo.Web = t.web
			fVideo.TorrentName = name
			fVideo.TorrentURL = v.Link
			fVideo.Type = "movie"

			// 原始数据
			//nolint:errchkjson
			bytes, _ := json.Marshal(v)
			//nolint:exhaustruct
			fVideo.RowData = sql.NullString{String: string(bytes)}

			// 片名
			matchArr := nameReg.FindStringSubmatch(name)
			if len(matchArr) > 0 {
				fVideo.Name = fVideo.FormatName(matchArr[1])
			} else {
				fVideo.Name = fVideo.FormatName(name)
			}
			// 年份
			submatch := yearReg.FindStringSubmatch(name)

			if len(matchArr) > 0 {
				if matchArr[2] != "" {
					fVideo.Year = matchArr[2]
				}
			} else {
				if len(submatch) > 2 {
					fVideo.Year = submatch[2]
				}
			}
			videos1 = append(videos1, &fVideo)
		}

		videos2 := t.fetchMagnetDownLoad(videos1)

		Videos = t.fetchMagnet(videos2)
		return Videos, nil
	}
	if t.typ == types.VideoTypeTV {
		fd, err := fp.ParseURL(urlTorlockTV)
		if err != nil {
			log.Error(err)
		}
		if fd == nil {
			return nil, errors.New("TORLOCK.tv: 没有feed数据")
		}
		log.Debugf("TORLOCK.tv Data: %#v", fd.String())
		if len(fd.Items) == 0 {
			return nil, errors.New("TORLOCK.tv: 没有feed数据")
		}
		compileRegex := regexp.MustCompile("(.*)\\.[sS][0-9][0-9]|[Ee][0-9][0-9]?\\.")

		var videos1 []*types.FeedVideo

		for _, v := range fd.Items {
			// 片名
			name := strings.ReplaceAll(v.Title, " ", ".")

			matchArr := compileRegex.FindStringSubmatch(name)
			var fVideo types.FeedVideo
			fVideo.TorrentName = fVideo.FormatName(name)
			fVideo.TorrentURL = v.Link
			fVideo.Type = "tv"
			// 原始数据
			//nolint:errchkjson
			bytes, _ := json.Marshal(v)
			//nolint:exhaustruct
			fVideo.RowData = sql.NullString{String: string(bytes)}
			fVideo.Web = t.web
			// 片名
			if len(matchArr) > 0 {
				fVideo.Name = fVideo.FormatName(matchArr[1])
			} else {
				fVideo.Name = fVideo.FormatName(name)
			}
			videos1 = append(videos1, &fVideo)
		}

		videos2 := t.fetchMagnetDownLoad(videos1)

		Videos = t.fetchMagnet(videos2)
		return Videos, nil
	}
	//nolint:nakedret
	return
}

func (t *torlock) fetchMagnet(videos []*types.FeedVideo) (feedVideos []*types.FeedVideo) {
	var wg sync.WaitGroup
	for _, video := range videos {
		wg.Add(1)
		magnet, err := magnetconvert.FetchMagnet(video.Magnet)
		if err != nil {
			log.Errorf("TORLOCK: get %s magnet download url is %s", video.Name, video.Magnet)
			wg.Done()
			continue
		}
		video.Magnet = magnet
		feedVideos = append(feedVideos, video)
		wg.Done()
	}
	wg.Wait()
	return feedVideos
}

func (t *torlock) fetchMagnetDownLoad(videos []*types.FeedVideo) []*types.FeedVideo {
	var wg sync.WaitGroup
	var videos2 []*types.FeedVideo
	for _, video := range videos {
		wg.Add(1)
		//nolint:noctx
		resp, err := httpclient.NewHTTPClient().Get(video.TorrentURL)
		if err != nil {
			log.Errorf("TORLOCK.%s %s http request url is %s, error:%s", video.Type, video.Name, video.TorrentURL, err)
			wg.Done()
			continue
		}
		defer resp.Body.Close()
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Errorf("TORLOCK.%s %s http goquery error:%s", video.Type, video.Name, err)
			continue
		}
		val, exists := doc.Find("body > article > div:nth-child(6) > div > div:nth-child(2) > a").Attr("href")
		if exists {
			video.Magnet = val
			videos2 = append(videos2, video)
		}
		wg.Done()
	}
	wg.Wait()
	return videos2
}
func (t *torlock) Run(ch chan *types.FeedVideo) {
	if t.scheduling == "" {
		log.Errorf("TORLOCK %s: Scheduling is null", t.typ)
		os.Exit(1)
	}
	log.Infof("TORLOCK %s: Scheduling is: [%s]", t.typ, t.scheduling)
	c := cron.New()
	_, _ = c.AddFunc(t.scheduling, func() {
		videos, err := t.Crawler()
		if err != nil {
			log.Error(err)
		}
		// model.ProxySaveVideo2DB(videos...)
		for _, video := range videos {
			ch <- video
		}
	})
	c.Start()
}
