package schema

import (
	"encoding/json"
	"testing"
)

func TestInferPlan_FromJSONSummary(t *testing.T) {
	summary := Summary{Format: "json", HasHeader: true, Columns: []string{"name", "age"}}

	plan, err := InferPlan("convert this json to csv", summary, "csv")
	if err != nil {
		t.Fatalf("InferPlan returned error: %v", err)
	}
	if plan.Operation != "convert" {
		t.Fatalf("expected convert operation, got %q", plan.Operation)
	}
	if plan.Source != "json" || plan.Target != "csv" {
		t.Fatalf("unexpected endpoints: %#v", plan)
	}
	if len(plan.Steps) == 0 {
		t.Fatal("expected steps to be generated")
	}
}

func TestInferPlan_FromCSVSummary(t *testing.T) {
	summary := Summary{Format: "csv", HasHeader: true, Columns: []string{"name", "age"}}

	plan, err := InferPlan("normalize this csv", summary, "json")
	if err != nil {
		t.Fatalf("InferPlan returned error: %v", err)
	}
	if plan.Operation != "normalize" {
		t.Fatalf("expected normalize operation, got %q", plan.Operation)
	}
}

func TestInferPlan_AmbiguousPrompt(t *testing.T) {
	summary := Summary{Format: "json"}
	_, err := InferPlan("do something useful", summary, "csv")
	if err == nil {
		t.Fatal("expected error for ambiguous prompt")
	}
}

func TestPlanJSONSerializable(t *testing.T) {
	plan := Plan{
		Objective: "convert json to csv",
		Source:    "json",
		Target:    "csv",
		Operation: "convert",
		Steps:     []string{"inspect json", "convert into csv"},
		Summary:   Summary{Format: "json", HasHeader: true, Columns: []string{"name"}},
	}

	data, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("json marshal: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected json output")
	}
}

func TestInferPlan_EmptySummary(t *testing.T) {
	_, err := InferPlan("convert json to csv", Summary{}, "csv")
	if err == nil {
		t.Fatal("expected error for empty summary")
	}
}
