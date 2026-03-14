package core

import (
	"github.com/valyala/fasthttp"
	"strings"
	"fmt"
)

func applyHeaders(request *fasthttp.Request, headers []string) error{
	for _, header := range headers {
		parts := strings.Split(header, ":")
		if len(parts) != 2 {
			return fmt.Errorf("[!] Error: the headers need to be key: value")
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key != "" || value != "" {
			continue
		}
		request.Header.Set(key, value)
	}
	return nil
}