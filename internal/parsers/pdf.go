package parsers

import (
	"bytes"
	"os"
	"strings"

	"morphgo/internal/schema"
)

// InspectPDF reads a PDF file and returns a minimal structural summary.
func InspectPDF(path string) (schema.Summary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schema.Summary{}, err
	}
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		return schema.Summary{}, os.ErrInvalid
	}

	summary := schema.Summary{Format: "pdf"}
	text := string(data)
	if strings.Contains(text, "/Type /Page") || strings.Contains(text, "/Type/Page") {
		summary.Rows = strings.Count(text, "/Type /Page") + strings.Count(text, "/Type/Page")
		if summary.Rows == 0 {
			summary.Rows = 1
		}
	} else {
		summary.Rows = 1
	}
	summary.HasHeader = true
	summary.Columns = []string{"pages"}

	return summary, nil
}
