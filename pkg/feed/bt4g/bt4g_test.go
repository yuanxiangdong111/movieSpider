package bt4g

import (
	"fmt"
	"movieSpider/pkg/config"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"testing"
)

func Test_bt4g_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	model.NewMovieDB()
	//https: //bt4g.org/search/%E9%BE%99%E4%B9%8B%E5%AE%B6%E6%97%8F?page=rss
	//https://bt4g.org/search/%E9%BE%99%E4%B9%8B%E5%AE%B6%E6%97%8F?page=rss
	b := NewFeedBt4g(config.Bt4G.Url, "杀手疾风号", types.ResolutionOther)

	gotVideos, err := b.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, v := range gotVideos {
		fmt.Println(v.Magnet)
	}
}
