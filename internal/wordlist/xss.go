package wordlist

import (
	_"embed"
	"sync"
	"strings"
)

//go:embed wordlists-txt/xss.txt
var xssRaw string
var xssSlice []string
var xssOnce sync.Once
func Xss() []string{
	xssOnce.Do(func() {
		xssSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(xssRaw, "\r\n", "\n")), "\n")
	})
	return xssSlice
}