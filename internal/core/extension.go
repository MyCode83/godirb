package core

import (
	"slices"
	"time"

	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/pkg/random"
)

func (c *Core) processExtensions(
	request *transport.RequestOptions,
	results chan<- Result,
	prefix string,
	debugName string,
	buildURL func(ext string) (string, error),
) bool {
	for _, ext := range c.Exts {
		urlWithExt, err := buildURL(ext)
		if err != nil {
			continue
		}

		request.URL = urlWithExt
		request.Method = c.nextRequestMethod()
		request.UserAgent = random.RandChoice(c.UserAgents)

		response, err := c.Client.Do(request)
		if err != nil {
			debug.Error(debugName, err)
			continue
		}

		debug.Printf(
			"%s response status=%d body=%d",
			debugName,
			response.StatusCode,
			response.Lenght,
		)

		statusCode := response.StatusCode
		length := response.Lenght

		if c.Calibration.Match(statusCode, length) {
			debug.Printf(
				"%s filtered calibration url=%s status=%d length=%d",
				debugName,
				urlWithExt,
				statusCode,
				length,
			)

			continue
		}

		if slices.Contains(c.IgnoreCodes, statusCode) {
			debug.Printf(
				"%s ignored url=%s status=%d",
				debugName,
				urlWithExt,
				statusCode,
			)
			continue
		}

		results <- Result{
			Prefix: prefix,
			URL:    urlWithExt,
			Size:   length,
			Status: statusCode,
		}

		if c.Delay <= 0 {
			continue
		}

		debug.Printf(
			"%s delay=%s url=%s",
			debugName,
			c.Delay,
			urlWithExt,
		)

		select {
		case <-time.After(c.Delay):
		case <-c.Ctx.Done():
			debug.Printf(
				"%s canceled during delay url=%s",
				debugName,
				urlWithExt,
			)
			return false
		}
	}

	return true
}