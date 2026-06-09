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

func TestInspectJSON_ValidArray(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "array.json")
	if err := os.WriteFile(path, []byte(`[1,2,3]`), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	summary, err := InspectJSON(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Format != "json" {
		t.Errorf("expected format json, got %q", summary.Format)
	}
	if summary.Rows != 3 {
		t.Errorf("expected 3 rows, got %d", summary.Rows)
	}
}

func TestInspectJSON_ArrayOfObjects(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "objects.json")
	if err := os.WriteFile(path, []byte(`[{"name":"ana"},{"name":"bob"}]`), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	summary, err := InspectJSON(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Rows != 2 {
		t.Fatalf("expected 2 rows, got %d", summary.Rows)
	}
	if !summary.HasHeader {
		t.Fatal("expected array of objects to be marked as headered")
	}
}

func TestInspectJSON_EmptyArray(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty-array.json")
	if err := os.WriteFile(path, []byte(`[]`), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	summary, err := InspectJSON(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Rows != 0 {
		t.Errorf("expected 0 rows, got %d", summary.Rows)
	}
}

func TestInspectJSON_Scalar(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "scalar.json")
	if err := os.WriteFile(path, []byte(`42`), 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	summary, err := InspectJSON(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(summary.Columns) != 1 || summary.Columns[0] != "value" {
		t.Errorf("unexpected columns for scalar: %#v", summary.Columns)
	}
}

func TestInspectJSON_FileNotFound(t *testing.T) {
	_, err := InspectJSON(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInspectJSON_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("write json: %v", err)
	}

	_, err := InspectJSON(path)
	if err == nil {
		t.Fatal("expected parse error for empty file")
	}
}
