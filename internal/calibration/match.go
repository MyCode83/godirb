package calibration

func (c *Calibration) Match(status int, length int) bool {
	if c == nil || !c.Stable {
		return false
	}

	if status != c.Status {
		return false
	}

	diff := length - c.Length
	if diff < 0 {
		diff = -diff
	}

	return diff <= c.Tolerance
}
