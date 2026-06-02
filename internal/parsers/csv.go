package parsers

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"morphgo/internal/schema"
)

// InspectCSV reads a CSV file and returns a minimal structural summary.
func InspectCSV(path string) (schema.Summary, error) {
	f, err := os.Open(path)
	if err != nil {
		return schema.Summary{}, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return schema.Summary{Format: "csv"}, nil
		}
		return schema.Summary{}, err
	}
	if len(records) == 0 {
		return schema.Summary{Format: "csv"}, nil
	}

	cols := append([]string(nil), records[0]...)
	return schema.Summary{
		Format:    "csv",
		Separator: ',',
		HasHeader: true,
		Columns:   cols,
		Rows:      len(records) - 1,
	}, nil
}
