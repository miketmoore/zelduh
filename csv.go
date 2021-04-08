package zelduh

import (
	"encoding/csv"
	"io"
	"log"
	"strings"
)

// ParseCSV translates a CSV string into a multi-dimensional array of strings
func ParseCSV(in string) [][]string {
	r := csv.NewReader(strings.NewReader(in))

	records := [][]string{}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		records = append(records, record)
	}
	return records
}
