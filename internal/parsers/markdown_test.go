package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInspectMarkdown_DetectsParagraphAndList(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.md")
	content := "# Title\n\nA short paragraph.\n- item one\n"
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
	if !summary.HasHeader {
		t.Fatal("expected markdown summary to mark structural content")
	}
	if len(summary.Columns) != 3 {
		t.Fatalf("expected 3 detected structures, got %#v", summary.Columns)
	}
}
