package fs

import (
	"context"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	debounceDuration = 60 * time.Second
)

type FileChangeHandler func(ctx context.Context, filePaths []string)

type Watcher struct {
	workspaceRoot    string
	filter           *FileFilter
	handler          FileChangeHandler
	fsWatcher        *fsnotify.Watcher
	debounceTimer    *time.Timer
	pendingFiles     map[string]bool
	mu               sync.RWMutex
	debounceDuration time.Duration
	ctx              context.Context
	cancel           context.CancelFunc
}

func NewWatcher(ctx context.Context, workspaceRoot string, supportedExts []string, handler FileChangeHandler) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)
	w := &Watcher{
		workspaceRoot:    workspaceRoot,
		filter:           NewFileFilter(workspaceRoot, supportedExts),
		handler:          handler,
		fsWatcher:        fsWatcher,
		pendingFiles:     map[string]bool{},
		debounceDuration: debounceDuration,
		ctx:              ctx,
		cancel:           cancel,
	}

	err = w.addWatchers()
	if err != nil {
		fsWatcher.Close()
		cancel()

		return nil, err
	}

	go w.watch()

	return w, nil
}

func (w *Watcher) addWatchers() error {
	var supportedExts []string
	for ext := range w.filter.supportedExts {
		supportedExts = append(supportedExts, ext)
	}

	return WalkSourceFiles(w.workspaceRoot, supportedExts, func(filePath string) error {
		dir := filepath.Dir(filepath.Join(w.workspaceRoot, filePath))
		return w.fsWatcher.Add(dir)
	})
}

func (w *Watcher) watch() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			w.handleEvent(event)
		case _, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.shouldIgnoreEvent(event) {
		return
	}

	relPath, err := filepath.Rel(w.workspaceRoot, event.Name)
	if err != nil {
		return
	}

	w.pendingFiles[relPath] = true

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(w.debounceDuration, w.processPendingFiles)
}

func (w *Watcher) shouldIgnoreEvent(event fsnotify.Event) bool {
	if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove) == 0 {
		return true
	}

	return w.filter.ShouldIgnore(event.Name)
}

func (w *Watcher) processPendingFiles() {
	w.mu.Lock()
	defer w.mu.Unlock()

	changes := make([]string, 0, len(w.pendingFiles))
	for filePath := range w.pendingFiles {
		changes = append(changes, filePath)
	}

	if len(changes) > 0 {
		w.handler(w.ctx, changes)
	}

	w.pendingFiles = map[string]bool{}
}

func (w *Watcher) FlushPending() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.debounceTimer != nil {
		w.debounceTimer.Reset(0)
	}
}

func (w *Watcher) PendingCount() int {
	w.mu.Lock()
	defer w.mu.Unlock()

	return len(w.pendingFiles)
}

func (w *Watcher) Close() error {
	w.cancel()

	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	return w.fsWatcher.Close()
}
