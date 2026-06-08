package agent

import (
	"fmt"
	"time"
	"path/filepath"
	"strings"

	"morphgo/internal/codeact"
	"morphgo/internal/schema"
	"morphgo/internal/telemetry"
)

// Run coordinates the full agent lifecycle: Inspect -> Plan -> Execute (Retry) -> Trace.
func Run(inputPath, task, target string) (telemetry.RunLog, error) {
	log := telemetry.RunLog{
		Timestamp: time.Now(),
		Status:    "error",
	}

	// 1. Inspect
	summary, err := Inspect(inputPath)
	if err != nil {
		return log, fmt.Errorf("inspection failed: %w", err)
	}

	// Normalize target format
	rawTarget := strings.ToLower(target)
	internalTarget := rawTarget
	outputExt := rawTarget
	switch rawTarget {
	case "md", "markdown":
		internalTarget = "markdown"
		outputExt = "md"
	case "xlsx", "excel":
		internalTarget = "excel"
		outputExt = "xlsx"
	case "cvs":
		internalTarget = "csv"
		outputExt = "csv"
	}

	// 2. Plan
	absInputPath, err := filepath.Abs(inputPath)
	if err != nil {
		return log, fmt.Errorf("failed to resolve absolute input path: %w", err)
	}
	outputDir := filepath.Dir(absInputPath)
	outputPath := filepath.Join(outputDir, fmt.Sprintf("output.%s", outputExt))
	plan, err := schema.InferPlan(task, summary, internalTarget, absInputPath, outputPath)
	if err != nil {
		return log, fmt.Errorf("planning failed: %w", err)
	}
	log.Plan = plan

	// 3. Execute with Retry Loop
	res, history, err := codeact.RetryRunPlan(plan, 30*time.Second, 3)
	log.Attempts = history
	log.Result = res

	if err != nil {
		log.Status = "failure"
		return log, fmt.Errorf("execution failed: %w", err)
	}

	fmt.Printf("Output saved to: %s\n", outputPath)
	log.Status = "success"
	return log, nil
}
