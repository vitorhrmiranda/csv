# vitorhrmiranda/csv

[![godoc](https://godoc.org/github.com/vitorhrmiranda/csv?status.svg)](https://godoc.org/github.com/vitorhrmiranda/csv)

Package `encoding/csv` extends the standard library's `encoding/csv` package with additional features.

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/vitorhrmiranda/csv
```

## Example

```go
package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/vitorhrmiranda/csv/io"
)

func main() {
	csv := "Sun,Mon,Tue,Wed,Thu,Fri,Sat\n0,1,2,3,4,5,6"
	reader := io.NewReaderWithHeaders(csv, ',')

	incrementOneDay := func(day string) string {
		d, _ := strconv.ParseInt(day, 10, 64)
		return fmt.Sprintf("%d", d+1)
	}

	output := &bytes.Buffer{}
	_ = reader.MapRows(output, func(row *io.Row) {
		for i, column := range row.Columns {
			_ = row.Columns.Set(i, incrementOneDay(column))
		}
	})

	fmt.Print(output.String())
}
```

## License
MIT licensed. See the LICENSE file for details.
