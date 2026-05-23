package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/MyCode83/godirb/internal/confirmation"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/duration"

	"github.com/fatih/color"
)

var (
	timeoutErr error
	delayErr   error
)

func ValidateFlags(cfg *Config) {
	debug.Printf("validate flags start")
	cfg.BaseURL = strings.TrimRight(cfg.URL, "/")
	if cfg.NoColor || cfg.Quiet {
		useColors = false
		debug.Printf("colors disabled no_color=%t quiet=%t", cfg.NoColor, cfg.Quiet)
	}
	if cfg.JSON && cfg.CSV {
		debug.Printf("invalid output flags: json and csv both enabled")
		fmt.Fprintln(os.Stderr, "[X] Error, use only one output format: --json or --csv")
		os.Exit(1)
	}

	cfg.Timeout, timeoutErr = duration.ParseDuration(cfg.RawTimeout, "s")
	cfg.Delay, delayErr = duration.ParseDuration(cfg.RawDelay, "ms")
	debug.Printf("parsed durations timeout=%s delay=%s", cfg.Timeout, cfg.Delay)
	// Delay error
	if delayErr != nil {
		debug.Error("delay parse", delayErr)
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'\n", cfg.RawDelay)
		os.Exit(1)
	}
	// Timeout error
	if timeoutErr != nil {
		debug.Error("timeout parse", timeoutErr)
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'\n", cfg.RawTimeout)
		os.Exit(1)
	}
	if strings.TrimSpace(cfg.URL) == "" {
		debug.Printf("missing url flag")
		fmt.Fprintln(os.Stderr, "[X] Error, missing '--url'(-u) flag")
		fmt.Fprintln(os.Stderr, "Run godirb --help for usage")
		os.Exit(2)
	}
	if cfg.Threads >= 250 && !cfg.ForceThreads {
		debug.Printf("high thread confirmation required threads=%d", cfg.Threads)

		if !confirmation.ThreadsConfirmation(fmt.Sprintf("Are you shure you want to send %d cfg.Threads", cfg.Threads)) {
			debug.Printf("high thread confirmation rejected")
			fmt.Println("Cancelled by the user")
			os.Exit(0)
		}
		debug.Printf("high thread confirmation accepted")

	}
	if cfg.Threads <= 0 {
		debug.Printf("invalid threads=%d", cfg.Threads)
		fmt.Fprintf(os.Stderr, "[X] Error, you can't send %d threads.\n", cfg.Threads)
		os.Exit(2)
	}
	if !useColors {
		color.NoColor = true
	}
	if cfg.Timeout <= 0 {
		debug.Printf("invalid timeout=%s", cfg.Timeout)
		fmt.Fprintln(os.Stderr, "[X] Error: timeout must be greater than 0")
		os.Exit(2)
	}
	debug.Printf("validate flags done")
}
