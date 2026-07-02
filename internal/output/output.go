package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/MyCode83/godirb/internal/core"
)

type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatCSV
)

func FromFlags(jsonOutput, csvOutput bool) Format {
	switch {
	case jsonOutput:
		return FormatJSON
	case csvOutput:
		return FormatCSV
	default:
		return FormatText
	}
}

type Stream struct {
	mu        sync.Mutex
	writer    io.Writer
	file      *os.File
	csvWriter *csv.Writer
	encoder   *json.Encoder
	format    Format
	quiet     bool
	closed    bool
}

func NewStream(format Format, outputPath string, quiet bool) (*Stream, error) {
	writer := io.Writer(os.Stdout)
	var file *os.File
	if strings.TrimSpace(outputPath) != "" {
		var err error
		file, err = os.Create(outputPath)
		if err != nil {
			return nil, err
		}
		writer = file
	}

	stream := &Stream{
		writer: writer,
		file:   file,
		format: format,
		quiet:  quiet,
	}

	switch format {
	case FormatJSON:
		stream.encoder = json.NewEncoder(writer)
	case FormatCSV:
		stream.csvWriter = csv.NewWriter(writer)
		if err := stream.csvWriter.Write(csvHeader()); err != nil {
			stream.closeFile()
			return nil, err
		}
		stream.csvWriter.Flush()
		if err := stream.csvWriter.Error(); err != nil {
			stream.closeFile()
			return nil, err
		}
	}

	return stream, nil
}

func (s *Stream) Write(result core.Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("output stream is closed")
	}

	switch s.format {
	case FormatJSON:
		return s.encoder.Encode(result)
	case FormatCSV:
		if err := s.csvWriter.Write(csvRecord(result)); err != nil {
			return err
		}
		s.csvWriter.Flush()
		return s.csvWriter.Error()
	default:
		_, err := fmt.Fprintln(s.writer, FormatTextResult(result, s.quiet))
		return err
	}
}

func (s *Stream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return nil
	}
	s.closed = true

	var err error
	if s.csvWriter != nil {
		s.csvWriter.Flush()
		err = s.csvWriter.Error()
	}
	if closeErr := s.closeFile(); err == nil {
		err = closeErr
	}
	return err
}

func (s *Stream) closeFile() error {
	if s.file == nil {
		return nil
	}
	err := s.file.Close()
	s.file = nil
	return err
}

func Write(results []core.Result, format Format, outputPath string, quiet bool) error {
	stream, err := NewStream(format, outputPath, quiet)
	if err != nil {
		return err
	}
	for _, result := range results {
		if err := stream.Write(result); err != nil {
			_ = stream.Close()
			return err
		}
	}
	return stream.Close()
}

func writeJSON(writer io.Writer, results []core.Result) error {
	encoder := json.NewEncoder(writer)
	for _, result := range results {
		if err := encoder.Encode(result); err != nil {
			return err
		}
	}
	return nil
}

func FormatTextResult(result core.Result, quiet bool) string {
	if quiet {
		return fmt.Sprintf("%d %s %d", result.Status, result.URL, result.Size)
	}
	if strings.TrimSpace(result.Extra) != "" {
		return fmt.Sprintf("[%s] %s ---> %d %s | %d", result.Prefix, result.URL, result.Status, result.Extra, result.Size)
	}
	return fmt.Sprintf("[%s] %s ---> %d | %d", result.Prefix, result.URL, result.Status, result.Size)
}

func writeCSV(writer io.Writer, results []core.Result) error {
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write(csvHeader()); err != nil {
		return err
	}
	for _, result := range results {
		if err := csvWriter.Write(csvRecord(result)); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

func csvHeader() []string {
	return []string{"prefix", "url", "status", "size", "extra"}
}

func csvRecord(result core.Result) []string {
	return []string{
		result.Prefix,
		result.URL,
		strconv.Itoa(result.Status),
		strconv.Itoa(result.Size),
		result.Extra,
	}
}
