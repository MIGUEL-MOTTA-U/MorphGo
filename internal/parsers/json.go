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
		// If it's a map with only one key and that key is an array, inspect the array
		if len(v) == 1 {
			for _, val := range v {
				if arr, ok := val.([]any); ok && len(arr) > 0 {
					if first, ok := arr[0].(map[string]any); ok {
						summary.Columns = make([]string, 0, len(first))
						for k := range first {
							summary.Columns = append(summary.Columns, k)
						}
					}
				}
			}
		}
	case []any:
		summary.Rows = len(v)
		if len(v) > 0 {
			if first, ok := v[0].(map[string]any); ok {
				summary.Columns = make([]string, 0, len(first))
				for k := range first {
					summary.Columns = append(summary.Columns, k)
				}
				summary.HasHeader = true
			}
		}
	default:
		summary.HasHeader = true
		summary.Columns = []string{"value"}
	}

	return summary, nil
}
