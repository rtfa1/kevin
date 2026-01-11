# CLI Commands & Initialization

## Core Philosophy: CLI/TUI Parity
Every action available in the TUI must be performable via a headless CLI command. This allows scripting, automation, and quick interactions without entering the full UI.

## File Structure: The `.kevin` Directory
All project-specific data lives in a hidden directory at the project root.

```text
ProjectRoot/
├── .kevin/
│   ├── config.yaml       # Project configuration (Executor, Agents)
│   ├── board/            # The Kanban Board data
│   │   ├── task-001.md
│   │   └── task-002.md
│   └── .cache/           # Temporary agent artifacts (logs, diffs)
└── (User Files)
```

## Command Reference

### 1. Initialization
**Command**: `kevin init`
*   **Action**:
    1.  Checks if `.kevin` exists (Error if yes).
    2.  Creates `.kevin/`.
    3.  Creates default `config.yaml` (asking simple questions or defaults).
    4.  Creates `.kevin/board/` with a welcome task.
    5.  Optionally adds `.kevin` to `.gitignore`.

### 2. TUI Mode
**Command**: `kevin` (or `kevin board`)
*   **Action**: Launches the Bubble Tea interface.

### 3. Task Management (Headless)
*   `kevin task ls`: List tasks (JSON or Table output).
*   `kevin task new "Title" --status todo`: Create a new task.
*   `kevin task show <ID>`: Print task details to stdout.
*   `kevin task move <ID> <COLUMN>`: Move a task (e.g., `kevin task move 1 done`).
*   `kevin task edit <ID>`: Opens `$EDITOR` on the markdown file.

### 4. Agent Interaction
*   `kevin run <ID>`: Manually trigger the assigned agent for a task.
*   `kevin logs <ID>`: Tail the logs of a running agent.
