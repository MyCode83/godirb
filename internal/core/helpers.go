package core

import (
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
