package detention

import (
	"net/url"
	"strings"
)

func urlPath(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	return u.Path
}

func removeFinalSlash(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return strings.TrimRight(raw, "/")
	}

	if u.Path != "/" {
		u.Path = strings.TrimRight(raw, "/")
	}
	
	return u.String()
}

func addFinalSlash(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return strings.TrimRight(raw, "")
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}

	return u.String()
}