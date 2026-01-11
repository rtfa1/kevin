package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rtfa/kevin/internal/cli"
	"github.com/rtfa/kevin/internal/core"
	"github.com/rtfa/kevin/internal/store"
)

func main() {
	// 1. Try Load Config
	// We look in current directory .kevin/config.yaml
	configPath := filepath.Join(".kevin", "config.yaml") // Relative for now as per PRD "Current Dir"

	// Intentionally ignored error here (cfg might be nil) because 'init' command needs to run without config
	cfg, _ := core.LoadConfig(configPath)

	// 2. Init Store if config exists
	var st store.Store
	if cfg != nil {
		var err error
		dataDir := filepath.Join(".kevin", "board")
		st, err = store.NewFileStore(dataDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing store: %v\n", err)
			os.Exit(1)
		}
	}

	// 3. Init Root CLI with dependencies
	rootCmd := cli.NewRootCmd(cfg, st)

	// 4. Execute
	if err := rootCmd.Execute(); err != nil {
		// Cobra prints error if SilenceErrors is false (default)
		os.Exit(1)
	}
}
