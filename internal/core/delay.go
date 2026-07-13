package core

import (
	"time"

	"github.com/MyCode83/godirb/internal/debug"
)

func (c *Core) applyDelay(debugName, url string) bool {
	if c.Delay <= 0 {
		return true
	}

	debug.Printf("%s delay=%s url=%s", debugName, c.Delay, url)

	timer := time.NewTimer(c.Delay)
	defer timer.Stop()

	if c.Ctx == nil {
		<-timer.C
		return true
	}

	select {
	case <-timer.C:
		return true
	case <-c.Ctx.Done():
		debug.Printf("%s canceled during delay url=%s", debugName, url)
		return false
	}
}
