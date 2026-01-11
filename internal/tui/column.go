package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/rtfa/kevin/internal/core"
)

// Column represents a single status column in the Kanban board
type Column struct {
	Status  core.TaskStatus
	Tasks   []core.Task
	Focused bool
	width   int
	height  int

	// We might use a bubbletea List component in the future,
	// but for now let's implement simple rendering to understand the basics
	cursor int
}

func NewColumn(status core.TaskStatus) Column {
	return Column{
		Status: status,
		Tasks:  []core.Task{},
	}
}

func (c *Column) SetSize(w, h int) {
	c.width = w
	c.height = h
}

// Update handles events for the column
func (c *Column) Update(msg tea.Msg) (Column, tea.Cmd) {
	// Simple navigation for now
	// In a real app we'd bubble key msgs here only if focused
	return *c, nil
}

// SelectNext moves the cursor down
func (c *Column) SelectNext() {
	if c.cursor < len(c.Tasks)-1 {
		c.cursor++
	}
}

// SelectPrev moves the cursor up
func (c *Column) SelectPrev() {
	if c.cursor > 0 {
		c.cursor--
	}
}

func (c Column) View() string {
	// Header
	header := TitleStyle.Render(string(c.Status))

	// Tasks
	var lines []string
	for i, t := range c.Tasks {
		style := TaskStyle
		if c.Focused && i == c.cursor {
			style = SelectedTaskStyle
		}

		// Truncate title if too long
		// This is naive; normally we'd check width
		title := t.Title
		lines = append(lines, style.Render(title))
	}

	if len(c.Tasks) == 0 {
		lines = append(lines, TaskStyle.Foreground(ColorDim).Render("(empty)"))
	}

	// Let's join manually
	var content string
	for _, l := range lines {
		content += l + "\n"
	}

	// Check height - strictly for MVP just clip?
	// For MVP let's just make it scrollable later.

	// Style the container
	style := ColumnStyle
	if c.Focused {
		style = FocusedColumnStyle
	}

	return style.
		Width(c.width).
		Height(c.height).
		Render(
			fmt.Sprintf("%s\n\n%s", header, content),
		)
}
