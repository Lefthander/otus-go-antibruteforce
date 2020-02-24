package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "abf-ctl",
	Short: "abf-ctl",
	Long:  "abf-ctl",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Usage: abf-ctl [command]")
	},
}

// Execute the Cobra root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Cannot execute the command -", err)
	}
}
