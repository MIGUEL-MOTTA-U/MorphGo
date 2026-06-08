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

	// Matrix of conversions
	if plan.Source == plan.Target {
		return generateIdentity(plan)
	}

	switch plan.Source {
	case "csv":
		switch plan.Target {
		case "json": return generateCsvToJson(plan)
		case "xml":  return generateCsvToXml(plan)
		case "yaml": return generateCsvToYaml(plan)
		case "markdown": return generateCsvToMarkdown(plan)
		case "excel": return generateCsvToExcel(plan)
		case "pdf": return generateCsvToPdf(plan)
		}
	case "json":
		switch plan.Target {
		case "csv": return generateJsonToCsv(plan)
		case "xml": return generateJsonToXml(plan)
		case "yaml": return generateJsonToYaml(plan)
		case "markdown": return generateJsonToMarkdown(plan)
		case "excel": return generateJsonToExcel(plan)
		case "pdf": return generateJsonToPdf(plan)
		}
	case "xml":
		switch plan.Target {
		case "json": return generateXmlToJson(plan)
		case "csv":  return generateXmlToCsv(plan)
		case "yaml": return generateXmlToYaml(plan)
		case "markdown": return generateXmlToMarkdown(plan)
		case "excel": return generateXmlToExcel(plan)
		case "pdf": return generateXmlToPdf(plan)
		}
	case "yaml":
		switch plan.Target {
		case "json": return generateYamlToJson(plan)
		case "csv":  return generateYamlToCsv(plan)
		case "xml":  return generateYamlToXml(plan)
		case "markdown": return generateYamlToMarkdown(plan)
		case "excel": return generateYamlToExcel(plan)
		case "pdf": return generateYamlToPdf(plan)
		}
	case "excel":
		switch plan.Target {
		case "json": return generateExcelToJson(plan)
		case "csv":  return generateExcelToCsv(plan)
		case "pdf":  return generateExcelToPdf(plan)
		}
	case "markdown":
		switch plan.Target {
		case "json": return generateMarkdownToJson(plan)
		case "csv":  return generateMarkdownToCsv(plan)
		case "xml":  return generateMarkdownToXml(plan)
		case "yaml": return generateMarkdownToYaml(plan)
		case "excel": return generateMarkdownToExcel(plan)
		case "pdf": return generateMarkdownToPdf(plan)
		}
	case "pdf":
		// Identity handled above, others fall through
	}

	// Fallback/Generic
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

// --- CSV Generators ---

func generateIdentity(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"io\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tsrc, err := os.Open(%q)\n", plan.InputPath)
	b.WriteString("\tif err != nil { return err }\n")
	b.WriteString("\tdefer src.Close()\n")
	fmt.Fprintf(&b, "\tdst, err := os.Create(%q)\n", plan.OutputPath)
	b.WriteString("\tif err != nil { return err }\n")
	b.WriteString("\tdefer dst.Close()\n")
	b.WriteString("\t_, err = io.Copy(dst, src)\n")
	b.WriteString("\treturn err\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
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
			if i > 0 { b.WriteString(", ") }
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

func generateCsvToXml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
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

	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tstartIndex := 0\n")
	b.WriteString("\tif len(records) > 0 {\n")
	b.WriteString("\t\theaders = records[0]\n")
	b.WriteString("\t\tstartIndex = 1\n")
	b.WriteString("\t}\n\n")

	b.WriteString("\txmlContent := \"<?xml version=\\\"1.0\\\" encoding=\\\"UTF-8\\\"?>\\n<root>\\n\"\n")
	b.WriteString("\tfor i := startIndex; i < len(records); i++ {\n")
	b.WriteString("\t\txmlContent += \"  <item>\\n\"\n")
	b.WriteString("\t\tfor j, val := range records[i] {\n")
	b.WriteString("\t\t\tkey := fmt.Sprintf(\"col_%d\", j)\n")
	b.WriteString("\t\t\tif j < len(headers) { key = headers[j] }\n")
	b.WriteString("\t\t\txmlContent += fmt.Sprintf(\"    <%s>%v</%s>\\n\", key, val, key)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t\txmlContent += \"  </item>\\n\"\n")
	b.WriteString("\t}\n")
	b.WriteString("\txmlContent += \"</root>\\n\"\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(xmlContent), 0644)\n")
	b.WriteString("}\n\n")

	b.WriteString("func main() {\n")
	b.WriteString("\tif err := run(); err != nil { os.Exit(1) }\n")
	b.WriteString("}\n")
	return b.String(), nil
}

