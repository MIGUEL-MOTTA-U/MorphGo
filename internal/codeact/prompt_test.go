package codeact

import (
	"strings"
	"testing"

	"morphgo/internal/schema"
)

func TestGenerateMain_KnownPlan(t *testing.T) {
	plan := schema.Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
	}

	code, err := GenerateMain(plan)
	if err != nil {
		t.Fatalf("GenerateMain returned error: %v", err)
	}
	if code == "" {
		t.Fatal("expected generated code")
	}
	if !strings.Contains(code, "package main") {
		t.Fatal("expected package main")
	}
	if !strings.Contains(code, "\"errors\"") || !strings.Contains(code, "\"fmt\"") {
		t.Fatal("expected valid imports")
	}
	if !strings.Contains(code, "func run() error") {
		t.Fatal("expected run helper")
	}
	if !strings.Contains(code, "if err := run(); err != nil") {
		t.Fatal("expected basic error handling")
	}
}

func TestGenerateMain_MissingData(t *testing.T) {
	_, err := GenerateMain(schema.Plan{})
	if err == nil {
		t.Fatal("expected error for incomplete plan")
	}
}
