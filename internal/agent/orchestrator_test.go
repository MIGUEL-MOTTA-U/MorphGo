package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun_UnsupportedInputStatus(t *testing.T) {
	log, err := Run("input.txt", "convert json to csv", "csv")
	if err == nil {
		t.Fatal("expected error for unsupported input extension")
	}
	if log.Status != "error" {
		t.Fatalf("expected error status, got %q", log.Status)
	}
}

func TestRun_SuccessfulJsonToCsv(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	if err := os.WriteFile(inputPath, []byte(`{"employees":[{"name":"ana","age":30}]}`), 0o600); err != nil {
		t.Fatalf("write input json: %v", err)
	}

	log, err := Run(inputPath, "convert json to csv", "csv")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	if log.Plan.OutputPath != "output.csv" {
		t.Fatalf("expected default output path, got %q", log.Plan.OutputPath)
	}
}
