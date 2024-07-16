package io_test

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/vitorhrmiranda/csv/io"
)

func ExampleReader_MapRows() {
	csv := "Sun,Mon,Tue,Wed,Thu,Fri,Sat\n0,1,2,3,4,5,6"
	reader := io.NewReaderWithHeaders(csv, ',')

	incrementOneDay := func(day string) string {
		d, _ := strconv.ParseInt(day, 10, 64)
		return fmt.Sprintf("%d", d+1)
	}

	output := &bytes.Buffer{}
	reader.MapRows(output, func(row *io.Row) {
		for i, column := range row.Columns {
			row.Columns.Set(i, incrementOneDay(column))
		}
	})

	fmt.Println(output.String())
	// Output:
	// Sun,Mon,Tue,Wed,Thu,Fri,Sat
	// 1,2,3,4,5,6,7
}
