package validate

import (
	"fmt"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	parseurl "net/url"
	"os"
)

func ValidateUrl(raw string, client *transport.Client, method transport.Method, methodMode transport.MethodMode, userAgent string) bool {
	debug.Printf("validating url=%q method=%q user_agent=%q", raw, method, userAgent)
	u, err := parseurl.Parse(raw)
	if err != nil {
		debug.Error("validate url parse", err)
		fmt.Fprintf(os.Stderr, "\n[X] Err: Invalid URL. Check if you spelled it correctly '%s'\n", raw)
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		debug.Printf("validate url invalid scheme=%q", u.Scheme)
		fmt.Fprintf(os.Stderr, "\n[X] Err: Invalid Scheme '%s'. Use 'http' or 'https'\n", u.Scheme)
		return false
	}
	if u.Hostname() == "" {
		debug.Printf("validate url empty hostname")
		fmt.Fprintf(os.Stderr, "\n[X] Err: Empty host '%s'. Cheack the URL\n", u.Hostname())
		return false
	}
	request := transport.RequestOptions{
		URL:        raw,
		Method:     method,
		MethodMode: methodMode,
		UserAgent:  userAgent,
	}
	fails := 0
	for range 5 {
		response, err := client.Do(&request)
		if err != nil {
			debug.Error("validate-url", err)
			fails++
		} else {
			debug.Printf("validate-url response status=%d body=%d", response.StatusCode, response.Lenght)
			debug.Printf("validate url success fails_before_success=%d", fails)
			return true
		}

	}
	debug.Printf("validate url failed attempts=%d", fails)
	fmt.Fprintf(os.Stderr, "\n[X] Err: Invalid URL. Check if you spelled it correctly '%s'\n", raw)
	return false
}
