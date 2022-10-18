package pkg

import (
	"github.com/pkg/errors"
	"movieSpider/pkg/log"
)

var (
	ErrDBExist               = errors.New("数据已存在")
	ErrDBNotFound            = errors.New("数据未找到")
	ErrBTBTMoviePageBody     = errors.New("BTBT电影页没有数据")
	ErrBTBTMagnet            = errors.New("BTBT获取磁链错误")
	ErrBTBTMagnetMeta        = errors.New("BTBT读取磁链meta信息错误")
	ErrBTBTMagnetMarshal     = errors.New("BTBT磁链解析错误")
	ErrGLODLSFeedNull        = errors.New("GLODLS没有feed数据")
	ErrGLODLSFeedMagnetFetch = errors.New("GLODLS磁链获取错误")
	ErrRARBFeedNull          = errors.New("RARBG没有feed数据")
	ErrRARBMovieFeedNull     = errors.New("RARBG没有movie feed数据")
	ErrRARBTVFeedNull        = errors.New("RARBG没有tv feed数据")
	ErrBT4GFeedNull          = errors.New("BT4G没有feed数据")
	ErrEZTVFeedNull          = errors.New("EZTV没有feed数据")
	ErrKNABENFeedNull        = errors.New("KNABEN没有feed数据")
	ErrWMDBSpiderNull        = errors.New("WMDB没有Spider数据")
	ErrAria2IsNil            = errors.New("aria2 is nil")
	ErrDownloaderAria2       = errors.New("aria2下载错误")
)

func CheckError(who string, err error) {
	switch errors.Unwrap(err) {
	case ErrDBNotFound:
		log.Warnf("%s: %s", who, err)
	case ErrDBExist:
		log.Warnf("%s: %s", who, err)
	case ErrBTBTMoviePageBody:
		log.Warnf("%s: %s", who, err)
	case ErrBTBTMagnet:
		log.Warnf("%s: %s", who, err)
	case ErrBTBTMagnetMeta:
		log.Warnf("%s: %s", who, err)
	case ErrBTBTMagnetMarshal:
		log.Warnf("%s: %s", who, err)
	case ErrGLODLSFeedNull:
		log.Warnf("%s: %s", who, err)
	case ErrRARBFeedNull:
		log.Warnf("%s: %s", who, err)
	case ErrRARBMovieFeedNull:
		log.Warnf("%s: %s", who, err)
	case ErrRARBTVFeedNull:
		log.Warnf("%s: %s", who, err)
	case ErrBT4GFeedNull:
		log.Warnf("%s: %s", who, err)
	case ErrEZTVFeedNull:
		log.Warnf("%s: %s", who, err)
	case ErrKNABENFeedNull:
		log.Warnf("%s: %s", who, err)
	case ErrGLODLSFeedMagnetFetch:
		log.Warnf("%s: %s", who, err)
	case ErrWMDBSpiderNull:
		log.Warnf("%s: %s", who, err)
	case ErrAria2IsNil:
		log.Warnf("%s: %s", who, err)
	case ErrDownloaderAria2:
		log.Warnf("%s: %s", who, err)
	case nil:

	default:
		log.Errorf("%s: %s", who, err)
	}
}
