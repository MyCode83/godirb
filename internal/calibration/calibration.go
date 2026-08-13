package calibration

import (
	"fmt"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/pkg/random"
)

const (
	defaultTries        = 3
	defaultRandomLength = 12
)

func Build(client *transport.Client, opts Options) error {
	if client == nil {
		return fmt.Errorf("nil transport client")
	}

	if opts.Tries <= 0 {
		opts.Tries = defaultTries
	}

	samples := make([]Sample, 0, opts.Tries)

	debug.Printf(
		"calibration start base_url=%q placeholder=%q tries=%d",
		opts.BaseURL,
		opts.Placeholder,
		opts.Tries,
	)

	for range opts.Tries {
		url, err := generateURL(opts.BaseURL, opts.Placeholder)
		if err != nil {
			return err
		}

		response, err := client.Do(&transport.RequestOptions{
			URL:        url,
			Method:     transport.MethodGET,
			MethodMode: transport.MethodModeFixed,
			UserAgent:  random.RandChoice(opts.UserAgents),
		})
		if err != nil {
			debug.Error("calibration", err)
			return err
		}

		debug.Printf(
			"calibration sample url=%q status=%d length=%d",
			url,
			response.StatusCode,
			response.Lenght,
		)

		samples = append(samples, Sample{
			URL:    url,
			Status: response.StatusCode,
			Length: response.Lenght,
		})
	}

	calibration := buildSignature(samples)

	entries = append(entries, entry{
		BaseURL:     opts.BaseURL,
		Placeholder: opts.Placeholder,
		Calibration: calibration,
	},
	)

	return nil
}
