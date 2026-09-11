package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/MyCode83/godirb/internal/core"
	// "github.com/MyCode83/godirb/internal/ui"
)

type Format int

const (
	FormatQuiet Format = iota
	FormatHuman
	FormatJSON
	FormatCSV
)

func FromFlags(jsonOutput, csvOutput, quietOutput bool) Format {
	switch {
	case jsonOutput:
		return FormatJSON
	case csvOutput:
		return FormatCSV
	case quietOutput:
		return FormatQuiet
	default:
		return FormatHuman
	}
}

type Stream struct {
	mu        sync.Mutex
	writer    io.Writer
	file      *os.File
	csvWriter *csv.Writer
	encoder   *json.Encoder
	format    Format
	closed    bool
}

func NewStream(format Format, outputPath string, noColor bool) (*Stream, error) {
	writer := io.Writer(os.Stdout)
	var file *os.File
	if strings.TrimSpace(outputPath) != "" {
		var err error
		file, err = os.Create(outputPath)
		if err != nil {
			return nil, err
		}

		writer = file

		if format == FormatHuman {
			writer = ansiStripWriter{
				w: file,
			}
		}
	}

	if noColor {
		writer = ansiStripWriter{
			w: writer,
		}
	}

	stream := &Stream{
		writer: writer,
		file:   file,
		format: format,
	}

	switch format {
	case FormatJSON:
		stream.encoder = json.NewEncoder(writer)
	case FormatCSV:
		csvWriter, err := NewCSVWriter(writer)

		if err != nil {
			stream.closeFile()
			return nil, err
		}

		stream.csvWriter = csvWriter
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
		return WriteCSV(s.csvWriter, result)
	case FormatQuiet:
		_, err := fmt.Fprintf(s.writer, "%d %s %d\n", result.Status, result.URL, result.Size)
		return err
	default:
		return WriteHuman(s.writer, result)
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

func tagName(tag string) string {
	name, _, _ := strings.Cut(tag, ",")
	return name
}
