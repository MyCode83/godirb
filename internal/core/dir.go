package core

import (
	"fmt"
	"time"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/detention"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/wildcard"

	"github.com/MyCode83/godirb/pkg/random"
	"os"
	"slices"
	"strings"
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

					defer c.WG.Done()

					defer func() { <-c.Limiter }()
					select {
					case <-c.Ctx.Done():
						debug.Printf("dir worker canceled word=%q dir=%q", word, dir)
						return
					default:

					}
					fullURL := fmt.Sprintf("%s/%s", dir, word)
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
						debug.Error("dir", err)
						return
					}
					debug.Printf("dir response status=%d body=%d", response.StatusCode, response.Lenght)
					status := response.StatusCode
					lenght := response.Lenght
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
							urlWithExt := fullURL + "." + ext
							request.URL = urlWithExt
							request.Method = c.nextRequestMethod()
							request.UserAgent = random.RandChoice(c.UserAgents)
							response2, err2 := c.Client.Do(&request)

							if err2 != nil {
								debug.Error("dir-ext", err2)
								continue
							}
							debug.Printf("dir-ext response status=%d body=%d", response2.StatusCode, response2.Lenght)
							statusCode2 := response2.StatusCode
							lenght2 := response2.Lenght
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
					DirDetention, err := detention.Detect(c.Client, transport.RequestOptions{
						URL:        fullURL,
						Method:     transport.MethodHEAD,
						MethodMode: transport.MethodModeFixed,
						UserAgent:  random.RandChoice(c.UserAgents),
						Headers:    headers,
					})

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
