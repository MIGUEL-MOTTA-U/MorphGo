package codeact

import "morphgo/internal/schema"

// PrepareTempMain generates and writes the temporary main.go for a run.
func PrepareTempMain(plan schema.Plan) (string, error) {
	code, err := GenerateMain(plan)
	if err != nil {
		return "", err
	}
	return WriteTempMain(code)
}
