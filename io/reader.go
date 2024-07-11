package csv

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
	header Headers
}

// Headers is a map of column names to indexes.
type Headers map[string]int

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

func (r *Reader) Headers() Headers {
	return r.header
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

func (r *Reader) readline() (Columns, error) {
	line, err := r.Read()
	if err != nil && !errors.Is(err, io.EOF) {
		return line, err
	}
	return line, nil
}

func (r *Reader) readRow() (Row, error) {
	line, err := r.readline()
	if err != nil {
		return Row{}, err
	}
	return Row{Columns: line, Headers: r.header, Comma: string(r.Comma)}, nil
}
