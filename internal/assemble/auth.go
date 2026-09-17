package assemble

import (
	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/debug"
)

func BuildAuth(cfg cli.Config) string {
	var auth string
	if cfg.Password != "" && cfg.Username != "" {
		auth = BuildBasicAuth(cfg.Username, cfg.Password)
		debug.Printf("basic auth enabled user=%q", cfg.Username)
	}

	return auth
}
