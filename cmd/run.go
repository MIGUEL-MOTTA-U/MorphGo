package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the agent run pipeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validation is handled by rootCmd.PersistentPreRunE
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
