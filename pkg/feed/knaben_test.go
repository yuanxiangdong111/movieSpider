package feed

import (
	"encoding/json"
	"fmt"
	"movieSpider/pkg/config"
	"movieSpider/pkg/model"
	"movieSpider/pkg/types"
	"testing"
)

func TestEztv_Crawler(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	model.NewMovieDB()

	f := NewFeedKnaben("https://rss.knaben.eu/", "House Of The Dragon", types.ResolutionOther)
	videos, err := f.Crawler()
	if err != nil {
		t.Error(err)
	}
	for _, v := range videos {
		bytes, _ := json.Marshal(v)

		fmt.Println(string(bytes))
		//fmt.Println(v.Name)

	}
}
