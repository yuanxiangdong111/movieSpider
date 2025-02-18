package feedspider

import (
	"movieSpider/internal/bus"
	"movieSpider/internal/config"
	"testing"
)

func Test_tgx_Run(t1 *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/core/bin/core/config.yaml")
	t := tgx{
		scheduling: "tt.fields.scheduling",
		url:        urlTgx,
		web:        "tgx",
	}
	t.Run(bus.FeedVideoChan)
}
