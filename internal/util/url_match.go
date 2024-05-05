package util

import "strings"

const (
	MatchAll = "/*"
)

func MatchUrl(pattern string, url string) bool {
	if len(url) == 0 || len(pattern) == 0 {
		return false
	}

	if pattern == url {
		return true
	}

	if pattern == MatchAll {
		return true
	}

	if pattern[len(pattern)-2:] == MatchAll {
		tmpUrl := url
		for true {
			if pattern == tmpUrl+"/*" {
				return true
			}
			index := strings.LastIndex(tmpUrl, "/")
			if index < 0 {
				break
			}

			tmpUrl = tmpUrl[:index]
		}
	}
	return false
}
