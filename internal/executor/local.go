package executor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// LocalExecutor runs commands directly on the host machine
type LocalExecutor struct {
	Stdout io.Writer
	Stderr io.Writer
}

// NewLocalExecutor creates a new LocalExecutor.
// If stdout/stderr are nil, defaults to os.Stdout/os.Stderr.
func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (e *LocalExecutor) Run(ctx context.Context, command []string, env []string, workDir string) error {
	if len(command) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Dir = workDir

	// Merge env: os.Environ() + env
	// We append 'env' last so it overrides system defaults if needed
	cmd.Env = append(os.Environ(), env...)

	cmd.Stdout = e.Stdout
	cmd.Stderr = e.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("local execution failed: %w", err)
	}

	return nil
}
