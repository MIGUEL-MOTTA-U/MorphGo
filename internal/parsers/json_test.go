package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectJSON_ValidObject(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	if err := os.WriteFile(path, []byte(`{"name":"ana","age":30}`), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	summary, err := InspectJSON(path)
	if err != nil {
		t.Fatalf("InspectJSON returned error: %v", err)
	}
	if summary.Format != "json" {
		t.Fatalf("expected format json, got %q", summary.Format)
	}
	if !summary.HasHeader {
		t.Fatal("expected object fields to be detected")
	}
	if len(summary.Columns) != 2 {
		t.Fatalf("expected 2 fields, got %#v", summary.Columns)
	}
}

func TestInspectJSON_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(path, []byte(`{"name":`), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	_, err := InspectJSON(path)
	if err == nil {
		t.Fatal("expected error for invalid json")
	}
}
