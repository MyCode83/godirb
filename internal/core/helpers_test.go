package core

import "testing"

func TestShouldUsePlaceholdersOnly(t *testing.T) {
	cases := []struct {
		name  string
		words []string
		want  bool
	}{
		{"empty", nil, false},
		{"few placeholders", []string{"a.%EXT%", "b", "c"}, false},
		{"ninety eight placeholders", append(extWords(98), "plain", "other"), false},
		{"ninety nine placeholders", append(extWords(99), "plain"), true},
	}

	for _, tc := range cases {
		if got := shouldUsePlaceholdersOnly(tc.words); got != tc.want {
			t.Fatalf("%s: shouldUsePlaceholdersOnly() = %t, want %t", tc.name, got, tc.want)
		}
	}
}

func extWords(n int) []string {
	words := make([]string, n)
	for i := range words {
		words[i] = "asset.%EXT%"
	}
	return words
}
