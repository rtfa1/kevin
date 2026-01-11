package tui

import (
	"os"
	"os/exec"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rtfa/kevin/internal/agent"
	"github.com/rtfa/kevin/internal/core"
	"github.com/rtfa/kevin/internal/store"
)

type KeyMap struct {
	Quit      key.Binding
	Left      key.Binding
	Right     key.Binding
	Up        key.Binding
	Down      key.Binding
	Enter     key.Binding
	MoveLeft  key.Binding
	MoveRight key.Binding
	Run       key.Binding
}

var Keys = KeyMap{
	Quit:      key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
	Left:      key.NewBinding(key.WithKeys("h", "left"), key.WithHelp("h", "left")),
	Right:     key.NewBinding(key.WithKeys("l", "right"), key.WithHelp("l", "right")),
	Up:        key.NewBinding(key.WithKeys("k", "up"), key.WithHelp("k", "up")),
	Down:      key.NewBinding(key.WithKeys("j", "down"), key.WithHelp("j", "down")),
	Enter:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "open task")),
	MoveLeft:  key.NewBinding(key.WithKeys("H"), key.WithHelp("H", "move left")),
	MoveRight: key.NewBinding(key.WithKeys("L"), key.WithHelp("L", "move right")),
	Run:       key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "run agent")),
}

type Model struct {
	config  *core.ProjectConfig
	store   store.Store
	columns []Column
	active  int // active column index
	width   int
	height  int
	err     error
}

func NewModel(cfg *core.ProjectConfig, s store.Store) Model {
	m := Model{
		config: cfg,
		store:  s,
		columns: []Column{
			NewColumn(core.StatusBacklog),
			NewColumn(core.StatusTodo),
			NewColumn(core.StatusDoing),
			NewColumn(core.StatusDone),
		},
		active: 0,
	}
	m.columns[0].Focused = true
	return m
}

// Reactivity
type TaskUpdateMsg store.TaskUpdateEvent

func waitForActivity(sub <-chan store.TaskUpdateEvent) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-sub
		if !ok {
			return nil
		}
		return TaskUpdateMsg(event)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.reloadTasksCmd,
		waitForActivity(m.store.Watch()),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateSizes()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, Keys.Left):
			m.columns[m.active].Focused = false
			if m.active > 0 {
				m.active--
			}
			m.columns[m.active].Focused = true
		case key.Matches(msg, Keys.Right):
			m.columns[m.active].Focused = false
			if m.active < len(m.columns)-1 {
				m.active++
			}
			m.columns[m.active].Focused = true
		case key.Matches(msg, Keys.Up):
			m.columns[m.active].SelectPrev()
		case key.Matches(msg, Keys.Down):
			m.columns[m.active].SelectNext()

		// --- Interactions ---
		case key.Matches(msg, Keys.Enter):
			selected := m.SelectedTask()
			if selected != nil {
				editor := os.Getenv("EDITOR")
				if editor == "" {
					editor = "vim"
				}
				c := exec.Command(editor, selected.FilePath)
				return m, tea.ExecProcess(c, func(err error) tea.Msg {
					return nil // We rely on FS watch to reload
				})
			}

		case key.Matches(msg, Keys.Run):
			selected := m.SelectedTask()
			if selected != nil && selected.Assignee != "" {
				// Find agent
				var agentCfg *core.AgentConfig
				for _, a := range m.config.Agents {
					if a.Name == selected.Assignee {
						agentCfg = &a
						break
					}
				}

				if agentCfg != nil {
					// Duplicate logic from run.go for MVP simplicity
					// In real app, refactor to shared service
					cwd, _ := os.Getwd()
					cmdSlice, envSlice, err := agent.Prepare(*agentCfg, *selected, cwd)
					if err == nil {
						// Create executable command
						c := exec.Command(cmdSlice[0], cmdSlice[1:]...)
						c.Env = append(os.Environ(), envSlice...)
						c.Dir = cwd // or where?

						return m, tea.ExecProcess(c, func(err error) tea.Msg {
							// After returning, maybe show a status msg?
							return nil
						})
					}
				}
			}

		case key.Matches(msg, Keys.MoveLeft):
			m.moveTask(-1)

		case key.Matches(msg, Keys.MoveRight):
			m.moveTask(1)
		}

	case TaskReloadMsg:
		m.reloadTasks(msg)

	case TaskUpdateMsg:
		// Reload and re-subscribe
		return m, tea.Batch(
			m.reloadTasksCmd,
			waitForActivity(m.store.Watch()),
		)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	cols := []string{}
	for _, c := range m.columns {
		cols = append(cols, c.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, cols...)
}

// Helpers

func (m *Model) SelectedTask() *core.Task {
	col := m.columns[m.active]
	if len(col.Tasks) > 0 && col.cursor < len(col.Tasks) {
		return &col.Tasks[col.cursor]
	}
	return nil
}

func (m *Model) moveTask(direction int) {
	task := m.SelectedTask()
	if task == nil {
		return
	}

	// Find current column index (should match m.active but for safety)
	// Actually we know what `m.active` is.
	// But we need to know what the next status is.

	newIdx := m.active + direction
	if newIdx < 0 || newIdx >= len(m.columns) {
		return // Can't move
	}

	task.Status = m.columns[newIdx].Status

	// Update via Store
	// Note: We ignore error for MVP, ideally show toast
	_ = m.store.Update(*task)
	// Reactivity will update UI
}

func (m *Model) updateSizes() {
	colWidth := m.width / len(m.columns)
	// subtract borders/margins from height?
	colHeight := m.height - 2

	for i := range m.columns {
		m.columns[i].SetSize(colWidth-2, colHeight) // -2 for margin
	}
}

// Data Loading

type TaskReloadMsg []core.Task

func (m Model) reloadTasksCmd() tea.Msg {
	tasks, err := m.store.List()
	if err != nil {
		return nil // TODO: handle error
	}
	return TaskReloadMsg(tasks)
}

func (m *Model) reloadTasks(tasks []core.Task) {
	// clear all
	for i := range m.columns {
		m.columns[i].Tasks = []core.Task{}
	}

	// distribute
	for _, t := range tasks {
		for i, c := range m.columns {
			if c.Status == t.Status {
				m.columns[i].Tasks = append(m.columns[i].Tasks, t)
				break
			}
		}
	}
}
