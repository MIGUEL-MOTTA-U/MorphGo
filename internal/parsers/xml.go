package parsers

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"

	"morphgo/internal/schema"
)

// InspectXML reads an XML file and returns a minimal structural summary.
func InspectXML(path string) (schema.Summary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return schema.Summary{}, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(data))
	summary := schema.Summary{Format: "xml"}
	nodes := make(map[string]struct{})
	attrs := make(map[string]struct{})
	seenToken := false

	for {
		tok, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				if !seenToken {
					return schema.Summary{}, errors.New("empty xml file")
				}
				break
			}
			return schema.Summary{}, err
		}
		seenToken = true
		switch t := tok.(type) {
		case xml.StartElement:
			nodes[t.Name.Local] = struct{}{}
			for _, attr := range t.Attr {
				attrs[attr.Name.Local] = struct{}{}
			}
		}
	}

	summary.Columns = make([]string, 0, len(nodes)+len(attrs))
	for name := range nodes {
		summary.Columns = append(summary.Columns, name)
	}
	for name := range attrs {
		summary.Columns = append(summary.Columns, "@"+name)
	}
	if len(summary.Columns) > 0 {
		summary.HasHeader = true
	}

	return summary, nil
}
