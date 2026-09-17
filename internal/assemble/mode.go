package assemble

import (
	"fmt"
	"os"
	"time"

	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/wordlist"
	"github.com/spf13/pflag"
)

func ApplyModeDefaults(cfg *cli.Config, wd *wordlist.Wordlist, mode core.Mode) {
	switch mode {
	case core.ModeFuzz:
		if !pflag.Lookup("ext").Changed {
			cfg.Exts = []string{}
			debug.Printf("fuzz mode without explicit --ext; extensions disabled")
		}
	case core.ModePort:
		if !pflag.Lookup("wordlist").Changed {
			wd.Wordlist = "ports"
			debug.Printf("port mode without explicit wordlist; using ports wordlist")
		}
		if !pflag.Lookup("timeout").Changed {
			cfg.Timeout = time.Duration(500) * time.Millisecond
			debug.Printf("port mode without explicit timeout; using %s", cfg.Timeout)
		}
		switch {
		case cfg.Timeout > time.Second:
			fmt.Fprintf(os.Stderr, "[!] High timeout (%s). Scan may be slow.\n", cfg.Timeout)
		case cfg.Timeout >= time.Duration(5)*time.Second:
			fmt.Fprintf(os.Stderr, "[!] Very high timeout (%s). Scan will be very slow.\nCTRL + C will take a while (up to 30s).\n", cfg.Timeout)
		}
	case core.ModeDir:
	}
}
