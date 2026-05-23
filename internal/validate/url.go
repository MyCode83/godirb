package validate

import (
	"fmt"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/valyala/fasthttp"
	parseurl "net/url"
	"os"
	"strings"
)

func ValidateUrl(raw string, client *fasthttp.Client, method string, userAgent string) bool {
	debug.Printf("validating url=%q method=%q user_agent=%q", raw, method, userAgent)
	methodSwitch := "GET"
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
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(request)
	defer fasthttp.ReleaseResponse(response)
	request.SetRequestURI(raw)
	request.Header.SetUserAgent(userAgent)
	switch method {
	case "GET", "HEAD", "get", "head":
		request.Header.SetMethod(strings.ToUpper(method))
	case "SWITCH", "switch", "swich":
		if methodSwitch == "GET" {
			methodSwitch = "HEAD"
		} else {
			methodSwitch = "GET"
		}
		request.Header.SetMethod(methodSwitch)
	}
	fails := 0
	for range 5 {
		debug.Request("validate-url", request)
		err := client.Do(request, response)
		if err != nil {
			debug.Error("validate-url", err)
			fails++
		} else {
			debug.Response("validate-url", response)
			debug.Printf("validate url success fails_before_success=%d", fails)
			return true
		}

	}
	debug.Printf("validate url failed attempts=%d", fails)
	fmt.Fprintf(os.Stderr, "\n[X] Err: Invalid URL. Check if you spelled it correctly '%s'\n", raw)
	return false
}
