package core

import (
	"strings"

	"github.com/MyCode83/godirb/internal/debug"
	"github.com/MyCode83/godirb/internal/transport"
)

func (c *Core) hasSignature(resp transport.Response) bool {
	if c.Signatures == nil {
		return false
	}

	matches := c.Signatures.MatchDefaultError(&resp)
	if len(matches) == 0 {
		return false
	}

	debug.Printf(
		"known signature matched url=%q status=%d signatures=%v",
		resp.URL,
		resp.StatusCode,
		matches,
	)

	return true
}

func (c *Core) hasLineExt(line string) bool {
	return strings.Contains(line, `%EXT%`)
}

func (c *Core) processExtPlaceholder(
	request transport.RequestOptions,
	results chan<- Result,
	fullURL string,
	prefix string,
	debugName string,
	buildCalibrationURL func(ext string) string,
	) {

	c.processExtensions(
		&request,
		results,
		prefix,
		debugName,
		func(ext string) string {
			return strings.ReplaceAll(
				fullURL,
				"%EXT%",
				strings.TrimPrefix(ext, "."),
			)
		},
		buildCalibrationURL,
	)
}

func shouldUsePlaceholdersOnly(words []string) bool {
	if len(words) <= 0 {
		return false
	}

	placeholders := 0
	for _, word := range words {
		if strings.Contains(word, "%EXT%") {
			placeholders++
		}
	}

	return placeholders >= 99 &&
		float64(placeholders)/float64(len(words)) >= 0.10
}
