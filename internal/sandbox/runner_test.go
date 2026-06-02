package sandbox

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunMain_Success(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module tempgen\ngo 1.26.3\n"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainPath, []byte("package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"ok\")}\n"), 0o600); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	res, err := RunMain(mainPath, 5*time.Second)
	if err != nil {
		t.Fatalf("RunMain returned error: %v", err)
	}
	if !strings.Contains(res.Stdout, "ok") {
		t.Fatalf("expected stdout to contain ok, got %q", res.Stdout)
	}
}

func TestRunMain_Timeout(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module tempgen\ngo 1.26.3\n"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainPath, []byte("package main\nfunc main(){for{}}\n"), 0o600); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	_, err := RunMain(mainPath, 100*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestRunMain_CompileError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module tempgen\ngo 1.26.3\n"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainPath, []byte("package main\nfunc main(){broken}\n"), 0o600); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	_, err := RunMain(mainPath, 5*time.Second)
	if err == nil {
		t.Fatal("expected compile error")
	}
}

func TestRunMain_RuntimeError(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module tempgen\ngo 1.26.3\n"), 0o600); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	mainPath := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainPath, []byte(`package main
import "os"
func main(){os.Exit(1)}
`), 0o600); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	res, err := RunMain(mainPath, 5*time.Second)
	if err == nil {
		t.Fatal("expected runtime error")
	}
	if res.ExitCode == 0 {
		t.Fatal("expected non-zero exit code")
	}
}
