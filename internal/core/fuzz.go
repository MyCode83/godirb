package core

import (
	"fmt"
	"os"
	"time"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
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

		if c.Baseline == nil {
			debug.Printf("fuzz run stopped: baseline is nil")
			fmt.Fprintf(os.Stderr, "[!] Baseline is nil")
			return
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
				if err != nil {
					debug.Error("fuzz", err)
					return
				}
				debug.Printf("fuzz response status=%d body=%d", response.StatusCode, response.Lenght)
				status := response.StatusCode
				lenght := response.Lenght

				if !c.Baseline.IsInteresting(status, lenght, c.Baseline.Tolerance) {
					debug.Printf("fuzz filtered baseline url=%s status=%d length=%d baseline_status=%d baseline_length=%d tolerance=%d",
						fullURL, status, lenght, c.Baseline.Status, c.Baseline.Lenght, c.Baseline.Tolerance)
					return
				}

				if len(c.Exts) > 0 {
					for _, ext := range c.Exts {
						urlWithExt := urlParts[0] + word + "." + ext + urlParts[1]
						request.URL = urlWithExt
						request.Method = c.nextRequestMethod()
						request.UserAgent = random.RandChoice(c.UserAgents)
						response2, err2 := c.Client.Do(&request)

						if err2 != nil {
							debug.Error("fuzz-ext", err2)
							continue
						}
						debug.Printf("fuzz-ext response status=%d body=%d", response2.StatusCode, response2.Lenght)
						statusCode2 := response2.StatusCode
						lenght2 := response2.Lenght
						if !c.Baseline.IsInteresting(statusCode2, lenght2, c.Baseline.Tolerance) {
							debug.Printf("fuzz-ext filtered baseline url=%s status=%d length=%d", urlWithExt, statusCode2, lenght2)
							continue
						}
						if slices.Contains(c.IgnoreCodes, statusCode2) {
							debug.Printf("fuzz-ext ignored url=%s status=%d", urlWithExt, statusCode2)
							continue
						}

						results <- Result{
							Prefix: prefix,
							URL:    urlWithExt,
							Size:   lenght2,
							Status: statusCode2,
						}

						if c.Delay > 0 {
							debug.Printf("fuzz-ext delay=%s url=%s", c.Delay, urlWithExt)
							select {
							case <-time.After(c.Delay):
							case <-c.Ctx.Done():
								debug.Printf("fuzz-ext canceled during delay url=%s", urlWithExt)
								return
							}
						}

					}
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

				if c.Delay > 0 {
					debug.Printf("fuzz delay=%s url=%s", c.Delay, fullURL)
					select {
					case <-time.After(c.Delay):
					case <-c.Ctx.Done():
						debug.Printf("fuzz canceled during delay url=%s", fullURL)
						return
					}
				}

			}(word)
		}
		c.WG.Wait()
	}()

	return results
}
