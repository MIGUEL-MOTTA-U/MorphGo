package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRun_UnsupportedInputStatus(t *testing.T) {
	log, err := Run("input.txt", "convert json to csv", "csv")
	if err == nil {
		t.Fatal("expected error for unsupported input extension")
	}
	if log.Status != "error" {
		t.Fatalf("expected error status, got %q", log.Status)
	}
}

func TestRun_SuccessfulJsonToCsv(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	if err := os.WriteFile(inputPath, []byte(`{"employees":[{"name":"ana","age":30}]}`), 0o600); err != nil {
		t.Fatalf("write input json: %v", err)
	}

	log, err := Run(inputPath, "convert json to csv", "csv")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	if filepath.Base(log.Plan.OutputPath) != "output.csv" {
		t.Fatalf("expected default output path, got %q", log.Plan.OutputPath)
	}
}

func TestRun_SuccessfulCsvToJson(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.csv")
	csvContent := "name,age\nana,30\nbeatriz,25"
	if err := os.WriteFile(inputPath, []byte(csvContent), 0o600); err != nil {
		t.Fatalf("write input csv: %v", err)
	}

	log, err := Run(inputPath, "convert csv to json", "json")
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	if filepath.Base(log.Plan.OutputPath) != "output.json" {
		t.Fatalf("expected default output path base name, got %q", log.Plan.OutputPath)
	}
}

func TestRun_SuccessfulYamlToJson(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.yaml")
	yamlContent := "employees:\n  - name: ana\n    age: 30\n  - name: beatriz\n    age: 25"
	if err := os.WriteFile(inputPath, []byte(yamlContent), 0o600); err != nil {
		t.Fatalf("write input yaml: %v", err)
	}

	log, err := Run(inputPath, "convert yaml to json", "json")
	if err != nil {
		for _, att := range log.Attempts {
			t.Logf("Attempt %d [%s] error: %s, stderr: %s", att.Number, att.Stage, att.Error, att.Stderr)
		}
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	if filepath.Base(log.Plan.OutputPath) != "output.json" {
		t.Fatalf("expected default output path base name, got %q", log.Plan.OutputPath)
	}
}

func TestRun_SuccessfulXmlToJson(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.xml")
	xmlContent := `<company><employee name="ana">30</employee><employee name="beatriz">25</employee></company>`
	if err := os.WriteFile(inputPath, []byte(xmlContent), 0o600); err != nil {
		t.Fatalf("write input xml: %v", err)
	}

	log, err := Run(inputPath, "convert xml to json", "json")
	if err != nil {
		for _, att := range log.Attempts {
			t.Logf("Attempt %d [%s] error: %s, stderr: %s", att.Number, att.Stage, att.Error, att.Stderr)
		}
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	if filepath.Base(log.Plan.OutputPath) != "output.json" {
		t.Fatalf("expected default output path base name, got %q", log.Plan.OutputPath)
	}
}

func TestRun_SuccessfulJsonToYaml(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	jsonContent := `{"employees":[{"name":"ana","age":30},{"name":"beatriz","age":25}]}`
	if err := os.WriteFile(inputPath, []byte(jsonContent), 0o600); err != nil {
		t.Fatalf("write input json: %v", err)
	}

	log, err := Run(inputPath, "convert json to yaml", "yaml")
	if err != nil {
		for _, att := range log.Attempts {
			t.Logf("Attempt %d [%s] error: %s, stderr: %s", att.Number, att.Stage, att.Error, att.Stderr)
		}
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	if filepath.Base(log.Plan.OutputPath) != "output.yaml" {
		t.Fatalf("expected default output path base name, got %q", log.Plan.OutputPath)
	}
}

func TestRun_SuccessfulJsonToXml(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	jsonContent := `{"company":"MorphGo","employees":[{"name":"ana","age":30}]}`
	if err := os.WriteFile(inputPath, []byte(jsonContent), 0o600); err != nil {
		t.Fatalf("write input json: %v", err)
	}

	log, err := Run(inputPath, "convert json to xml", "xml")
	if err != nil {
		for _, att := range log.Attempts {
			t.Logf("Attempt %d [%s] error: %s, stderr: %s", att.Number, att.Stage, att.Error, att.Stderr)
		}
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	if filepath.Base(log.Plan.OutputPath) != "output.xml" {
		t.Fatalf("expected default output path base name, got %q", log.Plan.OutputPath)
	}
}

// TestRun_JsonToYaml_DataIntegrity verifies that the data in the YAML output matches the input JSON data.
func TestRun_JsonToYaml_DataIntegrity(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.json")
	// Complex JSON with nested structures, null, empty arrays and objects.
	jsonContent := `{
		"name": "test",
		"values": [1, 2, 3],
		"nested": {
			"inner": "value",
			"list": ["a", "b"]
		},
		"null_field": null,
		"empty_list": [],
		"empty_object": {}
	}`
	if err := os.WriteFile(inputPath, []byte(jsonContent), 0o600); err != nil {
		t.Fatalf("write input json: %v", err)
	}

	log, err := Run(inputPath, "convert json to yaml", "yaml")
	if err != nil {
		for _, att := range log.Attempts {
			t.Logf("Attempt %d [%s] error: %s, stderr: %s", att.Number, att.Stage, att.Error, att.Stderr)
		}
		t.Fatalf("Run returned error: %v", err)
	}
	if log.Status != "success" {
		t.Fatalf("expected success status, got %q", log.Status)
	}
	// The output path should be set by the plan.
	outputPath := log.Plan.OutputPath
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(dir, outputPath)
	}

	// Read the generated YAML file.
	yamlData, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output yaml: %v", err)
	}

	// Parse the YAML into a map[string]any.
	var yamlMap map[string]any
	if err := yaml.Unmarshal(yamlData, &yamlMap); err != nil {
		t.Fatalf("parse yaml: %v", err)
	}

	// Parse the original JSON into a map[string]any for comparison.
	var jsonMap map[string]any
	if err := json.Unmarshal([]byte(jsonContent), &jsonMap); err != nil {
		t.Fatalf("parse original json: %v", err)
	}

	// Compare the two maps.
	if !deepEqualRobust(jsonMap, yamlMap) {
		t.Errorf("data mismatch\noriginal: %#v\noutput: %#v", jsonMap, yamlMap)
	}
}