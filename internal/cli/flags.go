package cli
import (
	"godirb/internal/wordlist"
	"github.com/spf13/pflag"
	"godirb/internal/help"

	"fmt"
)

var preUserAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
}

func ParseFlags() (Config, wordlist.Wordlist){
	// Objects
	cfg := Config{}
	wd := wordlist.Wordlist{}

	// pflags
	pflag.StringVarP(&cfg.URL, "cfg.URL", "u", "", "Target URL (e.g. https://localhost)")
	pflag.StringVarP(&wd.Wordlist, "wordlist", "w", "medium", "Path to the wordlist. Defaults to raft-medium-directories from SecLists")
	pflag.IntVarP(&cfg.Threads, "cfg.Threads", "t", 15, "Number of cfg.Threads(goroutines) to use. Default: 15")
	pflag.IntSliceVarP(&cfg.IgnoreCode, "ignore", "i", []int{404, 400, 405, 408}, "Comma-separated list of status codes to ignore")
	pflag.StringSliceVarP(&cfg.Exts, "ext", "x", nil, "")
	pflag.StringVarP(&cfg.RawTimeout, "timeout", "T", "5s", "Request timeout in seconds. Default: 5")
	pflag.StringVarP(&cfg.RawDelay, "delay", "d", "0", "Request timeout in miliseconds. Default: 0") // ! Current
	pflag.StringSliceVarP(&cfg.UserAgent, "user-agent", "a", preUserAgents, "Comma-separated list of User-Agents to rotate")
	// bools
	pflag.BoolVarP(&cfg.NoColor, "no-color", "n", false, "Disable colored output") // no color
	pflag.BoolVarP(&cfg.Recursive, "cfg.Recursive", "r", false, "Recursive mode")      // cfg.Recursive
	// forces
	pflag.BoolVarP(&cfg.ForceHead, "force-head", "", false, "Skip wd confirmation")
	pflag.BoolVarP(&cfg.ForceThreads, "force-cfg.Threads", "", false, "Skip cfg.Threads confirmation")
	pflag.BoolVarP(&cfg.ForceProxy, "force-proxy", "", false, "Skip Proxy confirmation")

	pflag.StringVarP(&cfg.Proxy, "cfg.Proxy", "p", "", "HTTP/S cfg.Proxy (e.g. http://127.0.0.1:8080)")
	pflag.StringVarP(&cfg.Method, "cfg.Method", "m", "GET", "HTTP cfg.Method to use: GET, HEAD, SWITCH (rotate)")
	pflag.StringVarP(&cfg.Username, "user", "U", "", "Username for Basic Auth")
	pflag.StringVarP(&cfg.Password, "cfg.Password", "P", "", "Password for Basic Auth")
	pflag.StringVar(&cfg.Placeholder, "cfg.Placeholder", "FUZZ", "Fuzzing cfg.Placeholder")
	pflag.StringSliceVarP(&cfg.Header, "cfg.Header", "H", nil, "Add cfg.Header")
	pflag.BoolVarP(&cfg.Insecure, "cfg.Insecure", "k", false, "Skip tls validation")
	pflag.BoolVarP(&cfg.Quiet, "cfg.Quiet", "q", false, "Luego")


	pflag.Usage = func() {
		fmt.Println(help.PrintHelp())
	}
	pflag.Parse()

	return cfg, wd
}