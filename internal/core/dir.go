package core

import (
	"github.com/MyCode83/godirb/internal/calibration"
	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/detection"
	"github.com/MyCode83/godirb/internal/transport"
	"github.com/MyCode83/godirb/internal/urlutil"

	"slices"
	"strings"

	"github.com/MyCode83/godirb/pkg/random"
)

const ExtPlaceholder = "GODIRB_EXT_PLACEHOLDER"

func (c *Core) RunDir(baseURL string) <-chan Result {
	results := make(chan Result)
	debug.Printf("dir run start base_url=%q recursive=%t depth=%d words=%d exts=%v", baseURL, c.Recursive, c.Depth, len(c.WL), c.Exts)

	go func() {
		defer close(results)
		c.WG.Add(1)
		c.DirsChan <- DirTask{URL: baseURL}

		go func() {

			c.WG.Wait()
			close(c.DirsChan)

		}()

		// Dirs loop
	dirLoop:
		for task := range c.DirsChan {
			dir := task.URL
			debug.Printf("dir queue item=%q", dir)

			cal, ok := calibration.Get(dir, "")
			if !ok {
				err := calibration.Build(c.Client, calibration.Options{
					BaseURL:     dir,
					Placeholder: "",
					Tries:       3,
					UserAgents:  c.UserAgents,
				})
				if err != nil {
					debug.Error("dir calibration build", err)
					c.WG.Done()
					continue
				}

				cal, ok = calibration.Get(dir, "")
				if !ok {
					debug.Printf("dir calibration missing after build base_url=%q", dir)
					c.WG.Done()
					continue
				}
			}

			if len(c.Exts) > 0 {
				for _, ext := range c.Exts {
					base := urlutil.JoinPath(dir, ExtPlaceholder)

					templateURL := urlutil.AddExtension(base, ext)

					if _, ok := calibration.Get(templateURL, ExtPlaceholder); ok {
						continue
					}

					if err := calibration.Build(c.Client, calibration.Options{
						BaseURL:     templateURL,
						Placeholder: ExtPlaceholder,
						Tries:       3,
						UserAgents:  c.UserAgents,
					}); err != nil {
						debug.Error("extension calibration build", err)
						continue
					}
				}
			}

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
				go func(word string, cal *calibration.Calibration) {
					dirPrefix := ""

					defer c.WG.Done()

					defer func() { <-c.Limiter }()
					select {
					case <-c.Ctx.Done():
						debug.Printf("dir worker canceled word=%q dir=%q", word, dir)
						return
					default:

					}
					fullURL := urlutil.JoinPath(dir, word)
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
					if !c.applyDelay("dir", fullURL) {
						return
					}

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
							func(ext string) string {
								return urlutil.AddExtension(fullURL, ext)
							},
							func(ext string) string {
								return urlutil.AddExtension(urlutil.JoinPath(dir, ExtPlaceholder), ext)
							},
						)

						if !ok {
							return
						}
					}

					if c.hasSignature(response) {
						return
					}

					if cal.MatchURL(status, lenght, fullURL) {
						debug.Printf(
							"dir filtered calibration url=%s status=%d length=%d calibration_status=%d calibration_length=%d tolerance=%d",
							fullURL,
							status,
							lenght,
							cal.Status,
							cal.Length,
							cal.Tolerance,
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
					if !c.applyDelay("dir-detection", fullURL) {
						return
					}

					if err == nil {

						switch {

						case dirDetection.IsDir:

							dirPrefix = "DIR"

							if c.Recursive && (c.Depth < 0 || task.Depth < c.Depth) {
								err := calibration.Build(c.Client, calibration.Options{
									BaseURL:     fullURL,
									Placeholder: "",
									Tries:       3,
									UserAgents:  c.UserAgents,
								})
								if err != nil {
									debug.Error("recursive calibration build", err)
									break
								}

								c.WG.Add(1)
								c.DirsChan <- DirTask{URL: fullURL, Depth: task.Depth + 1}
								debug.Printf("dir recursive enqueue url=%s depth=%d", fullURL, task.Depth+1)

							}

						case dirDetection.IsFile:
							dirPrefix = "FILE"
						default:
							dirPrefix = "UNKNOWN"
						}
						debug.Printf("dir detection classification url=%s prefix=%s", fullURL, dirPrefix)
					} else {
						debug.Error("dir detection", err)
					}

					results <- Result{
						Kind:   dirPrefix,
						Size:   lenght,
						Status: status,

						Method:        response.Method.String(),
						ContentType:   response.ContentType,
						ContentLength: response.ContentLenght,
						Location:      response.Location,
						Duration:      response.Duration.String(),
						Title:         response.Title,

						URL: fullURL,
					}

				}(word, cal)
			}

			c.WG.Done()

		}
	}()

	return results

}
