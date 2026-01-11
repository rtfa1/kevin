package store

import (
	"github.com/rtfa/kevin/internal/core"
)

// EventType defines the type of change observed
type EventType string

const (
	EventCreate EventType = "create"
	EventUpdate EventType = "update"
	EventDelete EventType = "delete"
)

// TaskUpdateEvent represents a change in the task board
type TaskUpdateEvent struct {
	TaskID string
	Type   EventType
}

// Store defines the interface for data persistence and reactivity
type Store interface {
	// CRUD
	List() ([]core.Task, error)
	Get(id string) (*core.Task, error)
	Create(task core.Task) error
	Update(task core.Task) error
	Delete(id string) error

	// Reactivity
	Watch() <-chan TaskUpdateEvent
}
