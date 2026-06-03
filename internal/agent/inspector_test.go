package agent

import "testing"

func TestInspect_UnsupportedExtension(t *testing.T) {
	_, err := Inspect("sample.txt")
	if err == nil {
		t.Fatal("expected error for unsupported file extension")
	}
}
