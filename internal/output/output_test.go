package output

import (
	"bytes"
	"strings"
	"testing"

	"godirb/internal/core"
)

func TestFormatTextResult(t *testing.T) {
	result := core.Result{
		Prefix: "DIR",
		URL:    "http://example.test/admin",
		Status: 200,
		Size:   123,
	}

	got := FormatTextResult(result, false)
	want := "[DIR] http://example.test/admin ---> 200 | 123"
	if got != want {
		t.Fatalf("FormatTextResult() = %q, want %q", got, want)
	}

	got = FormatTextResult(result, true)
	want = "200 http://example.test/admin 123"
	if got != want {
		t.Fatalf("FormatTextResult(quiet) = %q, want %q", got, want)
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	results := []core.Result{{
		Prefix: "DIR",
		URL:    "http://example.test/admin",
		Status: 200,
		Size:   123,
	}}

	if err := writeJSON(&buf, results); err != nil {
		t.Fatalf("writeJSON() error = %v", err)
	}

	got := buf.String()
	for _, want := range []string{`"prefix": "DIR"`, `"url": "http://example.test/admin"`, `"status": 200`, `"size": 123`} {
		if !strings.Contains(got, want) {
			t.Fatalf("writeJSON() = %q, missing %q", got, want)
		}
	}
}

func TestWriteCSV(t *testing.T) {
	var buf bytes.Buffer
	results := []core.Result{{
		Prefix: "DIR",
		URL:    "http://example.test/admin",
		Status: 200,
		Size:   123,
	}}

	if err := writeCSV(&buf, results); err != nil {
		t.Fatalf("writeCSV() error = %v", err)
	}

	want := "prefix,url,status,size,extra\nDIR,http://example.test/admin,200,123,\n"
	got := strings.ReplaceAll(buf.String(), "\r\n", "\n")
	if got != want {
		t.Fatalf("writeCSV() = %q, want %q", got, want)
	}
}
