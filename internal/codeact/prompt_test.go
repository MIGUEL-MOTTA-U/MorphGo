package codeact

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"morphgo/internal/schema"
)

func TestGenerateMain_KnownPlan(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
	}

	code, err := GenerateMain(plan)
	if err != nil {
		t.Fatalf("GenerateMain returned error: %v", err)
	}
	if code == "" {
		t.Fatal("expected generated code")
	}
	if !strings.Contains(code, "package main") {
		t.Fatal("expected package main")
	}
	if strings.Contains(code, "\"errors\"") {
		t.Fatal("expected no extra error import")
	}
	if !strings.Contains(code, "\"fmt\"") {
		t.Fatal("expected valid imports")
	}
	if !strings.Contains(code, "func run() error") {
		t.Fatal("expected run helper")
	}
	if !strings.Contains(code, "if err := run(); err != nil") {
		t.Fatal("expected basic error handling")
	}
	if count := strings.Count(code, "func "); count != 2 {
		t.Fatalf("expected exactly 2 functions, got %d", count)
	}
	if strings.Contains(code, "reflect") {
		t.Fatal("expected no reflection usage")
	}
}

func TestGenerateMain_MissingData(t *testing.T) {
	_, err := GenerateMain(schema.Plan{})
	if err == nil {
		t.Fatal("expected error for incomplete plan")
	}
}

func TestGenerateMain_NotEmpty(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
	}

	code, err := GenerateMain(plan)
	if err != nil {
		t.Fatalf("GenerateMain returned error: %v", err)
	}
	if strings.TrimSpace(code) == "" {
		t.Fatal("expected non-empty generated code")
	}
}

func TestGenerateMain_CompilesInHappyPath(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
	}

	code, err := GenerateMain(plan)
	if err != nil {
		t.Fatalf("GenerateMain returned error: %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(code), 0o600); err != nil {
		t.Fatalf("write generated code: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module tempgen\ngo 1.26.3\n"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	cmd := exec.Command("go", "test", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("generated code did not compile: %v\n%s", err, string(out))
	}
}
