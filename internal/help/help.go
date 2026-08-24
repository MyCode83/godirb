package help

func PrintHelp() string {
	const help = `
   ____ _  ____   ____/ /   (_)   _____   / /
  / __  / / __ \ / __  /   / /   / ___/  / __ \
 / /_/ / / /_/ // /_/ /   / /   / /     / /_/ /
 \__  /  \____/ \____/   /_/   /_/     /_____/
/____/
godirb - fast directory brute-forcer built in Go

USAGE:
	./godirb -u [TARGET] [OPTIONS] 

FLAGS: 
	-a   --user-agent slice     User-Agents to rotate; repeat to add multiple values
	-d   --delay string         Delay between requests; accepts 200ms or 2s, bare numbers are milliseconds (default: 0)
	     --depth int            Maximum recursive depth with --recursive; -1 means unlimited (default: 2)
	     --debug                Enable verbose debug output
	-H   --header slice         Add custom HTTP headers; repeat to add multiple values
	-h   --help                 Show this help message
	-i   --ignore slice         Comma-separated list of HTTP status codes to ignore (default: 404,400,405,408)
	     --csv                  Print results as CSV
	-k   --insecure             Skip TLS certificate verification
	     --json                 Print results as JSON Lines, one object per line
	-m   --method string        HTTP method to use: GET, HEAD, SWITCH (rotate)
	-n   --no-color             Disable colored output
	-o   --output string        Write results to file
	-p   --proxy string         HTTP, HTTPS, or SOCKS5 proxy (e.g. http://127.0.0.1:8080)
	-P   --password string      Password for Basic Auth
	     --placeholder string   Fuzzing placeholder (default: FUZZ)
	-q   --quiet                Print results in minimal, parse-friendly format
	-r   --recursive            Enable recursive directory enumeration
	-t   --threads int          Number of threads (goroutines) to use (default: 15)
	-T   --timeout string       Request timeout; accepts 500ms or 2s, bare numbers are seconds (default: 5s; port mode default: 500ms)
	-u   --url string           Target URL (e.g. http://localhost)
	-U   --user string          Username for Basic Auth
	     --version              Print version and exit
	-w   --wordlist string      Embedded wordlist name, path, or - for stdin (default: medium)
	-x   --ext slice            File extensions to append (comma-separated)
	     --force-head           Skip HEAD/SWITCH wildcard confirmation
	     --force-proxy          Skip unknown proxy scheme confirmation and continue direct
	     --force-threads        Skip high thread-count confirmation

EMBEDDED WORDLISTS:
	small	
	medium
	big
	ports
	payloads
	xss
	lfi

EXAMPLES:
	godirb -u http://localhost
	godirb -u http://localhost -t 5 -a BOT/1.1
	godirb -u http://localhost -a "Scanner/1.0" -a "Mozilla/5.0"
	godirb -u http://localhost -H "Authorization: Bearer TOKEN" -H "X-Test: value"
	godirb -u https://example.com --json -o results.json
	godirb -u https://example.com:FUZZ --csv -o ports.csv

	godirb -u http://localhost:FUZZ
	godirb -u "http://localhost?msg=FUZZ" -w xss

NOTES:
	-  If you do not specify a wordlist, godirb uses the embedded medium wordlist
	-  If you do not want colors or your terminal does not support them, use -n, --no-color or NO_COLOR=1
	-  If you want to disable colors in godirb permanently, set GODIRB_NO_COLOR=1

CREDITS:
	-  Credits to SecLists: https://github.com/danielmiessler/SecLists (MIT LICENSE)
	-  Inspired by dirb and gobuster
`
	return help
}
