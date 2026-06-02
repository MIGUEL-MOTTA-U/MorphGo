package codeact

import (
	"time"

	"morphgo/internal/sandbox"
	"morphgo/internal/schema"
)

// PrepareTempMain generates and writes the temporary main.go for a run.
func PrepareTempMain(plan schema.Plan) (string, error) {
	code, err := GenerateMain(plan)
	if err != nil {
		return "", err
	}
	return WriteTempMain(code)
}

// RunPlan prepares and executes the generated program through the sandbox.
func RunPlan(plan schema.Plan, timeout time.Duration) (sandbox.Result, error) {
	mainPath, err := PrepareTempMain(plan)
	if err != nil {
		return sandbox.Result{}, err
	}
	return sandbox.RunMain(mainPath, timeout)
}
