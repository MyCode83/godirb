package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"github.com/MyCode83/godirb/internal/core"
)

func WriteCSV(w *csv.Writer, result core.Result) error {
	if err := w.Write(csvRecord(result)); err != nil {
		return err
	}

	w.Flush()
	return w.Error()
}

func NewCSVWriter(w io.Writer) (*csv.Writer, error) {
	csvWriter := csv.NewWriter(w)
	if err := csvWriter.Write(csvHeader()); err != nil {
		return nil, err
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, err
	}

	return csvWriter, nil
}

func csvHeader() []string {
	return csvValues(core.Result{}, func(field reflect.StructField, _ reflect.Value) string {
		if name := csvFieldName(field); name != "" {
			return name
		}
		return field.Name
	})
}

func csvRecord(result core.Result) []string {
	return csvValues(result, func(_ reflect.StructField, value reflect.Value) string {
		return fmt.Sprint(value.Interface())
	})
}

func csvValues(result core.Result, valueFor func(reflect.StructField, reflect.Value) string) []string {
	resultType := reflect.TypeOf(result)
	resultValue := reflect.ValueOf(result)
	values := make([]string, 0, resultType.NumField())

	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		if !field.IsExported() || csvFieldName(field) == "-" {
			continue
		}
		values = append(values, valueFor(field, resultValue.Field(i)))
	}

	return values
}

func csvFieldName(field reflect.StructField) string {
	if name := tagName(field.Tag.Get("csv")); name != "" {
		return name
	}
	return tagName(field.Tag.Get("json"))
}
