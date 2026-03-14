package wordlist

import (
	_"embed"
	"sync"
	"strings"
)

//go:embed wordlists-txt/medium.txt
var mediumRaw string
var mediumSlice []string
var mediumOnce sync.Once
func Medium() []string{
	mediumOnce.Do(func() {
		mediumSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(mediumRaw, "\r\n", "\n")), "\n")
	})
	return mediumSlice
}