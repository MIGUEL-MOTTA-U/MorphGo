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

func TestInspectXML_ColumnsContent(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "sample.xml")
    if err := os.WriteFile(path, []byte(`<root id="1"><child name="ana"/></root>`), 0o600); err != nil {
        t.Fatalf("write xml: %v", err)
    }

    summary, err := InspectXML(path)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    expected := map[string]bool{"root": true, "child": true, "@id": true, "@name": true}
    if len(summary.Columns) != len(expected) {
        t.Fatalf("expected %d columns, got %d: %#v", len(expected), len(summary.Columns), summary.Columns)
    }
    for _, col := range summary.Columns {
        if !expected[col] {
            t.Errorf("unexpected column %q", col)
        }
    }
}

func TestInspectXML_NoAttributes(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "no-attrs.xml")
    if err := os.WriteFile(path, []byte(`<root><item/><item/></root>`), 0o600); err != nil {
        t.Fatalf("write xml: %v", err)
    }

    summary, err := InspectXML(path)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    for _, col := range summary.Columns {
        if len(col) > 0 && col[0] == '@' {
            t.Errorf("expected no attribute columns, got %q", col)
        }
    }
    if len(summary.Columns) != 2 { // root, item
        t.Fatalf("expected 2 node columns, got %#v", summary.Columns)
    }
}

func TestInspectXML_MinimalRoot(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "minimal.xml")
    if err := os.WriteFile(path, []byte(`<root/>`), 0o600); err != nil {
        t.Fatalf("write xml: %v", err)
    }

    summary, err := InspectXML(path)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !summary.HasHeader {
        t.Fatal("expected HasHeader true for non-empty XML")
    }
    if len(summary.Columns) != 1 || summary.Columns[0] != "root" {
        t.Fatalf("expected [root], got %#v", summary.Columns)
    }
}

func TestInspectXML_EmptyFile(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "empty.xml")
    if err := os.WriteFile(path, nil, 0o600); err != nil {
        t.Fatalf("write xml: %v", err)
    }

    _, err := InspectXML(path)
    if err == nil {
        t.Fatal("expected parse error for empty file")
    }
}

func TestInspectXML_FileNotFound(t *testing.T) {
    _, err := InspectXML(filepath.Join(t.TempDir(), "missing.xml"))
    if err == nil {
        t.Fatal("expected error for missing file")
    }
}