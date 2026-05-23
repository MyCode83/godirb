package main

import (
	// stdlib
	"context"
	"fmt"
	"log"
	"runtime"

	"time"

	"godirb/pkg/random"
	"os"
	"os/signal"

	"sync"
	"syscall"

	// Third-libs
	"github.com/spf13/pflag"
	"github.com/valyala/fasthttp"

	// Godirb-lib
	"godirb/internal/assemble"
	"godirb/internal/cli"
	"godirb/internal/confirmation"

	"godirb/internal/core" // core
	"godirb/internal/output"
	"godirb/internal/validate"

	"godirb/internal/baseline"
	"godirb/internal/wildcard"

	"godirb/internal/tui"

	"github.com/fatih/color"
)

const banner string = (`		                     
   ____ _  ____   ____/ /   (_)   _____   / /_
  / __  / / __ \ / __  /   / /   / ___/  / __ \
 / /_/ / / /_/ // /_/ /   / /   / /     / /_/ /
 \__  /  \____/ \____/   /_/   /_/     /_____/
/____/
`)

var (
	client       *fasthttp.Client
	wg           sync.WaitGroup
	tasksWG      sync.WaitGroup
	visitedMutex sync.Mutex
	mode         core.Mode = core.ModeDir
)

const version = "0.9.0"

var preUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
}

// others
var (
	auth          string
	contextCancel context.Context
	cancel        context.CancelFunc
)

func main() {
	log.SetOutput(os.Stderr)
	_, ok := os.LookupEnv("GODIRB_NO_COLOR")
	if ok {
		color.NoColor = true
	}

	contextCancel, cancel = context.WithCancel(context.Background())
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
		log.Println(runtime.NumGoroutine())
		// log.Println(": context canceled")
	}()
	cfg, wd := cli.ParseFlags()
	cli.ValidateFlags(&cfg)
	mode = cli.SelectMode(mode, cfg)

	// wd = instance
	// wl = wordlist slice
	client = assemble.BuildProxyAndClient(cfg.Proxy, cfg.Timeout, cfg.Insecure) // Fasthttp-Client
	switch mode {
	case core.ModeFuzz:
		if !pflag.Lookup("placeholder").Changed {
			cfg.Exts = []string{}
		}
	case core.ModePort:
		if !pflag.Lookup("wordlist").Changed {
			wd.Wordlist = "ports"
		}
		if !pflag.Lookup("timeout").Changed {
			cfg.Timeout = time.Duration(500) * time.Millisecond
		}
		log.Printf(": %s\n", wd.Wordlist)
		switch {
		case cfg.Timeout > time.Second:
			fmt.Fprintf(os.Stderr, "[!] High timeout (%s). Scan may be slow.\n", cfg.Timeout)
		case cfg.Timeout >= time.Duration(5)*time.Second:
			fmt.Fprintf(os.Stderr, "[!] Very high timeout (%s). Scan will be very slow.\nCTRL + C will take a while (up to 30s).\n", cfg.Timeout)
		}
	case core.ModeDir:
		if !validate.ValidateUrl(cfg.BaseURL, client, cfg.Method, random.RandChoice(cfg.UserAgent)) {
			os.Exit(1)
		}
	}

	wl := wd.LoadWordlist() // Load Wordlist

	// Basic-Auth
	if cfg.Password != "" && cfg.Username != "" {
		auth = assemble.BuildBasicAuth(cfg.Username, cfg.Password)
	}

	outputFormat := output.FromFlags(cfg.JSON, cfg.CSV)
	collectOutput := cfg.Output != "" || outputFormat != output.FormatText

	if !cfg.Quiet && !(collectOutput && cfg.Output == "") {
		fmt.Println("\n------------------")
		fmt.Println("[*] Url: ", cfg.BaseURL)
		fmt.Println("[*] Method: ", cfg.Method)
		fmt.Println("[*] Threads: ", cfg.Threads)
		fmt.Println("[*] Timeout: ", cfg.Timeout)
		fmt.Println("[*] Delay: ", cfg.Delay)
		fmt.Println("[*] UAs: ", len(cfg.UserAgent))
		fmt.Print("[*] Mode: ")
		switch mode {
		case core.ModeDir:
			fmt.Print("Dir\n")
		case core.ModeFuzz:
			fmt.Print("Fuzz\n")
		case core.ModePort:
			fmt.Print("Port\n")
		}
		fmt.Printf("------------------\n\n")
	}

	limiter := make(chan struct{}, cfg.Threads)
	var dirsChan chan string
	if mode == core.ModeDir {
		dirsChan = make(chan string, cfg.Threads*50)
	}

	engine := &core.Core{
		// Mode
		Mode: mode,

		// Bools

		Recursive: cfg.Recursive,

		// Context
		Ctx:    contextCancel,
		Cancel: cancel,
		// Config
		Timeout: cfg.Timeout,
		Delay:   cfg.Delay,
		Quiet:   cfg.Quiet,

		// HTTP
		Client:     client,
		Method:     cfg.Method,
		UserAgents: cfg.UserAgent,
		AuthHeader: auth,
		Header:     cfg.Header,
		// cfg.Placeholder
		Placeholder: cfg.Placeholder,
		// Control
		IgnoreCodes: cfg.IgnoreCode,
		Exts:        cfg.Exts,

		// Concurrency
		Limiter:  limiter,
		DirsChan: dirsChan,

		// WG
		WG: &wg,

		// WordList
		WL: wl,

		// State
		VisitedDirs: make(map[string]bool),

		// Output / colors
		Others: tui.Other,
		File:   tui.File,
	}
	// Wildcard
	switch mode {
	case core.ModeDir:
		wildcard, err := wildcard.DetectWildcard(client, cfg.BaseURL, cfg.Placeholder, cfg.UserAgent...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(2)
		}
		if wildcard.Active {
			fmt.Fprintf(os.Stderr, "[!] Wildcard detected: %d | %d bytes \n", wildcard.Status, wildcard.Lenght)

			if cfg.Method != "GET" && !cfg.ForceHead {
				fmt.Fprintf(os.Stderr, "[!] Wildcard-like behavior detected using HEAD/SWITCH requests.\n")
				fmt.Fprintf(os.Stderr, "You can skip this confirmation with puting '--force-head'\n")
				fmt.Fprintf(os.Stderr, "HEAD/SWITCH responses do not include a body, so wildcard filtering\ncannot be done reliably and may produce false positives.\n")
				fmt.Fprintf(os.Stderr, "\nSwitch cfg.Method to 'GET'? [y/N]: \n")

				if confirmation.WildcardConfirmation() {
					cfg.Method = "GET"
				}
			}
		}
		engine.Wildcard = wildcard
	case core.ModeFuzz:
		baseline, err := baseline.BuildBaseLine(cfg.BaseURL, client, cfg.Placeholder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		engine.Baseline = baseline

	}
	results := make([]core.Result, 0)
	for result := range engine.Run(cfg.BaseURL) {
		if collectOutput {
			results = append(results, result)
			continue
		}
		tui.Print(result, cfg.Quiet)
	}
	if collectOutput {
		if err := output.Write(results, outputFormat, cfg.Output, cfg.Quiet); err != nil {
			fmt.Fprintf(os.Stderr, "[X] Error writing output: %v\n", err)
			os.Exit(1)
		}
	}

}
