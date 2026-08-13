package core

import (
	"errors"
	"net"

	// "github.com/MyCode83/godirb/internal/assemble"

	"slices"
	"strings"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/pkg/random"
)

func looksLikeOpenService(err error) bool {
	if err == nil {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		debug.Printf("port error classified as timeout error=%q", err)
		return false
	}

	message := strings.ToLower(err.Error())

	closedOrUnknown := []string{
		"connection refused",
		"no route to host",
		"network is unreachable",
		"host is down",
		"context deadline exceeded",
		"i/o timeout",
		"handshake timed out",
	}

	for _, pattern := range closedOrUnknown {
		if strings.Contains(message, pattern) {
			debug.Printf(
				"port error classified as closed-or-filtered error=%q pattern=%q",
				err,
				pattern,
			)
			return false
		}
	}

	for _, pattern := range openServiceSignals {
		if strings.Contains(message, pattern) {
			debug.Printf(
				"port error classified as probable-service error=%q pattern=%q",
				err,
				pattern,
			)
			return true
		}
	}

	debug.Printf("port error classified as unknown error=%q", err)
	return false
}

func (c *Core) RunPorts(baseUrl string) <-chan Result {
	results := make(chan Result)
	debug.Printf("ports run start base_url=%q words=%d timeout=%s", baseUrl, len(c.WL), c.Timeout)

	go func() {
		defer close(results)

	launch:
		for _, word := range c.WL {

			select {
			case <-c.Ctx.Done():
				debug.Printf("ports run canceled before scheduling word=%q", word)
				break launch
			case c.Limiter <- struct{}{}:
			}

			word = strings.TrimLeft(word, "/")

			c.WG.Add(1)

			go func(word string) {
				defer c.WG.Done()
				defer func() { <-c.Limiter }()

				urlParts := strings.Split(baseUrl, c.Placeholder)
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
				if !c.applyDelay("ports", fullURL) {
					return
				}
				if c.Ctx.Err() != nil {
					debug.Printf("ports worker canceled url=%s", fullURL)
					return
				}
				status := response.StatusCode
				if err != nil {
					debug.Error("ports", err)
					if looksLikeOpenService(err) {
						results <- Result{
							Kind:   "UNKNOWN",
							URL:    fullURL,
							Status: status,
							Error:  err.Error(),

							Method: response.Method.String(),
							Duration: response.Duration.String(),
						}
					}
					return
				}
				debug.Printf("ports response status=%d body=%d", response.StatusCode, response.Lenght)

				lenght := response.Lenght

				if slices.Contains(c.IgnoreCodes, status) {
					debug.Printf("ports ignored url=%s status=%d", fullURL, status)
					return
				}
				results <- Result{
					Kind:   prefix,
					Size:   lenght,
					Status: status,
					URL:    fullURL,

					Method: response.Method.String(),
					ContentType: response.ContentType,
					ContentLength: response.ContentLenght,
					Location: response.Location,
					Duration: response.Duration.String(),
				}

			}(word)
		}
		c.WG.Wait()
	}()

	return results

}
