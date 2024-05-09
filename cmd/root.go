package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wado",
	Short: "A tool for mixing and munging Doom WAD files.",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: print help
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
