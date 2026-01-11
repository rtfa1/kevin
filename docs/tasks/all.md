# Master Task List - Project Kevin

This document consolidates all tasks required to build Kevin, derived from the [PRD](../PRD.md), [Implementation Plan](../implementation_plan.md), and [Brainstorming](../brainstorm/) documents.

## Phase 1: Core Skeleton & Data Layer
**Goal**: Establish the "Filesystem-as-Database" foundation and basic CLI interaction.

### 1.1. Project Initialization
- [ ] **Module Init**: Initialize `github.com/rtfa/kevin`.
- [ ] **Directory Structure**: Create:
    - `cmd/kevin/`: CLI entrypoint.
    - `internal/core/`: Domain models.
    - `internal/store/`: Persistence logic.
    - `internal/executor/`: Agent runners.
    - `internal/tui/`: UI logic.
    - `internal/cli/`: Cobra commands.

### 1.2. Domain Modeling (`internal/core`)
- [ ] **Config Model**: Define `ProjectConfig` and `AgentConfig` structs (Ref: `docs/brainstorm/configuration.md`).
    - [ ] Field: `SystemPrompt` (string).
    - [ ] Field: `EnvPass` ([]string).
- [ ] **Task Model**: Define `Task` struct (Ref: `docs/brainstorm/features.md`).
    - [ ] YAML tags for frontmatter (`id`, `status`, `assignee`, `tags`).
    - [ ] Internal fields for `Body` (content) and `FilePath`.

### 1.3. Data Layer (`internal/store`)
- [ ] **Interface Definition**: Define `Store` interface (`List`, `Get`, `Create`, `Update`, `Watch`).
- [ ] **FileStore Implementation**:
    - [ ] **Encoding**: Use `adrg/frontmatter` to parse/marshal MD files.
    - [ ] **CRUD**: Implement strict read/write to `.kevin/board/`.
    - [ ] **Watcher**: Implement `fsnotify` loop to watch `.kevin/board/`.
    - [ ] **Debounce**: Add 50-100ms debounce to prevent tearing on rapid writes.
    - [ ] **Event Channel**: Broadcast `TaskUpdateEvent` on change.

### 1.4. CLI Foundation (`internal/cli`)
- [ ] **Root Command**: Setup `cobra` root command.
- [ ] **Init Command** (`kevin init`):
    - [ ] Check if `.kevin` exists.
    - [ ] Create `.kevin/config.yaml` with defaults.
    - [ ] Create `.kevin/board/` directory.
    - [ ] Create a "Welcome" task.
- [ ] **Task Commands**:
    - [ ] `kevin task new <title>`: Create a new file with default frontmatter.
    - [ ] `kevin task list`: Print tabular list of tasks.
    - [ ] `kevin task move <id> <status>`: Update status field in file.

---

## Phase 2: The TUI (Bubble Tea)
**Goal**: A reactive Kanban board that updates instantly on file changes.

### 2.1. View Components (`internal/tui`)
- [ ] **Column Model**:
    - [ ] Render a list of tasks for a specific status (`todo`, `doing`, etc.).
    - [ ] Style using `lipgloss` (colors, borders).
- [ ] **Board Model**:
    - [ ] Layout 3-4 columns horizontally.
    - [ ] Handle window resize (responsive layout).

### 2.2. Interaction Logic
- [ ] **Navigation**:
    - [ ] `h/l`: Switch active column.
    - [ ] `j/k`: Move cursor in column.
- [ ] **Task Actions**:
    - [ ] `enter`: View/Edit task (open `$EDITOR`).
    - [ ] `H/L`: Move task status left/right (write to disk).
    - [ ] `n`: Create new task (input form).

### 2.3. Reactivity (The Loop)
- [ ] **Store Integration**: Inject `Store` into the Bubble Tea model.
- [ ] **Event Handling**:
    - [ ] Listen to `Store.Watch()` channel in a `tea.Cmd`.
    - [ ] On `TaskUpdateEvent`: Reload task list and re-render.

---

## Phase 3: Agent Orchestration
**Goal**: Wrap external CLI tools to execute work against tasks.

### 3.1. Executor System (`internal/executor`)
- [ ] **Interface**: Define `Executor.Run(ctx, cmd, env, workdir)`.
- [ ] **Local Executor**: Implement using `exec.Command` (direct host process).
- [ ] **Docker Executor**:
    - [ ] Implement using Docker CLI or SDK.
    - [ ] Handle volume mounting (`HostDir:MountDir`).

### 3.2. Agent Config & Context
- [ ] **Context Injection**:
    - [ ] Implement template parsing for command strings (`{{.TaskPath}}`, `{{.SystemPrompt}}`).
- [ ] **Environment**:
    - [ ] Read `env_pass` list from config.
    - [ ] Capture values from Host `os.Environ`.

### 3.3. Execution CLI
- [ ] **Command**: `kevin run <task_id>`.
    - [ ] Look up task and assignee.
    - [ ] Look up Agent config.
    - [ ] Prepare Command (Templates) and Env.
    - [ ] Invoke Executor.
    - [ ] Stream Stdout/Stderr to console.

---

## Phase 4: Polish & Parity
- [ ] **Log Streaming**: Show agent logs in a TUI pane (or "ghost" terminal).
- [ ] **Search/Filter**: Filter tasks by tags in TUI.
- [ ] **Cross-Platform**: Verify paths and Docker socket on Windows/Linux.
- [ ] **Documentation**: Generate `README.md` and `CONTRIBUTING.md`.
