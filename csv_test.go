package csv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vitorhrmiranda/csv"
	"github.com/vitorhrmiranda/csv/io"
)

func TestNewReader(t *testing.T) {
	fileName := CreateExampleFile(t, []byte("a,b,c\n1,2,3\n4,5,6\n7,8,9"))

	reader, err := csv.NewReader(fileName, ',')
	require.NoError(t, err)

	reader.ForEach(func(row io.Row) {
		t.Run("column", func(t *testing.T) {
			value, exists := row.Column("a")
			require.True(t, exists)
			require.NotEmpty(t, value)
		})
		t.Run("set", func(t *testing.T) {
			row.Set(0, "10")
			value, exists := row.Column("a")
			require.True(t, exists)
			require.Equal(t, "10", value)
		})
	})
}

func TestNewReader_Error(t *testing.T) {
	t.Run("read file", func(t *testing.T) {
		reader, err := csv.NewReader("", ',')
		assert.Error(t, err)
		assert.Nil(t, reader)
	})
	t.Run("parse header", func(t *testing.T) {
		fileName := CreateExampleFile(t, []byte(`§a∑""b,c`))
		reader, err := csv.NewReader(fileName, ',')
		assert.ErrorIs(t, err, io.ErrReadHeader)
		assert.Nil(t, reader)
	})
}

func CreateExampleFile(t *testing.T, csvContent []byte) string {
	tempDir := t.TempDir()
	fileName := tempDir + "/test.csv"

	file, err := os.Create(fileName)
	require.NoError(t, err)

	t.Cleanup(func() {
		file.Close()
		os.Remove(fileName)
	})

	_, err = file.Write(csvContent)
	require.NoError(t, err)
	return fileName
}
