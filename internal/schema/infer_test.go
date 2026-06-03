package schema

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestInferPlan_FromJSONSummary(t *testing.T) {
	summary := Summary{Format: "json", HasHeader: true, Columns: []string{"name", "age"}}

	plan, err := InferPlan("convert this json to csv", summary, "csv", "in.json", "out.csv")
	if err != nil {
		t.Fatalf("InferPlan returned error: %v", err)
	}
	if plan.Operation != "convert" {
		t.Fatalf("expected convert operation, got %q", plan.Operation)
	}
	if plan.Source != "json" || plan.Target != "csv" {
		t.Fatalf("unexpected endpoints: %#v", plan)
	}
	if plan.InputPath != "in.json" || plan.OutputPath != "out.csv" {
		t.Fatalf("unexpected paths: %#v", plan)
	}
	if len(plan.Steps) == 0 {
		t.Fatal("expected steps to be generated")
	}
}

func TestInferPlan_FromSpanishConvertPrompt(t *testing.T) {
	summary := Summary{Format: "json", HasHeader: true, Columns: []string{"name", "age"}}

	plan, err := InferPlan("convierte este json a csv", summary, "csv", "in.json", "out.csv")
	if err != nil {
		t.Fatalf("InferPlan returned error: %v", err)
	}
	if plan.Operation != "convert" {
		t.Fatalf("expected convert operation, got %q", plan.Operation)
	}
}

func TestInferPlan_FromCSVSummary(t *testing.T) {
	summary := Summary{Format: "csv", HasHeader: true, Columns: []string{"name", "age"}}

	plan, err := InferPlan("normalize this csv", summary, "json", "", "")
	if err != nil {
		t.Fatalf("InferPlan returned error: %v", err)
	}
	if plan.Operation != "normalize" {
		t.Fatalf("expected normalize operation, got %q", plan.Operation)
	}
}

func TestInferPlan_AmbiguousPrompt(t *testing.T) {
	summary := Summary{Format: "json"}
	_, err := InferPlan("do something useful", summary, "csv", "", "")
	if err == nil {
		t.Fatal("expected error for ambiguous prompt")
	}
}

func TestInferPlan_MultiOperationPrompt(t *testing.T) {
	summary := Summary{Format: "json"}
	_, err := InferPlan("convert and filter this json", summary, "csv", "", "")
	if err == nil {
		t.Fatal("expected error for multi-operation prompt")
	}
}

func TestPlanJSONSerializable(t *testing.T) {
	plan := Plan{
		Objective:  "convert json to csv",
		Source:     "json",
		Target:     "csv",
		Operation:  "convert",
		InputPath:  "in.json",
		OutputPath: "out.csv",
		Steps:      []string{"inspect json", "convert into csv"},
		Summary:    Summary{Format: "json", HasHeader: true, Columns: []string{"name"}},
	}

	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected json output")
	}
	if !strings.Contains(string(data), "\"input_path\":\"in.json\"") || !strings.Contains(string(data), "\"output_path\":\"out.csv\"") {
		t.Fatalf("expected paths in json output, got %s", string(data))
	}
}

func TestPlanMarkdownSerializable(t *testing.T) {
	plan := Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
		Steps:     []string{"inspect json", "convert into csv"},
	}

	md := plan.MarshalMarkdown()
	if md == "" {
		t.Fatal("expected markdown output")
	}
	if !containsAll(md, []string{"Objective:", "Source:", "Target:", "Operation:"}) {
		t.Fatalf("unexpected markdown output: %q", md)
	}
}

func TestInferPlan_EmptySummary(t *testing.T) {
	_, err := InferPlan("convert json to csv", Summary{}, "csv", "", "")
	if err == nil {
		t.Fatal("expected error for empty summary")
	}
}

func TestInferPlan_BlankPrompt(t *testing.T) {
	summary := Summary{Format: "json"}
	_, err := InferPlan("   ", summary, "csv", "", "")
	if err == nil {
		t.Fatal("expected error for blank prompt")
	}
}

func TestPlanEmptyStructure(t *testing.T) {
	var plan Plan

	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("json marshal empty plan: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected json output for empty plan")
	}

	if md := plan.MarshalMarkdown(); md == "" {
		t.Fatal("expected markdown output for empty plan")
	}
}

func containsAll(s string, parts []string) bool {
	for _, part := range parts {
		if !strings.Contains(s, part) {
			return false
		}
	}
	return true
}
