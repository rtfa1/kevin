# Configuration Schema (.kevin/config.yaml)

## Overview
The `config.yaml` file acts as the registry for Project Kevin. It defines **Project Settings**, **Board Layout**, and **Agent Capabilities**.

## Schema Definition

### 1. Global Settings
```yaml
version: "1.0"
project_name: "Kevin"
default_executor: "docker" # or "local"
```

### 2. Board Configuration
Customize the TUI workflow.
```yaml
board:
  columns:
    - id: "backlog"
      title: "Backlog"
    - id: "todo"
      title: "To Do"
    - id: "doing"
      title: "In Progress"
    - id: "done"
      title: "Done"
  tags:
    - "bug"
    - "feature"
    - "docs"
```

### 3. Agent Configuration
This is the core definition of external tools.

```yaml
agents:
  - name: "senior-coder"
    description: "Expert Go developer for complex logic"
    
    # SYSTEM PROMPT / CONTEXT
    # Can be a direct string or a path to a file (relative to .kevin/)
    system_prompt: "You are a Senior Go Engineer. Prefer idiomatic code. Use Table Driven Tests."
    # OR
    # system_prompt_file: "prompts/senior_go.md"

    # EXECUTION
    executor: "docker"
    image: "ghcr.io/opencode/opencode:latest"
    
    # COMMAND TEMPLATE
    # Available variables:
    # - {{.TaskPath}}: Path to the markdown file
    # - {{.SystemPrompt}}: The content defined above
    # - {{.Workspace}}: The mounted workspace path
    command: 
      - "opencode"
      - "--task"
      - "{{.TaskPath}}"
      - "--system"
      - "{{.SystemPrompt}}"
    
    # ENVIRONMENT
    # Variables to pass from the Host to the Container
    env_pass:
      - "OPENAI_API_KEY"
      - "GITHUB_TOKEN"
      - "ANTHROPIC_API_KEY"

  - name: "quick-fix"
    description: "Fast local LLM for typos"
    executor: "local"
    system_prompt: "Fix spelling and grammar only."
    command:
      - "local-llm"
      - "--prompt"
      - "{{.SystemPrompt}} \n\n Task: {{.TaskPath}}"
```

## Variable Interpolation
Kevin will verify that required variables are available before running.
*   `{{.SystemPrompt}}` is injected into the command line if the tool supports a flag for it.
*   If the tool *doesn't* support a system prompt flag, the User can include it in the `command` template (as shown in the `quick-fix` example) or Kevin can prepend it to the Task Markdown body temporarily.

## Secrets Management
Kevin **never logs** the values of `env_pass`. It simply looks them up in `os.Environ()` and passes them to the Executor.
