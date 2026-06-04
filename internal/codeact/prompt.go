package codeact

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"morphgo/internal/schema"
)

// GenerateMain renders a minimal temporary Go program for a given plan.
func GenerateMain(plan schema.Plan, history []Attempt) (string, error) {
	if plan.Objective == "" {
		return "", errors.New("plan objective is required")
	}
	if plan.Source == "" || plan.Target == "" || plan.Operation == "" {
		return "", errors.New("plan is incomplete")
	}

	if plan.Source == "json" && plan.Target == "csv" && plan.Operation == "convert" {
		if strings.TrimSpace(plan.InputPath) == "" || strings.TrimSpace(plan.OutputPath) == "" {
			return "", errors.New("json to csv plan requires input and output paths")
		}
		return generateJsonToCsv(plan)
	}

	if plan.Source == "csv" && plan.Target == "json" && plan.Operation == "convert" {
		if strings.TrimSpace(plan.InputPath) == "" || strings.TrimSpace(plan.OutputPath) == "" {
			return "", errors.New("csv to json plan requires input and output paths")
		}
		return generateCsvToJson(plan)
	}

	if plan.Source == "yaml" && plan.Target == "json" && plan.Operation == "convert" {
		if strings.TrimSpace(plan.InputPath) == "" || strings.TrimSpace(plan.OutputPath) == "" {
			return "", errors.New("yaml to json plan requires input and output paths")
		}
		return generateYamlToJson(plan)
	}

	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import \"fmt\"\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tfmt.Println(%q)\n", "execute "+plan.Operation+" from "+plan.Source+" to "+plan.Target)
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")
	b.WriteString("func main() {\n")
	b.WriteString("\tif err := run(); err != nil {\n")
	b.WriteString("\t\tfmt.Println(err)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")
	return b.String(), nil
}

func generateJsonToCsv(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)

	b.WriteString("\tdata, err := os.ReadFile(inputPath)\n")
	b.WriteString("\tif err != nil { return err }\n\n")

	b.WriteString("\tvar root map[string]interface{}\n")
	b.WriteString("\tif err := json.Unmarshal(data, &root); err != nil { return err }\n\n")

	b.WriteString("\tvar items []interface{}\n")
	b.WriteString("\tif list, ok := root[\"employees\"].([]interface{}); ok {\n")
	b.WriteString("\t\titems = list\n")
	b.WriteString("\t} else {\n")
	b.WriteString("\t\tfor _, v := range root {\n")
	b.WriteString("\t\t\tif list, ok := v.([]interface{}); ok {\n")
	b.WriteString("\t\t\t\titems = list; break\n")
	b.WriteString("\t\t\t}\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tif len(items) == 0 { return fmt.Errorf(\"no items found to convert\") }\n\n")

	b.WriteString("\tf, err := os.Create(outputPath)\n")
	b.WriteString("\tif err != nil { return err }\n")
	b.WriteString("\tdefer f.Close()\n\n")

	b.WriteString("\tw := csv.NewWriter(f)\n")
	b.WriteString("\tdefer w.Flush()\n\n")

	if len(plan.Summary.Columns) > 0 {
		b.WriteString("\theaders := []string{")
		for i, col := range plan.Summary.Columns {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%q", col)
		}
		b.WriteString("}\n")
		b.WriteString("\tif err := w.Write(headers); err != nil { return err }\n\n")

		b.WriteString("\tfor _, item := range items {\n")
		b.WriteString("\t\tm, ok := item.(map[string]interface{})\n")
		b.WriteString("\t\tif !ok { continue }\n")
		b.WriteString("\t\trow := make([]string, len(headers))\n")
		b.WriteString("\t\tfor i, h := range headers {\n")
		b.WriteString("\t\t\tval := m[h]\n")
		b.WriteString("\t\t\tif val == nil { val = \"\" }\n")
		b.WriteString("\t\t\trow[i] = fmt.Sprintf(\"%v\", val)\n")
		b.WriteString("\t\t}\n")
		b.WriteString("\t\tif err := w.Write(row); err != nil { return err }\n")
		b.WriteString("\t}\n")
	}

	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	b.WriteString("func main() {\n")
	b.WriteString("\tif err := run(); err != nil {\n")
	b.WriteString("\t\tfmt.Fprintf(os.Stderr, \"error: %v\\n\", err)\n")
	b.WriteString("\t\tos.Exit(1)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	return b.String(), nil
}

