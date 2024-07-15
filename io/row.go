package io

import (
	"fmt"
	"strings"
)

// Row is a CSV row.
type Row struct {
	Columns
	Headers
	Comma string
}

// Column returns the value of a column by name.
// Need Column values and Headers loaded.
func (row Row) Column(name string) (string, bool) {
	index, ok := row.Headers[name]
	return row.At(index), ok
}

// String converts a row to a string using CSV format with a comma separator.
func (row Row) String() string {
	s := make([]string, len(row.Columns)+1)
	format := strings.TrimSuffix(strings.Join(s, "%s"+row.Comma), row.Comma)
	args := make([]any, len(row.Columns))
	for i, v := range row.Columns {
		args[i] = v
	}
	return fmt.Sprintf(format, args...)
}
