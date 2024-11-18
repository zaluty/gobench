package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher handles file system events for test files
type Watcher struct {
	watcher *fsnotify.Watcher
	dirs    []string
	events  chan Event
}

// Event represents a file system event that should trigger a benchmark run
type Event struct {
	Path string
	Type string
}

// NewWatcher creates a new file system watcher
func NewWatcher(dirs []string) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	w := &Watcher{
		watcher: fsWatcher,
		dirs:    dirs,
		events:  make(chan Event),
	}

	return w, nil
}

// Start begins watching for file changes
func (w *Watcher) Start() error {
	fmt.Println("Setting up file watchers...")
	// Add all directories to the watcher
	for _, dir := range w.dirs {
		fmt.Printf("Adding directory to watch: %s\n", dir)
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				fmt.Printf("Watching directory: %s\n", path)
				return w.watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to add directory %s to watcher: %w", dir, err)
		}
	}

	// Start watching for events
	go w.watch()

	return nil
}

// Stop stops watching for file changes
func (w *Watcher) Stop() error {
	close(w.events)
	return w.watcher.Close()
}

// Events returns the channel of file system events
func (w *Watcher) Events() <-chan Event {
	return w.events
}

// watch processes file system events
func (w *Watcher) watch() {
	fmt.Println("Starting file watcher...")
	// Use a timer to debounce events
	var timer *time.Timer
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Only care about Go files
			if !strings.HasSuffix(event.Name, ".go") {
				continue
			}

			fmt.Printf("Detected change in file: %s\n", event.Name)

			// Reset or create timer
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(100*time.Millisecond, func() {
				w.events <- Event{
					Path: event.Name,
					Type: event.Op.String(),
				}
			})

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}
