package agent

import (
	"fmt"
	"time"

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

	// 2. Plan
	outputPath := fmt.Sprintf("output.%s", target)
	plan, err := schema.InferPlan(task, summary, target, inputPath, outputPath)
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

	log.Status = "success"
	return log, nil
}
