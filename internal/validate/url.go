package validate

import (
	"fmt"
	"github.com/valyala/fasthttp"
	parseurl "net/url"
	"os"
	"strings"
)

func ValidateUrl(raw string, client *fasthttp.Client, method string, userAgent string) bool {
	methodSwitch := "GET"
	u, err := parseurl.Parse(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\n[X] Err: Invalid URL. Check if you spelled it correctly '%s'\n", raw)
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		fmt.Fprintf(os.Stderr, "\n[X] Err: Invalid Scheme '%s'. Use 'http' or 'https'\n", u.Scheme)
		return false
	}
	if u.Hostname() == "" {
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
		err := client.Do(request, response)
		if err != nil {
			fails++
		} else {
			return true
		}

	}
	fmt.Fprintf(os.Stderr, "\n[X] Err: Invalid URL. Check if you spelled it correctly '%s'\n", raw)
	return false
}
