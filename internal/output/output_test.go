package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MyCode83/godirb/internal/core"
)

func TestFormatTextResult(t *testing.T) {
	result := core.Result{
		Kind:   "DIR",
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
	results := []core.Result{
		{
			Kind:   "DIR",
			URL:    "http://example.test/admin",
			Status: 200,
			Size:   123,
		},
		{
			Kind:   "DIR",
			URL:    "http://example.test/login",
			Status: 301,
			Size:   456,
		},
	}

	if err := writeJSON(&buf, results); err != nil {
		t.Fatalf("writeJSON() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("writeJSON() wrote %d lines, want 2: %q", len(lines), buf.String())
	}
	for i, line := range lines {
		var got core.Result
		if err := json.Unmarshal([]byte(line), &got); err != nil {
			t.Fatalf("writeJSON() line %d is not JSON: %v", i, err)
		}
		if got != results[i] {
			t.Fatalf("writeJSON() line %d = %+v, want %+v", i, got, results[i])
		}
	}
}

func TestStreamJSONWritesLineBeforeClose(t *testing.T) {
	result := core.Result{
		Kind:   "DIR",
		URL:    "http://example.test/admin",
		Status: 200,
		Size:   123,
	}

	path := filepath.Join(t.TempDir(), "results.jsonl")
	stream, err := NewStream(FormatJSON, path, false)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}
	defer stream.Close()

	if err := stream.Write(result); err != nil {
		t.Fatalf("Stream.Write() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var got core.Result
	if err := json.Unmarshal(bytes.TrimSpace(data), &got); err != nil {
		t.Fatalf("streamed JSON = %q, unmarshal error = %v", data, err)
	}
	if got != result {
		t.Fatalf("streamed JSON = %+v, want %+v", got, result)
	}
}

func TestWriteCSV(t *testing.T) {
	var buf bytes.Buffer
	results := []core.Result{{
		Kind:   "DIR",
		URL:    "http://example.test/admin",
		Status: 200,
		Size:   123,
	}}

	if err := writeCSV(&buf, results); err != nil {
		t.Fatalf("writeCSV() error = %v", err)
	}

	records := parseCSV(t, buf.String())
	want := [][]string{
		csvHeader(),
		csvRecord(results[0]),
	}
	if !equalCSV(records, want) {
		t.Fatalf("writeCSV() = %#v, want %#v", records, want)
	}
}

func TestStreamCSVWritesHeaderAndRowsBeforeClose(t *testing.T) {
	result := core.Result{
		Kind:   "DIR",
		URL:    "http://example.test/admin",
		Status: 200,
		Size:   123,
	}

	path := filepath.Join(t.TempDir(), "results.csv")
	stream, err := NewStream(FormatCSV, path, false)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}
	defer stream.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	records := parseCSV(t, string(data))
	if !equalCSV(records, [][]string{csvHeader()}) {
		t.Fatalf("CSV header before close = %#v, want %#v", records, [][]string{csvHeader()})
	}

	if err := stream.Write(result); err != nil {
		t.Fatalf("Stream.Write() error = %v", err)
	}

	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	records = parseCSV(t, string(data))
	want := [][]string{
		csvHeader(),
		csvRecord(result),
	}
	if !equalCSV(records, want) {
		t.Fatalf("CSV before close = %#v, want %#v", records, want)
	}
}

func parseCSV(t *testing.T, data string) [][]string {
	t.Helper()

	records, err := csv.NewReader(strings.NewReader(data)).ReadAll()
	if err != nil {
		t.Fatalf("parse CSV %q: %v", data, err)
	}
	return records
}

func equalCSV(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
