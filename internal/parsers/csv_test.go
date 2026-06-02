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
