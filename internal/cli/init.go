package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func (app *CLI) newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Kevin project",
		Long:  "Creates the .kevin directory structure and a default config.yaml in the current directory.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Create .kevin directory
			kevinDir := ".kevin"
			boardDir := filepath.Join(kevinDir, "board")

			if err := os.MkdirAll(boardDir, 0755); err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}

			// 2. Create config.yaml if not exists
			configPath := filepath.Join(kevinDir, "config.yaml")
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				defaultConfig := `project:
  name: "My Project"

agents:
  - name: "coder"
    executor: "local" # or docker
    command: ["echo", "I am a placeholder agent"]
`
				if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
					return fmt.Errorf("failed to write config: %w", err)
				}
				fmt.Println("Initialized .kevin/config.yaml")
			} else {
				fmt.Println("Config already exists, skipping creation.")
			}

			fmt.Println("Kevin Project Initialized! 🍌")
			return nil
		},
	}
}
