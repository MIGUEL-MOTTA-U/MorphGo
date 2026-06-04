package sandbox

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Result captures the outcome of running generated code.
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

// RunMain compiles and executes a temporary main.go with a timeout.
func RunMain(mainPath string, timeout time.Duration) (Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	dir := filepath.Dir(mainPath)
	binPath := filepath.Join(dir, "morphgo-runner.exe")
	if err := buildBinary(ctx, dir, binPath); err != nil {
		return Result{Stderr: err.Error(), ExitCode: 1}, err
	}
	defer os.Remove(binPath)

	res, err := execBinary(ctx, binPath, dir)
	if err != nil {
		return res, err
	}
	return res, nil
}

func buildBinary(ctx context.Context, dir, binPath string) error {
	// Ensure dependencies are resolved
	tidyCmd := exec.CommandContext(ctx, "go", "mod", "tidy")
	tidyCmd.Dir = dir
	_ = tidyCmd.Run()

	cmd := exec.CommandContext(ctx, "go", "build", "-o", binPath, ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return ctx.Err()
		}
		if len(out) > 0 {
			return errors.New(string(out))
		}
		return err
	}
	return nil
}

func execBinary(ctx context.Context, binPath, dir string) (Result, error) {
	cmd := exec.CommandContext(ctx, binPath)
	cmd.Dir = dir
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	res := Result{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: 0}
	if err != nil {
		res.ExitCode = 1
		if ctx.Err() == context.DeadlineExceeded {
			return res, ctx.Err()
		}
		return res, err
	}
	return res, nil
}
