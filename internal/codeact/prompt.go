package codeact

import (
	"errors"
	"fmt"
	"strings"

	"morphgo/internal/schema"
)

// GenerateMain renders a minimal temporary Go program for a given plan.
func GenerateMain(plan schema.Plan) (string, error) {
	if plan.Objective == "" {
		return "", errors.New("plan objective is required")
	}
	if plan.Source == "" || plan.Target == "" || plan.Operation == "" {
		return "", errors.New("plan is incomplete")
	}

	var b strings.Builder
	b.WriteString("package main\n\n")
	b.WriteString("import (\n\t\"errors\"\n\t\"fmt\"\n)\n\n")
	b.WriteString("func run() error {\n")
	fmt.Fprintf(&b, "\tfmt.Println(%q)\n", "execute "+plan.Operation+" from "+plan.Source+" to "+plan.Target)
	b.WriteString("\treturn nil\n")
	b.WriteString("}\n\n")
	b.WriteString("func main() {\n")
	b.WriteString("\tif err := run(); err != nil {\n")
	b.WriteString("\t\tfmt.Println(errors.New(err.Error()))\n")
	b.WriteString("\t}\n")
	b.WriteString("}\n")
	return b.String(), nil
}
