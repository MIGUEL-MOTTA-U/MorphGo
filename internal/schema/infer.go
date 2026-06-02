package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// InferPlan translates a user prompt and a file summary into a small execution plan.
func InferPlan(prompt string, summary Summary, target string) (Plan, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return Plan{}, errors.New("prompt is required")
	}
	if summary.Format == "" {
		return Plan{}, errors.New("summary is required")
	}
	if target == "" {
		return Plan{}, errors.New("target is required")
	}

	op := detectOperation(prompt)
	if op == "" {
		return Plan{}, errors.New("prompt is ambiguous")
	}

	steps := []string{
		"inspect " + summary.Format,
		op + " into " + target,
	}
	return Plan{
		Objective: prompt,
		Source:    summary.Format,
		Target:    target,
		Operation: op,
		Steps:     steps,
		Summary:   summary,
	}, nil
}

func detectOperation(prompt string) string {
	p := strings.ToLower(prompt)
	matches := 0
	op := ""
	switch {
	case strings.Contains(p, "convert") || strings.Contains(p, "convertir") || strings.Contains(p, "transform"):
		op = "convert"
		matches++
	}
	if strings.Contains(p, "filter") || strings.Contains(p, "filtrar") {
		op = "filter"
		matches++
	}
	if strings.Contains(p, "extract") || strings.Contains(p, "extraer") {
		op = "extract"
		matches++
	}
	if strings.Contains(p, "normalize") || strings.Contains(p, "normalizar") {
		op = "normalize"
		matches++
	}
	if matches != 1 {
		return ""
	}
	return op
}

// MarshalJSON keeps the plan serializable with a stable shape.
func (p Plan) MarshalJSON() ([]byte, error) {
	type alias Plan
	return json.Marshal(alias(p))
}

// MarshalMarkdown serializes the plan into a compact Markdown representation.
func (p Plan) MarshalMarkdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Plan\n")
	fmt.Fprintf(&b, "- Objective: %s\n", p.Objective)
	fmt.Fprintf(&b, "- Source: %s\n", p.Source)
	fmt.Fprintf(&b, "- Target: %s\n", p.Target)
	fmt.Fprintf(&b, "- Operation: %s\n", p.Operation)
	if len(p.Steps) > 0 {
		b.WriteString("- Steps:\n")
		for _, step := range p.Steps {
			fmt.Fprintf(&b, "  - %s\n", step)
		}
	}
	return b.String()
}
