package csv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitorhrmiranda/csv"
	csvio "github.com/vitorhrmiranda/csv/io"
)

func TestNewReader(t *testing.T) {
	fileName := CreateExampleFile(t, []byte("a,b,c\n1,2,3\n4,5,6\n7,8,9"))

	reader, err := csv.NewReader(fileName, ',')
	require.NoError(t, err)

	reader.ForEach(func(row csvio.Row) {
		value, ok := row.Column("a")
		require.True(t, ok)
		require.NotEmpty(t, value)
		t.Log(value)
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
