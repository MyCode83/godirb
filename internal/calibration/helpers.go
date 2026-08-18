package calibration

import (
	"net/url"

	"github.com/MyCode83/godirb/internal/urlutil"
	"github.com/MyCode83/godirb/pkg/maths"
	"github.com/MyCode83/godirb/pkg/random"

	"fmt"
	"strings"
)

func generateURL(BaseURL, placeholder string, randomLength int) (string, error) {
	randomText := random.RandomString(randomLength)

	if placeholder != "" {
		parts := strings.Split(BaseURL, placeholder)
		if len(parts) != 2 {
			return "", fmt.Errorf("placeholder must appear once")
		}

		return parts[0] + randomText + parts[1], nil
	}

	return urlutil.JoinPath(BaseURL, randomText), nil
}

func buildSignature(samples []Sample) *Calibration {
	if len(samples) == 0 {
		return &Calibration{}
	}

	status := samples[0].Status
	length := samples[0].Length

	lengths := make([]int, 0, len(samples))
	adjustedLengths := make([]int, 0, len(samples))
	stable := true

	for _, sample := range samples {
		lengths = append(lengths, sample.Length)
		adjustedLengths = append(adjustedLengths, sample.Length-sample.PathLength)

		if sample.Status != status {
			stable = false
		}
	}

	min, max := maths.MinMax(lengths...)
	tolerance := max - min
	if tolerance <= 0 {
		tolerance = 5
	}

	adjustedMin, adjustedMax := maths.MinMax(adjustedLengths...)
	adjustedTolerance := adjustedMax - adjustedMin
	if adjustedTolerance <= 0 {
		adjustedTolerance = 5
	}

	pathLengthAdjusted := stable &&
		tolerance > 5 &&
		adjustedTolerance <= 5 &&
		adjustedTolerance < tolerance

	return &Calibration{
		Status:    status,
		Length:    length,
		Tolerance: tolerance,

		PathLengthAdjusted: pathLengthAdjusted,
		AdjustedLength:     adjustedLengths[0],
		AdjustedTolerance:  adjustedTolerance,

		Stable:   stable,
		Wildcard: stable && status != 404,
		Samples:  samples,
	}
}

func decodedURLPathLength(rawURL string) int {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return 0
	}

	return len(parsed.Path)
}
