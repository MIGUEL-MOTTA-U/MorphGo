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
	summary.Columns = append(summary.Columns, sheets...)

	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return schema.Summary{}, err
		}
		if len(rows) > 0 {
			summary.Rows += len(rows)
		}
	}

	return summary, nil
}
