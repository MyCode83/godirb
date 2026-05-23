package suggest

import (
	"strings"

	"github.com/agnivade/levenshtein"
)

func absolute(number int) int {
	if number < 0 {
		return -number
	}
	return number
}

func SuggestClosest(maxDistance int, input string, options ...string) string {
	input = strings.TrimSpace(input)
	best := ""
	bestDist := maxDistance + 1
	for _, opt := range options {
		if absolute(len(input)-len(opt)) > maxDistance {
			continue
		}
		distance := levenshtein.ComputeDistance(input, opt)
		if distance < bestDist {
			bestDist = distance
			best = opt
		}
	}
	if bestDist <= maxDistance {
		return best
	}
	return ""
}
