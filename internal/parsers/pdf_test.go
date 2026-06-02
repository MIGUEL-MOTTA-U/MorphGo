package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectPDF_ValidMinimalPDF(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.pdf")
	content := []byte("%PDF-1.4\n1 0 obj\n<< /Type /Page >>\nendobj\ntrailer\n<<>>\n%%EOF")
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write pdf: %v", err)
	}

	summary, err := InspectPDF(path)
	if err != nil {
		t.Fatalf("InspectPDF returned error: %v", err)
	}
	if summary.Format != "pdf" {
		t.Fatalf("expected format pdf, got %q", summary.Format)
	}
	if !summary.HasHeader {
		t.Fatal("expected PDF structure to be detected")
	}
	if len(summary.Columns) != 1 || summary.Columns[0] != "pages" {
		t.Fatalf("unexpected columns: %#v", summary.Columns)
	}
}

func TestInspectPDF_Invalid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "broken.pdf")
	if err := os.WriteFile(path, []byte("not a pdf"), 0o600); err != nil {
		t.Fatalf("write pdf: %v", err)
	}

	_, err := InspectPDF(path)
	if err == nil {
		t.Fatal("expected error for invalid pdf")
	}
}
