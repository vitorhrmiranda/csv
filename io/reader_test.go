package io_test

import (
	"bytes"
	"fmt"
	"testing"

	csv "github.com/vitorhrmiranda/csv/io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewReadeable(t *testing.T) {
	want := [][]string{{"a", "b", "c"}}
	ActAndAssertNewReadeable(t, "a,b,c", want)
	ActAndAssertNewReadeable(t, []byte("a,b,c"), want)
	ActAndAssertNewReadeable(t, "", nil)
	ActAndAssertNewReadeable(t, []byte{}, nil)
}

func ActAndAssertNewReadeable[T csv.Readeable](t *testing.T, input T, want [][]string) {
	t.Helper()
	reader := csv.NewReadeable(input)
	content, err := reader.ReadAll()
	require.NoError(t, err)
	require.Equal(t, want, content)
}

func TestNewReaderWithHeaders(t *testing.T) {
	want := [][]string{{"1", "2", "3"}}
	ActAndAssertNewReaderWithHeaders(t, "a,b,c\n1,2,3", want)
	ActAndAssertNewReaderWithHeaders(t, []byte("a,b,c\n1,2,3"), want)
	ActAndAssertNewReaderWithHeaders(t, "", nil)
	ActAndAssertNewReaderWithHeaders(t, []byte{}, nil)
}

func ActAndAssertNewReaderWithHeaders[T csv.Readeable](t *testing.T, input T, want [][]string) {
	t.Helper()
	reader := csv.NewReaderWithHeaders(input, ',')
	content, err := reader.ReadAll()
	require.NoError(t, err)
	require.Equal(t, want, content)
}

func TestReader_ParseHeader(t *testing.T) {
	tcases := []struct {
		csv    string
		reader func(string) *csv.Reader
		assert func(*testing.T, *csv.Reader, error)
	}{
		{
			csv:    "a,b,c\n1,2,3",
			reader: csv.NewReadeable[string],
			assert: func(t *testing.T, reader *csv.Reader, err error) {
				require.NoError(t, err)
				AssertHeaderRow(t, csv.Header{"a": 0, "b": 1, "c": 2}, "a,b,c", reader.Header())
			},
		},
		{
			csv:    "a,b,c",
			reader: csv.NewReadeable[string],
			assert: func(t *testing.T, reader *csv.Reader, err error) {
				require.NoError(t, err)
				AssertHeaderRow(t, csv.Header{"a": 0, "b": 1, "c": 2}, "a,b,c", reader.Header())
			},
		},
		{
			csv:    "",
			reader: csv.NewReadeable[string],
			assert: func(t *testing.T, reader *csv.Reader, err error) {
				require.NoError(t, err)
				AssertHeaderRow(t, csv.Header{}, "", reader.Header())
			},
		},
		{
			csv: "a;b;c\n1,2,3",
			reader: func(content string) *csv.Reader {
				reader := csv.NewReadeable(content)
				reader.Read()
				return reader
			},
			assert: func(t *testing.T, reader *csv.Reader, err error) {
				require.ErrorIs(t, err, csv.ErrReadHeader)
				AssertHeaderRow(t, csv.Header(nil), "", reader.Header())
			},
		},
	}
	for _, tcase := range tcases {
		t.Run(tcase.csv, func(t *testing.T) {
			reader := tcase.reader(tcase.csv)
			err := reader.ParseHeader()
			tcase.assert(t, reader, err)
		})
	}
}

func AssertHeaderRow(t *testing.T, want csv.Header, wantString string, got csv.Row) {
	t.Helper()
	require.Equal(t, want, got.Header)
	require.Equal(t, wantString, got.String())
}

func TestReader_ForEach(test *testing.T) {
	tcases := []struct {
		csv    string
		reader func(string) *csv.Reader
		assert func(*testing.T, int, error)
		exec   func(*testing.T, csv.Row)
	}{
		{
			csv:    "a,b,c\n1,2,3",
			reader: csv.NewReadeable[string],
			exec: func(t *testing.T, row csv.Row) {
				t.Helper()
				expect := []csv.Columns{{"a", "b", "c"}, {"1", "2", "3"}}
				AssertRows(t, expect, row.Columns)
			},
			assert: func(t *testing.T, count int, err error) {
				require.NoError(t, err)
				require.Equal(t, 2, count)
			},
		},
		{
			csv:    "a,b,c\n\n1,2,3",
			reader: csv.NewReadeable[string],
			exec: func(t *testing.T, row csv.Row) {
				t.Helper()
				expect := []csv.Columns{{"a", "b", "c"}, {"1", "2", "3"}}
				AssertRows(t, expect, row.Columns)
			},
			assert: func(t *testing.T, count int, err error) {
				require.NoError(t, err)
				require.Equal(t, 2, count)
			},
		},
		{
			csv:    "a,b,c\n1,2,3,4\n5,6,7",
			reader: csv.NewReadeable[string],
			exec: func(t *testing.T, row csv.Row) {
				t.Helper()
				expect := []csv.Columns{{"a", "b", "c"}}
				AssertRows(t, expect, row.Columns)
			},
			assert: func(t *testing.T, count int, err error) {
				require.ErrorIs(t, err, csv.ErrReadLine)
				require.Equal(t, 1, count)
			},
		},
	}
	for _, tcase := range tcases {
		test.Run(tcase.csv, func(t *testing.T) {
			reader := tcase.reader(tcase.csv)
			var count int
			err := reader.ForEach(func(r csv.Row) {
				tcase.exec(t, r)
				count++
			})
			tcase.assert(t, count, err)
		})
	}
}

func AssertRows(t *testing.T, oneOf []csv.Columns, got csv.Columns) {
	t.Helper()
	for _, want := range oneOf {
		if assert.ObjectsAreEqualValues(want, got) {
			return
		}
	}
	assert.Fail(t, fmt.Sprintf("unexpected row: %s", got))
}

func TestReader_MapRows(t *testing.T) {
	transform := func(r *csv.Row) { r.Columns.Set(0, "T") }

	writer := bytes.NewBuffer(nil)
	err := csv.NewReadeable("a,b\n1,2\n").MapRows(writer, transform)
	assert.NoError(t, err)
	assert.Equal(t, "T,b\nT,2\n", writer.String())

	t.Run("parse header", func(t *testing.T) {
		writer := bytes.NewBuffer(nil)
		err := csv.NewReaderWithHeaders("a,b\n1,2\n", ',').MapRows(writer, transform)
		assert.NoError(t, err)
		assert.Equal(t, "a,b\nT,2\n", writer.String())
	})

	t.Run("reader err", func(t *testing.T) {
		wantErr := fmt.Errorf("error")
		err := csv.NewReader(&ioWithErr{wantErr}).MapRows(nil, nil)
		assert.ErrorIs(t, err, wantErr)
	})

	t.Run("write error", func(t *testing.T) {
		wantErr := fmt.Errorf("error")
		err := csv.NewReadeable("a,b\n1,2\n").MapRows(&ioWithErr{wantErr}, transform)
		assert.ErrorIs(t, err, wantErr)
		assert.ErrorIs(t, err, csv.ErrMapRows)
	})
}

type ioWithErr struct{ error }

func (w *ioWithErr) Write([]byte) (int, error) {
	return 0, w.error
}

func (w *ioWithErr) Read([]byte) (int, error) {
	return 0, w.error
}
