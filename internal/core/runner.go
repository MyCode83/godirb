package core
func (c *Core) Run(baseUrl string) {
	switch c.Mode {
	case ModeDir:
		c.RunDir(baseUrl)
	case ModeFuzz:
		c.RunFuzz(baseUrl)
	case ModePort:
		c.RunPorts(baseUrl)	
	}

}