package calibration

import (
	"github.com/MyCode83/godirb/internal/urlutil"
	"github.com/MyCode83/godirb/pkg/maths"
	"github.com/MyCode83/godirb/pkg/random"

	"fmt"
	"strings"
)

func generateURL(BaseURL, placeholder string) (string, error) {
	randomText := random.RandomString(12)

	if placeholder != "" {
		parts := strings.Split(BaseURL, placeholder)
		if len(parts) != 2 {
			return "", fmt.Errorf("placeholder must appear once")
		}

		return parts[0] + randomText + parts[1], nil
	}

	return urlutil.JoinPath(BaseURL, randomText)
}

func buildSignature(samples []Sample) *Calibration {
	if len(samples) == 0 {
		return &Calibration{}
	}

	status := samples[0].Status
	length := samples[0].Length

	lengths := make([]int, 0, len(samples))
	stable := true

	for _, sample := range samples {
		lengths = append(lengths, sample.Length)

		if sample.Status != status {
			stable = false
		}
	}

	min, max := maths.MinMax(lengths...)
	tolerance := max - min
	if tolerance <= 0 {
		tolerance = 5
	}

	return &Calibration{
		Status:    status,
		Length:    length,
		Tolerance: tolerance,
		Stable:    stable,
		Wildcard:  stable && status != 404,
		Samples:   samples,
	}
}
