package cmd

import (
	"fmt"
	"os"

	"morphgo/internal/agent"
	"morphgo/internal/telemetry"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Execute the agent run pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Starting run: %s -> %s (task: %s)\n", inputPath, targetFormat, taskStr)

		log, err := agent.Run(inputPath, taskStr, targetFormat)

		// Always try to save the trace
		traceDir := "runs"
		if outputPath != "" {
			traceDir = outputPath
		}

		saveDir, saveErr := telemetry.SaveRun(traceDir, log)
		if saveErr != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to save run trace: %v\n", saveErr)
		} else {
			fmt.Printf("Trace saved to: %s\n", saveDir)
		}

		if err != nil {
			return fmt.Errorf("run failed: %w", err)
		}

		fmt.Println("Run completed successfully!")
		return nil
		},
	}
}
