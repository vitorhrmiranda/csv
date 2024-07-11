package csv

type Columns []string

// At returns the value of a column by index.
// If the index is out of range, it returns an empty string.
func (columns Columns) At(index int) string {
	if index < 0 || index >= len(columns) {
		return ""
	}
	return columns[index]
}

// IsEmpty returns true if the columns are empty.
func (columns Columns) IsEmpty() bool {
	return len(columns) == 0
}
