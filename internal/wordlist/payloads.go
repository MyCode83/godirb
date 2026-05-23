package wordlist

import (
	_ "embed"
	"strings"
	"sync"
)

//go:embed wordlists-txt/payloads.txt
var payloadsRaw string
var payloadsSlice []string
var payloadsOnce sync.Once

func Payloads() []string {
	payloadsOnce.Do(func() {
		payloadsSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(payloadsRaw, "\r\n", "\n")), "\n")
	})
	return payloadsSlice
}
