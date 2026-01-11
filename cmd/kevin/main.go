package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "kevin",
	Short: "Kevin is a TUI-first Task Orchestrator",
	Long: `Kevin connects your existing tools to a structured workflow,
using the filesystem as the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: TUI Launch logic will go here
		fmt.Println("Kevin is ready to help!")
	},
}

func init() {
	// Global flags will be defined here
}
