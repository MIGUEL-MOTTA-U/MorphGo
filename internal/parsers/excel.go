package parsers

import (
	"os"

	"github.com/xuri/excelize/v2"

	"morphgo/internal/schema"
)

// InspectExcel reads an Excel workbook and returns a minimal structural summary.
func InspectExcel(path string) (schema.Summary, error) {
	if _, err := os.Stat(path); err != nil {
		return schema.Summary{}, err
	}

	f, err := excelize.OpenFile(path)
	if err != nil {
		return schema.Summary{}, err
	}
	defer func() {
		_ = f.Close()
	}()

	sheets := f.GetSheetList()
	summary := schema.Summary{Format: "excel"}
	if len(sheets) == 0 {
		return summary, nil
	}

	summary.HasHeader = true
	// Use headers from the first sheet as columns
	if rows, err := f.GetRows(sheets[0]); err == nil && len(rows) > 0 {
		summary.Columns = rows[0]
	}

	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return schema.Summary{}, err
		}
		if len(rows) > 1 {
			summary.Rows += len(rows) - 1
		}
	}

	return summary, nil
}
