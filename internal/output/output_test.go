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
		name       string
		jsonOutput bool
		csvOutput  bool
		want       Format
	}{
		{name: "text by default", want: FormatText},
		{name: "json", jsonOutput: true, want: FormatJSON},
		{name: "csv", csvOutput: true, want: FormatCSV},
		{name: "json wins over csv", jsonOutput: true, csvOutput: true, want: FormatJSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromFlags(tt.jsonOutput, tt.csvOutput)
			if got != tt.want {
				t.Fatalf("FromFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStreamTextWritesResult(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.txt")
	result := testResult()

	stream, err := NewStream(FormatText, path, false)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}

	if err := stream.Write(result); err != nil {
		t.Fatalf("Stream.Write() error = %v", err)
	}
	if err := stream.Close(); err != nil {
		t.Fatalf("Stream.Close() error = %v", err)
	}

	got := readFile(t, path)
	want := "[DIR] http://example.test/admin ---> 200 | 123\n"
	if got != want {
		t.Fatalf("text stream = %q, want %q", got, want)
	}
}

func TestStreamTextQuietWritesMinimalResult(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.txt")
	result := testResult()

	stream, err := NewStream(FormatText, path, true)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}

	if err := stream.Write(result); err != nil {
		t.Fatalf("Stream.Write() error = %v", err)
	}
	if err := stream.Close(); err != nil {
		t.Fatalf("Stream.Close() error = %v", err)
	}

	got := readFile(t, path)
	want := "200 http://example.test/admin 123\n"
	if got != want {
		t.Fatalf("quiet text stream = %q, want %q", got, want)
	}
}

func TestStreamJSONWritesLineBeforeClose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.jsonl")
	result := testResult()

	stream, err := NewStream(FormatJSON, path, false)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}
	defer stream.Close()

	if err := stream.Write(result); err != nil {
		t.Fatalf("Stream.Write() error = %v", err)
	}

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

	stream, err := NewStream(FormatCSV, path, false)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}
	defer stream.Close()

	got := parseCSV(t, readFile(t, path))
	want := [][]string{csvHeader()}
	if !equalCSV(got, want) {
		t.Fatalf("CSV header = %#v, want %#v", got, want)
	}
}

func TestStreamCSVWritesRowsBeforeClose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "results.csv")
	result := testResult()

	stream, err := NewStream(FormatCSV, path, false)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}
	defer stream.Close()

	if err := stream.Write(result); err != nil {
		t.Fatalf("Stream.Write() error = %v", err)
	}

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
	stream, err := NewStream(FormatText, "", false)
	if err != nil {
		t.Fatalf("NewStream() error = %v", err)
	}
	if err := stream.Close(); err != nil {
		t.Fatalf("Stream.Close() error = %v", err)
	}

	if err := stream.Write(testResult()); err == nil {
		t.Fatal("Stream.Write() after Close() succeeded, want error")
	}
}

func TestWriteUsesSelectedFormat(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		quiet  bool
		assert func(t *testing.T, data string)
	}{
		{
			name:   "text",
			format: FormatText,
			assert: func(t *testing.T, data string) {
				t.Helper()
				if want := "[DIR] http://example.test/admin ---> 200 | 123\n"; data != want {
					t.Fatalf("Write() text = %q, want %q", data, want)
				}
			},
		},
		{
			name:   "quiet text",
			format: FormatText,
			quiet:  true,
			assert: func(t *testing.T, data string) {
				t.Helper()
				if want := "200 http://example.test/admin 123\n"; data != want {
					t.Fatalf("Write() quiet text = %q, want %q", data, want)
				}
			},
		},
		{
			name:   "json",
			format: FormatJSON,
			assert: func(t *testing.T, data string) {
				t.Helper()
				var got core.Result
				if err := json.Unmarshal(bytes.TrimSpace([]byte(data)), &got); err != nil {
					t.Fatalf("Write() JSON = %q, unmarshal error = %v", data, err)
				}
				if want := testResult(); got != want {
					t.Fatalf("Write() JSON = %+v, want %+v", got, want)
				}
			},
		},
		{
			name:   "csv",
			format: FormatCSV,
			assert: func(t *testing.T, data string) {
				t.Helper()
				got := parseCSV(t, data)
				want := [][]string{csvHeader(), csvRecord(testResult())}
				if !equalCSV(got, want) {
					t.Fatalf("Write() CSV = %#v, want %#v", got, want)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "results")
			if err := Write([]core.Result{testResult()}, tt.format, path, tt.quiet); err != nil {
				t.Fatalf("Write() error = %v", err)
			}
			tt.assert(t, readFile(t, path))
		})
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
