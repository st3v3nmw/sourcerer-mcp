package analyzer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/st3v3nmw/sourcerer-mcp/internal/fs"
	"github.com/st3v3nmw/sourcerer-mcp/internal/index"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

type Analyzer struct {
	workspaceRoot string
	parsers       map[Language]parser.Parser
	watcher       *fs.Watcher

	index         *index.Index
	indexMu       sync.RWMutex
	nPendingFiles int
	lastIndexedAt time.Time
}

func New(ctx context.Context, workspaceRoot string) (*Analyzer, error) {
	index, err := index.New(ctx, workspaceRoot)
	if err != nil {
		return nil, err
	}

	analyzer := &Analyzer{
		workspaceRoot: workspaceRoot,
		parsers:       map[Language]parser.Parser{},
		index:         index,
	}

	go analyzer.IndexWorkspace(ctx)

	w, err := fs.NewWatcher(
		ctx,
		workspaceRoot,
		languages.supportedExts(),
		analyzer.handleFileChange,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	analyzer.watcher = w
	return analyzer, nil
}

func (a *Analyzer) IndexWorkspace(ctx context.Context) {
	a.flushPendingChanges()

	var filesToProcess []string
	fs.WalkSourceFiles(a.workspaceRoot, languages.supportedExts(), func(filePath string) error {
		if a.index.IsStale(filePath) {
			filesToProcess = append(filesToProcess, filePath)
		}

		return nil
	})

	a.processFiles(ctx, filesToProcess)
}

func (a *Analyzer) handleFileChange(ctx context.Context, filePaths []string) {
	a.processFiles(ctx, filePaths)
}

func (a *Analyzer) processFiles(ctx context.Context, filePaths []string) {
	if len(filePaths) == 0 {
		return
	}

	a.nPendingFiles = len(filePaths)
	for _, filePath := range filePaths {
		a.chunk(ctx, filePath)

		a.indexMu.Lock()
		a.nPendingFiles = max(a.nPendingFiles-1, 0)
		a.indexMu.Unlock()
	}

	a.indexMu.Lock()
	a.lastIndexedAt = time.Now()
	a.indexMu.Unlock()
}

func (a *Analyzer) getParser(filePath string) (parser.Parser, error) {
	lang := languages.detect(filepath.Ext(filePath))
	parser, exists := a.parsers[lang]
	if exists {
		return parser, nil
	}

	return languages.createParser(a.workspaceRoot, lang)
}

func (a *Analyzer) chunk(ctx context.Context, filePath string) error {
	parser, err := a.getParser(filePath)
	if err != nil {
		return err
	}

	file, err := parser.Chunk(filePath)
	if err != nil {
		return err
	}

	err = a.index.Index(ctx, file)
	if err != nil {
		return err
	}

	return nil
}

func (a *Analyzer) SemanticSearch(ctx context.Context, query string) ([]string, error) {
	a.flushPendingChanges()
	return a.index.Search(ctx, query)
}

func (a *Analyzer) flushPendingChanges() {
	if a.watcher != nil {
		a.watcher.FlushPending()
	}
}

func (a *Analyzer) GetChunkSources(ctx context.Context, ids []string) string {
	chunks := ""
	for _, id := range ids {
		chunks += a.getChunkSource(ctx, id)
	}

	return chunks
}

func (a *Analyzer) getChunkSource(ctx context.Context, id string) string {
	parts := strings.SplitN(id, "::", 2)
	if len(parts) != 2 {
		return fmt.Sprintf("== %s ==\n\n<invalid chunk id>\n\n", id)
	}

	err := a.chunk(ctx, parts[0])
	if err != nil {
		return fmt.Sprintf("== %s ==\n\n<processing error: %v>\n\n", id, err)
	}

	chunk, err := a.index.GetChunk(ctx, id)
	if err != nil {
		return fmt.Sprintf("== %s ==\n\n<source not found for chunk>\n\n", id)
	}

	return fmt.Sprintf("== %s ==\n\n%s\n\n", id, chunk.Source)
}

func (a *Analyzer) GetIndexStatus() (int, time.Time) {
	a.indexMu.RLock()
	pendingFiles := a.nPendingFiles
	lastIndexedAt := a.lastIndexedAt
	a.indexMu.RUnlock()

	if a.watcher != nil {
		pendingFiles += a.watcher.PendingCount()
	}

	return pendingFiles, lastIndexedAt
}

func (a *Analyzer) Close() {
	if a.watcher != nil {
		a.watcher.Close()
	}

	for _, parser := range a.parsers {
		parser.Close()
	}
}
