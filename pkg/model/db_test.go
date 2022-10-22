package model

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"movieSpider/pkg/config"
	"movieSpider/pkg/types"
	"strings"
	"testing"
)

func Test_movieDB_FetchMagnetByName(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	NewMovieDB()

	var names = []string{"Andor"}
	videos, err := MovieDB.FetchMovieMagnetByName(names)
	if err != nil {
		t.Error(err)
	}
	for _, v := range videos {
		fmt.Println(v.Magnet)
	}

}

func Test_movieDB_UpdateFeedVideoDownloadByID(t *testing.T) {
	//config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	//NewMovieDB()
	//
	//err := MovieDB.UpdateFeedVideoDownloadByID(55)
	//if err != nil {
	//	t.Error(err)
	//}
	isURL := govalidator.IsURL("https://api.wmdb.tv/movie/api?id=")
	//parse, err := url.ParseRequestURI("http://dsadsa.com")
	//if err != nil {
	//	t.Error(err)
	//}
	fmt.Println(isURL)
}

func Test_movieDB_RandomOneDouBanVideo(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	NewMovieDB()
	MovieDB.RandomOneDouBanVideo()
}

func Test_movieDB_FetchTVMagnetByName(t *testing.T) {

	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	NewMovieDB()

	var names = []string{"Andor"}
	videos, err := MovieDB.FetchTVMagnetByName(names)
	if err != nil {
		t.Error(err)
	}
	is1, is3 := sotByResolution(videos)
	fmt.Println(is1)
	fmt.Println(is3)

}
func sotByResolution(videos []*types.FeedVideo) (downloadIs1 []*types.FeedVideo, downloadIs3 []*types.FeedVideo) {
	var Videos2160P []*types.FeedVideo
	var Videos1080P []*types.FeedVideo
	for _, v := range videos {
		switch {
		case strings.Contains(v.TorrentName, "2160"):
			Videos2160P = append(Videos2160P, v)
		case strings.Contains(v.TorrentName, "1080"):
			Videos1080P = append(Videos1080P, v)
		}
	}
	if len(Videos2160P) >= 0 {
		if len(Videos2160P) >= 2 {
			downloadIs1 = append(downloadIs1, Videos2160P[0:2]...)
			downloadIs3 = append(downloadIs3, Videos2160P[2:]...)
			downloadIs3 = append(downloadIs3, Videos1080P...)
		} else {
			downloadIs1 = append(downloadIs1, Videos2160P...)
			downloadIs3 = append(downloadIs3, Videos1080P...)
		}

	} else {
		if len(Videos2160P) >= 2 {
			downloadIs1 = append(downloadIs1, Videos1080P[0:2]...)
			downloadIs3 = append(downloadIs3, Videos1080P[2:]...)
		} else {
			downloadIs1 = append(downloadIs1, Videos1080P...)
		}

	}
	return
}
