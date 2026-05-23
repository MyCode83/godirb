package duration

import (
	"strings"
	"time"
	"unicode"
)

func ParseDuration(input string, defaultTime string) (time.Duration, error) {
	input = strings.ToLower(strings.TrimSpace(input))
	input = strings.ReplaceAll(input, ",", ".")
	input = strings.ReplaceAll(input, " ", "")
	if strings.Contains(input, "mls") {
		input = strings.ReplaceAll(input, "mls", "ms")
	}
	hasUnit := false
	for _, c := range input {
		if unicode.IsLetter(c) {
			hasUnit = true
			break
		}
	}
	if !hasUnit {
		input += defaultTime
	}
	timeout, err := time.ParseDuration(input)
	return timeout, err
}
