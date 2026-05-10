package main

import (
	// stdlib
	"context"
	"fmt"
	"log"
	"runtime"

	"time"

	"os"
	"os/signal"

	"strings"

	"godirb/pkg/parse"
	"godirb/pkg/random"

	"sync"
	"syscall"

	// Third-libs
	"github.com/spf13/pflag"
	"github.com/valyala/fasthttp"

	// Godirb-lib
	"godirb/internal/assemble"
	"godirb/internal/cli"
	"godirb/internal/confirmation"
	"godirb/internal/duration"

	"godirb/internal/core" // core
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
	log.SetOutput(os.Stdout)
	var useColors = true
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

	// wd = instance
	// wl = wordlist slice
	
	if !cfg.Quiet {
		fmt.Printf(banner)
		log.Println("Godirb v", version)
		fmt.Println("Author: MyCode83")

	}
	BaseURL := strings.TrimRight(cfg.URL, "/")
	if cfg.NoColor || cfg.Quiet{
			useColors = false
	}
	timeout, timeoutErr := duration.ParseDuration(cfg.RawTimeout, "s")
	delay, delayErr := duration.ParseDuration(cfg.RawDelay, "ms")
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
	if timeout <= 0 {
		fmt.Printf("[X] Error: timeout must be greater than 0\n")
		os.Exit(2)
	}
	if parse.ExtractPort(BaseURL) == cfg.Placeholder{
		mode = core.ModePort
	} else if strings.Contains(BaseURL, cfg.Placeholder) {
		mode = core.ModeFuzz
	}
	client = assemble.BuildProxyAndClient(cfg.Proxy, timeout, cfg.Insecure) // Fasthttp-Client
	switch mode {
	case core.ModeFuzz:
		if !pflag.Lookup("cfg.Placeholder").Changed {
			cfg.Exts = []string{}
		}
	case core.ModePort:
		if !pflag.Lookup("wordlist").Changed {
			wd.Wordlist = "ports"
		}
		if !pflag.Lookup("timeout").Changed {
			timeout = time.Duration(500) * time.Millisecond
		}
		log.Printf(": %s\n", wd.Wordlist)
		switch  {
		case timeout > time.Second:
			fmt.Printf("[!] High timeout (%s). Scan may be slow.\n", timeout)
		case timeout >= time.Duration(5) * time.Second:
			fmt.Printf("[!] Very high timeout (%s). Scan will be very slow.\nCTRL + C will take a while (up to 30s).\n", timeout)
		}
	case core.ModeDir:
		if !validate.ValidateUrl(BaseURL, client, cfg.Method, random.RandChoice(cfg.UserAgent)){
	 		os.Exit(1)
		}
	}

	wl := wd.LoadWordlist() // Load Wordlist

	// Basic-Auth
	if cfg.Password != "" && cfg.Username != "" {
		auth = assemble.BuildBasicAuth(cfg.Username, cfg.Password)
	}

	
if !cfg.Quiet {
	fmt.Println("\n------------------")
	fmt.Println("[*] Url: ", BaseURL)
	fmt.Println("[*] Method: ", cfg.Method)
	fmt.Println("[*] Threads: ", cfg.Threads)
	fmt.Println("[*] Timeout: ", timeout)
	fmt.Println("[*] UAs: ", len(cfg.UserAgent))
	fmt.Print("[*] Mode: ")
	switch mode{
	case core.ModeDir:
		fmt.Print("Dir\n")
	case core.ModeFuzz:
		fmt.Print("Fuzz\n")
	case core.ModePort:
		fmt.Print("Port\n")
	}
	// nolint
	fmt.Println("------------------\n")
}


	limiter := make(chan struct{}, cfg.Threads)
	var dirsChan chan string
	if mode == core.ModeDir {
		dirsChan = make(chan string, cfg.Threads * 50)
	}

	engine := &core.Core{
		// Mode
		Mode: mode,

		// Bools

		Recursive:  cfg.Recursive,

		// Context
		Ctx:    contextCancel,
		Cancel: cancel,
		// Config
		Timeout: timeout,
		Delay: delay,
		Quiet: cfg.Quiet,

		// HTTP
		Client:      client,
		Method:      cfg.Method,
		UserAgents:  cfg.UserAgent,
		AuthHeader:  auth,
		Header:      cfg.Header,
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
		File: tui.File,

	}
	// Wildcard
	switch mode {
	case core.ModeDir:
		wildcard, err := wildcard.DetectWildcard(client, BaseURL, cfg.Placeholder, cfg.UserAgent...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(2)
		}
		if wildcard.Active {
			fmt.Printf("[!] Wildcard detected: %d | %d bytes \n", wildcard.Status, wildcard.Lenght)
			
			if cfg.Method != "GET" && !cfg.ForceHead{
				fmt.Printf("[!] Wildcard-like behavior detected using HEAD/SWITCH requests.\n")
				fmt.Printf("You can skip this confirmation with puting '--force-head'\n")
				fmt.Printf("HEAD/SWITCH responses do not include a body, so wildcard filtering\ncannot be done reliably and may produce false positives.\n")
				fmt.Printf("\nSwitch cfg.Method to 'GET'? [y/N]: \n")
				
				if confirmation.WildcardConfirmation() {
					cfg.Method = "GET"
				}
			}
		}
		engine.Wildcard = wildcard
	case core.ModeFuzz:
		baseline, err := baseline.BuildBaseLine(BaseURL, client, cfg.Placeholder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		engine.Baseline = baseline
	
	}
	engine.Run(BaseURL)
	wg.Wait()

}