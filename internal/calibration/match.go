package calibration

func (c *Calibration) Match(status int, length int) bool {
	return c.match(status, length, c.Length, c.Tolerance)
}

func (c *Calibration) MatchURL(status int, length int, rawURL string) bool {
	if c == nil || !c.PathLengthAdjusted {
		return c.Match(status, length)
	}

	return c.match(
		status,
		length-decodedURLPathLength(rawURL),
		c.AdjustedLength,
		c.AdjustedTolerance,
	)
}

func (c *Calibration) match(status int, length int, expectedLength int, tolerance int) bool {
	if c == nil || !c.Stable {
		return false
	}

	if status != c.Status {
		return false
	}

	diff := length - expectedLength
	if diff < 0 {
		diff = -diff
	}

	return diff <= tolerance
}
