package codeact

import (
	"errors"
	"strings"
	"time"

	"morphgo/internal/sandbox"
	"morphgo/internal/schema"
)

var (
	ErrMaxAttempts   = errors.New("max attempts reached")
	ErrRepeatedError = errors.New("repeated error")
)

// PrepareTempMain generates and writes the temporary main.go for a run.
func PrepareTempMain(plan schema.Plan, history []Attempt) (string, error) {
	code, err := GenerateMain(plan, history)
	if err != nil {
		return "", err
	}
	return WriteTempMain(code)
}

// RunPlan prepares and executes the generated program through the sandbox.
func RunPlan(plan schema.Plan, timeout time.Duration) (sandbox.Result, error) {
	mainPath, err := PrepareTempMain(plan, nil)
	if err != nil {
		return sandbox.Result{}, err
	}
	return sandbox.RunMain(mainPath, timeout)
}

// RetryRunPlan retries generation and sandbox execution a limited number of times.
func RetryRunPlan(plan schema.Plan, timeout time.Duration, maxAttempts int) (sandbox.Result, []Attempt, error) {
	return retryRunPlan(plan, timeout, maxAttempts, GenerateMain, WriteTempMain, sandbox.RunMain)
}

func retryRunPlan(
	plan schema.Plan,
	timeout time.Duration,
	maxAttempts int,
	generate func(schema.Plan, []Attempt) (string, error),
	write func(string) (string, error),
	run func(string, time.Duration) (sandbox.Result, error),
) (sandbox.Result, []Attempt, error) {
	if maxAttempts < 1 {
		return sandbox.Result{}, nil, errors.New("maxAttempts must be at least 1")
	}

	var history []Attempt
	var lastFingerprint string
	for i := 0; i < maxAttempts; i++ {
		attempt := Attempt{Number: i + 1}
		code, err := generate(plan, history)
		if err != nil {
			attempt.Stage = "generate"
			attempt.Error = err.Error()
			history = append(history, attempt)
			if attempt.fingerprint() == lastFingerprint {
				return sandbox.Result{}, history, ErrRepeatedError
			}
			lastFingerprint = attempt.fingerprint()
			continue
		}

		mainPath, err := write(code)
		if err != nil {
			attempt.Stage = "write"
			attempt.Error = err.Error()
			history = append(history, attempt)
			if attempt.fingerprint() == lastFingerprint {
				return sandbox.Result{}, history, ErrRepeatedError
			}
			lastFingerprint = attempt.fingerprint()
			continue
		}

		res, err := run(mainPath, timeout)
		if err == nil {
			attempt.Stage = "run"
			attempt.Stdout = res.Stdout
			attempt.Stderr = res.Stderr
			attempt.ExitCode = res.ExitCode
			history = append(history, attempt)
			return res, history, nil
		}

		attempt.Stage = classifyRunError(res, err)
		attempt.Error = err.Error()
		attempt.Stdout = res.Stdout
		attempt.Stderr = res.Stderr
		attempt.ExitCode = res.ExitCode
		history = append(history, attempt)
		if attempt.fingerprint() == lastFingerprint {
			return res, history, ErrRepeatedError
		}
		lastFingerprint = attempt.fingerprint()
	}

	return sandbox.Result{}, history, ErrMaxAttempts
}

func classifyRunError(res sandbox.Result, err error) string {
	msg := err.Error()
	if strings.Contains(strings.ToLower(res.Stderr), "undefined") || strings.Contains(strings.ToLower(msg), "compile") {
		return "compile"
	}
	return "runtime"
}

// Attempt records one retry step for traceability.
type Attempt struct {
	Number   int
	Stage    string
	Error    string
	Stdout   string
	Stderr   string
	ExitCode int
}

func (a Attempt) fingerprint() string {
	return a.Stage + ":" + strings.TrimSpace(a.Error) + ":" + strings.TrimSpace(a.Stderr)
}
