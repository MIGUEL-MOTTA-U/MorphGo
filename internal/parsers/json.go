package parsers

import (
	"encoding/json"
	"os"

	"morphgo/internal/schema"
)

// InspectJSON reads a JSON file and returns a minimal structural summary.
func InspectJSON(path string) (schema.Summary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schema.Summary{}, err
	}

	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return schema.Summary{}, err
	}

	summary := schema.Summary{Format: "json"}
	switch v := value.(type) {
	case map[string]any:
		summary.Columns = make([]string, 0, len(v))
		for key := range v {
			summary.Columns = append(summary.Columns, key)
		}
		summary.HasHeader = true
	case []any:
		summary.Rows = len(v)
	default:
		summary.HasHeader = true
		summary.Columns = []string{"value"}
	}

	return summary, nil
}
