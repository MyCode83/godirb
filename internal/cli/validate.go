package cli
import (
	"os"
	"fmt"
	"strings"

	"godirb/internal/duration"
	"godirb/internal/confirmation"
	"godirb/internal/core"
	"godirb/pkg/parse"

	"github.com/fatih/color"
)

var (
	timeoutErr error
	delayErr error
)

func ValidateFlags(cfg *Config) {
	cfg.BaseURL = strings.TrimRight(cfg.URL, "/")
	if cfg.NoColor || cfg.Quiet{
			useColors = false
	}

	cfg.Timeout, timeoutErr = duration.ParseDuration(cfg.RawTimeout, "s")
	cfg.Delay, delayErr = duration.ParseDuration(cfg.RawDelay, "ms")
	// Delay error
	if delayErr != nil {
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'", cfg.RawDelay)
		os.Exit(1)		
	}
	// Timeout error
	if timeoutErr != nil {
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'", cfg.RawTimeout)
		os.Exit(1)
	}
	if strings.TrimSpace(cfg.URL) == "" {
		fmt.Fprintln(os.Stderr, "[X] Error, missing '--cfg.URL'(-u) flag")
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
		fmt.Fprintf(os.Stderr, "[X] Error, you can't send %d cfg.Threads.\n", cfg.Threads)
		os.Exit(2)
	}
	if !useColors {
		color.NoColor = true
	}
	if cfg.Timeout <= 0 {
		fmt.Printf("[X] Error: timeout must be greater than 0\n")
		os.Exit(2)
	}
	if parse.ExtractPort(cfg.BaseURL) == cfg.Placeholder{
		mode = core.ModePort
	} else if strings.Contains(cfg.BaseURL, cfg.Placeholder) {
		mode = core.ModeFuzz
	}
}