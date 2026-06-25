package core

import (
	"fmt"

	// "github.com/MyCode83/godirb/internal/assemble"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/pkg/random"
	"slices"
	"strings"
	"time"
)

func looksLikeService(err error) bool {
	s := err.Error()
	result := strings.Contains(s, "tls") ||
		strings.Contains(s, "EOF") ||
		strings.Contains(s, "reset")
	debug.Printf("port service error check error=%q result=%t", s, result)
	return result
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
					Method:     c.Method,
					MethodMode: c.MethodMode,
					UserAgent:  random.RandChoice(c.UserAgents),
					Headers:    headers,
				}
				response, err := c.Client.Do(&request)
				if c.Ctx.Err() != nil {
					debug.Printf("ports worker canceled url=%s", fullURL)
					return
				}
				status := response.StatusCode
				if err != nil {
					debug.Error("ports", err)
					if looksLikeService(err) {
						results <- Result{
							Prefix: "?",
							URL:    fullURL,
							Status: status,
							Extra:  fmt.Sprintf("(error: %v)", err),
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
					Prefix: prefix,
					Size:   lenght,
					Status: status,
					URL:    fullURL,
				}

				if c.Delay > 0 {
					debug.Printf("ports delay=%s url=%s", c.Delay, fullURL)
					select {
					case <-time.After(c.Delay):
					case <-c.Ctx.Done():
						debug.Printf("ports canceled during delay url=%s", fullURL)
						return
					}
				}

			}(word)
		}
		c.WG.Wait()
	}()

	return results

}
