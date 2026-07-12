package core

import (
	"time"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/detection"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/urlutil"

	"github.com/MyCode83/godirb/pkg/random"
	"slices"
	"strings"
)

func (c *Core) RunDir(baseURL string) <-chan Result {
	results := make(chan Result)
	debug.Printf("dir run start base_url=%q recursive=%t words=%d exts=%v", baseURL, c.Recursive, len(c.WL), c.Exts)

	go func() {
		defer close(results)
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
					fullURL, err := urlutil.JoinPath(dir, word)
					if err != nil {
						debug.Printf("fullURL error in urlutil.JoinPath(dir, word)")
						return
					}
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

					if len(c.Exts) > 0 {
						ok := c.processExtensions(
							&request,
							results,
							"FILE",
							"dir-ext",
							func(ext string) (string, error) {
								return urlutil.AddExtension(fullURL, ext)
							},
						)

						if !ok {
							return
						}
					}

					if c.Calibration.Match(status, lenght) {
						debug.Printf(
							"dir filtered calibration url=%s status=%d length=%d calibration_status=%d calibration_length=%d tolerance=%d",
							fullURL,
							status,
							lenght,
							c.Calibration.Status,
							c.Calibration.Length,
							c.Calibration.Tolerance,
						)
						return
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

					debug.Printf("dir detection url=%s path=%s", fullURL, pathOnly)
					dirDetection, err := detection.Detect(c.Client, transport.RequestOptions{
						URL:        fullURL,
						Method:     transport.MethodHEAD,
						MethodMode: transport.MethodModeFixed,
						UserAgent:  random.RandChoice(c.UserAgents),
						Headers:    headers,
					})

					if err == nil {

						switch {

						case dirDetection.IsDir:

							dirPrefix = "DIR"

							if c.Recursive {

								c.WG.Add(1)
								c.DirsChan <- fullURL
								debug.Printf("dir recursive enqueue url=%s", fullURL)

							}

						case dirDetection.IsFile:
							dirPrefix = "FILE"
						default:
							dirPrefix = "Unknown"
						}
						debug.Printf("dir detection classification url=%s prefix=%s", fullURL, dirPrefix)
					} else {
						debug.Error("dir detection", err)
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
