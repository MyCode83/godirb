package core

import (
	"fmt"

	// "github.com/MyCode83/godirb/internal/assemble"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/pkg/random"
	"slices"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
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
				methodSwitch := "GET"
				request := fasthttp.AcquireRequest()
				response := fasthttp.AcquireResponse()

				defer fasthttp.ReleaseRequest(request)
				defer fasthttp.ReleaseResponse(response)

				switch c.Method {
				case "GET", "HEAD", "get", "head":
					request.Header.SetMethod(strings.ToUpper(c.Method))
				case "SWITCH", "switch", "swich":
					if methodSwitch == "GET" {
						methodSwitch = "HEAD"
					} else {
						methodSwitch = "GET"
					}
					request.Header.SetMethod(methodSwitch)
				}

				urlParts := strings.Split(baseUrl, c.Placeholder)
				fullURL := urlParts[0] + word + urlParts[1]
				request.SetRequestURI(fullURL)

				request.Header.SetUserAgent(random.RandChoice(c.UserAgents))
				if c.AuthHeader != "" {
					request.Header.Set("Authorization", c.AuthHeader)
				}
				debug.Request("ports", request)
				err := c.Client.DoTimeout(request, response, c.Timeout)
				if c.Ctx.Err() != nil {
					debug.Printf("ports worker canceled url=%s", fullURL)
					return
				}
				status := response.StatusCode()
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
				debug.Response("ports", response)

				lenght := len(response.Body())

				request.Reset()
				response.Reset()
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
