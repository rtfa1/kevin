package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/rtfa/kevin/internal/agent"
	"github.com/rtfa/kevin/internal/core"
	"github.com/rtfa/kevin/internal/executor"
	"github.com/spf13/cobra"
)

func (app *CLI) newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run [task-id|title]",
		Short: "Run the assigned agent for a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.ensureInit(); err != nil {
				return err
			}

			// 1. Context with cancellation
			ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
			defer stop()

			// 2. Resolve Task
			query := args[0]
			task, err := app.resolveTask(query)
			if err != nil {
				return err
			}
			fmt.Printf("Found Task: %s (%s)\n", task.Title, task.ID)

			// 3. Resolve Agent
			if task.Assignee == "" {
				return fmt.Errorf("task %s has no assignee", task.ID)
			}
			agentCfg, err := app.resolveAgent(task.Assignee)
			if err != nil {
				// helpful error for user
				return fmt.Errorf("agent '%s' not found in config.yaml. Available agents: %v",
					task.Assignee, app.listAgentNames())
			}
			fmt.Printf("Using Agent: %s (Executor: %s)\n", agentCfg.Name, agentCfg.Executor)

			// 4. Prepare Execution
			cwd, _ := os.Getwd() // Use current directory as project root for now
			cmdSlice, envSlice, err := agent.Prepare(*agentCfg, *task, cwd)
			if err != nil {
				return fmt.Errorf("failed to prepare agent: %w", err)
			}

			// 5. Execute
			// TODO: switch on agentCfg.Executor when we have Docker support
			exec := executor.NewLocalExecutor()

			fmt.Printf("Running: %s\n", strings.Join(cmdSlice, " "))
			fmt.Println("--- Agent Output ---")

			if err := exec.Run(ctx, cmdSlice, envSlice, cwd); err != nil {
				if errors.Is(err, context.Canceled) {
					fmt.Println("\nStopped by user.")
					return nil
				}
				return err
			}

			return nil
		},
	}
}

// Helpers

func (app *CLI) resolveTask(query string) (*core.Task, error) {
	// 1. Try exact ID match logic (if we had efficient lookup)
	// But since we only have List(), let's list and search.
	tasks, err := app.Store.List()
	if err != nil {
		return nil, err
	}

	for _, t := range tasks {
		// Exact ID match
		if t.ID == query {
			return &t, nil
		}
	}

	// 2. Substring Title match
	var matches []core.Task
	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Title), strings.ToLower(query)) {
			matches = append(matches, t)
		}
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no task found matching '%s'", query)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("multiple tasks found matching '%s'. Please use ID", query)
	}

	return &matches[0], nil
}

func (app *CLI) resolveAgent(name string) (*core.AgentConfig, error) {
	for _, a := range app.Config.Agents {
		if a.Name == name {
			return &a, nil
		}
	}
	return nil, fmt.Errorf("agent not found")
}

func (app *CLI) listAgentNames() []string {
	names := []string{}
	for _, a := range app.Config.Agents {
		names = append(names, a.Name)
	}
	return names
}
