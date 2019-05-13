package cmd

import (
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "start",
	Short: "Start monitoring.",
	Long:  "Start monitoring.",
	Run: func(cmd *cobra.Command, args []string) {
		startMonitoring()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func startMonitoring() {
}
