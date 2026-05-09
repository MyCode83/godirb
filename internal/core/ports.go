package core

import (
	"fmt"

	// "godirb/internal/assemble"
	"godirb/internal/tui"
	"godirb/pkg/random"
	"slices"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)
func looksLikeService(err error) bool {
    s := err.Error()
    return strings.Contains(s, "tls") ||
           strings.Contains(s, "EOF") ||
           strings.Contains(s, "reset")
}

func (c *Core) RunPorts(baseUrl string) {
		result := tui.Result{}
		for _, word := range c.WL {

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
			err := c.Client.DoTimeout(request, response, c.Timeout)
			if c.Ctx.Err() != nil {
				return
			}
			status := response.StatusCode()
			if err != nil {
				if looksLikeService(err) {
					result.Prefix = "?"
					result.URL = fullURL
					result.Status = status
					result.Extra = fmt.Sprintf("(error: %v)", err)
					tui.Print(result, c.Quiet)
					result.Extra = ""
				}
				return
			}
			
			lenght := len(response.Body())

			
			
			request.Reset()
			response.Reset()
			if slices.Contains(c.IgnoreCodes, status) {
				return
			}
			result.Prefix = prefix
			result.Size = lenght
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