package codeact

import (
	"errors"
	"strings"
	"time"

	"morphgo/internal/sandbox"
	"morphgo/internal/schema"
)

// PrepareTempMain generates and writes the temporary main.go for a run.
func PrepareTempMain(plan schema.Plan) (string, error) {
	code, err := GenerateMain(plan)
	if err != nil {
		return "", err
	}
	return WriteTempMain(code)
}

// RunPlan prepares and executes the generated program through the sandbox.
func RunPlan(plan schema.Plan, timeout time.Duration) (sandbox.Result, error) {
	mainPath, err := PrepareTempMain(plan)
	if err != nil {
		return sandbox.Result{}, err
	}
	return sandbox.RunMain(mainPath, timeout)
}

// RetryRunPlan retries generation and sandbox execution a limited number of times.
func RetryRunPlan(plan schema.Plan, timeout time.Duration, maxAttempts int) (sandbox.Result, []string, error) {
	return retryRunPlan(plan, timeout, maxAttempts, GenerateMain, WriteTempMain, sandbox.RunMain)
}

func retryRunPlan(
	plan schema.Plan,
	timeout time.Duration,
	maxAttempts int,
	generate func(schema.Plan) (string, error),
	write func(string) (string, error),
	run func(string, time.Duration) (sandbox.Result, error),
) (sandbox.Result, []string, error) {
	if maxAttempts < 1 {
		return sandbox.Result{}, nil, errors.New("maxAttempts must be at least 1")
	}

	var history []string
	var lastError string
	for i := 0; i < maxAttempts; i++ {
		code, err := generate(plan)
		if err != nil {
			msg := "generate: " + err.Error()
			history = append(history, msg)
			if msg == lastError {
				return sandbox.Result{}, history, errors.New("repeated error")
			}
			lastError = msg
			continue
		}

		mainPath, err := write(code)
		if err != nil {
			msg := "write: " + err.Error()
			history = append(history, msg)
			if msg == lastError {
				return sandbox.Result{}, history, errors.New("repeated error")
			}
			lastError = msg
			continue
		}

		res, err := run(mainPath, timeout)
		if err == nil {
			return res, history, nil
		}

		msg := classifyRunError(res, err)
		history = append(history, msg)
		if msg == lastError {
			return res, history, errors.New("repeated error")
		}
		lastError = msg
	}

	return sandbox.Result{}, history, errors.New("max attempts reached")
}

func classifyRunError(res sandbox.Result, err error) string {
	msg := err.Error()
	if strings.Contains(strings.ToLower(res.Stderr), "undefined") || strings.Contains(strings.ToLower(msg), "compile") {
		return "compile: " + msg
	}
	return "runtime: " + msg
}
