package analyzer

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/st3v3nmw/sourcerer-mcp/internal/fs"
	"github.com/st3v3nmw/sourcerer-mcp/internal/index"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

type Analyzer struct {
	workspaceRoot string
	parsers       map[Language]parser.Parser
	index         *index.Index
	watcher       *fs.Watcher
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

	analyzer.indexWorkspace(ctx)

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

func (a *Analyzer) indexWorkspace(ctx context.Context) {
	var filesToProcess []string

	fs.WalkSourceFiles(a.workspaceRoot, languages.supportedExts(), func(filePath string) error {
		if a.index.ShouldIndex(filePath) {
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
	for _, filePath := range filePaths {
		a.chunk(ctx, filePath)
	}
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

func (a *Analyzer) Close() {
	if a.watcher != nil {
		a.watcher.Close()
	}

	for _, parser := range a.parsers {
		parser.Close()
	}
}
