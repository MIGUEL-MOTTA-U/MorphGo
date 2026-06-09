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
		Objective:  "convert json to csv",
		Source:     "json",
		Target:     "csv",
		Operation:  "convert",
		InputPath:  "input.json",
		OutputPath: "output.csv",
		Summary:    schema.Summary{Columns: []string{"id", "name"}},
	}

	code, err := GenerateMain(plan, nil)
	if err != nil {
		t.Fatalf("GenerateMain returned error: %v", err)
	}
	if code == "" {
		t.Fatal("expected generated code")
	}
	if !strings.Contains(code, "package main") {
		t.Fatal("expected package main")
	}
	if !strings.Contains(code, "\"encoding/json\"") {
		t.Fatal("expected json import in real generator")
	}
	if !strings.Contains(code, "func run() error") {
		t.Fatal("expected run helper")
	}
	if count := strings.Count(code, "func "); count != 3 {
		t.Fatalf("expected exactly 3 functions, got %d", count)
	}
}

func TestGenerateMain_MissingData(t *testing.T) {
	_, err := GenerateMain(schema.Plan{}, nil)
	if err == nil {
		t.Fatal("expected error for incomplete plan")
	}
}

func TestGenerateMain_JsonToCsvRequiresPaths(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
	}

	_, err := GenerateMain(plan, nil)
	if err == nil {
		t.Fatal("expected error when json to csv paths are missing")
	}
}

func TestGenerateMain_CompilesInHappyPath(t *testing.T) {
	plan := schema.Plan{
		Objective:  "convert json to csv",
		Source:     "json",
		Target:     "csv",
		Operation:  "convert",
		InputPath:  "input.json",
		OutputPath: "output.csv",
		Summary:    schema.Summary{Columns: []string{"id", "name"}},
	}

	code, err := GenerateMain(plan, nil)
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
	if _, err := cmd.CombinedOutput(); err != nil {
		// Since we don't have a go.sum or dependencies, just build is enough
		cmdBuild := exec.Command("go", "build", "-o", "main.exe", "main.go")
		cmdBuild.Dir = dir
		if outB, errB := cmdBuild.CombinedOutput(); errB != nil {
			t.Fatalf("generated code did not compile: %v\n%s", errB, string(outB))
		}
	}
}
