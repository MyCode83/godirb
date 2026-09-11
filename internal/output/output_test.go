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

func TestFromFlags(t *testing.T) {
	tests := []struct {
		name        string
		jsonOutput  bool
		csvOutput   bool
		quietOutput bool
		want        Format
	}{
		{name: "human by default", want: FormatHuman},
		{name: "quiet", quietOutput: true, want: FormatQuiet},
		{name: "json", jsonOutput: true, want: FormatJSON},
		{name: "csv", csvOutput: true, want: FormatCSV},
		{name: "json wins over csv and quiet", jsonOutput: true, csvOutput: true, quietOutput: true, want: FormatJSON},
		{name: "csv wins over quiet", csvOutput: true, quietOutput: true, want: FormatCSV},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromFlags(tt.jsonOutput, tt.csvOutput, tt.quietOutput)
			if got != tt.want {
				t.Fatalf("FromFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamHumanWritesPlainResultToFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.txt")

	stream := newTestStream(t, FormatHuman, path, false)
	writeTestResult(t, stream, testResult())
	closeTestStream(t, stream)

	want := "DIR      200       123 B  http://example.test/admin\n"
	if got := readFile(t, path); got != want {
		t.Fatalf("human stream = %q, want %q", got, want)
	}
}

func TestStreamHumanIncludesError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.txt")
	result := testResult()
	result.Status = 500
	result.Error = "request failed"

	stream := newTestStream(t, FormatHuman, path, false)
	writeTestResult(t, stream, result)
	closeTestStream(t, stream)

	want := "DIR      500       123 B  http://example.test/admin  request failed\n"
	if got := readFile(t, path); got != want {
		t.Fatalf("human stream with error = %q, want %q", got, want)
	}
}

func TestStreamQuietWritesMinimalResult(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.txt")

	stream := newTestStream(t, FormatQuiet, path, false)
	writeTestResult(t, stream, testResult())
	closeTestStream(t, stream)

	want := "200 http://example.test/admin 123\n"
	if got := readFile(t, path); got != want {
		t.Fatalf("quiet stream = %q, want %q", got, want)
	}
}

func TestStreamJSONWritesLineBeforeClose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.jsonl")
	result := testResult()

	stream := newTestStream(t, FormatJSON, path, false)
	defer closeTestStream(t, stream)

	writeTestResult(t, stream, result)

	var got core.Result
	data := readFile(t, path)
	if err := json.Unmarshal(bytes.TrimSpace([]byte(data)), &got); err != nil {
		t.Fatalf("streamed JSON = %q, unmarshal error = %v", data, err)
	}
	if got != result {
		t.Fatalf("streamed JSON = %+v, want %+v", got, result)
	}
}

func TestStreamCSVWritesHeaderBeforeClose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.csv")

	stream := newTestStream(t, FormatCSV, path, false)
	defer closeTestStream(t, stream)

	got := parseCSV(t, readFile(t, path))
	want := [][]string{csvHeader()}
	if !equalCSV(got, want) {
		t.Fatalf("CSV header = %#v, want %#v", got, want)
	}
}

func TestStreamCSVWritesRowsBeforeClose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.csv")
	result := testResult()

	stream := newTestStream(t, FormatCSV, path, false)
	defer closeTestStream(t, stream)

	writeTestResult(t, stream, result)

	got := parseCSV(t, readFile(t, path))
	want := [][]string{
		csvHeader(),
		csvRecord(result),
	}
	if !equalCSV(got, want) {
		t.Fatalf("CSV rows = %#v, want %#v", got, want)
	}
}

func TestStreamWriteAfterCloseFails(t *testing.T) {
	stream := newTestStream(t, FormatHuman, "", true)
	closeTestStream(t, stream)

	if err := stream.Write(testResult()); err == nil {
		t.Fatal("Stream.Write() after Close() succeeded, want error")
	}
}

func newTestStream(t *testing.T, format Format, path string, noColor bool) *Stream {
	t.Helper()

	stream, err := NewStream(format, path, noColor)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}
	return stream
}

func writeTestResult(t *testing.T, stream *Stream, result core.Result) {
	t.Helper()

	if err := stream.Write(result); err != nil {
		t.Fatalf("Stream.Write() error = %v", err)
	}
}

func closeTestStream(t *testing.T, stream *Stream) {
	t.Helper()

	if err := stream.Close(); err != nil {
		t.Fatalf("Stream.Close() error = %v", err)
	}
}

func testResult() core.Result {
	return core.Result{
		Kind:   "DIR",
		URL:    "http://example.test/admin",
		Status: 200,
		Size:   123,
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	return string(data)
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
