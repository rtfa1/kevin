# Agent Workflow & Lifecycle

## 1. Core Concept: The "Wrapped" Agent
In Kevin, an "Agent" is simply a configuration that describes how to run an external CLI tool. Kevin does not "think"; it orchestrates.

### The Flow
1.  **Trigger**: User runs `kevin run 42` (or presses `r` in TUI).
2.  **Context**: Kevin gathers:
    *   Task Description (`.kevin/board/task-42.md`).
    *   Project Config (`.kevin/config.yaml`).
    *   User Auth/Env (passed through from host).
3.  **Execution**: Kevin invokes the configured `Processor`.
4.  **Feedback**: Output streams to TUI; artifacts are written to disk.

## 2. Agent Configuration
Agents are defined in `.kevin/config.yaml`.

```yaml
agents:
  - name: "coder"
    description: "General purpose coding agent"
    # The command template
    command: ["opencode", "--task", "{{.TaskPath}}", "--workspace", "{{.Workspace}}"]
    # Environment variables to pass through
    env_pass: ["OPENAI_API_KEY", "GITHUB_TOKEN"]
    # How to run it
    executor: "docker"
    image: "ghcr.io/opencode/opencode:latest"

  - name: "helper"
    command: ["claude", "--prompt-file", "{{.TaskPath}}"]
    executor: "local"
```

## 3. Workflow Stages

### Stage A: Planning (Optional)
*   **Role**: Break down a high-level request into tracking items.
*   **Tool**: `opencode` (or similar high-capability agent).
*   **Action**:
    1.  Kevin runs the planner agent with the prompt.
    2.  Agent writes new markdown files to `.kevin/board/`.
    3.  Kevin's **File Watcher** detects new files and updates the board.

### Stage B: Execution (The Loop)
*   **Role**: Complete a specific task.
*   **Tool**: `claude-cli`, `codex`, custom scripts.
*   **Action**:
    1.  User assigns `agent: coder` to Task #42.
    2.  Kevin prepares the environment (mounts `.kevin/board` and source code).
    3.  Kevin executes the tool.
    4.  Tool modifies code files and updates `task-42.md` (e.g., appending logs or checking boxes).

## 4. Operational Context in `.kevin`
The hidden folder acts as the bridge.

*   **Input**: The External Tool reads `.kevin/board/task-XYZ.md` to know what to do.
*   **Output**: The External Tool writes to `src/` (code) and back to `.kevin/board/task-XYZ.md` (status updates).
*   **State**: Kevin monitors `.kevin/board/*.md` to reflect changes in the UI.

## 5. Lifecycle Hooks (Future)
*   `on_start`: e.g., git pull.
*   `on_success`: e.g., git commit, move card to "Done".
*   `on_fail`: Add "Blocked" tag.
