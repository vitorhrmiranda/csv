package csv

import (
	"os"

	csv "github.com/vitorhrmiranda/csv/io"
)

// NewReader creates a new Reader from a file path.
func NewReader(filePath string, comma rune) (reader *csv.Reader, err error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return
	}

	reader = csv.NewReader(file)
	reader.Comma = comma
	if err := reader.ParseHeader(); err != nil {
		return nil, err
	}

	return reader, nil
}
