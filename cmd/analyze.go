package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyzeCmd)
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze the difficulty of a WAD",
	Long: `Analyzes each level in a WAD by looking at thing counts and
calculates a estimated difficulty score along with other metrics. Optionally
generates a graphical report.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("This function is not yet implemented.")
	},
}
