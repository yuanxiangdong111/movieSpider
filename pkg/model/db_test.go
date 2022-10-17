package model

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"movieSpider/pkg/config"
	"testing"
)

func Test_movieDB_FetchMovieDouBanVideo(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	NewMovieDB()

	video, err := MovieDB.FetchDouBanMovies()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(video)
}

func Test_movieDB_FetchMagnetByName(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	NewMovieDB()

	var names = []string{"Moon"}
	videos, err := MovieDB.FetchMagnetByName(names)
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
