package convert

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"io"
	"movieSpider/pkg"
)

func FileToMagnet(file string) (string, error) {
	mi, err := metainfo.LoadFromFile(file)
	if err != nil {
		return "", fmt.Errorf("Cannot read the metainfo from file: %s. %v", file, err)
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return "", fmt.Errorf("Cannot unmarshal the metainfo from file: %s. %v", file, err)
	}
	hs := mi.HashInfoBytes()

	if info.Name == "" {
		return "", nil
	}

	return mi.Magnet(&hs, &info).String(), nil
}

func IO2Magnet(r io.Reader) (string, error) {
	mi, err := metainfo.Load(r)
	if err != nil {
		return "", pkg.ErrBTBTMagnetMeta
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return "", pkg.ErrBTBTMagnetMarshal
	}
	hs := mi.HashInfoBytes()

	if info.Name == "" {
		return "", nil
	}

	return mi.Magnet(&hs, &info).String(), nil
}
