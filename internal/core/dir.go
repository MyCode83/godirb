package core

import (
	"fmt"
	"time"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/detention"
	"github.com/MyCode83/godirb/internal/wildcard"

	"github.com/MyCode83/godirb/pkg/random"
	"os"
	"slices"
	"strings"

	"github.com/valyala/fasthttp"
)

func (c *Core) RunDir(baseURL string) <-chan Result {
	results := make(chan Result)
	debug.Printf("dir run start base_url=%q recursive=%t words=%d exts=%v", baseURL, c.Recursive, len(c.WL), c.Exts)

	go func() {
		defer close(results)

		if c.Wildcard == nil {
			debug.Printf("dir run stopped: wildcard is nil")
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
			debug.Printf("dir queue item=%q", dir)

			// Wordlist loop
			for _, word := range c.WL {
				select {
				case <-c.Ctx.Done():
					debug.Printf("dir run canceled before scheduling word=%q dir=%q", word, dir)
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
						debug.Printf("dir worker canceled word=%q dir=%q", word, dir)
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

					debug.Request("dir", request)
					err := c.Client.Do(request, response)
					if err != nil {
						debug.Error("dir", err)
						return
					}
					debug.Response("dir", response)
					status := response.StatusCode()
					lenght := len(response.Body())
					if c.Wildcard.Active {
						if status == c.Wildcard.Status && wildcard.IsSimilarSize(lenght, c.Wildcard.Lenght, c.Wildcard.Tolerance) {
							debug.Printf("dir filtered wildcard url=%s status=%d length=%d wildcard_status=%d wildcard_length=%d tolerance=%d",
								fullURL, status, lenght, c.Wildcard.Status, c.Wildcard.Lenght, c.Wildcard.Tolerance)
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
							debug.Request("dir-ext", request2)
							err2 := c.Client.Do(request2, response2)

							if err2 != nil {
								debug.Error("dir-ext", err2)
								fasthttp.ReleaseRequest(request2)
								fasthttp.ReleaseResponse(response2)
								continue
							}
							debug.Response("dir-ext", response2)
							statusCode2 := response2.StatusCode()
							lenght2 := len(response2.Body())
							if c.Wildcard.Active && statusCode2 == c.Wildcard.Status && wildcard.IsSimilarSize(lenght2, c.Wildcard.Lenght, c.Wildcard.Tolerance) {
								debug.Printf("dir-ext filtered wildcard url=%s status=%d length=%d", urlWithExt, statusCode2, lenght2)

								continue
							}
							if slices.Contains(c.IgnoreCodes, statusCode2) {
								debug.Printf("dir-ext ignored url=%s status=%d", urlWithExt, statusCode2)

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
								debug.Printf("dir-ext delay=%s url=%s", c.Delay, urlWithExt)
								select {

								case <-time.After(c.Delay):

								case <-c.Ctx.Done():
									debug.Printf("dir-ext canceled during delay url=%s", urlWithExt)
									return
								}

							}

						}
					}
					if slices.Contains(c.IgnoreCodes, status) {
						debug.Printf("dir ignored url=%s status=%d", fullURL, status)
						return

					}
					c.VisitedMutex.Lock()

					if c.VisitedDirs[fullURL] {
						debug.Printf("dir skipped visited url=%s", fullURL)

						c.VisitedMutex.Unlock()

						return
					}
					c.VisitedDirs[fullURL] = true

					c.VisitedMutex.Unlock()

					pathOnly := strings.TrimPrefix(fullURL, baseURL)

					debug.Printf("dir detention url=%s path=%s", fullURL, pathOnly)
					DirDetention, err := detention.Detect(c.Client, baseURL, pathOnly, c.Method)

					if err == nil {

						switch {

						case DirDetention.IsDir:

							dirPrefix = "DIR"

							if c.Recursive {

								c.WG.Add(1)
								c.DirsChan <- fullURL
								debug.Printf("dir recursive enqueue url=%s", fullURL)

							}

						case DirDetention.IsFile:
							dirPrefix = "FILE"
						default:
							dirPrefix = "Unknown"
						}
						debug.Printf("dir detention classification url=%s prefix=%s", fullURL, dirPrefix)
					} else {
						debug.Error("dir detention", err)
					}

					results <- Result{
						Prefix: dirPrefix,
						Size:   lenght,
						Status: status,
						URL:    fullURL,
					}

					if c.Delay > 0 {
						debug.Printf("dir delay=%s url=%s", c.Delay, fullURL)
						select {
						case <-time.After(c.Delay):
						case <-c.Ctx.Done():
							debug.Printf("dir canceled during delay url=%s", fullURL)
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
