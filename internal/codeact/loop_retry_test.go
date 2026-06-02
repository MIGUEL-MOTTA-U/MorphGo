package codeact

import (
	"errors"
	"testing"
	"time"

	"morphgo/internal/sandbox"
	"morphgo/internal/schema"
)

func TestRetryRunPlan_CompileCorrection(t *testing.T) {
	plan := schema.Plan{Objective: "x", Source: "a", Target: "b", Operation: "convert"}
	attempts := 0

	res, history, err := retryRunPlan(
		plan,
		time.Second,
		3,
		func(schema.Plan) (string, error) {
			attempts++
			if attempts == 1 {
				return "bad code", nil
			}
			return "good code", nil
		},
		func(code string) (string, error) { return code, nil },
		func(mainPath string, timeout time.Duration) (sandbox.Result, error) {
			if mainPath == "bad code" {
				return sandbox.Result{Stderr: "compile: undefined"}, errors.New("compile failed")
			}
			return sandbox.Result{Stdout: "ok"}, nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 1 {
		t.Fatalf("expected 1 failed attempt, got %#v", history)
	}
	if res.Stdout != "ok" {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestRetryRunPlan_AbandonsAfterNAttempts(t *testing.T) {
	plan := schema.Plan{Objective: "x", Source: "a", Target: "b", Operation: "convert"}
	_, history, err := retryRunPlan(
		plan,
		time.Second,
		2,
		func(schema.Plan) (string, error) { return "code", nil },
		func(code string) (string, error) { return code, nil },
		func(mainPath string, timeout time.Duration) (sandbox.Result, error) {
			return sandbox.Result{Stderr: "runtime panic"}, errors.New("runtime failed")
		},
	)
	if err == nil {
		t.Fatal("expected error after max attempts")
	}
	if len(history) != 2 {
		t.Fatalf("expected 2 attempts, got %#v", history)
	}
}

func TestRetryRunPlan_RepeatedErrorStopsEarly(t *testing.T) {
	plan := schema.Plan{Objective: "x", Source: "a", Target: "b", Operation: "convert"}
	calls := 0

	_, history, err := retryRunPlan(
		plan,
		time.Second,
		5,
		func(schema.Plan) (string, error) { return "code", nil },
		func(code string) (string, error) { return code, nil },
		func(mainPath string, timeout time.Duration) (sandbox.Result, error) {
			calls++
			return sandbox.Result{Stderr: "runtime panic"}, errors.New("runtime failed")
		},
	)
	if err == nil {
		t.Fatal("expected repeated error")
	}
	if len(history) != 2 {
		t.Fatalf("expected two repeated attempts, got %#v", history)
	}
	if calls != 2 {
		t.Fatalf("expected early stop after repeated error, got %d calls", calls)
	}
}
