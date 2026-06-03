package agent

import (
	"strings"
	"testing"
)

func TestRun_UnsupportedInputExtension(t *testing.T) {
	_, err := Run("input.txt", "convert json to csv", "csv")
	if err == nil {
		t.Fatal("expected error for unsupported input extension")
	}
	if !strings.Contains(err.Error(), "inspection failed") {
		t.Fatalf("expected inspection failure, got %v", err)
	}
}
