package calibration

import "testing"

func TestMatchURLUsesPathLengthAdjustedCalibration(t *testing.T) {
	cal := buildSignature([]Sample{
		{
			URL:        "http://example.test/aaaaaaaa",
			Status:     200,
			Length:     108,
			PathLength: len("/aaaaaaaa"),
		},
		{
			URL:        "http://example.test/bbbbbbbbbbbb",
			Status:     200,
			Length:     112,
			PathLength: len("/bbbbbbbbbbbb"),
		},
		{
			URL:        "http://example.test/cccccccccccccccccccc",
			Status:     200,
			Length:     120,
			PathLength: len("/cccccccccccccccccccc"),
		},
		{
			URL:        "http://example.test/dddddddddddddddddddddddddddddddd",
			Status:     200,
			Length:     132,
			PathLength: len("/dddddddddddddddddddddddddddddddd"),
		},
	})

	if !cal.PathLengthAdjusted {
		t.Fatal("PathLengthAdjusted = false, want true")
	}

	if !cal.MatchURL(200, 104, "http://example.test/css") {
		t.Fatal("MatchURL() = false for shorter reflected wildcard path, want true")
	}

	if !cal.MatchURL(200, 118, "http://example.test/administrator-long") {
		t.Fatal("MatchURL() = false for longer reflected wildcard path, want true")
	}

	if cal.MatchURL(200, 180, "http://example.test/real-page") {
		t.Fatal("MatchURL() = true for different body length, want false")
	}
}

func TestMatchURLUsesDecodedPathLengthForUnicodePaths(t *testing.T) {
	cal := buildSignature([]Sample{
		{
			URL:        "http://example.test/aaaaaaaa",
			Status:     200,
			Length:     108,
			PathLength: len("/aaaaaaaa"),
		},
		{
			URL:        "http://example.test/bbbbbbbbbbbb",
			Status:     200,
			Length:     112,
			PathLength: len("/bbbbbbbbbbbb"),
		},
		{
			URL:        "http://example.test/cccccccccccccccccccc",
			Status:     200,
			Length:     120,
			PathLength: len("/cccccccccccccccccccc"),
		},
		{
			URL:        "http://example.test/dddddddddddddddddddddddddddddddd",
			Status:     200,
			Length:     132,
			PathLength: len("/dddddddddddddddddddddddddddddddd"),
		},
	})

	path := "/除投票"
	if !cal.MatchURL(200, 100+len(path), "http://example.test"+path) {
		t.Fatal("MatchURL() = false for reflected unicode path, want true")
	}
}

func TestMatchURLFallsBackToRawLengthCalibration(t *testing.T) {
	cal := buildSignature([]Sample{
		{
			URL:        "http://example.test/aaaaaaaa",
			Status:     200,
			Length:     100,
			PathLength: len("/aaaaaaaa"),
		},
		{
			URL:        "http://example.test/bbbbbbbbbbbb",
			Status:     200,
			Length:     100,
			PathLength: len("/bbbbbbbbbbbb"),
		},
		{
			URL:        "http://example.test/cccccccccccccccccccc",
			Status:     200,
			Length:     100,
			PathLength: len("/cccccccccccccccccccc"),
		},
		{
			URL:        "http://example.test/dddddddddddddddddddddddddddddddd",
			Status:     200,
			Length:     100,
			PathLength: len("/dddddddddddddddddddddddddddddddd"),
		},
	})

	if cal.PathLengthAdjusted {
		t.Fatal("PathLengthAdjusted = true for fixed-size wildcard, want false")
	}

	if !cal.MatchURL(200, 100, "http://example.test/any-length") {
		t.Fatal("MatchURL() = false for fixed-size wildcard, want true")
	}
}
