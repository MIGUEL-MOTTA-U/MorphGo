package codeact

import (
	"strings"
	"testing"
	"time"

	"morphgo/internal/sandbox"
	"morphgo/internal/schema"
)

func TestRunPlan_Success(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "placeholder",
	}

	res, err := RunPlan(plan, 5*time.Second)
	if err != nil {
		t.Fatalf("RunPlan returned error: %v", err)
	}
	if !strings.Contains(res.Stdout, "execute placeholder from json to csv") {
		t.Fatalf("unexpected stdout: %q", res.Stdout)
	}
}

func TestRunPlan_MissingData(t *testing.T) {
	_, err := RunPlan(schema.Plan{}, 5*time.Second)
	if err == nil {
		t.Fatal("expected error for incomplete plan")
	}
}

func TestRetryRunPlan_InvalidAttempts(t *testing.T) {
	_, _, err := RetryRunPlan(schema.Plan{}, 5*time.Second, 0)
	if err == nil {
		t.Fatal("expected error for invalid max attempts")
	}
}

func TestRetryRunPlan_RepeatedGenerationError(t *testing.T) {
	plan := schema.Plan{Objective: "convert json to csv"}
	generate := func(schema.Plan, []Attempt) (string, error) {
		return "", ErrRepeatedError
	}
	write := func(string) (string, error) {
		t.Fatal("write should not be called when generation fails")
		return "", nil
	}
	run := func(string, time.Duration) (sandbox.Result, error) {
		t.Fatal("run should not be called when generation fails")
		return sandbox.Result{}, nil
	}

	_, history, err := retryRunPlan(plan, 5*time.Second, 2, generate, write, run)
	if err != ErrRepeatedError {
		t.Fatalf("expected ErrRepeatedError, got %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("expected 2 attempts in history, got %d", len(history))
	}
	if history[0].Stage != "generate" {
		t.Fatalf("expected generate stage, got %q", history[0].Stage)
	}
	if history[1].Stage != "generate" {
		t.Fatalf("expected repeated generate stage, got %q", history[1].Stage)
	}
}
