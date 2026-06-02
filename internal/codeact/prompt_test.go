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
	if !strings.Contains(code, "fmt.Println") {
		t.Fatal("expected fmt.Println call")
	}
}

func TestGenerateMain_MissingData(t *testing.T) {
	_, err := GenerateMain(schema.Plan{})
	if err == nil {
		t.Fatal("expected error for incomplete plan")
	}
}
