package executor

import (
	"context"
)

// Executor defines the interface for running external tools
type Executor interface {
	// Run executes a command and streams output.
	// workDir: The directory the agent should start in (on host or in container).
	// env: "KEY=VALUE" strings.
	Run(ctx context.Context, cmd []string, env []string, workDir string) error
}
