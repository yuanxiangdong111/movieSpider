package utils

import "strings"

// ExcludeVideo 排除  480p 720p  dvsux  hdr 视频源
func ExcludeVideo(name string) bool {
	lowerTorrentName := strings.ToLower(name)
	if strings.Contains(lowerTorrentName, "720p") || strings.Contains(lowerTorrentName, "dvsux") || strings.Contains(lowerTorrentName, "480p") || strings.Contains(lowerTorrentName, "hdr") {
		return true
	} else {
		return false
	}
}
