package schema

import (
	"encoding/json"
	"errors"
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
	switch {
	case strings.Contains(p, "convert") || strings.Contains(p, "convertir") || strings.Contains(p, "transform"):
		return "convert"
	case strings.Contains(p, "filter") || strings.Contains(p, "filtrar"):
		return "filter"
	case strings.Contains(p, "extract") || strings.Contains(p, "extraer"):
		return "extract"
	case strings.Contains(p, "normalize") || strings.Contains(p, "normalizar"):
		return "normalize"
	default:
		return ""
	}
}

// MarshalJSON keeps the plan serializable with a stable shape.
func (p Plan) MarshalJSON() ([]byte, error) {
	type alias Plan
	return json.Marshal(alias(p))
}
