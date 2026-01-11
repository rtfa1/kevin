# Features & Requirements (Finalized)

## 1. Interaction Modes: TUI & CLI Parity
*   **Philosophy**: "Headless First". Everything possible in the UI is possible via flags.
*   **TUI (Bubble Tea)**:
    *   Horizontal Kanban Layout (Backlog, Todo, In Progress, Done).
    *   Vim-style navigation (`h`, `j`, `k`, `l`).
    *   Live log streaming from agents.
*   **CLI (Cobra)**:
    *   `kevin init`: Bootstrap `.kevin/` folder.
    *   `kevin task new/move/edit`.
    *   `kevin run <ID>`: Trigger agent execution headless.

## 2. Agent Orchestration
*   **External Tools**: We do not implement AI. We wrap tools like `opencode`, `aider`, `claude-cli`.
*   **Execution**:
    *   **Pluggable Runtimes**: Docker, Local, or SSH.
    *   **Environment**: Secrets (`OPENAI_API_KEY`) are passed from host to agent.
    *   **Context**: Workspaces are mounted R/W to the agent.
*   **Configuration**:
    *   Defined in `.kevin/config.yaml`.
    *   Supports **System Prompts** per agent profile.

## 3. Data & Storage
*   **Location**: `.kevin/board/*.md`.
*   **Format**: Markdown with Frontmatter.
*   **Reactivity**:
    *   App watches file changes (`fsnotify`).
    *   Manual edits in `vim` update the TUI instantly.
    *   Agent writes update the TUI instantly.

## 4. Task Schema (Frontmatter)
Example `.kevin/board/task-001.md`:

```markdown
---
id: "task-001"
title: "Implement Login Handler"
status: "todo"        # backlog, todo, doing, done
assignee: "coder"     # references agent name in config.yaml
tags: ["backend", "auth"]
priority: "high"
created: "2026-01-10T10:00:00Z"
---

# Description
We need a standard login handler using OAuth2.

## Context
- [ ] Returns 200 OK on success
- [ ] Returns JWT token
```
