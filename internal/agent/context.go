package agent

// ExecutionContext defines the variables available to agent command templates
type ExecutionContext struct {
	TaskID     string
	TaskTitle  string
	TaskStatus string
	TaskPath   string // Absolute path to the task file
	ProjectDir string // Root of the project
}
