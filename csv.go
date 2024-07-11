package csv

import (
	"os"

	csv "github.com/vitorhrmiranda/csv/io"
)

// NewReader creates a new Reader from a file path.
func NewReader(filePath string, comma rune) (*csv.Reader, error) {
	reader := new(csv.Reader)
	file, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return reader, err
	}

	reader = csv.NewReader(file)
	reader.Comma = comma
	if err := reader.ParseHeader(); err != nil {
		return reader, err
	}

	return reader, nil
}
