package wordlist

import (
	_"embed"
	"sync"
	"strings"
)

//go:embed wordlists-txt/ports.txt
var portsRaw string
var portsSlice []string
var portsOnce sync.Once
func Ports() []string{
	portsOnce.Do(func() {
		portsSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(portsRaw, "\r\n", "\n")), "\n")
	})
	return portsSlice
}