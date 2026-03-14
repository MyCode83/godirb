package assemble
import (
	"github.com/valyala/fasthttp/fasthttpproxy"
	"strings"
	"fmt"
	"os"
	"time"
	"github.com/valyala/fasthttp"
	"godirb/internal/confirmation"
	"crypto/tls"
)

func BuildProxyAndClient(proxy string, timeout time.Duration, insecure bool) *fasthttp.Client {
	TLS := &tls.Config{
		InsecureSkipVerify: insecure,
	}
	client := &fasthttp.Client{
		ReadTimeout:     timeout,
		WriteTimeout:    timeout,
		MaxConnDuration: timeout,
		TLSConfig: TLS,
	}
	if proxy != "" {
		switch{
		case strings.HasPrefix(proxy, "https://"):
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		case strings.HasPrefix(proxy, "http://"):
			client.Dial = fasthttpproxy.FasthttpHTTPDialer(proxy)
		case strings.HasPrefix(proxy, "sock5://"):
			client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
		case strings.HasPrefix(proxy, "socks5://"):
			client.Dial = fasthttpproxy.FasthttpSocksDialer(proxy)
		default:
			fmt.Fprintf(os.Stderr,"\n[!] Unkown proxy scheme %s", proxy)
			if !confirmation.ProxyConfirmation() {
				os.Exit(2)
			}
		}
	}
	return client
}