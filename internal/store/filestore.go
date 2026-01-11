package store

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/fsnotify/fsnotify"
	"github.com/rtfa/kevin/internal/core"
	"gopkg.in/yaml.v3"
)

// FileStore implements Store using the filesystem
type FileStore struct {
	baseDir string
	watcher *fsnotify.Watcher
	updates chan TaskUpdateEvent
	mu      sync.RWMutex // Protects concurrent reads if needed, though most OS ops are safe
}

// NewFileStore creates a new store watching the given directory
func NewFileStore(dir string) (*FileStore, error) {
	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create board directory: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	fs := &FileStore{
		baseDir: dir,
		watcher: watcher,
		updates: make(chan TaskUpdateEvent, 100),
	}

	if err := watcher.Add(dir); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to watch directory: %w", err)
	}

	go fs.watchLoop()

	return fs, nil
}

func (s *FileStore) Watch() <-chan TaskUpdateEvent {
	return s.updates
}

func (s *FileStore) List() ([]core.Task, error) {
	entries, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}

	var tasks []core.Task
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		task, err := s.readTask(filepath.Join(s.baseDir, entry.Name()))
		if err != nil {
			// Skip malformed files but maybe log them?
			continue
		}
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

func (s *FileStore) Get(id string) (*core.Task, error) {
	// We assume ID matches filename for simplicity in MVP: "task-001" -> "task-001.md"
	// Or we scan. For performance, let's scan or enforce a convention.
	// The PRD implies .kevin/board/task-001.md
	filename := fmt.Sprintf("%s.md", id)
	path := filepath.Join(s.baseDir, filename)

	return s.readTask(path)
}

func (s *FileStore) Create(task core.Task) error {
	filename := fmt.Sprintf("%s.md", task.ID)
	path := filepath.Join(s.baseDir, filename)
	return s.writeTask(path, task)
}

func (s *FileStore) Update(task core.Task) error {
	return s.Create(task) // Same as create for file overwrite
}

func (s *FileStore) Delete(id string) error {
	filename := fmt.Sprintf("%s.md", id)
	path := filepath.Join(s.baseDir, filename)
	return os.Remove(path)
}

func (s *FileStore) Close() error {
	return s.watcher.Close()
}

// Internal helpers

func (s *FileStore) readTask(path string) (*core.Task, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var task core.Task
	rest, err := frontmatter.Parse(f, &task)
	if err != nil {
		return nil, err
	}

	task.Content = string(rest)
	task.FilePath = path

	// Ensure ID is set from filename if missing or to guarantee consistency
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	task.ID = strings.TrimSuffix(base, ext)

	return &task, nil
}

func (s *FileStore) writeTask(path string, task core.Task) error {
	data, err := yaml.Marshal(&task)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	buffer.WriteString("---\n")
	buffer.Write(data)
	buffer.WriteString("---\n\n")
	buffer.WriteString(task.Content)

	return os.WriteFile(path, buffer.Bytes(), 0644)
}

func (s *FileStore) watchLoop() {
	// Simple map to debounce events per file
	timers := make(map[string]*time.Timer)
	var timersMu sync.Mutex
	const debounceDuration = 100 * time.Millisecond

	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			// Only care about markdown files
			if !strings.HasSuffix(event.Name, ".md") {
				continue
			}

			// Determine event type
			var et EventType
			if event.Op&fsnotify.Create == fsnotify.Create {
				et = EventCreate
			} else if event.Op&fsnotify.Write == fsnotify.Write {
				et = EventUpdate
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				et = EventDelete
			} else if event.Op&fsnotify.Rename == fsnotify.Rename {
				et = EventDelete // Treated as delete for the old name
			} else {
				continue
			}

			timersMu.Lock()
			if t, exists := timers[event.Name]; exists {
				t.Stop()
			}

			// Get task ID from filename
			base := filepath.Base(event.Name)
			id := strings.TrimSuffix(base, filepath.Ext(base))

			timers[event.Name] = time.AfterFunc(debounceDuration, func() {
				timersMu.Lock()
				delete(timers, event.Name)
				timersMu.Unlock()

				s.updates <- TaskUpdateEvent{
					TaskID: id,
					Type:   et,
				}
			})
			timersMu.Unlock()

		case _, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
			// log error?
		}
	}
}
