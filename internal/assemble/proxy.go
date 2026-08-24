package assemble

import (
	"crypto/tls"
	"fmt"
	"github.com/MyCode83/godirb/internal/confirmation"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"os"
	"strings"
	"time"
)

func BuildProxyAndClient(proxy string, timeout time.Duration, insecure bool, forceProxy bool) *fasthttp.Client {
	debug.Printf("building client timeout=%s insecure=%t proxy=%q", timeout, insecure, proxy)
	TLS := &tls.Config{
		InsecureSkipVerify: insecure,
	}
	client := &fasthttp.Client{

		ReadTimeout:     timeout,
		WriteTimeout:    timeout,
		MaxConnDuration: timeout,
		TLSConfig:       TLS,
	}
	if proxy != "" {
		switch proxyScheme(proxy) {
		case "https":
			debug.Printf("using HTTP proxy dialer for HTTPS proxy scheme")
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(httpProxyDialerAddress(proxy))
		case "http":
			debug.Printf("using HTTP proxy dialer")
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(httpProxyDialerAddress(proxy))
		case "sock5":
			debug.Printf("using SOCKS proxy dialer")
			client.Dial = fasthttpproxy.FasthttpSocksDialer(socksProxyDialerAddress(proxy))
		case "socks5":
			debug.Printf("using SOCKS proxy dialer")
			client.Dial = fasthttpproxy.FasthttpSocksDialer(socksProxyDialerAddress(proxy))
		default:
			debug.Printf("unknown proxy scheme proxy=%q", proxy)
			fmt.Fprintf(os.Stderr, "\n[!] Unknown proxy scheme %s. Continue without proxy", proxy)
			if !confirmUnknownProxy(forceProxy, confirmation.ProxyConfirmation) {
				os.Exit(2)
			}
			if forceProxy {
				fmt.Fprintln(os.Stderr, ".")
			}
		}
	}
	return client
}

func proxyScheme(proxy string) string {
	scheme, _, ok := strings.Cut(proxy, "://")
	if !ok {
		return ""
	}
	return strings.ToLower(scheme)
}

func proxyAllowed(proxy string) bool {
	switch proxyScheme(proxy) {
	case "http", "https", "socks5", "sock5":
		return true
	default:
		return false
	}
}

func httpProxyDialerAddress(proxy string) string {
	switch proxyScheme(proxy) {
	case "http":
		_, address, _ := strings.Cut(proxy, "://")
		return address
	case "https":
		_, address, _ := strings.Cut(proxy, "://")
		return address
	default:
		return proxy
	}
}

func socksProxyDialerAddress(proxy string) string {
	if proxyScheme(proxy) != "sock5" {
		return proxy
	}
	_, address, _ := strings.Cut(proxy, "://")
	return "socks5://" + address
}

func confirmUnknownProxy(forceProxy bool, confirm func() bool) bool {
	return forceProxy || confirm()
}
