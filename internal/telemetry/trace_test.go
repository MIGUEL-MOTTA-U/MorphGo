package telemetry

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"morphgo/internal/codeact"
	"morphgo/internal/sandbox"
	"morphgo/internal/schema"
)

func TestSaveRun_CreatesDirectoryAndFile(t *testing.T) {
	baseDir := t.TempDir()
	log := RunLog{
		ID:        "test-run",
		Timestamp: time.Date(2026, 6, 3, 10, 0, 0, 0, time.UTC),
		Plan:      schema.Plan{Objective: "test"},
		Attempts: []codeact.Attempt{
			{Number: 1, Stage: "run", Stdout: "hello"},
		},
		Result: sandbox.Result{Stdout: "hello", ExitCode: 0},
		Status: "success",
	}

	runDir, err := SaveRun(baseDir, log)
	if err != nil {
		t.Fatalf("SaveRun failed: %v", err)
	}

	expectedDirName := "20260603-100000-test-run"
	if filepath.Base(runDir) != expectedDirName {
		t.Errorf("expected directory name %q, got %q", expectedDirName, filepath.Base(runDir))
	}

	logPath := filepath.Join(runDir, "run.json")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Fatal("run.json was not created")
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read run.json: %v", err)
	}

	var savedLog RunLog
	if err := json.Unmarshal(data, &savedLog); err != nil {
		t.Fatalf("failed to unmarshal run.json: %v", err)
	}

	if savedLog.ID != log.ID {
		t.Errorf("expected ID %q, got %q", log.ID, savedLog.ID)
	}
	if len(savedLog.Attempts) != 1 {
		t.Errorf("expected 1 attempt, got %d", len(savedLog.Attempts))
	}
}

func TestSaveRun_GeneratesIDAndTimestamp(t *testing.T) {
	baseDir := t.TempDir()
	log := RunLog{}

	runDir, err := SaveRun(baseDir, log)
	if err != nil {
		t.Fatalf("SaveRun failed: %v", err)
	}

	logPath := filepath.Join(runDir, "run.json")
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("failed to read run.json: %v", err)
	}

	var savedLog RunLog
	if err := json.Unmarshal(data, &savedLog); err != nil {
		t.Fatalf("failed to unmarshal run.json: %v", err)
	}

	if savedLog.ID == "" {
		t.Error("expected generated ID, got empty string")
	}
	if savedLog.Timestamp.IsZero() {
		t.Error("expected generated timestamp, got zero")
	}
}
