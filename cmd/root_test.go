package cmd

import "testing"

func TestRootValidation_MissingFlags(t *testing.T) {
    cmd := newTestRootCmd()
    cmd.SetArgs([]string{})
    _, err := cmd.ExecuteC()
    if err == nil {
        t.Fatal("expected error when flags missing")
    }
}

func TestRootValidation_WithFlags(t *testing.T) {
    cmd := newTestRootCmd()
    cmd.SetArgs([]string{"--input", "file.txt", "--task", "do"})
    _, err := cmd.ExecuteC()
    if err != nil {
        t.Fatalf("expected no error when flags present: %v", err)
    }
}
