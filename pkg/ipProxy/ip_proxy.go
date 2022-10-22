package ipProxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"movieSpider/pkg/config"
	"movieSpider/pkg/log"
	"net/http"
)

type proxyData struct {
	Anonymous  string `json:"anonymous"`
	CheckCount int    `json:"check_count"`
	FailCount  int    `json:"fail_count"`
	Https      bool   `json:"https"`
	LastStatus bool   `json:"last_status"`
	LastTime   string `json:"last_time"`
	Proxy      string `json:"proxy"`
	Region     string `json:"region"`
	Source     string `json:"source"`
}

func FetchProxy() string {
	if config.ProxyPool != "" {
		resp, err := http.Get(fmt.Sprintf("%s/get", config.ProxyPool))
		if err != nil {
			return ""
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		defer resp.Body.Close()

		var data proxyData
		err = json.Unmarshal(buf.Bytes(), &data)
		if err != nil {
			log.Warnf("FetchProxy: %s.", err.Error())
			return ""
		}
		return fmt.Sprintf("http://%s", data.Proxy)

	}
	log.Warn("FetchProxy: Global.ProxyPool没有配置.")
	return ""
}
