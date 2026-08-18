package calibration

var entries []entry

type Calibration struct {
	Status    int
	Length    int
	Tolerance int

	PathLengthAdjusted bool
	AdjustedLength     int
	AdjustedTolerance  int

	Stable   bool
	Wildcard bool

	Samples []Sample
}

type Sample struct {
	URL        string
	Status     int
	Length     int
	PathLength int
}

type Options struct {
	BaseURL     string
	Placeholder string
	Tries       int
	UserAgents  []string
}

type entry struct {
	BaseURL     string
	Placeholder string
	Calibration *Calibration
}
