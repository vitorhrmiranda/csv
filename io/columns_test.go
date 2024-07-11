package csv_test

import (
	"fmt"
	"testing"

	csv "github.com/vitorhrmiranda/csv/io"

	"github.com/stretchr/testify/assert"
)

func TestColumns_At(t *testing.T) {
	testCases := []struct {
		columns csv.Columns
		index   int
		want    string
	}{
		{csv.Columns{"a", "b", "c"}, 0, "a"},
		{csv.Columns{"a", "b", "c"}, -1, ""},
		{csv.Columns{"a", "b", "c"}, 3, ""},
		{csv.Columns{"a", "b", "c"}, 1, "b"},
		{csv.Columns{"a", "b", "c"}, 2, "c"},
		{csv.Columns{"a", "b", "c", ""}, 3, ""},
	}
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("%v", testCase.columns), func(t *testing.T) {
			got := testCase.columns.At(testCase.index)
			assert.Equal(t, testCase.want, got)
		})
	}
}
