package csv_test

import (
	"testing"

	"github.com/miketmoore/zelduh/assert"
	"github.com/miketmoore/zelduh/csv"
)

func TestParse(t *testing.T) {
	a := "86,86,86,86,\n86,86,86,86,"
	got := csv.Parse(a)
	expected := [][]string{
		[]string{"86,86,86,86"},
		[]string{"86,86,86,86"},
	}
	assert.Ok(t, got != nil)
	assert.Ok(t, len(got) == len(expected))
}
