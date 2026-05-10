package core

import (
	"fmt"
	"os"
	"time"

	"godirb/pkg/random"
	"slices"
	"strings"

	"github.com/valyala/fasthttp"
)

const prefix = "+"

func (c *Core) RunFuzz(baseURL string) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		if c.Baseline == nil {
			fmt.Fprintf(os.Stderr, "[!] Baseline is nil")
			return
		}

	launch:
		for _, word := range c.WL {

			select {
			case <-c.Ctx.Done():
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
					return
				default:

				}
				methodSwitch := "GET"
				request := fasthttp.AcquireRequest()
				response := fasthttp.AcquireResponse()
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

				urlParts := strings.Split(baseURL, c.Placeholder)
				fullURL := urlParts[0] + word + urlParts[1]

				request.SetRequestURI(fullURL)

				request.Header.SetUserAgent(random.RandChoice(c.UserAgents))
				if c.AuthHeader != "" {
					request.Header.Set("Authorization", c.AuthHeader)
				}

				err := c.Client.Do(request, response)
				if err != nil {
					fasthttp.ReleaseRequest(request)
					fasthttp.ReleaseResponse(response)
					return
				}
				status := response.StatusCode()
				lenght := len(response.Body())

				if !c.Baseline.IsInteresting(status, lenght, c.Baseline.Tolerance) {
					fasthttp.ReleaseRequest(request)
					fasthttp.ReleaseResponse(response)
					return
				}

				fasthttp.ReleaseRequest(request)
				fasthttp.ReleaseResponse(response)
				if len(c.Exts) > 0 {
					for _, ext := range c.Exts {
						urlWithExt := urlParts[0] + word + "." + ext + urlParts[1]
						request2 := fasthttp.AcquireRequest()
						response2 := fasthttp.AcquireResponse()
						switch c.Method {
						case "GET", "HEAD", "get", "head":
							request2.Header.SetMethod(strings.ToUpper(c.Method))
						case "SWITCH", "switch", "swich":
							if methodSwitch == "GET" {
								methodSwitch = "HEAD"
							} else {
								methodSwitch = "GET"
							}
							request2.Header.SetMethod(strings.ToUpper(methodSwitch))

						}
						request2.SetRequestURI(urlWithExt)
						request2.Header.SetUserAgent(random.RandChoice(c.UserAgents))
						if c.AuthHeader != "" {
							request2.Header.Set("Authorization", c.AuthHeader)
						}
						err2 := c.Client.Do(request2, response2)

						if err2 != nil {
							fasthttp.ReleaseRequest(request2)
							fasthttp.ReleaseResponse(response2)
							continue
						}
						statusCode2 := response2.StatusCode()
						lenght2 := len(response2.Body())
						if !c.Baseline.IsInteresting(statusCode2, lenght2, c.Baseline.Tolerance) {
							fasthttp.ReleaseRequest(request2)
							fasthttp.ReleaseResponse(response2)
							continue
						}
						if slices.Contains(c.IgnoreCodes, statusCode2) {
							fasthttp.ReleaseRequest(request2)
							fasthttp.ReleaseResponse(response2)
							continue
						}
						fasthttp.ReleaseRequest(request2)
						fasthttp.ReleaseResponse(response2)

						results <- Result{
							Prefix: prefix,
							URL:    urlWithExt,
							Size:   lenght2,
							Status: statusCode2,
						}

						if c.Delay > 0 {
							select {
							case <-time.After(c.Delay):
							case <-c.Ctx.Done():
								return
							}
						}

					}
				}
				if slices.Contains(c.IgnoreCodes, status) {
					return
				}
				results <- Result{
					Prefix: prefix,
					URL:    fullURL,
					Size:   lenght,
					Status: status,
				}

				if c.Delay > 0 {
					select {
					case <-time.After(c.Delay):
					case <-c.Ctx.Done():
						return
					}
				}

			}(word)
		}
		c.WG.Wait()
	}()

	return results
}
