package cmd

import "testing"

func TestRootValidation_MissingFlags(t *testing.T) {
    rootCmd.SetArgs([]string{})
    _, err := rootCmd.ExecuteC()
    if err == nil {
        t.Fatal("expected error when flags missing")
    }
}

func TestRootValidation_WithFlags(t *testing.T) {
    rootCmd.SetArgs([]string{"--input", "file.txt", "--task", "do"})
    _, err := rootCmd.ExecuteC()
    if err != nil {
        t.Fatalf("expected no error when flags present: %v", err)
    }
}
