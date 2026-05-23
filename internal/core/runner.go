package core

import "github.com/MyCode83/godirb/internal/debug"

func (c *Core) Run(baseUrl string) <-chan Result {
	debug.Printf("core run mode=%d base_url=%q debug=%t words=%d", c.Mode, baseUrl, c.Debug, len(c.WL))
	switch c.Mode {
	case ModeDir:
		return c.RunDir(baseUrl)
	case ModeFuzz:
		return c.RunFuzz(baseUrl)
	case ModePort:
		return c.RunPorts(baseUrl)
	}

	results := make(chan Result)
	close(results)
	return results
}
