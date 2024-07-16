package io

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Reader is a CSV reader.
// Extend csv.Reader to add a header map.
type Reader struct {
	*csv.Reader
	header Header
}

// Header is a map of column names to indexes.
type Header map[string]int

// Readeable is a type that can be converted to a string.
type Readeable interface {
	~string | ~[]byte
}

// NewReadeable creates a new Reader from a Readeable.
func NewReadeable[T Readeable](r T) *Reader {
	content := fmt.Sprintf("%s", r)
	reader := strings.NewReader(content)
	return NewReader(reader)
}

// NewReaderWithHeaders creates a new Reader from a Readeable and parses first line as headers.
func NewReaderWithHeaders[T Readeable](r T, comma rune) *Reader {
	reader := NewReadeable(r)
	reader.Comma = comma
	_ = reader.ParseHeader()
	return reader
}

// NewReader creates a new Reader from an io.Reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{Reader: csv.NewReader(r)}
}

// Header returns the header of the CSV as Row.
// The header is a map of column names to indexes.
func (r *Reader) Header() Row {
	columns := make(Columns, len(r.header))
	for name, index := range r.header {
		columns[index] = name
	}
	return Row{Columns: columns, Header: r.header, Comma: string(r.Comma)}
}

// ParseHeader reads a line of the CSV and creates a header map.
func (r *Reader) ParseHeader() error {
	header, err := r.readline()
	if err != nil {
		return errors.Join(ErrReadHeader, err)
	}
	r.header = make(map[string]int)
	for i, name := range header {
		r.header[name] = i
	}
	return nil
}

// ForEach reads each row of the CSV and executes a function.
func (r *Reader) ForEach(exec func(Row)) error {
	for {
		line, err := r.readRow()
		if err != nil {
			return err
		}
		if line.IsEmpty() {
			return nil
		}
		exec(line)
	}
}

// MapRows reads each row of the CSV and maps a function to it.
// Flushes the writer after mapping all rows.
func (r *Reader) MapRows(writer io.Writer, mapper func(*Row)) (err error) {
	csvMapped := csv.NewWriter(writer)
	csvMapped.Comma = r.Comma

	if !r.Header().IsEmpty() {
		_ = csvMapped.Write(r.Header().Columns)
	}

	if err := r.ForEach(func(row Row) {
		mapper(&row)
		_ = csvMapped.Write(row.Columns)
	}); err != nil {
		return err
	}

	csvMapped.Flush()
	if err := csvMapped.Error(); err != nil {
		return errors.Join(ErrMapRows, err)
	}
	return nil
}

func (r *Reader) readline() (Columns, error) {
	line, err := r.Read()
	if err != nil && !errors.Is(err, io.EOF) {
		return line, errors.Join(ErrReadLine, err)
	}
	return line, nil
}

func (r *Reader) readRow() (Row, error) {
	line, err := r.readline()
	if err != nil {
		return Row{}, err
	}
	return Row{Columns: line, Header: r.header, Comma: string(r.Comma)}, nil
}
