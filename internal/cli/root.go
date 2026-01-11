package cli

import (
	"fmt"

	"github.com/rtfa/kevin/internal/core"
	"github.com/rtfa/kevin/internal/store"
	"github.com/spf13/cobra"
)

// CLI holds the dependencies for the application commands
type CLI struct {
	Config *core.ProjectConfig
	Store  store.Store
}

// NewRootCmd creates the root kevin command and injects dependencies
func NewRootCmd(cfg *core.ProjectConfig, s store.Store) *cobra.Command {
	app := &CLI{Config: cfg, Store: s}

	cmd := &cobra.Command{
		Use:          "kevin",
		Short:        "The Helpful Agent Orchestrator",
		Long:         "Kevin connects your existing tools to a structured workflow using the filesystem as the database.",
		SilenceUsage: true, // Don't show usage on error
	}

	cmd.AddCommand(app.newInitCmd())
	cmd.AddCommand(app.newTaskCmd())
	// cmd.AddCommand(app.newRunCmd()) // Future

	return cmd
}

func (app *CLI) ensureInit() error {
	if app.Config == nil || app.Store == nil {
		return fmt.Errorf("project not initialized. Run 'kevin init' first")
	}
	return nil
}
