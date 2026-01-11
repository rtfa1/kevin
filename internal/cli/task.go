package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/rtfa/kevin/internal/core"
	"github.com/spf13/cobra"
)

func (app *CLI) newTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}

	// kevin task new <title>
	cmd.AddCommand(&cobra.Command{
		Use:   "new [title]",
		Short: "Create a new task",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			if err := app.ensureInit(); err != nil {
				return err
			}

			// Generate ID (timestamp based for now, or simple counter if store supported it)
			// For MVP, let's use task-YYYYMMDD-HHMMSS
			id := fmt.Sprintf("task-%s", time.Now().Format("20060102-150405"))

			task := core.Task{
				ID:      id,
				Title:   args[0],
				Status:  core.StatusTodo,
				Created: time.Now(),
				Content: "\nAdd details here...",
			}

			if err := app.Store.Create(task); err != nil {
				return fmt.Errorf("failed to create task: %w", err)
			}

			fmt.Printf("Created task %s: %s\n", task.ID, task.Title)
			return nil
		},
	})

	// kevin task list
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE: func(c *cobra.Command, args []string) error {
			if err := app.ensureInit(); err != nil {
				return err
			}

			tasks, err := app.Store.List()
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			// Simple Table Output
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tSTATUS\tTITLE")
			for _, t := range tasks {
				fmt.Fprintf(w, "%s\t%s\t%s\n", t.ID, t.Status, t.Title)
			}
			w.Flush()
			return nil
		},
	})

	// kevin task json (helper for scripts)
	cmd.AddCommand(&cobra.Command{
		Use:   "json",
		Short: "Dump tasks as JSON",
		RunE: func(c *cobra.Command, args []string) error {
			if err := app.ensureInit(); err != nil {
				return err
			}
			tasks, err := app.Store.List()
			if err != nil {
				return err
			}

			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(tasks)
		},
	})

	return cmd
}
