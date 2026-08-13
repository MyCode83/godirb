package calibration

func Get(baseURL, placeholder string) (*Calibration, bool) {
	for _, e := range entries {
		if e.BaseURL == baseURL &&
			e.Placeholder == placeholder {
			return e.Calibration, true
		}
	}

	return nil, false
}
