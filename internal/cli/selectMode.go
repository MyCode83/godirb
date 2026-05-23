package cli

import (
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/pkg/parse"
	"strings"
)

func SelectMode(mode core.Mode, cfg Config) core.Mode {
	if parse.ExtractPort(cfg.BaseURL) == cfg.Placeholder {
		mode = core.ModePort
		return mode
	}
	if strings.Contains(cfg.BaseURL, cfg.Placeholder) {
		mode = core.ModeFuzz
		return mode
	}

	mode = core.ModeDir
	return mode
}
