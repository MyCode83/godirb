package wordlist

import (
	_ "embed"
	"strings"
	"sync"
)

//go:embed wordlists-txt/lfi.txt
var lfiRaw string
var lfiSlice []string
var lfiOnce sync.Once

func Lfi() []string {
	lfiOnce.Do(func() {
		lfiSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(lfiRaw, "\r\n", "\n")), "\n")
	})
	return lfiSlice
}
