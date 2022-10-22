package aria2

import (
	"fmt"
	"movieSpider/pkg/config"
	"testing"
)

func Test_aria2_DownloadList(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")

	newAria2, err := NewAria2(config.Downloader.Aria2Label)
	if err != nil {
		t.Error(err)
	}
	//info, err := newAria2.aria2Client.GetGlobalStat()
	//if err != nil {
	//	t.Error(err)
	//}
	//sessionInfo, err := newAria2.aria2Client.TellActive()
	//if err != nil {
	//	t.Error(err)
	//}
	//
	//for _, v := range sessionInfo {
	//	fmt.Println(v.TotalLength)
	//	fmt.Println(v.Gid)
	//	//fmt.Printf("%#v\n", v)
	//	infos, err := newAria2.aria2Client.GetServers(v.Gid)
	//
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	fmt.Println(infos)
	//	//for _, f := range v.Files {
	//	//	fmt.Printf("%#v\n", f)
	//	//
	//	//}
	//}
	sessionInfo, err := newAria2.aria2Client.TellStopped(0, 1)
	if err != nil {
		t.Error(err)
	}

	for _, v := range sessionInfo {
		fmt.Println(v.TotalLength)
		fmt.Println(v.Gid)
		fmt.Println(v.Files)
		//fmt.Printf("%#v\n", v)

		//for _, f := range v.Files {
		//	fmt.Printf("%#v\n", f)
		//
		//}
	}

}
