package ipProxy

import (
	"fmt"
	"movieSpider/pkg/config"
	"testing"
)

func TestFetchProxy(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	proxy := FetchProxy()
	fmt.Println(proxy)
}
