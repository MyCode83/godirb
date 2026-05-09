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

	"godirb/internal/confirmation"
	"godirb/internal/duration"
	"godirb/internal/wordlist"

	"godirb/internal/core" // core
	"godirb/internal/help"
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

// Flags
var (
	url        string
	threads    int
	ignoreCode []int
	exts       []string
	rawTimeout  string
	rawDelay	string
	userAgent  []string
	// bool
	noColor   bool
	recursive bool
	//forces
	forceHead    bool
	forceThreads bool
	forceProxy   bool

	proxy  string
	method string
	// Basic Auth
	username string
	password string
	// Placeholder
	placeholder string

	header []string
	// TLS
	insecure bool

	quiet bool
)

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
	wd := wordlist.Wordlist{}
	// pflags
	pflag.StringVarP(&url, "url", "u", "", "Target URL (e.g. https://localhost)")
	pflag.StringVarP(&wd.Wordlist, "wordlist", "w", "medium", "Path to the wordlist. Defaults to raft-medium-directories from SecLists")
	pflag.IntVarP(&threads, "threads", "t", 15, "Number of threads(goroutines) to use. Default: 15")
	pflag.IntSliceVarP(&ignoreCode, "ignore", "i", []int{404, 400, 405, 408}, "Comma-separated list of status codes to ignore")
	pflag.StringSliceVarP(&exts, "ext", "x", nil, "")
	pflag.StringVarP(&rawTimeout, "timeout", "T", "5s", "Request timeout in seconds. Default: 5")
	pflag.StringVarP(&rawDelay, "delay", "d", "0", "Request timeout in miliseconds. Default: 0") // ! Current
	pflag.StringSliceVarP(&userAgent, "user-agent", "a", preUserAgents, "Comma-separated list of User-Agents to rotate")
	// bools
	pflag.BoolVarP(&noColor, "no-color", "n", false, "Disable colored output") // no color
	pflag.BoolVarP(&recursive, "recursive", "r", false, "Recursive mode")      // recursive
	// forces
	pflag.BoolVarP(&forceHead, "force-head", "", false, "Skip wd confirmation")
	pflag.BoolVarP(&forceThreads, "force-threads", "", false, "Skip threads confirmation")
	pflag.BoolVarP(&forceProxy, "force-proxy", "", false, "Skip proxy confirmation")

	pflag.StringVarP(&proxy, "proxy", "p", "", "HTTP/S proxy (e.g. http://127.0.0.1:8080)")
	pflag.StringVarP(&method, "method", "m", "GET", "HTTP method to use: GET, HEAD, SWITCH (rotate)")
	pflag.StringVarP(&username, "user", "U", "", "Username for Basic Auth")
	pflag.StringVarP(&password, "password", "P", "", "Password for Basic Auth")
	pflag.StringVar(&placeholder, "placeholder", "FUZZ", "Fuzzing Placeholder")
	pflag.StringSliceVarP(&header, "header", "H", nil, "Add Header")
	pflag.BoolVarP(&insecure, "insecure", "k", false, "Skip tls validation")
	pflag.BoolVarP(&quiet, "quiet", "q", false, "Luego")


	pflag.Usage = func() {
		fmt.Println(help.PrintHelp())
	}
	pflag.Parse()
	
	if !quiet {
		fmt.Printf(banner)
		log.Println("Godirb v", version)
		fmt.Println("Author: MyCode83")

	}
	BaseURL := strings.TrimRight(url, "/")
	if noColor || quiet{
			useColors = false
	}
	timeout, timeoutErr := duration.ParseDuration(rawTimeout, "s")
	delay, delayErr := duration.ParseDuration(rawDelay, "ms")
	// Delay error
	if delayErr != nil {
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'", rawDelay)
		os.Exit(1)		
	}
	// Timeout error
	if timeoutErr != nil {
		fmt.Fprintf(os.Stderr, "[X] An error ocurred during parsing durantion.\nCheck if you spelled it correctly '%s'", rawTimeout)
		os.Exit(1)
	}
	if strings.TrimSpace(url) == "" {
		fmt.Fprintln(os.Stderr, "[X] Error, missing '--url'(-u) flag")
		fmt.Fprintln(os.Stderr, "Run godirb --help for usage")
		os.Exit(2)
	}
	if threads >= 250 && !forceThreads {

		if !confirmation.ThreadsConfirmation(fmt.Sprintf("Are you shure you want to send %d threads", threads)) {
			fmt.Println("Cancelled by the user")
			os.Exit(0)
		}

	}
	if threads <= 0 {
		fmt.Fprintf(os.Stderr, "[X] Error, you can't send %d threads.\n", threads)
		os.Exit(2)
	}
	if !useColors {
		color.NoColor = true
	}
	if timeout <= 0 {
		fmt.Printf("[X] Error: timeout must be greater than 0\n")
		os.Exit(2)
	}
	if parse.ExtractPort(BaseURL) == placeholder{
		mode = core.ModePort
	} else if strings.Contains(BaseURL, placeholder) {
		mode = core.ModeFuzz
	}
	client = assemble.BuildProxyAndClient(proxy, timeout, insecure) // Fasthttp-Client
	switch mode {
	case core.ModeFuzz:
		if !pflag.Lookup("placeholder").Changed {
			exts = []string{}
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
		if !validate.ValidateUrl(BaseURL, client, method, random.RandChoice(userAgent)){
	 		os.Exit(1)
		}
	}

	wl := wd.LoadWordlist() // Load Wordlist

	// Basic-Auth
	if password != "" && username != "" {
		auth = assemble.BuildBasicAuth(username, password)
	}

	
if !quiet {
	fmt.Println("\n\n------------------")
	fmt.Println("[*] Url: ", BaseURL)
	fmt.Println("[*] Method: ", method)
	fmt.Println("[*] Threads: ", threads)
	fmt.Println("[*] Timeout: ", timeout)
	fmt.Println("[*] UAs: ", len(userAgent))
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


	limiter := make(chan struct{}, threads)
	var dirsChan chan string
	if mode == core.ModeDir {
		dirsChan = make(chan string, threads * 50)
	}

	engine := &core.Core{
		// Mode
		Mode: mode,

		// Bools

		Recursive:  recursive,

		// Context
		Ctx:    contextCancel,
		Cancel: cancel,
		// Config
		Timeout: timeout,
		Delay: delay,
		Quiet: quiet,

		// HTTP
		Client:      client,
		Method:      method,
		UserAgents:  userAgent,
		AuthHeader:  auth,
		Header:      header,
		// Placeholder
		Placeholder: placeholder,
		// Control
		IgnoreCodes: ignoreCode,
		Exts:        exts,

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
		wildcard, err := wildcard.DetectWildcard(client, BaseURL, placeholder, userAgent...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(2)
		}
		if wildcard.Active {
			fmt.Printf("[!] Wildcard detected: %d | %d bytes \n", wildcard.Status, wildcard.Lenght)
			
			if method != "GET" && !forceHead{
				fmt.Printf("[!] Wildcard-like behavior detected using HEAD/SWITCH requests.\n")
				fmt.Printf("You can skip this confirmation with puting '--force-head'\n")
				fmt.Printf("HEAD/SWITCH responses do not include a body, so wildcard filtering\ncannot be done reliably and may produce false positives.\n")
				fmt.Printf("\nSwitch method to 'GET'? [y/N]: \n")
				
				if confirmation.WildcardConfirmation() {
					method = "GET"
				}
			}
		}
		engine.Wildcard = wildcard
	case core.ModeFuzz:
		baseline, err := baseline.BuildBaseLine(BaseURL, client, placeholder)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		engine.Baseline = baseline
	
	}
	engine.Run(BaseURL)
	wg.Wait()

}