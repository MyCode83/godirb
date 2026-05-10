package cli
import (
	"godirb/internal/wordlist"
	"github.com/spf13/pflag"
	"godirb/internal/help"

	"fmt"
)

func ParseFlags() (Config, wordlist.Wordlist){
	// Objects
	cfg := Config{}
	wd := wordlist.Wordlist{}

	// pflags
	pflag.StringVarP(&cfg.URL, "url", "u", "", "Target URL (e.g. https://localhost)")
	pflag.StringVarP(&wd.Wordlist, "wordlist", "w", "medium", "Path to the wordlist. Defaults to raft-medium-directories from SecLists")
	pflag.IntVarP(&cfg.Threads, "threads", "t", 15, "Number of cfg.Threads(goroutines) to use. Default: 15")
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
	pflag.BoolVarP(&cfg.ForceThreads, "force-threads", "", false, "Skip cfg.Threads confirmation")
	pflag.BoolVarP(&cfg.ForceProxy, "force-proxy", "", false, "Skip Proxy confirmation")

	pflag.StringVarP(&cfg.Proxy, "proxy", "p", "", "HTTP/S cfg.Proxy (e.g. http://127.0.0.1:8080)")
	pflag.StringVarP(&cfg.Method, "method", "m", "GET", "HTTP cfg.Method to use: GET, HEAD, SWITCH (rotate)")
	pflag.StringVarP(&cfg.Username, "user", "U", "", "Username for Basic Auth")
	pflag.StringVarP(&cfg.Password, "password", "P", "", "Password for Basic Auth")
	pflag.StringVar(&cfg.Placeholder, "placeholder", "FUZZ", "Fuzzing cfg.Placeholder")
	pflag.StringSliceVarP(&cfg.Header, "header", "H", nil, "Add cfg.Header")
	pflag.BoolVarP(&cfg.Insecure, "insecure", "k", false, "Skip tls validation")
	pflag.BoolVarP(&cfg.Quiet, "quiet", "q", false, "Luego")


	pflag.Usage = func() {
		fmt.Println(help.PrintHelp())
	}
	pflag.Parse()

	return cfg, wd
}