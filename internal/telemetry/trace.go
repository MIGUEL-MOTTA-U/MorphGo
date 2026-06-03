package telemetry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"morphgo/internal/codeact"
	"morphgo/internal/sandbox"
	"morphgo/internal/schema"
)

// RunLog captures the complete state and history of an agent execution.
type RunLog struct {
	ID        string            `json:"id"`
	Timestamp time.Time         `json:"timestamp"`
	Plan      schema.Plan       `json:"plan"`
	Attempts  []codeact.Attempt `json:"attempts"`
	Result    sandbox.Result    `json:"result"`
	Status    string            `json:"status"` // "success", "failure", "error"
}

// SaveRun persists the RunLog into a dedicated directory under baseDir.
// It returns the path to the created directory.
func SaveRun(baseDir string, log RunLog) (string, error) {
	if log.ID == "" {
		log.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	runDirName := fmt.Sprintf("%s-%s", log.Timestamp.Format("20060102-150405"), log.ID)
	runDir := filepath.Join(baseDir, runDirName)

	if err := os.MkdirAll(runDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create run directory: %w", err)
	}

	logPath := filepath.Join(runDir, "run.json")
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal run log: %w", err)
	}

	if err := os.WriteFile(logPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write run log: %w", err)
	}

	return runDir, nil
}
