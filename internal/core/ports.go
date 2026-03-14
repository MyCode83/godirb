package core

import (
	"fmt"

	// "godirb/internal/assemble"
	"godirb/internal/tui"
	"godirb/internal/wordlist"
	"godirb/pkg/random"
	"slices"
	"strings"
	"time"

	"net/http"
	"io"
)
func looksLikeService(err error) bool {
    s := err.Error()
    return strings.Contains(s, "tls") ||
           strings.Contains(s, "EOF") ||
           strings.Contains(s, "reset")
}
func getHostOnly(u string) string {
    parts := strings.SplitN(u, "://", 2)
    if len(parts) == 2 {
        u = parts[1]
    }

    host := strings.SplitN(u, "/", 2)[0]
    host = strings.SplitN(host, ":", 2)[0]

    return host
}
func (c *Core) RunPorts(baseUrl string) {

	client := &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			DisableKeepAlives:     true,
			MaxIdleConns:          0,
			MaxIdleConnsPerHost:   0,
			IdleConnTimeout:       0,
			TLSHandshakeTimeout:   c.Timeout,
			ResponseHeaderTimeout: c.Timeout,
			TLSClientConfig:       c.Client.TLSConfig,
		},
	}

	urlParts := strings.Split(baseUrl, c.Placeholder)

	for _, word := range wordlist.ListSlice {

		select {
		case <-c.Ctx.Done():
			return
		case c.Limiter <- struct{}{}:
		}

		word = strings.TrimLeft(word, "/")
		c.WG.Add(1)

		go func(word string) {

			defer c.WG.Done()
			defer func() { <-c.Limiter }()

			result := tui.Result{}
			methodSwitch := "GET"

			switch c.Method {
			case "GET", "HEAD", "get", "head":
				methodSwitch = strings.ToUpper(c.Method)

			case "SWITCH", "switch", "swich":
				if methodSwitch == "GET" {
					methodSwitch = "HEAD"
				} else {
					methodSwitch = "GET"
				}
			}

			fullURL := urlParts[0] + word + urlParts[1]

			req, err := http.NewRequest(methodSwitch, fullURL, nil)
			if err != nil {
				return
			}

			req.Header.Set("User-Agent", random.RandChoice(c.UserAgents))

			if c.AuthHeader != "" {
				req.Header.Set("Authorization", c.AuthHeader)
			}

			resp, err := client.Do(req)

			if c.Ctx.Err() != nil {
				return
			}

			if err != nil {
				if looksLikeService(err) {
					result.Prefix = "?"
					result.URL = fullURL
					result.Status = 0
					result.Extra = fmt.Sprintf("(error: %v)", err)
					tui.Print(result, c.Quiet)
					result.Extra = ""
				}
				return
			}

			defer resp.Body.Close()

			body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
			length := len(body)
			status := resp.StatusCode

			if slices.Contains(c.IgnoreCodes, status) {
				return
			}

			result.Prefix = prefix
			result.Size = length
			result.Status = status
			result.URL = fullURL

			tui.Print(result, c.Quiet)

			if c.Delay > 0 {
				select {
				case <-time.After(c.Delay):
				case <-c.Ctx.Done():
					return
				}
			}

		}(word)

	}

}