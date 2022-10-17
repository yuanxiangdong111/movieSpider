package httpClient

import (
	"movieSpider/pkg/config"
	"net/http"
	"net/url"
	"sync"
)

var (
	httpClient *http.Client
	once       = &sync.Once{}
)

func GetHttpClient() *http.Client {
	if httpClient == nil {
		once.Do(func() {
			if config.Global.Proxy != "" {
				proxyUrl, _ := url.Parse(config.Global.Proxy)
				proxy := http.ProxyURL(proxyUrl)
				transport := &http.Transport{Proxy: proxy}
				httpClient = &http.Client{Transport: transport}
			} else {
				httpClient = &http.Client{}
			}
		})
	}

	return httpClient
}
