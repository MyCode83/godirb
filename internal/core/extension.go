package core

import (
	"slices"

	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/pkg/random"
)

func (c *Core) processExtensions(
	request *transport.RequestOptions,
	results chan<- Result,
	prefix string,
	debugName string,
	buildURL func(ext string) string,
	buildCalibrationURL func(ext string) string,
) bool {
	for _, ext := range c.Exts {
		cal, ok := calibration.Get(buildCalibrationURL(ext), ExtPlaceholder)
		if !ok {
			debug.Printf("extension %q calibration not exists continue", ext)
			continue
		}

		urlWithExt := buildURL(ext)

		request.URL = urlWithExt
		request.Method = c.nextRequestMethod()
		request.UserAgent = random.RandChoice(c.UserAgents)

		response, err := c.Client.Do(request)
		if !c.applyDelay(debugName, urlWithExt) {
			return false
		}
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

		if cal.Match(statusCode, length) {
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
	}

	return true
}