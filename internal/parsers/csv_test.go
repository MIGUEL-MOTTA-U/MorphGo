package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectCSV_WithHeader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.csv")
	if err := os.WriteFile(path, []byte("name,age\nana,30\n"), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	summary, err := InspectCSV(path)
	if err != nil {
		t.Fatalf("InspectCSV returned error: %v", err)
	}
	if summary.Format != "csv" {
		t.Fatalf("expected format csv, got %q", summary.Format)
	}
	if !summary.HasHeader {
		t.Fatal("expected header to be detected")
	}
	if len(summary.Columns) != 2 || summary.Columns[0] != "name" || summary.Columns[1] != "age" {
		t.Fatalf("unexpected columns: %#v", summary.Columns)
	}
	if summary.Rows != 1 {
		t.Fatalf("expected 1 data row, got %d", summary.Rows)
	}
}

func TestInspectCSV_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.csv")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("write empty csv: %v", err)
	}

	summary, err := InspectCSV(path)
	if err != nil {
		t.Fatalf("InspectCSV returned error: %v", err)
	}
	if summary.Format != "csv" {
		t.Fatalf("expected format csv, got %q", summary.Format)
	}
	if summary.HasHeader {
		t.Fatal("expected no header for empty file")
	}
	if len(summary.Columns) != 0 {
		t.Fatalf("expected no columns, got %#v", summary.Columns)
	}
	if summary.Rows != 0 {
		t.Fatalf("expected 0 rows, got %d", summary.Rows)
	}
}

func TestInspectCSV_HeaderOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "header-only.csv")
	if err := os.WriteFile(path, []byte("a,b\n"), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	summary, err := InspectCSV(path)
	if err != nil {
		t.Fatalf("InspectCSV returned error: %v", err)
	}
	if !summary.HasHeader {
		t.Fatal("expected header to be detected")
	}
	if summary.Rows != 0 {
		t.Fatalf("expected 0 rows, got %d", summary.Rows)
	}
}

func TestInspectCSV_FileNotFound(t *testing.T) {
	_, err := InspectCSV(filepath.Join(t.TempDir(), "missing.csv"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestInspectCSV_SingleField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "single-field.csv")
	if err := os.WriteFile(path, []byte("name\nana\n"), 0o600); err != nil {
		t.Fatalf("write csv: %v", err)
	}

	summary, err := InspectCSV(path)
	if err != nil {
		t.Fatalf("InspectCSV returned error: %v", err)
	}
	if len(summary.Columns) != 1 || summary.Columns[0] != "name" {
		t.Fatalf("unexpected columns: %#v", summary.Columns)
	}
	if summary.Rows != 1 {
		t.Fatalf("expected 1 row, got %d", summary.Rows)
	}
}
