package aria2

import (
	"context"
	"errors"
	"fmt"
	"github.com/zyxar/argo/rpc"
	"movieSpider/pkg/config"
	"movieSpider/pkg/log"
	"movieSpider/pkg/types"
	"path"
	"strconv"
	"strings"
)

var (
	Aria2 *aria2
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
			Aria2 = &aria2{client}
			return Aria2, nil
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

func (a *aria2) CompletedFiles() (completedFiles []*types.ReportCompletedFiles) {
	sessionInfo, err := a.aria2Client.TellStopped(0, 100)
	if err != nil {
		log.Error(err)
		return nil
	}

	for _, v := range sessionInfo {
		if len(v.Files) > 0 {
			if strings.Contains(v.Files[0].Path, "[METADATA]") {
				continue
			} else {
				// 下载了多少
				CompletedLength, err := strconv.Atoi(v.Files[0].CompletedLength)
				if err != nil {
					log.Error(err)
				}
				// 文件大小
				Length, err := strconv.Atoi(v.Files[0].Length)
				if err != nil {
					log.Error(err)
				}
				//文件完成度百分比
				completed := CompletedLength / Length * 100
				//if completed != 100 {
				//	continue
				//}
				f := new(types.ReportCompletedFiles)
				f.GID = v.Gid
				f.Completed = fmt.Sprintf("%d%%", completed)
				f.Size = fmt.Sprintf("%.2fGB", float32(Length)/1024/1024/1024)
				_, file := path.Split(v.Files[0].Path)
				f.FileName = file
				completedFiles = append(completedFiles, f)
			}
		}

	}
	return
}
