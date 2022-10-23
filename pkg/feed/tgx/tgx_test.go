package tgx

import (
	"movieSpider/pkg/config"
	"testing"
)

func Test_tgx_Run(t1 *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	t := tgx{
		scheduling: "tt.fields.scheduling",
		url:        urlStr,
		web:        "tgx",
	}
	t.Run()
}
