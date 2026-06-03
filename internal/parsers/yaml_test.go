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

func TestInspectYAML_ListOfMaps(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "list-of-maps.yaml")
	content := "- name: ana\n  age: 30\n- name: bob\n  age: 28\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	summary, err := InspectYAML(path)
	if err != nil {
		t.Fatalf("InspectYAML returned error: %v", err)
	}
	if summary.Rows != 2 {
		t.Fatalf("expected 2 items, got %d", summary.Rows)
	}
	if len(summary.Columns) != 2 {
		t.Fatalf("expected 2 columns, got %#v", summary.Columns)
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

func TestInspectYAML_MapColumns(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.yaml")
	if err := os.WriteFile(path, []byte("name: ana\nage: 30\n"), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	summary, err := InspectYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cols := map[string]bool{}
	for _, c := range summary.Columns {
		cols[c] = true
	}
	if !cols["name"] || !cols["age"] {
		t.Fatalf("expected columns name and age, got %#v", summary.Columns)
	}
}

func TestInspectYAML_EmptyList(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty-list.yaml")
	if err := os.WriteFile(path, []byte("[]\n"), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	summary, err := InspectYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Rows != 0 {
		t.Fatalf("expected 0 rows, got %d", summary.Rows)
	}
}

func TestInspectYAML_Scalar(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scalar.yaml")
	if err := os.WriteFile(path, []byte("42\n"), 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	summary, err := InspectYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summary.Columns) != 1 || summary.Columns[0] != "value" {
		t.Fatalf("expected [value] for scalar, got %#v", summary.Columns)
	}
}

func TestInspectYAML_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.yaml")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	summary, err := InspectYAML(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.HasHeader {
		t.Fatal("expected no header for empty file")
	}
	if len(summary.Columns) != 0 {
		t.Fatalf("expected no columns for empty file, got %#v", summary.Columns)
	}
}

func TestInspectYAML_FileNotFound(t *testing.T) {
	_, err := InspectYAML(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
