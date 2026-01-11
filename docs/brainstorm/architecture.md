# System Architecture (Finalized)

## Overview
Kevin is a **TUI-first, CLI-compatible Task Orchestrator**. It manages a local Kanban board stored as Markdown files and orchestrates external "Agent" tools (CLI binaries/containers) to execute work against those tasks.

## Core Philosophy
1.  **Filesystem as Database**: The `.kevin/board/` directory is the single source of truth.
2.  **External Intelligence**: Kevin has no internal AI. It wraps tools like `opencode`, `claude`, or `aider`.
3.  **Config Driven**: Behavior is defined in `.kevin/config.yaml`.
4.  **Agnostic Execution**: Agents can run in Docker, locally, or via SSH.

## Component Design

### 1. The Controller (The Brain)
*   **Startup**: Reads `.kevin/config.yaml`.
*   **File Watcher**: Uses `fsnotify` to monitor `.kevin/board/*.md`.
    *   *Event Loop*: File Change -> Debounce -> Parse -> Update State -> Notify TUI.
*   **Executor Service**: Implementation of the `Executor` interface.
    *   `Run(cmd, env, workdir)` -> Streams Output.

### 2. The Interfaces (The Face)
*   **TUI (Bubble Tea)**:
    *   Visual Board (Columns).
    *   Real-time logs view.
    *   Command palette for triggering agents.
*   **CLI (Cobra)**:
    *   Headed-less operations (`kevin task move`, `kevin run`).
    *   Reuse the same internal `Controller` logic as the TUI.

### 3. The Data Layer
*   **Storage**: `.kevin/board/`
*   **Format**: Markdown with YAML Frontmatter.
    ```markdown
    ---
    id: "task-123"
    status: "todo"
    assignee: "coder"
    ---
    # Fix Login
    ...
    ```

### 4. The Agent Integration
*   **Definition**: A mapping in `config.yaml` of `Name -> Command Template`.
*   **Execution Flow**:
    1.  Kevin injects variables (`{{.TaskPath}}`, `{{.SystemPrompt}}`).
    2.  Kevin injects allow-listed ENV vars (`OPENAI_API_KEY`).
    3.  Kevin invokes the `Executor`.
    4.  The External Tool reads the Task File + Codebase, does work, and writes back to the Task File.

## Directory Structure
```text
kevin/
  cmd/kevin/       # Usage: kevin [flags] [command]
  internal/
    core/          # Domain models (Task, Config)
    store/         # fsnotify + parser
    executor/      # Docker/Local implementations
    tui/           # Bubble Tea views
    cli/           # Cobra commands
```
