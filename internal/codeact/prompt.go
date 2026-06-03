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
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module morphgo-temp\ngo 1.26.3\n"), 0o600); err != nil {
		return "", err
	}
	return path, nil
}
