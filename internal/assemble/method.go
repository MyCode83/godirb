package assemble

import (
	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/transport"
)

func BuildMethod(cfg cli.Config) (transport.Method, transport.MethodMode, error) {
	method, methodMode, err := transport.ParseMethod(cfg.Method)
	if err != nil {
		return method, methodMode, err
	}

	return method, methodMode, nil
}
