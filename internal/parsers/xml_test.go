package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectXML_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.xml")
	if err := os.WriteFile(path, []byte(`<root id="1"><child name="ana"/></root>`), 0o600); err != nil {
		t.Fatalf("write xml: %v", err)
	}

	summary, err := InspectXML(path)
	if err != nil {
		t.Fatalf("InspectXML returned error: %v", err)
	}
	if summary.Format != "xml" {
		t.Fatalf("expected format xml, got %q", summary.Format)
	}
	if !summary.HasHeader {
		t.Fatal("expected nodes or attributes to be detected")
	}
	if len(summary.Columns) == 0 {
		t.Fatal("expected extracted nodes or attributes")
	}
}

func TestInspectXML_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.xml")
	if err := os.WriteFile(path, []byte(`<root><child></root>`), 0o600); err != nil {
		t.Fatalf("write xml: %v", err)
	}

	_, err := InspectXML(path)
	if err == nil {
		t.Fatal("expected error for invalid xml")
	}
}
