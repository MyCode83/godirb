package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"godirb/internal/core"
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

func Write(results []core.Result, format Format, outputPath string, quiet bool) error {
	writer := io.Writer(os.Stdout)
	var file *os.File
	if strings.TrimSpace(outputPath) != "" {
		var err error
		file, err = os.Create(outputPath)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
	}

	switch format {
	case FormatJSON:
		return writeJSON(writer, results)
	case FormatCSV:
		return writeCSV(writer, results)
	default:
		for _, result := range results {
			if _, err := fmt.Fprintln(writer, FormatTextResult(result, quiet)); err != nil {
				return err
			}
		}
		return nil
	}
}

func writeJSON(writer io.Writer, results []core.Result) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
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
	if err := csvWriter.Write([]string{"prefix", "url", "status", "size", "extra"}); err != nil {
		return err
	}
	for _, result := range results {
		if err := csvWriter.Write([]string{
			result.Prefix,
			result.URL,
			strconv.Itoa(result.Status),
			strconv.Itoa(result.Size),
			result.Extra,
		}); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}
