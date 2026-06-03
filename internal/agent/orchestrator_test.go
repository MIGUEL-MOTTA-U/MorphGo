package agent

import "testing"

func TestRun_UnsupportedInputStatus(t *testing.T) {
	log, err := Run("input.txt", "convert json to csv", "csv")
	if err == nil {
		t.Fatal("expected error for unsupported input extension")
	}
	if log.Status != "error" {
		t.Fatalf("expected error status, got %q", log.Status)
	}
}
