package cli

import (
	"github.com/MyCode83/godirb/internal/help"
	"github.com/MyCode83/godirb/internal/wordlist"
	"github.com/spf13/pflag"

	"fmt"
)

func ParseFlags() (Config, wordlist.Wordlist) {
	// Objects
	cfg := Config{}
	wd := wordlist.Wordlist{}

	// pflags
	pflag.StringVarP(&cfg.URL, "url", "u", "", "Target URL (e.g. https://localhost)")
	pflag.StringVarP(&wd.Wordlist, "wordlist", "w", "medium", "Embedded wordlist name or path to a custom wordlist. Default: medium")
	pflag.IntVarP(&cfg.Threads, "threads", "t", 15, "Number of worker goroutines to use. Default: 15")
	pflag.IntSliceVarP(&cfg.IgnoreCode, "ignore", "i", []int{404, 400, 405, 408}, "Comma-separated list of status codes to ignore")
	pflag.StringSliceVarP(&cfg.Exts, "ext", "x", nil, "Comma-separated list of file extensions to append")
	pflag.StringVarP(&cfg.RawTimeout, "timeout", "T", "5s", "Request timeout in seconds. Default: 5")
	pflag.StringVarP(&cfg.RawDelay, "delay", "d", "0", "Delay between requests in milliseconds. Default: 0")
	pflag.StringSliceVarP(&cfg.UserAgent, "user-agent", "a", preUserAgents, "Comma-separated list of User-Agents to rotate")
	// bools
	pflag.BoolVarP(&cfg.NoColor, "no-color", "n", false, "Disable colored output") // no color
	pflag.BoolVarP(&cfg.Recursive, "recursive", "r", false, "Enable recursive directory enumeration")
	pflag.BoolVar(&cfg.Debug, "debug", false, "Enable verbose debug output")
	// forces
	pflag.BoolVarP(&cfg.ForceHead, "force-head", "", false, "Skip HEAD/SWITCH wildcard confirmation")
	pflag.BoolVarP(&cfg.ForceThreads, "force-threads", "", false, "Skip high thread-count confirmation")
	pflag.BoolVarP(&cfg.ForceProxy, "force-proxy", "", false, "Skip proxy confirmation")

	pflag.StringVarP(&cfg.Proxy, "proxy", "p", "", "HTTP/S proxy (e.g. http://127.0.0.1:8080)")
	pflag.StringVarP(&cfg.Method, "method", "m", "GET", "HTTP method to use: GET, HEAD, SWITCH (rotate)")
	pflag.StringVarP(&cfg.Username, "user", "U", "", "Username for Basic Auth")
	pflag.StringVarP(&cfg.Password, "password", "P", "", "Password for Basic Auth")
	pflag.StringVar(&cfg.Placeholder, "placeholder", "FUZZ", "Fuzzing placeholder")
	pflag.StringSliceVarP(&cfg.Header, "header", "H", nil, "Add custom HTTP headers")
	pflag.BoolVarP(&cfg.Insecure, "insecure", "k", false, "Skip TLS certificate verification")
	pflag.BoolVarP(&cfg.Quiet, "quiet", "q", false, "Print results in minimal, parse-friendly format")
	pflag.BoolVar(&cfg.JSON, "json", false, "Print results as JSON")
	pflag.BoolVar(&cfg.CSV, "csv", false, "Print results as CSV")
	pflag.StringVarP(&cfg.Output, "output", "o", "", "Write results to file")

	pflag.Usage = func() {
		fmt.Println(help.PrintHelp())
	}
	pflag.Parse()

	return cfg, wd
}
