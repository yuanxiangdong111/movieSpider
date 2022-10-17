package aria2

import (
	"context"
	"github.com/zyxar/argo/rpc"
	"movieSpider/pkg"
	"movieSpider/pkg/config"
	"movieSpider/pkg/log"
)

type aria2 struct {
	aria2Client rpc.Client
}

func NewAria2(label string) (*aria2, error) {
	for _, v := range config.Aria2cList {
		if v.Label == label {
			client, err := rpc.New(context.TODO(), v.Url, v.Token, 0, nil)
			if err != nil {
				return nil, err
			}
			log.Debug(config.Aria2cList)
			return &aria2{client}, nil
		}
	}
	return nil, pkg.ErrAria2IsNil
}

func (a *aria2) DownloadByUrl(url string) (gid string, err error) {
	return a.aria2Client.AddURI([]string{url})
}
