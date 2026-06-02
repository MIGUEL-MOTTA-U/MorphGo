package codeact

import (
	"strings"
	"testing"
	"time"

	"morphgo/internal/schema"
)

func TestRunPlan_Success(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
	}

	res, err := RunPlan(plan, 5*time.Second)
	if err != nil {
		t.Fatalf("RunPlan returned error: %v", err)
	}
	if !strings.Contains(res.Stdout, "execute convert from json to csv") {
		t.Fatalf("unexpected stdout: %q", res.Stdout)
	}
}

func TestRunPlan_MissingData(t *testing.T) {
	_, err := RunPlan(schema.Plan{}, 5*time.Second)
	if err == nil {
		t.Fatal("expected error for incomplete plan")
	}
}
