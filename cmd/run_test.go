package cmd

import "testing"

func TestRun_Help(t *testing.T) {
	rootCmd.SetArgs([]string{"run", "--help"})
	_, err := rootCmd.ExecuteC()
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
	rootCmd.SetArgs([]string{"run"})
	_, err := rootCmd.ExecuteC()
	if err == nil {
		t.Fatal("expected error when flags missing")
	}
}

func TestRun_WithFlags(t *testing.T) {
	rootCmd.SetArgs([]string{"run", "--input", "in.txt", "--task", "do", "--output", "out", "--target", "json"})
	_, err := rootCmd.ExecuteC()
	if err != nil {
		t.Fatalf("expected no error when flags present: %v", err)
	}
}
