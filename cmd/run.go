package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute the agent run pipeline",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.PersistentPreRunE(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
