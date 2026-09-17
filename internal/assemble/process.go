package assemble

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/wordlist"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func ConfigureProcess() {
	log.SetOutput(os.Stderr)
	_, ok := os.LookupEnv("GODIRB_NO_COLOR")
	if ok {
		lipgloss.SetColorProfile(termenv.Ascii)
	}
}

func SetupSignals() (context.Context, context.CancelFunc) {
	contextCancel, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()

		go func() {
			time.Sleep(1 * time.Second)
			os.Exit(1)
		}()
	}()
	go func() {
		<-contextCancel.Done()

		// log.Println(": context canceled")
	}()

	return contextCancel, cancel
}

func LogParsedFlags(cfg cli.Config, wd wordlist.Wordlist) {
	debug.Set(cfg.Debug)
	debug.Printf("parsed flags url=%q wordlist=%q threads=%d depth=%d timeout=%q delay=%q method=%q recursive=%t quiet=%t json=%t csv=%t output=%q",
		cfg.URL, wd.Wordlist, cfg.Threads, cfg.Depth, cfg.RawTimeout, cfg.RawDelay, cfg.Method, cfg.Recursive, cfg.Quiet, cfg.JSON, cfg.CSV, cfg.Output)
}

func LogValidatedFlags(cfg cli.Config) {
	debug.Printf("validated flags base_url=%q timeout=%s delay=%s ignore=%v exts=%v headers=%d proxy=%q insecure=%t",
		cfg.BaseURL, cfg.Timeout, cfg.Delay, cfg.IgnoreCode, cfg.Exts, len(cfg.Header), cfg.Proxy, cfg.Insecure)
}

func LogSelectedMode(mode core.Mode) {
	debug.Printf("selected mode=%d", mode)
}
