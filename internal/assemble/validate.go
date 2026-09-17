package assemble

import (
	"os"

	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/validate"
	"github.com/MyCode83/godirb/pkg/random"
)

func ValidateURL(cfg cli.Config, client *transport.Client, mode core.Mode, method transport.Method, methodMode transport.MethodMode) {
	if mode == core.ModeDir {
		if !validate.ValidateUrl(cfg.BaseURL, client, method, methodMode, random.RandChoice(cfg.UserAgent)) {
			os.Exit(1)
		}
	}
}
