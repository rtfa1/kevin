# Technology Stack Recommendations (Finalized)

## 1. CLI / TUI Framework
**Recommendation: [Bubble Tea](https://github.com/charmbracelet/bubbletea)**

*   **Why**: Industry standard for "modern" CLI tools. The MVU (Model-View-Update) architecture fits our "Reactive" requirements perfectly.
*   **Styling**: **[Lip Gloss](https://github.com/charmbracelet/lipgloss)**.
*   **CLI Commands**: **[Cobra](https://github.com/spf13/cobra)**.
    *   *Integration*: We will use Cobra for the command structure (`kevin init`, `kevin task`) and launch Bubble Tea programs from within those commands.

## 2. Execution / Isolation Strategy
**Recommendation: Pluggable `Executor` Interface**

*   **Concept**: A generic `Executor` interface allowing different runtimes.
*   **Implementations**:
    *   **Docker**: For isolated, reproducible agent runs (using `docker run`).
    *   **Local**: For running simple scripts or local tools directly (`os/exec`).
*   **Configuration**: Defined in `.kevin/config.yaml`.

## 3. Agent Integration Strategy
**Recommendation: External CLI Orchestration**

*   **Approach**: "Kevin as Orchestrator". No internal AI logic.
*   **Tools**:
    *   `opencode`
    *   `claude-cli`
    *   `aider`
    *   `codex-cli`
*   **Mechanism**:
    *   Kevin constructs the command string (e.g., `opencode --task task.md`).
    *   Kevin passes environment variables (Auth tokens) from the host to the process.

## 4. Data Layer & Persistence
**Recommendation: Reactive Filesystem Store**

*   **Source of Truth**: `.kevin/board/` (Markdown files).
*   **Library**: **[fsnotify](https://github.com/fsnotify/fsnotify)**.
    *   *Why*: Immediate TUI updates when files are changed manually by the user or by an agent.
*   **Parsing**: **[frontmatter](https://github.com/adrg/frontmatter)** for metadata, standard Go string manipulation for body.
*   **Structure**:
    ```text
    .kevin/
      config.yaml
      board/
        task-001.md
    ```
