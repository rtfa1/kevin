# Data Layer & File Watching Strategy (Finalized)

## Core Decision: The Reactive Filesystem Store
Kevin uses the filesystem (`.kevin/board/`) as the Database of Record.
*   **Single Source of Truth**: The Markdown files.
*   **Bi-directional**: Changes can happen from *within* (Kevin/Agents) or *without* (User/Vim).

## 1. Directory Structure
```text
.kevin/
  board/          # Watched Directory
    task-001.md
    task-002.md
```

## 2. Technology Stack
*   **Watcher**: `github.com/fsnotify/fsnotify` (Recursive watch on `.kevin/board`).
*   **Parser**: `github.com/adrg/frontmatter` (YAML Header + Body).

## 3. Synchronization Logic

### A. The "Vim" Path (External Edit)
1.  User edits `task-001.md` in Vim and saves.
2.  `fsnotify` emits `WRITE` event for `.kevin/board/task-001.md`.
3.  **Controller**:
    *   Receives event.
    *   Debounces (wait 50-100ms to allow flush).
    *   Reads file from disk.
    *   Parses Frontmatter into `Task` struct.
    *   Updates internal State.
    *   `Program.Send(MsgTaskUpdated)` -> TUI re-renders.

### B. The "Agent" Path (Internal Edit)
1.  Agent (running via Executor) appends a log line to `task-001.md`.
2.  `fsnotify` emits `WRITE` event.
3.  **Controller** treats this exactly like a User edit (See A).
    *   *Benefit*: We don't need complex internal event buses for agents. We just watch the file they are writing to.

## 4. Conflict handling
*   **Strategy**: Last Write Wins (Filesystem standard).
*   **Optimistic Updates**: If the TUI modifies a task (e.g., dragging column), it:
    1.  Reads the file to ensure freshness.
    2.  Updates the field.
    3.  Writes back to disk.
*   If the user was editing the *Description* while the TUI updated the *Status*, a rewrite might overwrite the user's unsaved buffer in Vim (Vim usually warns "File changed on disk").
*   *Mitigation*: We only parse/rewrite Frontmatter for status changes, appending to body where possible.

## 5. Interface Definition
```go
type Store interface {
    // Operations
    List() ([]Task, error)
    Get(id string) (Task, error)
    Update(task Task) error // Writes to disk
    Create(task Task) error // Creates new file

    // Reactive
    Watch(ctx context.Context) (<-chan TaskEvent, error)
}
```
