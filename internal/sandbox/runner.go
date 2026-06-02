package sandbox

import (
	"context"
	"os/exec"
	"path/filepath"
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
	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	res := Result{Stdout: string(out), ExitCode: 0}
	if err != nil {
		res.ExitCode = 1
		if ctx.Err() == context.DeadlineExceeded {
			res.Stderr = ctx.Err().Error()
			return res, ctx.Err()
		}
		res.Stderr = err.Error()
		return res, err
	}
	return res, nil
}
