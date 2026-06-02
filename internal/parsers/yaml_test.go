package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectYAML_Map(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.yaml")
	if err := os.WriteFile(path, []byte("name: ana\nage: 30\n"), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	summary, err := InspectYAML(path)
	if err != nil {
		t.Fatalf("InspectYAML returned error: %v", err)
	}
	if summary.Format != "yaml" {
		t.Fatalf("expected format yaml, got %q", summary.Format)
	}
	if !summary.HasHeader {
		t.Fatal("expected map keys to be detected")
	}
	if len(summary.Columns) != 2 {
		t.Fatalf("expected 2 keys, got %#v", summary.Columns)
	}
}

func TestInspectYAML_List(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "list.yaml")
	if err := os.WriteFile(path, []byte("- a\n- b\n"), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	summary, err := InspectYAML(path)
	if err != nil {
		t.Fatalf("InspectYAML returned error: %v", err)
	}
	if summary.Rows != 2 {
		t.Fatalf("expected 2 items, got %d", summary.Rows)
	}
}

func TestInspectYAML_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.yaml")
	if err := os.WriteFile(path, []byte("name: ["), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	_, err := InspectYAML(path)
	if err == nil {
		t.Fatal("expected error for invalid yaml")
	}
}
