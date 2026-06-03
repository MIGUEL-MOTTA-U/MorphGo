package agent

import (
	"fmt"
	"path/filepath"
	"strings"

	"morphgo/internal/parsers"
	"morphgo/internal/schema"
)

// Inspect identifies the file type by extension and delegates to the appropriate parser.
func Inspect(path string) (schema.Summary, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json":
		return parsers.InspectJSON(path)
	case ".xml":
		return parsers.InspectXML(path)
	case ".csv":
		return parsers.InspectCSV(path)
	case ".yaml", ".yml":
		return parsers.InspectYAML(path)
	case ".xlsx":
		return parsers.InspectExcel(path)
	case ".md":
		return parsers.InspectMarkdown(path)
	case ".pdf":
		return parsers.InspectPDF(path)
	default:
		return schema.Summary{}, fmt.Errorf("unsupported file extension: %s", ext)
	}
}
