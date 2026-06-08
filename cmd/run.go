package cmd

import (
	"fmt"
	"os"
	"strings"

	"morphgo/internal/agent"
	"morphgo/internal/telemetry"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Execute the agent run pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			targets := strings.Split(targetFormat, ",")
			var lastErr error

			for _, t := range targets {
				t = strings.TrimSpace(t)
				if t == "" {
					continue
				}

				fmt.Printf("\n--- Starting run: %s -> %s (task: %s) ---\n", inputPath, t, taskStr)

				log, err := agent.Run(inputPath, taskStr, t)

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
					fmt.Fprintf(os.Stderr, "Error: run failed for target %s: %v\n", t, err)
					lastErr = err
					continue
				}

				fmt.Printf("Run completed successfully for target %s!\n", t)
			}

			if lastErr != nil {
				return fmt.Errorf("one or more runs failed (last error: %w)", lastErr)
			}

			return nil
		},
	}
}
