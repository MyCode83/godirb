package core

import (
	"fmt"
	"time"

	"godirb/internal/detention"
	"godirb/internal/wildcard"

	"godirb/pkg/random"
	"os"
	"slices"
	"strings"

	"github.com/valyala/fasthttp"
)

func (c *Core) RunDir(baseURL string) <-chan Result {
	results := make(chan Result)

	go func() {
		defer close(results)

		if c.Wildcard == nil {
			fmt.Fprintf(os.Stderr, "[X] Wildcard is nil")
			return
		}

		c.WG.Add(1)
		c.DirsChan <- baseURL

		go func() {

			c.WG.Wait()
			close(c.DirsChan)

		}()

		// Dirs loop
	dirLoop:
		for dir := range c.DirsChan {

			// Wordlist loop
			for _, word := range c.WL {
				select {
				case <-c.Ctx.Done():
					c.WG.Done()
					c.WG.Wait()
					break dirLoop
				case c.Limiter <- struct{}{}:
				}
				word = strings.TrimLeft(word, "/")
				c.WG.Add(1)
				go func(word string) {
					dirPrefix := ""

					// Request/Response
					request := fasthttp.AcquireRequest()
					response := fasthttp.AcquireResponse()
					// Release
					defer fasthttp.ReleaseRequest(request)
					defer fasthttp.ReleaseResponse(response)

					// Request/Response of extensions
					request2 := fasthttp.AcquireRequest()
					response2 := fasthttp.AcquireResponse()
					// Release
					defer fasthttp.ReleaseRequest(request2)
					defer fasthttp.ReleaseResponse(response2)

					// Reset Headers, Methods... without release
					request.Reset()
					response.Reset()

					defer c.WG.Done()

					defer func() { <-c.Limiter }()
					select {
					case <-c.Ctx.Done():
						return
					default:

					}
					methodSwitch := "GET"

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

					fullURL := fmt.Sprintf("%s/%s", dir, word)
					request.SetRequestURI(fullURL)

					request.Header.SetUserAgent(random.RandChoice(c.UserAgents))
					if c.AuthHeader != "" {
						request.Header.Set("Authorization", c.AuthHeader)
					}
					if c.Header != nil {
						err := applyHeaders(request, c.Header)
						if err != nil {
							fmt.Fprintln(os.Stderr, err)
						}
					}

					err := c.Client.Do(request, response)
					if err != nil {
						return
					}
					status := response.StatusCode()
					lenght := len(response.Body())
					if c.Wildcard.Active {
						if status == c.Wildcard.Status && wildcard.IsSimilarSize(lenght, c.Wildcard.Lenght, c.Wildcard.Tolerance) {
							return
						}
					}

					if len(c.Exts) > 0 {
						for _, ext := range c.Exts {
							// Reset
							request2.Reset()
							response2.Reset()

							urlWithExt := fullURL + "." + ext

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
							if c.Wildcard.Active && statusCode2 == c.Wildcard.Status && wildcard.IsSimilarSize(lenght2, c.Wildcard.Lenght, c.Wildcard.Tolerance) {

								continue
							}
							if slices.Contains(c.IgnoreCodes, statusCode2) {

								continue
							}

							dirPrefix = "FILE"

							results <- Result{
								Prefix: dirPrefix,
								Size:   lenght2,
								URL:    urlWithExt,
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
					c.VisitedMutex.Lock()

					if c.VisitedDirs[fullURL] {

						c.VisitedMutex.Unlock()

						return
					}
					c.VisitedDirs[fullURL] = true

					c.VisitedMutex.Unlock()

					pathOnly := strings.TrimPrefix(fullURL, baseURL)

					DirDetention, err := detention.Detect(c.Client, baseURL, pathOnly, c.Method)

					if err == nil {

						switch {

						case DirDetention.IsDir:

							dirPrefix = "DIR"

							if c.Recursive {

								c.WG.Add(1)
								c.DirsChan <- fullURL

							}

						case DirDetention.IsFile:
							dirPrefix = "FILE"
						default:
							dirPrefix = "Unknown"
						}
					}

					results <- Result{
						Prefix: dirPrefix,
						Size:   lenght,
						Status: status,
						URL:    fullURL,
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

			c.WG.Done()

		}
	}()

	return results

}
