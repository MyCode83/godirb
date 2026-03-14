package wordlist

import (
	_"embed"
	"sync"
	"strings"
)

//go:embed wordlists-txt/common.txt
var commonRaw string
var commonSlice []string
var commonOnce sync.Once
func Common() []string{
	commonOnce.Do(func() {
		commonSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(commonRaw, "\r\n", "\n")), "\n")
	})
	return commonSlice
}