package codeact

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"morphgo/internal/schema"
)

func TestPrepareTempMain_WritesGeneratedCode(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
	}

	path, err := PrepareTempMain(plan, nil)
	if err != nil {
		t.Fatalf("PrepareTempMain returned error: %v", err)
	}
	if !strings.HasSuffix(path, string(filepath.Separator)+"main.go") {
		t.Fatalf("expected main.go path, got %q", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read generated file: %v", err)
	}
	if !strings.Contains(string(data), "package main") {
		t.Fatalf("unexpected generated content: %q", string(data))
	}
}

func TestPrepareTempMain_MissingData(t *testing.T) {
	_, err := PrepareTempMain(schema.Plan{}, nil)
	if err == nil {
		t.Fatal("expected error for incomplete plan")
	}
}
