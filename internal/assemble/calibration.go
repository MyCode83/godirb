package assemble

import (
	"fmt"
	"os"

	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/cli"
	"github.com/MyCode83/godirb/internal/confirmation"
	"github.com/MyCode83/godirb/internal/core"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
)

func BuildCalibration(
	engine *core.Core,
	cfg *cli.Config,
	mode core.Mode,
	client *transport.Client,
	method *transport.Method,
	methodMode *transport.MethodMode,
) error {
	if mode == core.ModeDir || mode == core.ModeFuzz {
		debug.Printf("building calibration")

		calibrationPlaceholder := ""
		if mode == core.ModeFuzz {
			calibrationPlaceholder = cfg.Placeholder
		}

		err := calibration.Build(client, calibration.Options{
			BaseURL:     cfg.BaseURL,
			Placeholder: calibrationPlaceholder,
			Tries:       3,
			UserAgents:  cfg.UserAgent,
		})
		
		if err != nil {
			return err
		}

		cal, ok := calibration.Get(cfg.BaseURL, calibrationPlaceholder)

		if !ok {
			err := fmt.Errorf(
				"calibration not found after build: base_url=%q placeholder=%q",
				cfg.BaseURL,
				calibrationPlaceholder,
			)
			debug.Error("calibration get", err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}

		debug.Printf(
			"calibration result stable=%t wildcard=%t status=%d length=%d tolerance=%d",
			cal.Stable,
			cal.Wildcard,
			cal.Status,
			cal.Length,
			cal.Tolerance,
		)

		if mode == core.ModeDir && cal.Wildcard {

			if cfg.Method != "GET" && !cfg.ForceHead {
				fmt.Fprintf(os.Stderr, "[!] Wildcard-like behavior detected using HEAD/SWITCH requests.\n")
				fmt.Fprintf(os.Stderr, "You can skip this confirmation with '--force-head'\n")
				fmt.Fprintf(os.Stderr, "HEAD/SWITCH responses do not include a body, so wildcard filtering\ncannot be done reliably and may produce false positives.\n")
				fmt.Fprintf(os.Stderr, "\nSwitch cfg.Method to 'GET'? [y/N]: \n")

				if confirmation.WildcardConfirmation() {
					cfg.Method = "GET"
					*method = transport.MethodGET
					*methodMode = transport.MethodModeFixed
					engine.Method = *method
					engine.MethodMode = *methodMode
				}
			}
		}

		engine.Calibration = cal
	}

	return nil
}
