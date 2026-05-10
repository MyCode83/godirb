package cli

import (
	"godirb/internal/core"
	"godirb/pkg/parse"
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