package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectXML_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.xml")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("write xml: %v", err)
	}

	_, err := InspectXML(path)
	if err == nil {
		t.Fatal("expected error for empty xml file")
	}
}
