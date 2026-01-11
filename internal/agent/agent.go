package agent

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/rtfa/kevin/internal/core"
)

// Prepare resolves the command templates and environment for an agent execution
func Prepare(agentCfg core.AgentConfig, task core.Task, projectDir string) ([]string, []string, error) {
	ctx := ExecutionContext{
		TaskID:     task.ID,
		TaskTitle:  task.Title,
		TaskStatus: string(task.Status),
		TaskPath:   task.FilePath,
		ProjectDir: projectDir,
	}

	// 1. Render Command Arguments
	cmd := make([]string, len(agentCfg.Command))
	for i, arg := range agentCfg.Command {
		val, err := render(arg, ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to render arg '%s': %w", arg, err)
		}
		cmd[i] = val
	}

	// 2. Build Envs
	env := buildEnv(agentCfg)

	return cmd, env, nil
}

func render(tmpl string, ctx ExecutionContext) (string, error) {
	t, err := template.New("cmd").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func buildEnv(cfg core.AgentConfig) []string {
	// Start with empty or base env?
	// The LocalExecutor appends os.Environ(), so we only need to return overrides/additions here.
	// However, if we want to restrict, we should do it here or in executor.
	// For LocalExecutor, it merges. For Docker, we probably want to be explicit.
	// Let's implement the "EnvPass" logic to be: these are variables we want to explicitely capture from Host and pass down.

	currentEnv := toMap(os.Environ())
	var finalEnv []string

	for _, key := range cfg.EnvPass {
		if val, ok := currentEnv[key]; ok {
			finalEnv = append(finalEnv, fmt.Sprintf("%s=%s", key, val))
		}
	}

	return finalEnv
}

func toMap(env []string) map[string]string {
	m := make(map[string]string)
	for _, e := range env {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			m[parts[0]] = parts[1]
		}
	}
	return m
}