func generateCsvToJson(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)

	b.WriteString("\tf, err := os.Open(inputPath)\n")
	b.WriteString("\tif err != nil { return err }\n")
	b.WriteString("\tdefer f.Close()\n\n")

	b.WriteString("\tr := csv.NewReader(f)\n")
	if plan.Summary.Separator != 0 && plan.Summary.Separator != ',' {
		fmt.Fprintf(&b, "\tr.Comma = %q\n", string(plan.Summary.Separator))
	}

	b.WriteString("\trecords, err := r.ReadAll()\n")
	b.WriteString("\tif err != nil { return err }\n\n")

	b.WriteString("\tif len(records) == 0 { return fmt.Errorf(\"csv file is empty\") }\n\n")

	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tstartIndex := 0\n")
	if plan.Summary.HasHeader {
		b.WriteString("\tif len(records) > 0 {\n")
		b.WriteString("\t\theaders = records[0]\n")
		b.WriteString("\t\tstartIndex = 1\n")
		b.WriteString("\t}\n")
	} else if len(plan.Summary.Columns) > 0 {
		b.WriteString("\theaders = []string{")
		for i, col := range plan.Summary.Columns {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%q", col)
		}
		b.WriteString("}\n")
	} else {
		b.WriteString("\tif len(records) > 0 {\n")
		b.WriteString("\t\theaders = make([]string, len(records[0]))\n")
		b.WriteString("\t\tfor i := range headers { headers[i] = fmt.Sprintf(\"column_%d\", i) }\n")
		b.WriteString("\t}\n")
	}

	b.WriteString("\tvar data []map[string]interface{}\n")
	b.WriteString("\tfor i := startIndex; i < len(records); i++ {\n")
	b.WriteString("\t\trow := records[i]\n")
	b.WriteString("\t\tobj := make(map[string]interface{})\n")
	b.WriteString("\t\tfor j, val := range row {\n")
	b.WriteString("\t\t\tif j < len(headers) {\n")
	b.WriteString("\t\t\t\tobj[headers[j]] = val\n")
	b.WriteString("\t\t\t}\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t\tdata = append(data, obj)\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\tjsonData, err := json.MarshalIndent(data, \"\", \"  \")\n")
	b.WriteString("\tif err != nil { return err }\n\n")

	b.WriteString("\tif err := os.WriteFile(outputPath, jsonData, 0644); err != nil { return err }\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	b.WriteString("func main() {\n")
	b.WriteString("\tif err := run(); err != nil {\n")
	b.WriteString("\t\tfmt.Fprintf(os.Stderr, \"error: %v\\n\", err)\n")
	b.WriteString("\t\tos.Exit(1)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	return b.String(), nil
}

// WriteTempMain stores generated code in a temporary main.go file for the current run.
func WriteTempMain(code string) (string, error) {
	if strings.TrimSpace(code) == "" {
		return "", errors.New("generated code is required")
	}

	dir, err := os.MkdirTemp("", "morphgo-*")
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, "main.go")
	if err := os.WriteFile(path, []byte(code), 0o600); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module morphgo-temp\ngo 1.26.3\n\nrequire gopkg.in/yaml.v3 v3.0.1\n"), 0o600); err != nil {
		return "", err
	}

	// Copy project go.sum to ensure dependencies can be resolved from cache
	sumData, err := os.ReadFile("go.sum")
	if err == nil {
		_ = os.WriteFile(filepath.Join(dir, "go.sum"), sumData, 0o600)
	}

	return path, nil
}

func generateYamlToJson(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)

	b.WriteString("\tdata, err := os.ReadFile(inputPath)\n")
	b.WriteString("\tif err != nil { return err }\n\n")

	b.WriteString("\tvar value any\n")
	b.WriteString("\tif err := yaml.Unmarshal(data, &value); err != nil { return err }\n\n")

	b.WriteString("\tjsonData, err := json.MarshalIndent(value, \"\", \"  \")\n")
	b.WriteString("\tif err != nil { return err }\n\n")

	b.WriteString("\tif err := os.WriteFile(outputPath, jsonData, 0644); err != nil { return err }\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")

	b.WriteString("func main() {\n")
	b.WriteString("\tif err := run(); err != nil {\n")
	b.WriteString("\t\tfmt.Fprintf(os.Stderr, \"error: %v\\n\", err)\n")
	b.WriteString("\t\tos.Exit(1)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")

	return b.String(), nil
}
