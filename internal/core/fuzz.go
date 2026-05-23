package core

import (
	"fmt"
	"os"
	"time"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/pkg/random"
	"slices"
	"strings"

	"github.com/valyala/fasthttp"
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

				debug.Request("fuzz", request)
				err := c.Client.Do(request, response)
				if err != nil {
					debug.Error("fuzz", err)
					fasthttp.ReleaseRequest(request)
					fasthttp.ReleaseResponse(response)
					return
				}
				debug.Response("fuzz", response)
				status := response.StatusCode()
				lenght := len(response.Body())

				if !c.Baseline.IsInteresting(status, lenght, c.Baseline.Tolerance) {
					debug.Printf("fuzz filtered baseline url=%s status=%d length=%d baseline_status=%d baseline_length=%d tolerance=%d",
						fullURL, status, lenght, c.Baseline.Status, c.Baseline.Lenght, c.Baseline.Tolerance)
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
						debug.Request("fuzz-ext", request2)
						err2 := c.Client.Do(request2, response2)

						if err2 != nil {
							debug.Error("fuzz-ext", err2)
							fasthttp.ReleaseRequest(request2)
							fasthttp.ReleaseResponse(response2)
							continue
						}
						debug.Response("fuzz-ext", response2)
						statusCode2 := response2.StatusCode()
						lenght2 := len(response2.Body())
						if !c.Baseline.IsInteresting(statusCode2, lenght2, c.Baseline.Tolerance) {
							debug.Printf("fuzz-ext filtered baseline url=%s status=%d length=%d", urlWithExt, statusCode2, lenght2)
							fasthttp.ReleaseRequest(request2)
							fasthttp.ReleaseResponse(response2)
							continue
						}
						if slices.Contains(c.IgnoreCodes, statusCode2) {
							debug.Printf("fuzz-ext ignored url=%s status=%d", urlWithExt, statusCode2)
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
