package aria2

import (
	"context"
	"errors"
	"fmt"
	"github.com/zyxar/argo/rpc"
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
	return nil, errors.New("aria2 is nil")
}

func (a *aria2) DownloadByUrl(url string) (gid string, err error) {
	return a.aria2Client.AddURI([]string{url})
}
func (a *aria2) DownloadList(url string) (gid string, err error) {
	info, err := a.aria2Client.GetSessionInfo()
	if err != nil {
		return "", err
	}
	fmt.Println(info)
	return
}
