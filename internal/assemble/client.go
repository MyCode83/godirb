package assemble

import (
	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
)

func BuildClient(cfg cli.Config) *transport.Client {
	rawClient := BuildProxyAndClient(cfg.Proxy, cfg.Timeout, cfg.Insecure)
	client := transport.New(rawClient)
	debug.Printf("http client ready proxy=%t timeout=%s insecure=%t", cfg.Proxy != "", cfg.Timeout, cfg.Insecure)

	return client
}
