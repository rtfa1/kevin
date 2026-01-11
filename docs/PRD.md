# Product Requirements Document (PRD) - Project Kevin

## 1. Executive Summary
**Kevin** is a TUI-first, CLI-compatible Task Orchestrator and Kanban board for developers. Inspired by the helpful Minions, Kevin manages the tedious tracking of tasks and orchestrates external "Agent" tools to execute work.

Unlike traditional AI coding tools that lock you into their ecosystem, Kevin is an **Orchestrator**: he connects your existing tools (OpenCode, Claude CLI, Aider) to a structured workflow, using the filesystem as the database.

## 2. Product Principles
1.  **Filesystem as Database**: All state lives in `.kevin/board/*.md`. No hidden SQlite, no cloud sync logic in the core. If you delete the file, the task is gone.
2.  **Headless Parity**: Every action available in the TUI (Bubble Tea) must be executable via CLI flags (Cobra). This enables scripting.
3.  **No Internal AI**: Kevin does not import LLM SDKs. He wraps external binaries/containers.
4.  **Config Driven**: Project-specific behavior is defined in `.kevin/config.yaml`.
5.  **Reactive**: The UI updates instantly when files change on disk (bi-directional sync).

## 3. Functional Requirements

### 3.1. The Kanban Board (TUI)
*   **Views**: A standard Kanban layout with columns: `Backlog`, `Todo`, `Doing`, `Done`.
*   **Navigation**: Vim-style keys (`h`, `j`, `k`, `l`) to navigate columns and items.
*   **Actions**:
    *   Create new task.
    *   Move task between columns.
    *   Edit task details (opens default `$EDITOR`).
    *   Stream logs from running agents.

### 3.2. CLI Commands
The application must provide the following commands:
*   `kevin init`: Bootstrap the `.kevin/` directory and default config.
*   `kevin task new <title>`: Create a task.
*   `kevin task list`: Dump tasks as JSON/Table.
*   `kevin task move <id> <status>`: Update status.
*   `kevin run <id>`: Trigger the assigned agent for a task.

### 3.3. Agent Orchestration
Kevin acts as a wrapper around external tools.
*   **Definition**: Agents are configured in `config.yaml` with a command template.
*   **Execution**:
    *   Supports **Docker** (isolated) and **Local** (shell) execution.
    *   Injects context via template variables: `{{.TaskPath}}`, `{{.SystemPrompt}}`.
    *   Passes allow-listed Environment Variables (e.g., `OPENAI_API_KEY`) from the host.
*   **System Prompts**: Users can define per-agent system prompts in the config to specialize generic tools (e.g., "Senior Go Dev" vs "Documentation Expert").

### 3.4. Data Layer & Persistence
*   **Storage**: Tasks are Markdown files in `.kevin/board/`.
*   **Metadata**: stored in YAML Frontmatter.
*   **Reactivity**: The application must watch the directory for changes using `fsnotify` to handle concurrent edits (e.g., User editing in Vim while Kevin is running).

## 4. Technical Architecture

### 4.1. Directory Structure
```text
ProjectRoot/
├── .kevin/
│   ├── config.yaml       # Project & Agent Config
│   ├── board/            # The Database
│   │   ├── task-001.md
│   │   └── task-002.md
│   └── context/          # (Optional) Shared context/prompts
└── ... (User Code)
```

### 4.2. Task Schema
```yaml
---
id: "task-001"
title: "Fix Login Handler"
status: "todo"        # backlog, todo, doing, done
assignee: "coder"     # references agent key in config
priority: "high"
tags: ["bug", "auth"]
created: "2026-01-10T12:00:00Z"
sys_prompt_override: "Use strict TDD." # (Optional)
---

# Description
...
```

### 4.3. Configuration Schema
```yaml
project:
  name: "My App"

agents:
  - name: "coder"
    executor: "docker"
    image: "ghcr.io/opencode/opencode:latest"
    command: ["opencode", "--task", "{{.TaskPath}}"]
    env_pass: ["OPENAI_API_KEY"]
    system_prompt: "You are an expert Go developer."
```

## 5. Non-Functional Requirements
*   **Performance**: TUI startup < 100ms. File watcher latency < 200ms.
*   **Compatibility**: macOS, Linux, Windows (via WSL or PowerShell).
*   **Dependencies**: Minimal runtime dependencies (single binary preferred, relies on system Docker/Git).

## 6. Success Metrics (MVP)
1.  User can initialize a repo with `kevin init`.
2.  User can create a task via CLI.
3.  User can see the task in TUI.

## 7. Development Guidelines
*   **Testing**: Create tests when possible and makes sense.

