package core

import (
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/pkg/random"
	"slices"
	"strings"
)

const prefix = "+"

func (c *Core) RunFuzz(baseURL string) <-chan Result {
	results := make(chan Result)
	debug.Printf("fuzz run start base_url=%q placeholder=%q words=%d exts=%v", baseURL, c.Placeholder, len(c.WL), c.Exts)

	go func() {
		defer close(results)

		if len(c.Exts) > 0 {
			for _, ext := range c.Exts {
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}

				urlParts := strings.Split(baseURL, c.Placeholder)
				templateURL := urlParts[0] + ExtPlaceholder + ext + urlParts[1]

				if _, ok := calibration.Get(templateURL, ExtPlaceholder); ok {
					continue
				}

				if err := calibration.Build(c.Client, calibration.Options{
					BaseURL:     templateURL,
					Placeholder: ExtPlaceholder,
					Tries:       3,
					UserAgents:  c.UserAgents,
				}); err != nil {
					debug.Error("fuzz extension calibration build", err)
					continue
				}
			}
		}

	launch:
		for _, word := range c.WL {

			select {
			case <-c.Ctx.Done():
				debug.Printf("fuzz run canceled before scheduling word=%q", word)
				break launch
			case c.Limiter <- struct{}{}:
			}
			word = strings.TrimLeft(word, "/")

			c.WG.Add(1)

			go func(word string) {

				defer func() { <-c.Limiter }()
				defer c.WG.Done()

				select {
				case <-c.Ctx.Done():
					debug.Printf("fuzz worker canceled word=%q", word)
					return
				default:

				}
				urlParts := strings.Split(baseURL, c.Placeholder)
				fullURL := urlParts[0] + word + urlParts[1]
				headers := c.Header
				if c.AuthHeader != "" {
					headers = append(append([]string{}, headers...), "Authorization: "+c.AuthHeader)
				}
				request := transport.RequestOptions{
					URL:        fullURL,
					Method:     c.nextRequestMethod(),
					MethodMode: transport.MethodModeFixed,
					UserAgent:  random.RandChoice(c.UserAgents),
					Headers:    headers,
				}

				response, err := c.Client.Do(&request)
				if !c.applyDelay("fuzz", fullURL) {
					return
				}
				if err != nil {
					debug.Error("fuzz", err)
					return
				}
				debug.Printf("fuzz response status=%d body=%d", response.StatusCode, response.Lenght)
				status := response.StatusCode
				lenght := response.Lenght

				if len(c.Exts) > 0 {
					ok := c.processExtensions(
						&request,
						results,
						prefix,
						"fuzz-ext",
						func(ext string) string {
							if !strings.HasPrefix(ext, ".") {
								ext = "." + ext
							}

							return urlParts[0] + word + ext + urlParts[1]
						},

						func(ext string) string {
							if !strings.HasPrefix(ext, ".") {
								ext = "." + ext
							}
							
							return urlParts[0] + ExtPlaceholder + ext + urlParts[1]
						},
					)

					if !ok {
						return
					}
				}

				if c.Calibration.Match(status, lenght) {
					debug.Printf("fuzz filtered calibration url=%s status=%d length=%d calibration_status=%d calibration_length=%d tolerance=%d",
						fullURL, status, lenght, c.Calibration.Status, c.Calibration.Length, c.Calibration.Tolerance,
					)
					return
				}
				if slices.Contains(c.IgnoreCodes, status) {
					debug.Printf("fuzz ignored url=%s status=%d", fullURL, status)
					return
				}
				results <- Result{
					Prefix: prefix,
					URL:    fullURL,
					Size:   lenght,
					Status: status,
				}

			}(word)
		}
		c.WG.Wait()
	}()

	return results
}
