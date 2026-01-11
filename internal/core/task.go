package core

import (
	"time"
)

// TaskStatus defines the possible states of a task
type TaskStatus string

const (
	StatusBacklog TaskStatus = "backlog"
	StatusTodo    TaskStatus = "todo"
	StatusDoing   TaskStatus = "doing"
	StatusDone    TaskStatus = "done"
)

// Task represents a unit of work stored in .kevin/board/*.md
type Task struct {
	ID                string     `yaml:"id"`
	Title             string     `yaml:"title"`
	Status            TaskStatus `yaml:"status"`
	Assignee          string     `yaml:"assignee,omitempty"`
	Priority          string     `yaml:"priority,omitempty"`
	Tags              []string   `yaml:"tags,omitempty"`
	Created           time.Time  `yaml:"created"`
	SysPromptOverride string     `yaml:"sys_prompt_override,omitempty"`

	// Content holds the markdown body content (excluding frontmatter)
	Content string `yaml:"-"`

	// FilePath is the location of the task on disk
	FilePath string `yaml:"-"`
}
