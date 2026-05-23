package assemble

import (
	"crypto/tls"
	"fmt"
	"github.com/MyCode83/godirb/internal/confirmation"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"os"
	"strings"
	"time"
)

func BuildProxyAndClient(proxy string, timeout time.Duration, insecure bool) *fasthttp.Client {
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
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		case strings.HasPrefix(proxy, "http://"):
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		case strings.HasPrefix(proxy, "sock5://"):
			client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
		case strings.HasPrefix(proxy, "socks5://"):
			client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
		default:
			fmt.Fprintf(os.Stderr, "\n[!] Unkown proxy scheme %s", proxy)
			if !confirmation.ProxyConfirmation() {
				os.Exit(2)
			}
		}
	}
	return client
}