func generateCsvToYaml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)

	b.WriteString("\tf, err := os.Open(inputPath)\n")
	b.WriteString("\tif err != nil { return err }\n")
	b.WriteString("\tdefer f.Close()\n\n")

	b.WriteString("\tr := csv.NewReader(f)\n")
	b.WriteString("\trecords, _ := r.ReadAll()\n")
	b.WriteString("\tvar data []map[string]string\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tstart := 0\n")
	b.WriteString("\tif len(records) > 0 { headers = records[0]; start = 1 }\n")
	b.WriteString("\tfor i := start; i < len(records); i++ {\n")
	b.WriteString("\t\tm := make(map[string]string)\n")
	b.WriteString("\t\tfor j, v := range records[i] {\n")
	b.WriteString("\t\t\tif j < len(headers) { m[headers[j]] = v }\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t\tdata = append(data, m)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tout, _ := yaml.Marshal(data)\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateCsvToMarkdown(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tr := csv.NewReader(f)\n")
	b.WriteString("\trecords, _ := r.ReadAll()\n")
	b.WriteString("\tvar md strings.Builder\n")
	b.WriteString("\tfor i, row := range records {\n")
	b.WriteString("\t\tmd.WriteString(\"| \" + strings.Join(row, \" | \") + \" |\\n\")\n")
	b.WriteString("\t\tif i == 0 {\n")
	b.WriteString("\t\t\tmd.WriteString(\"| \")\n")
	b.WriteString("\t\t\tfor j := range row {\n")
	b.WriteString("\t\t\t\tmd.WriteString(\"---\")\n")
	b.WriteString("\t\t\t\tif j < len(row)-1 { md.WriteString(\" | \") }\n")
	b.WriteString("\t\t\t}\n")
	b.WriteString("\t\t\tmd.WriteString(\" |\\n\")\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(md.String()), 0644)\n")
	b.WriteString("}\n\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateCsvToExcel(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString(")\n\n")

	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tr := csv.NewReader(f)\n")
	b.WriteString("\trecords, _ := r.ReadAll()\n")
	b.WriteString("\tx := excelize.NewFile()\n")
	b.WriteString("\tfor i, row := range records {\n")
	b.WriteString("\t\tfor j, v := range row {\n")
	b.WriteString("\t\t\tc, _ := excelize.CoordinatesToCellName(j+1, i+1)\n")
	b.WriteString("\t\t\tx.SetCellValue(\"Sheet1\", c, v)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn x.SaveAs(outputPath)\n")
	b.WriteString("}\n\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateCsvToPdf(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"github.com/jung-kurt/gofpdf\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tr := csv.NewReader(f)\n")
	b.WriteString("\trecords, _ := r.ReadAll()\n")
	b.WriteString("\tpdf := gofpdf.New(\"L\", \"mm\", \"A4\", \"\")\n")
	b.WriteString("\tpdf.AddPage(); pdf.SetFont(\"Arial\", \"B\", 10)\n")
	b.WriteString("\tfor _, row := range records {\n")
	b.WriteString("\t\tfor _, col := range row { pdf.Cell(30, 10, col) }\n")
	b.WriteString("\t\tpdf.Ln(-1)\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn pdf.OutputFileAndClose(outputPath)\n")
	b.WriteString("}\n\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

// --- JSON Generators ---

func generateJsonToCsv(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar raw interface{}\n")
	b.WriteString("\tjson.Unmarshal(data, &raw)\n")
	b.WriteString("\titems := extractList(raw)\n")
	b.WriteString("\tif len(items) == 0 { return fmt.Errorf(\"no list found\") }\n")
	b.WriteString("\tf, _ := os.Create(outputPath)\n")
	b.WriteString("\tw := csv.NewWriter(f)\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor k := range first { headers = append(headers, k) }\n")
	b.WriteString("\t\tw.Write(headers)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\trow := make([]string, len(headers))\n")
	b.WriteString("\t\t\tfor i, h := range headers { row[i] = fmt.Sprintf(\"%v\", m[h]) }\n")
	b.WriteString("\t\t\tw.Write(row)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tw.Flush(); return nil\n")
	b.WriteString("}\n")
	b.WriteString("func extractList(v interface{}) []interface{} {\n")
	b.WriteString("\tswitch val := v.(type) {\n")
	b.WriteString("\tcase []interface{}: return val\n")
	b.WriteString("\tcase map[string]interface{}:\n")
	b.WriteString("\t\tfor _, sub := range val {\n")
	b.WriteString("\t\t\tif list, ok := sub.([]interface{}); ok { return list }\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateJsonToXml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func marshalXML(v interface{}, name string) string {\n")
	b.WriteString("\tswitch val := v.(type) {\n")
	b.WriteString("\tcase map[string]interface{}:\n")
	b.WriteString("\t\tres := \"<\" + name + \">\\n\"\n")
	b.WriteString("\t\tfor k, v := range val { res += marshalXML(v, k) }\n")
	b.WriteString("\t\treturn res + \"</\" + name + \">\\n\"\n")
	b.WriteString("\tcase []interface{}:\n")
	b.WriteString("\t\tres := \"\"\n")
	b.WriteString("\t\tfor _, item := range val { res += marshalXML(item, name) }\n")
	b.WriteString("\t\treturn res\n")
	b.WriteString("\tdefault: return fmt.Sprintf(\"<%s>%v</%s>\\n\", name, val, name)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; json.Unmarshal(data, &v)\n")
	b.WriteString("\txml := \"<?xml version=\\\"1.0\\\"?>\\n\" + marshalXML(v, \"root\")\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(xml), 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateJsonToYaml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; json.Unmarshal(data, &v)\n")
	b.WriteString("\tout, _ := yaml.Marshal(v)\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateJsonToMarkdown(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; json.Unmarshal(data, &v)\n")
	b.WriteString("\tvar items []interface{}\n")
	b.WriteString("\tif list, ok := v.([]interface{}); ok { items = list } else {\n")
	b.WriteString("\t\tif m, ok := v.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor _, sub := range m { if l, ok := sub.([]interface{}); ok { items = l; break } }\n")
	b.WriteString("\t\t}\n\t}\n")
	b.WriteString("\tif len(items) == 0 { return nil }\n")
	b.WriteString("\tvar md strings.Builder\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor k := range first { headers = append(headers, k) }\n")
	b.WriteString("\t\tmd.WriteString(\"| \" + strings.Join(headers, \" | \") + \" |\\n\")\n")
	b.WriteString("\t\tmd.WriteString(\"| \" + strings.Repeat(\"--- | \", len(headers)) + \"\\n\")\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor _, h := range headers { md.WriteString(fmt.Sprintf(\"| %v \", m[h])) }\n")
	b.WriteString("\t\t\tmd.WriteString(\"|\\n\")\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(md.String()), 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateJsonToExcel(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; json.Unmarshal(data, &v)\n")
	b.WriteString("\tvar items []interface{}\n")
	b.WriteString("\tif l, ok := v.([]interface{}); ok { items = l } else if m, ok := v.(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor _, sub := range m { if l, ok := sub.([]interface{}); ok { items = l; break } }\n")
	b.WriteString("\t}\n")
	b.WriteString("\tx := excelize.NewFile()\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString("\t\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor k := range first {\n")
	b.WriteString("\t\t\t\theaders = append(headers, k)\n")
	b.WriteString("\t\t\t\tc, _ := excelize.CoordinatesToCellName(len(headers), 1)\n")
	b.WriteString("\t\t\t\tx.SetCellValue(\"Sheet1\", c, k)\n")
	b.WriteString("\t\t\t}\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor i, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor j, h := range headers {\n")
	b.WriteString("\t\t\t\tc, _ := excelize.CoordinatesToCellName(j+1, i+2)\n")
	b.WriteString("\t\t\t\tx.SetCellValue(\"Sheet1\", c, m[h])\n")
	b.WriteString("\t\t\t}\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn x.SaveAs(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateJsonToPdf(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"github.com/jung-kurt/gofpdf\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; json.Unmarshal(data, &v)\n")
	b.WriteString("\tvar items []interface{}\n")
	b.WriteString("\tif l, ok := v.([]interface{}); ok { items = l } else if m, ok := v.(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor _, sub := range m { if l, ok := sub.([]interface{}); ok { items = l; break } }\n")
	b.WriteString("\t}\n")
	b.WriteString("\tpdf := gofpdf.New(\"L\", \"mm\", \"A4\", \"\")\n")
	b.WriteString("\tpdf.AddPage(); pdf.SetFont(\"Arial\", \"B\", 10)\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString("\t\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor k := range first { headers = append(headers, k); pdf.Cell(30, 10, k) }\n")
	b.WriteString("\t\t\tpdf.Ln(-1)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tpdf.SetFont(\"Arial\", \"\", 10)\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor _, h := range headers { pdf.Cell(30, 10, fmt.Sprintf(\"%v\", m[h])) }\n")
	b.WriteString("\t\t\tpdf.Ln(-1)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn pdf.OutputFileAndClose(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

// --- Other Basic Generators ---

func generateXmlToJson(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"encoding/xml\"\n")
	b.WriteString("\t\"io\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func xmlToMap(d *xml.Decoder) (map[string]interface{}, error) {\n")
	b.WriteString("\tm := make(map[string]interface{})\n")
	b.WriteString("\tfor {\n")
	b.WriteString("\t\tt, err := d.Token()\n")
	b.WriteString("\t\tif err == io.EOF { break }\n")
	b.WriteString("\t\tswitch el := t.(type) {\n")
	b.WriteString("\t\tcase xml.StartElement:\n")
	b.WriteString("\t\t\tv, _ := xmlToMap(d)\n")
	b.WriteString("\t\t\tm[el.Name.Local] = v\n")
	b.WriteString("\t\tcase xml.EndElement: return m, nil\n")
	b.WriteString("\t\tcase xml.CharData: if s := strings.TrimSpace(string(el)); s != \"\" { m[\"#text\"] = s }\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn m, nil\n")
	b.WriteString("}\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\td := xml.NewDecoder(f)\n")
	b.WriteString("\tvar root interface{}\n")
	b.WriteString("\tfor {\n")
	b.WriteString("\t\tt, _ := d.Token()\n")
	b.WriteString("\t\tif start, ok := t.(xml.StartElement); ok {\n")
	b.WriteString("\t\t\tv, _ := xmlToMap(d)\n")
	b.WriteString("\t\t\troot = map[string]interface{}{start.Name.Local: v}\n")
	b.WriteString("\t\t\tbreak\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tout, _ := json.MarshalIndent(root, \"\", \"  \")\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateYamlToJson(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; yaml.Unmarshal(data, &v)\n")
	b.WriteString("\tout, _ := json.MarshalIndent(v, \"\", \"  \")\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateExcelToJson(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tx, _ := excelize.OpenFile(inputPath)\n")
	b.WriteString("\trows, _ := x.GetRows(x.GetSheetList()[0])\n")
	b.WriteString("\tvar data []map[string]string\n")
	b.WriteString("\theaders := rows[0]\n")
	b.WriteString("\tfor i := 1; i < len(rows); i++ {\n")
	b.WriteString("\t\tm := make(map[string]string)\n")
	b.WriteString("\t\tfor j, v := range rows[i] { if j < len(headers) { m[headers[j]] = v } }\n")
	b.WriteString("\t\tdata = append(data, m)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tout, _ := json.MarshalIndent(data, \"\", \"  \")\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateMarkdownToJson(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/json\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tlines := strings.Split(string(data), \"\\n\")\n")
	b.WriteString("\tvar dataOut []map[string]string\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tfor _, l := range lines {\n")
	b.WriteString("\t\tl = strings.TrimSpace(l)\n")
	b.WriteString("\t\tif !strings.HasPrefix(l, \"|\") || strings.Contains(l, \"---\") { continue }\n")
	b.WriteString("\t\tparts := strings.Split(strings.Trim(l, \"|\"), \"|\")\n")
	b.WriteString("\t\tfor i := range parts { parts[i] = strings.TrimSpace(parts[i]) }\n")
	b.WriteString("\t\tif headers == nil { headers = parts } else {\n")
	b.WriteString("\t\t\tm := make(map[string]string)\n")
	b.WriteString("\t\t\tfor i, v := range parts { if i < len(headers) { m[headers[i]] = v } }\n")
	b.WriteString("\t\t\tdataOut = append(dataOut, m)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tout, _ := json.MarshalIndent(dataOut, \"\", \"  \")\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateExcelToCsv(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tx, _ := excelize.OpenFile(inputPath)\n")
	b.WriteString("\trows, _ := x.GetRows(x.GetSheetList()[0])\n")
	b.WriteString("\tf, _ := os.Create(outputPath)\n")
	b.WriteString("\tw := csv.NewWriter(f)\n")
	b.WriteString("\tfor _, row := range rows { w.Write(row) }\n")
	b.WriteString("\tw.Flush(); return nil\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateExcelToPdf(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString("\t\"github.com/jung-kurt/gofpdf\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tx, _ := excelize.OpenFile(inputPath)\n")
	b.WriteString("\trows, _ := x.GetRows(x.GetSheetList()[0])\n")
	b.WriteString("\tpdf := gofpdf.New(\"L\", \"mm\", \"A4\", \"\")\n")
	b.WriteString("\tpdf.AddPage(); pdf.SetFont(\"Arial\", \"B\", 10)\n")
	b.WriteString("\tfor _, row := range rows {\n")
	b.WriteString("\t\tfor _, col := range row { pdf.Cell(30, 10, col) }\n")
	b.WriteString("\t\tpdf.Ln(-1)\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn pdf.OutputFileAndClose(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateXmlToCsv(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"encoding/xml\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"io\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")
	b.WriteString(xmlHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tvar root interface{}; d := xml.NewDecoder(f)\n")
	b.WriteString("\tfor { t, _ := d.Token(); if start, ok := t.(xml.StartElement); ok { v, _ := xmlToMap(d); root = map[string]interface{}{start.Name.Local: v}; break } }\n")
	b.WriteString("\titems := extractList(root)\n")
	b.WriteString("\tif len(items) == 0 { return fmt.Errorf(\"no list found\") }\n")
	b.WriteString("\tf2, _ := os.Create(outputPath)\n")
	b.WriteString("\tw := csv.NewWriter(f2)\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor k := range first { headers = append(headers, k) }\n")
	b.WriteString("\t\tw.Write(headers)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\trow := make([]string, len(headers))\n")
	b.WriteString("\t\t\tfor i, h := range headers { row[i] = fmt.Sprintf(\"%v\", m[h]) }\n")
	b.WriteString("\t\t\tw.Write(row)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tw.Flush(); return nil\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateXmlToYaml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/xml\"\n")
	b.WriteString("\t\"io\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")
	b.WriteString(xmlHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tvar root interface{}; d := xml.NewDecoder(f)\n")
	b.WriteString("\tfor { t, _ := d.Token(); if start, ok := t.(xml.StartElement); ok { v, _ := xmlToMap(d); root = map[string]interface{}{start.Name.Local: v}; break } }\n")
	b.WriteString("\tout, _ := yaml.Marshal(root)\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateXmlToMarkdown(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/xml\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"io\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")
	b.WriteString(xmlHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tvar root interface{}; d := xml.NewDecoder(f)\n")
	b.WriteString("\tfor { t, _ := d.Token(); if start, ok := t.(xml.StartElement); ok { v, _ := xmlToMap(d); root = map[string]interface{}{start.Name.Local: v}; break } }\n")
	b.WriteString("\titems := extractList(root)\n")
	b.WriteString("\tif len(items) == 0 { return nil }\n")
	b.WriteString("\tvar md strings.Builder\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor k := range first { headers = append(headers, k) }\n")
	b.WriteString("\t\tmd.WriteString(\"| \" + strings.Join(headers, \" | \") + \" |\\n\")\n")
	b.WriteString("\t\tmd.WriteString(\"| \" + strings.Repeat(\"--- | \", len(headers)) + \"\\n\")\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor _, h := range headers { md.WriteString(fmt.Sprintf(\"| %v \", m[h])) }\n")
	b.WriteString("\t\t\tmd.WriteString(\"|\\n\")\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(md.String()), 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateXmlToExcel(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/xml\"\n")
	b.WriteString("\t\"io\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString(")\n\n")
	b.WriteString(xmlHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tvar root interface{}; d := xml.NewDecoder(f)\n")
	b.WriteString("\tfor { t, _ := d.Token(); if start, ok := t.(xml.StartElement); ok { v, _ := xmlToMap(d); root = map[string]interface{}{start.Name.Local: v}; break } }\n")
	b.WriteString("\titems := extractList(root)\n")
	b.WriteString("\tx := excelize.NewFile()\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString("\t\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor k := range first { headers = append(headers, k); c, _ := excelize.CoordinatesToCellName(len(headers), 1); x.SetCellValue(\"Sheet1\", c, k) }\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor i, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor j, h := range headers { c, _ := excelize.CoordinatesToCellName(j+1, i+2); x.SetCellValue(\"Sheet1\", c, m[h]) }\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn x.SaveAs(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateXmlToPdf(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/xml\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"io\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString("\t\"github.com/jung-kurt/gofpdf\"\n")
	b.WriteString(")\n\n")
	b.WriteString(xmlHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tf, _ := os.Open(inputPath)\n")
	b.WriteString("\tvar root interface{}; d := xml.NewDecoder(f)\n")
	b.WriteString("\tfor { t, _ := d.Token(); if start, ok := t.(xml.StartElement); ok { v, _ := xmlToMap(d); root = map[string]interface{}{start.Name.Local: v}; break } }\n")
	b.WriteString("\titems := extractList(root)\n")
	b.WriteString("\tpdf := gofpdf.New(\"L\", \"mm\", \"A4\", \"\")\n")
	b.WriteString("\tpdf.AddPage(); pdf.SetFont(\"Arial\", \"B\", 10)\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString("\t\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor k := range first { headers = append(headers, k); pdf.Cell(30, 10, k) }\n")
	b.WriteString("\t\t\tpdf.Ln(-1)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tpdf.SetFont(\"Arial\", \"\", 10)\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor _, h := range headers { pdf.Cell(30, 10, fmt.Sprintf(\"%v\", m[h])) }\n")
	b.WriteString("\t\t\tpdf.Ln(-1)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn pdf.OutputFileAndClose(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateYamlToCsv(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")
	b.WriteString(jsonHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; yaml.Unmarshal(data, &v)\n")
	b.WriteString("\titems := extractList(v)\n")
	b.WriteString("\tif len(items) == 0 { return fmt.Errorf(\"no list found\") }\n")
	b.WriteString("\tf, _ := os.Create(outputPath)\n")
	b.WriteString("\tw := csv.NewWriter(f)\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor k := range first { headers = append(headers, k) }\n")
	b.WriteString("\t\tw.Write(headers)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\trow := make([]string, len(headers))\n")
	b.WriteString("\t\t\tfor i, h := range headers { row[i] = fmt.Sprintf(\"%v\", m[h]) }\n")
	b.WriteString("\t\t\tw.Write(row)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tw.Flush(); return nil\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateYamlToXml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func marshalXML(v interface{}, name string) string {\n")
	b.WriteString("\tswitch val := v.(type) {\n")
	b.WriteString("\tcase map[string]interface{}:\n")
	b.WriteString("\t\tres := \"<\" + name + \">\\n\"\n")
	b.WriteString("\t\tfor k, v := range val { res += marshalXML(v, k) }\n")
	b.WriteString("\t\treturn res + \"</\" + name + \">\\n\"\n")
	b.WriteString("\tcase []interface{}:\n")
	b.WriteString("\t\tres := \"\"\n")
	b.WriteString("\t\tfor _, item := range val { res += marshalXML(item, name) }\n")
	b.WriteString("\t\treturn res\n")
	b.WriteString("\tdefault: return fmt.Sprintf(\"<%s>%v</%s>\\n\", name, val, name)\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; yaml.Unmarshal(data, &v)\n")
	b.WriteString("\txml := \"<?xml version=\\\"1.0\\\"?>\\n\" + marshalXML(v, \"root\")\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(xml), 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateYamlToMarkdown(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")
	b.WriteString(jsonHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; yaml.Unmarshal(data, &v)\n")
	b.WriteString("\titems := extractList(v)\n")
	b.WriteString("\tif len(items) == 0 { return nil }\n")
	b.WriteString("\tvar md strings.Builder\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\tfor k := range first { headers = append(headers, k) }\n")
	b.WriteString("\t\tmd.WriteString(\"| \" + strings.Join(headers, \" | \") + \" |\\n\")\n")
	b.WriteString("\t\tmd.WriteString(\"| \" + strings.Repeat(\"--- | \", len(headers)) + \"\\n\")\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor _, h := range headers { md.WriteString(fmt.Sprintf(\"| %v \", m[h])) }\n")
	b.WriteString("\t\t\tmd.WriteString(\"|\\n\")\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(md.String()), 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateYamlToExcel(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString(")\n\n")
	b.WriteString(jsonHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; yaml.Unmarshal(data, &v)\n")
	b.WriteString("\titems := extractList(v)\n")
	b.WriteString("\tx := excelize.NewFile()\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString("\t\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor k := range first { headers = append(headers, k); c, _ := excelize.CoordinatesToCellName(len(headers), 1); x.SetCellValue(\"Sheet1\", c, k) }\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tfor i, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor j, h := range headers { c, _ := excelize.CoordinatesToCellName(j+1, i+2); x.SetCellValue(\"Sheet1\", c, m[h]) }\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn x.SaveAs(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateYamlToPdf(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString("\t\"github.com/jung-kurt/gofpdf\"\n")
	b.WriteString(")\n\n")
	b.WriteString(jsonHelpers)
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tvar v interface{}; yaml.Unmarshal(data, &v)\n")
	b.WriteString("\titems := extractList(v)\n")
	b.WriteString("\tpdf := gofpdf.New(\"L\", \"mm\", \"A4\", \"\")\n")
	b.WriteString("\tpdf.AddPage(); pdf.SetFont(\"Arial\", \"B\", 10)\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString("\t\tif first, ok := items[0].(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor k := range first { headers = append(headers, k); pdf.Cell(30, 10, k) }\n")
	b.WriteString("\t\t\tpdf.Ln(-1)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tpdf.SetFont(\"Arial\", \"\", 10)\n")
	b.WriteString("\tfor _, item := range items {\n")
	b.WriteString("\t\tif m, ok := item.(map[string]interface{}); ok {\n")
	b.WriteString("\t\t\tfor _, h := range headers { pdf.Cell(30, 10, fmt.Sprintf(\"%v\", m[h])) }\n")
	b.WriteString("\t\t\tpdf.Ln(-1)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn pdf.OutputFileAndClose(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateMarkdownToCsv(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"encoding/csv\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func parseMarkdown(data []byte) []map[string]string {\n")
	b.WriteString("\tlines := strings.Split(string(data), \"\\n\")\n")
	b.WriteString("\tvar dataOut []map[string]string\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tfor _, l := range lines {\n")
	b.WriteString("\t\tl = strings.TrimSpace(l)\n")
	b.WriteString("\t\tif !strings.HasPrefix(l, \"|\") || strings.Contains(l, \"---\") { continue }\n")
	b.WriteString("\t\tparts := strings.Split(strings.Trim(l, \"|\"), \"|\")\n")
	b.WriteString("\t\tfor i := range parts { parts[i] = strings.TrimSpace(parts[i]) }\n")
	b.WriteString("\t\tif headers == nil { headers = parts } else {\n")
	b.WriteString("\t\t\tm := make(map[string]string)\n")
	b.WriteString("\t\t\tfor i, v := range parts { if i < len(headers) { m[headers[i]] = v } }\n")
	b.WriteString("\t\t\tdataOut = append(dataOut, m)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn dataOut\n")
	b.WriteString("}\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\titems := parseMarkdown(data)\n")
	b.WriteString("\tf, _ := os.Create(outputPath)\n")
	b.WriteString("\tw := csv.NewWriter(f)\n")
	b.WriteString("\tif len(items) > 0 {\n")
	b.WriteString("\t\tvar headers []string\n")
	b.WriteString("\t\tfor k := range items[0] { headers = append(headers, k) }\n")
	b.WriteString("\t\tw.Write(headers)\n")
	b.WriteString("\t\tfor _, m := range items {\n")
	b.WriteString("\t\t\trow := make([]string, len(headers))\n")
	b.WriteString("\t\t\tfor i, h := range headers { row[i] = m[h] }\n")
	b.WriteString("\t\t\tw.Write(row)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\tw.Flush(); return nil\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateMarkdownToXml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"fmt\"\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func parseMarkdown(data []byte) []map[string]string {\n")
	b.WriteString("\tlines := strings.Split(string(data), \"\\n\")\n")
	b.WriteString("\tvar dataOut []map[string]string\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tfor _, l := range lines {\n")
	b.WriteString("\t\tl = strings.TrimSpace(l)\n")
	b.WriteString("\t\tif !strings.HasPrefix(l, \"|\") || strings.Contains(l, \"---\") { continue }\n")
	b.WriteString("\t\tparts := strings.Split(strings.Trim(l, \"|\"), \"|\")\n")
	b.WriteString("\t\tfor i := range parts { parts[i] = strings.TrimSpace(parts[i]) }\n")
	b.WriteString("\t\tif headers == nil { headers = parts } else {\n")
	b.WriteString("\t\t\tm := make(map[string]string)\n")
	b.WriteString("\t\t\tfor i, v := range parts { if i < len(headers) { m[headers[i]] = v } }\n")
	b.WriteString("\t\t\tdataOut = append(dataOut, m)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn dataOut\n")
	b.WriteString("}\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\titems := parseMarkdown(data)\n")
	b.WriteString("\txml := \"<?xml version=\\\"1.0\\\"?>\\n<root>\\n\"\n")
	b.WriteString("\tfor _, m := range items {\n")
	b.WriteString("\t\txml += \"  <item>\\n\"\n")
	b.WriteString("\t\tfor k, v := range m { xml += fmt.Sprintf(\"    <%s>%v</%s>\\n\", k, v, k) }\n")
	b.WriteString("\t\txml += \"  </item>\\n\"\n")
	b.WriteString("\t}\n")
	b.WriteString("\txml += \"</root>\"\n")
	b.WriteString("\treturn os.WriteFile(outputPath, []byte(xml), 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateMarkdownToYaml(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString("\t\"gopkg.in/yaml.v3\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func parseMarkdown(data []byte) []map[string]string {\n")
	b.WriteString("\tlines := strings.Split(string(data), \"\\n\")\n")
	b.WriteString("\tvar dataOut []map[string]string\n")
	b.WriteString("\tvar headers []string\n")
	b.WriteString("\tfor _, l := range lines {\n")
	b.WriteString("\t\tl = strings.TrimSpace(l)\n")
	b.WriteString("\t\tif !strings.HasPrefix(l, \"|\") || strings.Contains(l, \"---\") { continue }\n")
	b.WriteString("\t\tparts := strings.Split(strings.Trim(l, \"|\"), \"|\")\n")
	b.WriteString("\t\tfor i := range parts { parts[i] = strings.TrimSpace(parts[i]) }\n")
	b.WriteString("\t\tif headers == nil { headers = parts } else {\n")
	b.WriteString("\t\t\tm := make(map[string]string)\n")
	b.WriteString("\t\t\tfor i, v := range parts { if i < len(headers) { m[headers[i]] = v } }\n")
	b.WriteString("\t\t\tdataOut = append(dataOut, m)\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn dataOut\n")
	b.WriteString("}\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\titems := parseMarkdown(data)\n")
	b.WriteString("\tout, _ := yaml.Marshal(items)\n")
	b.WriteString("\treturn os.WriteFile(outputPath, out, 0644)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateMarkdownToExcel(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString("\t\"github.com/xuri/excelize/v2\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tlines := strings.Split(string(data), \"\\n\")\n")
	b.WriteString("\tx := excelize.NewFile()\n")
	b.WriteString("\trowIdx := 1\n")
	b.WriteString("\tfor _, l := range lines {\n")
	b.WriteString("\t\tl = strings.TrimSpace(l)\n")
	b.WriteString("\t\tif !strings.HasPrefix(l, \"|\") || strings.Contains(l, \"---\") { continue }\n")
	b.WriteString("\t\tparts := strings.Split(strings.Trim(l, \"|\"), \"|\")\n")
	b.WriteString("\t\tfor i, v := range parts { c, _ := excelize.CoordinatesToCellName(i+1, rowIdx); x.SetCellValue(\"Sheet1\", c, strings.TrimSpace(v)) }\n")
	b.WriteString("\t\trowIdx++\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn x.SaveAs(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

func generateMarkdownToPdf(plan schema.Plan) (string, error) {
	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"os\"\n")
	b.WriteString("\t\"strings\"\n")
	b.WriteString("\t\"github.com/jung-kurt/gofpdf\"\n")
	b.WriteString(")\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tinputPath := %q\n", plan.InputPath)
	fmt.Fprintf(&b, "\toutputPath := %q\n", plan.OutputPath)
	b.WriteString("\tdata, _ := os.ReadFile(inputPath)\n")
	b.WriteString("\tlines := strings.Split(string(data), \"\\n\")\n")
	b.WriteString("\tpdf := gofpdf.New(\"L\", \"mm\", \"A4\", \"\")\n")
	b.WriteString("\tpdf.AddPage(); pdf.SetFont(\"Arial\", \"B\", 10)\n")
	b.WriteString("\tfirst := true\n")
	b.WriteString("\tfor _, l := range lines {\n")
	b.WriteString("\t\tl = strings.TrimSpace(l)\n")
	b.WriteString("\t\tif !strings.HasPrefix(l, \"|\") || strings.Contains(l, \"---\") { continue }\n")
	b.WriteString("\t\tparts := strings.Split(strings.Trim(l, \"|\"), \"|\")\n")
	b.WriteString("\t\tfor _, v := range parts { pdf.Cell(30, 10, strings.TrimSpace(v)) }\n")
	b.WriteString("\t\tpdf.Ln(-1)\n")
	b.WriteString("\t\tif first { pdf.SetFont(\"Arial\", \"\", 10); first = false }\n")
	b.WriteString("\t}\n")
	b.WriteString("\treturn pdf.OutputFileAndClose(outputPath)\n")
	b.WriteString("}\n")
	b.WriteString("func main() { run() }\n")
	return b.String(), nil
}

const xmlHelpers = `
func xmlToMap(d *xml.Decoder) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for {
		t, err := d.Token()
		if err == io.EOF { break }
		switch el := t.(type) {
		case xml.StartElement:
			v, _ := xmlToMap(d)
			m[el.Name.Local] = v
		case xml.EndElement: return m, nil
		case xml.CharData: if s := strings.TrimSpace(string(el)); s != "" { m["#text"] = s }
		}
	}
	return m, nil
}

func extractList(v interface{}) []interface{} {
	switch val := v.(type) {
	case []interface{}: return val
	case map[string]interface{}:
		for _, sub := range val {
			if list, ok := sub.([]interface{}); ok { return list }
			if m, ok := sub.(map[string]interface{}); ok {
				if list := extractList(m); list != nil { return list }
			}
		}
	}
	return nil
}
`

const jsonHelpers = `
func extractList(v interface{}) []interface{} {
	switch val := v.(type) {
	case []interface{}: return val
	case map[string]interface{}:
		for _, sub := range val {
			if list, ok := sub.([]interface{}); ok { return list }
		}
	}
	return nil
}
`

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
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module morphgo-temp\ngo 1.26.3\n\nrequire (\n\tgopkg.in/yaml.v3 v3.0.1\n\tgithub.com/xuri/excelize/v2 v2.8.1\n\tgithub.com/jung-kurt/gofpdf v1.16.2\n)\n"), 0o600); err != nil {
		return "", err
	}

	sumData, err := os.ReadFile("go.sum")
	if err == nil {
		_ = os.WriteFile(filepath.Join(dir, "go.sum"), sumData, 0o600)
	}

	return path, nil
}
