package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootCmd() *cobra.Command {
	cmd := newRootCmd()
	cmd.AddCommand(newRunCmd())
	return cmd
}

func TestRun_Help(t *testing.T) {
	cmd := newTestRootCmd()
	cmd.SetArgs([]string{"run", "--help"})
	_, err := cmd.ExecuteC()
	if err != nil {
		t.Fatalf("unexpected error from help: %v", err)
	}
}

func TestRun_MissingFlags(t *testing.T) {
	// ensure flag variables are reset to simulate missing flags
	inputPath = ""
	taskStr = ""
	outputPath = ""
	targetFormat = ""
	cmd := newTestRootCmd()
	cmd.SetArgs([]string{"run"})
	_, err := cmd.ExecuteC()
	if err == nil {
		t.Fatal("expected error when flags missing")
	}
}

func TestRun_WithFlags(t *testing.T) {
	cmd := newTestRootCmd()
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "in.json")
	if err := os.WriteFile(inputPath, []byte(`{"items":[{"id":1,"name":"a"}]}`), 0o600); err != nil {
		t.Fatalf("write input file: %v", err)
	}
	outputDir := filepath.Join(dir, "out")
	cmd.SetArgs([]string{"run", "--input", inputPath, "--task", "convert to json", "--output", outputDir, "--target", "json"})
	_, err := cmd.ExecuteC()
	if err != nil {
		t.Fatalf("expected no error when flags present: %v", err)
	}
}
