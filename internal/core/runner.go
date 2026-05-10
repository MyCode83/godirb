package core

func (c *Core) Run(baseUrl string) <-chan Result {
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
