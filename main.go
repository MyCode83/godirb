package main

import (
	// stdlib
	"context"
	"fmt"
	"log"

	"time"

	"github.com/MyCode83/godirb/pkg/random"
	"os"
	"os/signal"

	"sync"
	"syscall"

	// Third-libs
	"github.com/spf13/pflag"

	// Godirb-lib
	"github.com/MyCode83/godirb/internal/assemble"
	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/confirmation"

	"github.com/MyCode83/godirb/internal/core" // core
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/output"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/validate"

	"github.com/MyCode83/godirb/internal/calibration"

	"github.com/MyCode83/godirb/internal/tui"

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
	wg           sync.WaitGroup
	tasksWG      sync.WaitGroup
	visitedMutex sync.Mutex
	mode         core.Mode = core.ModeDir
)

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

		// log.Println(": context canceled")
	}()
	cfg, wd := cli.ParseFlags()
	if cfg.Version {
		fmt.Println(cli.Version)
		return
	}
	debug.Set(cfg.Debug)
	debug.Printf("parsed flags url=%q wordlist=%q threads=%d timeout=%q delay=%q method=%q recursive=%t quiet=%t json=%t csv=%t output=%q",
		cfg.URL, wd.Wordlist, cfg.Threads, cfg.RawTimeout, cfg.RawDelay, cfg.Method, cfg.Recursive, cfg.Quiet, cfg.JSON, cfg.CSV, cfg.Output)
	cli.ValidateFlags(&cfg)
	debug.Printf("validated flags base_url=%q timeout=%s delay=%s ignore=%v exts=%v headers=%d proxy=%q insecure=%t",
		cfg.BaseURL, cfg.Timeout, cfg.Delay, cfg.IgnoreCode, cfg.Exts, len(cfg.Header), cfg.Proxy, cfg.Insecure)
	mode = cli.SelectMode(mode, cfg)
	debug.Printf("selected mode=%d", mode)

	// wd = instance
	// wl = wordlist slice
	method, methodMode, err := transport.ParseMethod(cfg.Method)
	if err != nil {
		debug.Error("method parse", err)
		fmt.Fprintf(os.Stderr, "[X] Error: invalid method '%s'\n", cfg.Method)
		os.Exit(2)
	}
	switch mode {
	case core.ModeFuzz:
		if !pflag.Lookup("placeholder").Changed {
			cfg.Exts = []string{}
			debug.Printf("fuzz mode without explicit placeholder; extensions disabled")
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
		log.Printf(": %s\n", wd.Wordlist)
		switch {
		case cfg.Timeout > time.Second:
			fmt.Fprintf(os.Stderr, "[!] High timeout (%s). Scan may be slow.\n", cfg.Timeout)
		case cfg.Timeout >= time.Duration(5)*time.Second:
			fmt.Fprintf(os.Stderr, "[!] Very high timeout (%s). Scan will be very slow.\nCTRL + C will take a while (up to 30s).\n", cfg.Timeout)
		}
	case core.ModeDir:
	}
	rawClient := assemble.BuildProxyAndClient(cfg.Proxy, cfg.Timeout, cfg.Insecure)
	client := transport.New(rawClient)
	debug.Printf("http client ready proxy=%t timeout=%s insecure=%t", cfg.Proxy != "", cfg.Timeout, cfg.Insecure)
	if mode == core.ModeDir {
		if !validate.ValidateUrl(cfg.BaseURL, client, method, methodMode, random.RandChoice(cfg.UserAgent)) {
			os.Exit(1)
		}
	}

	wl := wd.LoadWordlist() // Load Wordlist
	debug.Printf("loaded wordlist entries=%d source=%q", len(wl), wd.Wordlist)

	// Basic-Auth
	if cfg.Password != "" && cfg.Username != "" {
		auth = assemble.BuildBasicAuth(cfg.Username, cfg.Password)
		debug.Printf("basic auth enabled user=%q", cfg.Username)
	}

	outputFormat := output.FromFlags(cfg.JSON, cfg.CSV)
	collectOutput := cfg.Output != "" || outputFormat != output.FormatText
	debug.Printf("output format=%d collect_output=%t", outputFormat, collectOutput)

	if !cfg.Quiet && !(collectOutput && cfg.Output == "") {
		fmt.Printf(banner)
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
		Debug:   cfg.Debug,

		// HTTP
		Client:     client,
		Method:     method,
		MethodMode: methodMode,
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
	// Calibration
	debug.Printf("building calibration")

	calibrationPlaceholder := ""
	if mode == core.ModeFuzz {
		calibrationPlaceholder = cfg.Placeholder
	}

	cal, err := calibration.Build(client, calibration.Options{
		BaseURL:     cfg.BaseURL,
		Placeholder: calibrationPlaceholder,
		Tries:      3,
		UserAgents: cfg.UserAgent,
	})
	if err != nil {
		debug.Error("calibration build", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(2)
	}

	debug.Printf(
		"calibration result stable=%t wildcard=%t status=%d length=%d tolerance=%d",
		cal.Stable,
		cal.Wildcard,
		cal.Status,
		cal.Length,
		cal.Tolerance,
	)

	if mode == core.ModeDir && cal.Wildcard {
		if cfg.Method != "GET" && !cfg.ForceHead {
			fmt.Fprintf(os.Stderr, "[!] Wildcard-like behavior detected using HEAD/SWITCH requests.\n")
			fmt.Fprintf(os.Stderr, "You can skip this confirmation with '--force-head'\n")
			fmt.Fprintf(os.Stderr, "HEAD/SWITCH responses do not include a body, so wildcard filtering\ncannot be done reliably and may produce false positives.\n")
			fmt.Fprintf(os.Stderr, "\nSwitch cfg.Method to 'GET'? [y/N]: \n")

			if confirmation.WildcardConfirmation() {
				cfg.Method = "GET"
				method = transport.MethodGET
				methodMode = transport.MethodModeFixed
				engine.Method = method
				engine.MethodMode = methodMode
			}
		}
	}

	engine.Calibration = cal
	results := make([]core.Result, 0)
	for result := range engine.Run(cfg.BaseURL) {
		debug.Printf("result prefix=%s status=%d size=%d url=%s extra=%q", result.Prefix, result.Status, result.Size, result.URL, result.Extra)
		if collectOutput {
			results = append(results, result)
			continue
		}
		tui.Print(result, cfg.Quiet)
	}
	if collectOutput {
		debug.Printf("writing collected results count=%d", len(results))
		if err := output.Write(results, outputFormat, cfg.Output, cfg.Quiet); err != nil {
			debug.Error("output write", err)
			fmt.Fprintf(os.Stderr, "[X] Error writing output: %v\n", err)
			os.Exit(1)
		}
	}
	debug.Printf("scan finished")

}
