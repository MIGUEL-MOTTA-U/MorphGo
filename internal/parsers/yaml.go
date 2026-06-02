package parsers

import (
	"os"

	"gopkg.in/yaml.v3"

	"morphgo/internal/schema"
)

// InspectYAML reads a YAML file and returns a minimal structural summary.
func InspectYAML(path string) (schema.Summary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schema.Summary{}, err
	}

	var value any
	if err := yaml.Unmarshal(data, &value); err != nil {
		return schema.Summary{}, err
	}

	summary := schema.Summary{Format: "yaml"}
	switch v := value.(type) {
	case map[string]any:
		summary.HasHeader = true
		summary.Columns = make([]string, 0, len(v))
		for key := range v {
			summary.Columns = append(summary.Columns, key)
		}
	case []any:
		summary.Rows = len(v)
		summary.HasHeader = true
	default:
		summary.HasHeader = true
		summary.Columns = []string{"value"}
	}

	return summary, nil
}
