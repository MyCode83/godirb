package assemble

import (
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/wordlist"
)

func BuildWordlist(wd wordlist.Wordlist) []string {
	wl := wd.LoadWordlist()
	debug.Printf("loaded wordlist entries=%d source=%q", len(wl), wd.Wordlist)

	return wl
}
