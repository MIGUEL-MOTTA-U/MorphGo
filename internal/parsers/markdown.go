package parsers

import (
	"bufio"
	"os"
	"strings"

	"morphgo/internal/schema"
)

// InspectMarkdown reads a Markdown or text file and returns a minimal structural summary.
func InspectMarkdown(path string) (schema.Summary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schema.Summary{}, err
	}
	if len(data) == 0 {
		return schema.Summary{Format: "markdown"}, nil
	}

	summary := schema.Summary{Format: "markdown"}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	seen := map[string]struct{}{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		switch {
		case strings.HasPrefix(line, "#"):
			seen["heading"] = struct{}{}
		case strings.HasPrefix(line, "- "), strings.HasPrefix(line, "* "), strings.HasPrefix(line, "1. "):
			seen["list"] = struct{}{}
		default:
			seen["paragraph"] = struct{}{}
		}
	}
	if err := scanner.Err(); err != nil {
		return schema.Summary{}, err
	}

	summary.Columns = make([]string, 0, len(seen))
	for k := range seen {
		summary.Columns = append(summary.Columns, k)
	}
	if len(summary.Columns) > 0 {
		summary.HasHeader = true
	}
	return summary, nil
}
