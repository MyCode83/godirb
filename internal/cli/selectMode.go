package cli

import (
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/pkg/parse"
	"strings"
)

func SelectMode(mode core.Mode, cfg Config) core.Mode {
	if parse.ExtractPort(cfg.BaseURL) == cfg.Placeholder {
		mode = core.ModePort
		debug.Printf("select mode=port base_url=%q placeholder=%q", cfg.BaseURL, cfg.Placeholder)
		return mode
	}
	if strings.Contains(cfg.BaseURL, cfg.Placeholder) {
		mode = core.ModeFuzz
		debug.Printf("select mode=fuzz base_url=%q placeholder=%q", cfg.BaseURL, cfg.Placeholder)
		return mode
	}

	mode = core.ModeDir
	debug.Printf("select mode=dir base_url=%q placeholder=%q", cfg.BaseURL, cfg.Placeholder)
	return mode
}
