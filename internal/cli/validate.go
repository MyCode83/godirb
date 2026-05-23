package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/MyCode83/godirb/internal/confirmation"
	"github.com/MyCode83/godirb/internal/duration"

	"github.com/fatih/color"
)

var (
	timeoutErr error
	delayErr   error
)

func ValidateFlags(cfg *Config) {
	cfg.BaseURL = strings.TrimRight(cfg.URL, "/")
	if cfg.NoColor || cfg.Quiet {
		useColors = false
	}
	if cfg.JSON && cfg.CSV {
		fmt.Fprintln(os.Stderr, "[X] Error, use only one output format: --json or --csv")
		os.Exit(1)
	}

	cfg.Timeout, timeoutErr = duration.ParseDuration(cfg.RawTimeout, "s")
	cfg.Delay, delayErr = duration.ParseDuration(cfg.RawDelay, "ms")
	// Delay error
	if delayErr != nil {
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'\n", cfg.RawDelay)
		os.Exit(1)
	}
	// Timeout error
	if timeoutErr != nil {
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'\n", cfg.RawTimeout)
		os.Exit(1)
	}
	if strings.TrimSpace(cfg.URL) == "" {
		fmt.Fprintln(os.Stderr, "[X] Error, missing '--url'(-u) flag")
		fmt.Fprintln(os.Stderr, "Run godirb --help for usage")
		os.Exit(2)
	}
	if cfg.Threads >= 250 && !cfg.ForceThreads {

		if !confirmation.ThreadsConfirmation(fmt.Sprintf("Are you shure you want to send %d cfg.Threads", cfg.Threads)) {
			fmt.Println("Cancelled by the user")
			os.Exit(0)
		}

	}
	if cfg.Threads <= 0 {
		fmt.Fprintf(os.Stderr, "[X] Error, you can't send %d threads.\n", cfg.Threads)
		os.Exit(2)
	}
	if !useColors {
		color.NoColor = true
	}
	if cfg.Timeout <= 0 {
		fmt.Fprintln(os.Stderr, "[X] Error: timeout must be greater than 0")
		os.Exit(2)
	}
}
