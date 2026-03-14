package wordlist


import (
	_"embed"
	"sync"
	"strings"
)

//go:embed wordlists-txt/big.txt
var bigRaw string
var bigSlice []string
var bigOnce sync.Once
func Big() []string{
	bigOnce.Do(func() {
		bigSlice = strings.Split(strings.TrimSpace(strings.ReplaceAll(bigRaw, "\r\n", "\n")), "\n")
	})
	return bigSlice
}