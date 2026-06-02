package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectMarkdown_HeadingAndList(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.md")
	content := "# Title\n\n- item one\nplain paragraph\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}

	summary, err := InspectMarkdown(path)
	if err != nil {
		t.Fatalf("InspectMarkdown returned error: %v", err)
	}
	if summary.Format != "markdown" {
		t.Fatalf("expected format markdown, got %q", summary.Format)
	}
	if len(summary.Columns) != 3 {
		t.Fatalf("expected 3 patterns, got %#v", summary.Columns)
	}
}

func TestInspectMarkdown_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.md")
	if err := os.WriteFile(path, nil, 0o600); err != nil {
		t.Fatalf("write markdown: %v", err)
	}

	summary, err := InspectMarkdown(path)
	if err != nil {
		t.Fatalf("InspectMarkdown returned error: %v", err)
	}
	if len(summary.Columns) != 0 {
		t.Fatalf("expected no structural patterns, got %#v", summary.Columns)
	}
}

func TestInspectMarkdown_FileNotFound(t *testing.T) {
	_, err := InspectMarkdown(filepath.Join(t.TempDir(), "missing.md"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
