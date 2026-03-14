package wordlist

import (
	_"embed"
	"sync"
	"strings"
)

//go:embed wordlists-txt/small.txt
var smallRaw string
var smallSlice []string
var smallOnce sync.Once
func Small() []string{
	smallOnce.Do(func() {
		smallSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(smallRaw, "\r\n", "\n")), "\n")
	})
	return smallSlice
}