package io_test

import (
	"testing"

	csv "github.com/vitorhrmiranda/csv/io"

	"github.com/stretchr/testify/require"
)

func TestRow_Column(t *testing.T) {
	t.Run("no headers", func(t *testing.T) {
		r := csv.Row{}
		value, ok := r.Column("name")
		require.Empty(t, value)
		require.False(t, ok)
	})
	t.Run("with headers", func(t *testing.T) {
		r := csv.Row{
			Columns: []string{"csv"},
			Headers: csv.Headers{"name": 0},
		}
		value, ok := r.Column("name")
		require.Equal(t, value, "csv")
		require.True(t, ok)
	})
}

func TestRow_String(t *testing.T) {
	t.Run("no columns", func(t *testing.T) {
		r := csv.Row{}
		value := r.String()
		require.Equal(t, "", value)
	})
	t.Run("with columns", func(t *testing.T) {
		r := csv.Row{
			Columns: []string{"a", "b", "c"},
		}
		r.Comma = ","
		value := r.String()
		require.Equal(t, "a,b,c", value)
	})
}
