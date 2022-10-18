package wmdb

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
	"io"
	"movieSpider/pkg"
	"movieSpider/pkg/httpClient"
	"movieSpider/pkg/log"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"net/http"
	"os"
	"strings"
)

type wmdb struct {
	url        string
	scheduling string
}

func NewSpiderWmdb(url, scheduling string) *wmdb {
	return &wmdb{
		url,
		scheduling,
	}
}

// crawlerImdb 30s 内只允许一个请求
func (d *wmdb) crawler(doubanID string) (video *types.DouBanVideo, err error) {
	video = new(types.DouBanVideo)

	request, err := http.NewRequest(http.MethodGet, d.url+doubanID, nil)
	if err != nil {
		return nil, nil
	}
	client := httpClient.GetHttpClient()
	resp, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if resp == nil {
		log.Warn("未能正常获取wmdb数据")
		return nil, errors.New("未能正常获取wmdb数据")
	}

	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Debugf("WMDB Config: %#v", string(all))
	rowData := string(all)

	if strings.Contains(rowData, "Too Many Requests") {
		return nil, pkg.ErrWMDBSpiderNull
	}
	if strings.Contains(rowData, "your requests today are full") {
		return nil, pkg.ErrWMDBSpiderNull
	}

	video.ImdbID = gjson.Get(rowData, "imdbId").String()
	video.Type = strings.ToLower(gjson.Get(rowData, "type").String())
	video.RowData = rowData

	var ns []string
	array := gjson.Get(rowData, "data").Array()
	for _, v := range array {
		name := gjson.Get(v.String(), "name").String()
		replace := strings.ReplaceAll(name, " ", ".")
		ns = append(ns, replace)
	}
	marshal, _ := json.Marshal(ns)
	video.Names = string(marshal)

	return
}

func (d *wmdb) Run() {
	if d.scheduling == "" {
		log.Error("WMDB: Scheduling is null")
		os.Exit(1)
	}
	log.Infof("WMDB Scheduling is: [%s]", d.scheduling)
	c := cron.New()
	_, err := c.AddFunc(d.scheduling, func() {
		video, err := model.MovieDB.RandomOneDouBanVideo()
		if err != nil {
			log.Warn("WMDB: 没有可爬取的豆瓣数据")
			return
		}

		v, err := d.crawler(video.DoubanID)
		if err != nil {
			log.Error("WMDB", err)
			return
		}
		video.ImdbID = v.ImdbID
		video.RowData = v.RowData
		if v.Type == "TVSeries" {
			video.Type = "tv"
		} else {
			video.Type = v.Type
		}

		video.Names = v.Names
		err = model.MovieDB.UpDateDouBanVideo(video)
		if err != nil {
			log.Error(err)
			return
		}
		log.Warnf("WMDB: %s 更新完毕", video.Names)
	})
	if err != nil {
		log.Error("WMDB: AddFunc is null")
		os.Exit(1)
	}
	c.Start()
}
