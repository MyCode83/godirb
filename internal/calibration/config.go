package calibration

type Calibration struct {
	Status    int
	Length    int
	Tolerance int

	Stable   bool
	Wildcard bool

	Samples []Sample
}

type Sample struct {
	URL    string
	Status int
	Length int
}

type Options struct {
	BaseURL     string
	Placeholder string
	Tries       int
	UserAgents  []string
}
