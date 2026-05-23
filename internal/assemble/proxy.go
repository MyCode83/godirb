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

func BuildProxyAndClient(proxy string, timeout time.Duration, insecure bool) *fasthttp.Client {
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
		switch {
		case strings.HasPrefix(proxy, "https://"):
			debug.Printf("using HTTP proxy dialer for HTTPS proxy scheme")
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		case strings.HasPrefix(proxy, "http://"):
			debug.Printf("using HTTP proxy dialer")
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		case strings.HasPrefix(proxy, "sock5://"):
			debug.Printf("using SOCKS proxy dialer")
			client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
		case strings.HasPrefix(proxy, "socks5://"):
			debug.Printf("using SOCKS proxy dialer")
			client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
		default:
			debug.Printf("unknown proxy scheme proxy=%q", proxy)
			fmt.Fprintf(os.Stderr, "\n[!] Unkown proxy scheme %s", proxy)
			if !confirmation.ProxyConfirmation() {
				os.Exit(2)
			}
		}
	}
	return client
}
