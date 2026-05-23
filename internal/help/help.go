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
	-a   --user-agent slice     Comma-separated list of User-Agents to rotate
	-d   --delay string         Delay between requests in milliseconds (default: 0)
	-H   --header slice         Add custom HTTP headers (repeat or comma-separate)
	-h   --help                 Show this help message
	-i   --ignore slice         Comma-separated list of HTTP status codes to ignore (default: 404,400,405,408)
	     --csv                  Print results as CSV
	-k   --insecure             Skip TLS certificate verification
	     --json                 Print results as JSON
	-m   --method string        HTTP method to use: GET, HEAD, SWITCH (rotate)
	-n   --no-color             Disable colored output
	-o   --output string        Write results to file
	-p   --proxy string         HTTP/S proxy (e.g. http://127.0.0.1:8080)
	-P   --password string      Password for Basic Auth
	     --placeholder string   Fuzzing placeholder (default: FUZZ)
	-q   --quiet                Print results in minimal, parse-friendly format
	-r   --recursive            Enable recursive directory enumeration
	-t   --threads int          Number of threads (goroutines) to use (default: 15)
	-T   --timeout string       Request timeout (default: 5s)
	-u   --url string           Target URL (e.g. http://localhost)
	-U   --user string          Username for Basic Auth
	-w   --wordlist string      Embedded wordlist name or path to custom wordlist (default: medium)
	-x   --ext slice            File extensions to append (comma-separated)
	     --force-head           Skip HEAD/SWITCH wildcard confirmation
	     --force-proxy          Skip proxy confirmation
	     --force-threads        Skip high thread-count confirmation

EMBEDDED WORDLISTS:
	small	
	common
	medium
	big
	ports
	payloads
	xss
	lfi

EXAMPLES:
	godirb -u http://localhost
	godirb -u http://localhost -t 5 -a BOT/1.1

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
