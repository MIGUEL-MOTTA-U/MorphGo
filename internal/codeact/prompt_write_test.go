package codeact

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteTempMain_CreatesMainGo(t *testing.T) {
	path, err := WriteTempMain("package main\n")
	if err != nil {
		t.Fatalf("WriteTempMain returned error: %v", err)
	}
	if !strings.HasSuffix(path, string(filepath.Separator)+"main.go") {
		t.Fatalf("expected main.go path, got %q", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	if string(data) != "package main\n" {
		t.Fatalf("unexpected file content: %q", string(data))
	}
}

func TestWriteTempMain_EmptyCode(t *testing.T) {
	_, err := WriteTempMain("")
	if err == nil {
		t.Fatal("expected error for empty code")
	}
}
