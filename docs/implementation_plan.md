# Implementation Plan - Project Kevin

# Goal
Build **Kevin**, a TUI-first, CLI-compatible Task Orchestrator. The system is built on a "Filesystem-as-Database" philosophy, orchestrating external tools via a pluggable Executor interface.

## User Review Required
> [!IMPORTANT]
> **Dependencies**: Confirmed usage of `spf13/cobra` (CLI), `charmbracelet/bubbletea` (TUI), `fsnotify/fsnotify` (Watcher), and `moby/moby` (Docker Client, optional/later).
> **Layout**: Standard Go layout (`cmd/`, `internal/`) as per PRD.

## Phase 1: Core, Data Layer & CLI Skeleton
**Objective**: Establish the project structure, data persistence, and basic CLI commands.

### 1.1. Project Initialization
#### [NEW] [go.mod](file:///Users/rtfa/lab/kevin/go.mod)
- Module: `github.com/rtfa/kevin`

#### [NEW] [internal/core/config.go](file:///Users/rtfa/lab/kevin/internal/core/config.go)
- Structs for `ProjectConfig`, `AgentConfig`.
- Loader logic for `.kevin/config.yaml`.

#### [NEW] [internal/core/task.go](file:///Users/rtfa/lab/kevin/internal/core/task.go)
- Struct `Task` with frontmatter tags (`yaml:"status"`).

### 1.2. The Data Layer (Store)
#### [NEW] [internal/store/interface.go](file:///Users/rtfa/lab/kevin/internal/store/interface.go)
- `type Store interface { List(), Create(), Update(), Watch() }`

#### [NEW] [internal/store/filestore.go](file:///Users/rtfa/lab/kevin/internal/store/filestore.go)
- **Watcher**: Implement `fsnotify` loop on `.kevin/board/`.
- **Parser**: Use `adrg/frontmatter` to read/write Markdown.
- **Debounce**: Implement simple debounce for file events.

### 1.3. CLI Foundation
#### [NEW] [cmd/kevin/main.go](file:///Users/rtfa/lab/kevin/cmd/kevin/main.go)
- Root command using `cobra`.

#### [NEW] [internal/cli/init.go](file:///Users/rtfa/lab/kevin/internal/cli/init.go)
- Command `kevin init` to bootstrap `.kevin/` structure.

#### [NEW] [internal/cli/task.go](file:///Users/rtfa/lab/kevin/internal/cli/task.go)
- Commands `kevin task new`, `kevin task list`.

---

## Phase 2: The TUI (Bubble Tea)
**Objective**: A visual, reactive Kanban board that syncs with `FileStore`.

### 2.1. TUI Model
#### [NEW] [internal/tui/model.go](file:///Users/rtfa/lab/kevin/internal/tui/model.go)
- Main `tea.Model` struct.
- Holds `[]Task` and `Store` reference.

#### [NEW] [internal/tui/view.go](file:///Users/rtfa/lab/kevin/internal/tui/view.go)
- Render Columns (`Backlog`, `Todo`, `Doing`, `Done`).
- Styles using `lipgloss`.

### 2.2. Reactivity
#### [MODIFY] [internal/tui/model.go](file:///Users/rtfa/lab/kevin/internal/tui/model.go)
- Listen to `Store.Watch()` channel.
- Handle `TaskUpdateMsg` to reload the TUI state instantly.

---

## Phase 3: Agent Orchestration (Executor)
**Objective**: Run external tools via Docker/Local execution.

### 3.1. Executor Interface
#### [NEW] [internal/executor/interface.go](file:///Users/rtfa/lab/kevin/internal/executor/interface.go)
- `Run(ctx, cmd []string, env []string, workDir string) error`

### 3.2. Implementations
#### [NEW] [internal/executor/local.go](file:///Users/rtfa/lab/kevin/internal/executor/local.go)
- Wraps `exec.Command`.

#### [NEW] [internal/executor/docker.go](file:///Users/rtfa/lab/kevin/internal/executor/docker.go)
- Wraps Docker CLI or SDK.

### 3.3. Integration
#### [NEW] [internal/core/agent.go](file:///Users/rtfa/lab/kevin/internal/core/agent.go)
- Logic to resolve `{{.TaskPath}}` and `{{.SystemPrompt}}`.
- Logic to merge host ENV variables.

#### [NEW] [internal/cli/run.go](file:///Users/rtfa/lab/kevin/internal/cli/run.go)
- Command `kevin run <task-id>`.
- Finds agent from Config -> Prepares Executor -> Runs.

## Verification Plan

### Core & Data
1.  Run `kevin init` -> Verify `.kevin` created.
2.  Run `kevin task new "Test"` -> Verify file created.
3.  Edit file manually -> Verify `kevin task list` shows changes.

### TUI
1.  Run `kevin` -> Verify board renders.
2.  Edit file in 2nd terminal -> Verify TUI updates.

### Agents
1.  Configure a "echo" agent.
2.  Run `kevin run <id>` -> Verify "echo" output appears/logs.
