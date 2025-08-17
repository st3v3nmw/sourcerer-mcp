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
	files         map[string][]string // file path -> chunk paths

	index               *index.Index
	initialIndexingDone chan struct{}
}

func New(ctx context.Context, workspaceRoot string) (*Analyzer, error) {
	index, err := index.New(workspaceRoot)
	if err != nil {
		return nil, err
	}

	analyzer := &Analyzer{
		workspaceRoot:       workspaceRoot,
		parsers:             map[Language]parser.Parser{},
		files:               map[string][]string{},
		index:               index,
		initialIndexingDone: make(chan struct{}),
	}

	fileChunks := index.GetAllChunkIDs(ctx)
	for filePath, chunkIDs := range fileChunks {
		chunkPaths := make([]string, len(chunkIDs))
		for i, id := range chunkIDs {
			parts := strings.SplitN(id, "::", 2)
			if len(parts) == 2 {
				chunkPaths[i] = parts[1]
			}
		}
		analyzer.files[filePath] = chunkPaths
	}

	analyzer.indexWorkspace(ctx)

	return analyzer, nil
}

func (a *Analyzer) indexWorkspace(ctx context.Context) {
	fs.WalkSourceFiles(a.workspaceRoot, languages.supportedExts(), func(filePath string) error {
		a.chunk(ctx, filePath)
		return nil
	})

	close(a.initialIndexingDone)
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
	if paths, exists := a.files[filePath]; exists {
		allFresh := true
		for _, chunkPath := range paths {
			id := filePath + "::" + chunkPath
			isStale, err := a.index.IsChunkStale(ctx, id)
			if err != nil || isStale {
				allFresh = false
				break
			}
		}

		if allFresh {
			return nil
		}
	}

	parser, err := a.getParser(filePath)
	if err != nil {
		return err
	}

	file, err := parser.Chunk(filePath)
	if err != nil {
		return err
	}

	err = a.index.Upsert(ctx, file)
	if err != nil {
		return err
	}

	paths := make([]string, len(file.Chunks))
	for i, chunk := range file.Chunks {
		paths[i] = chunk.Path
	}

	a.files[filePath] = paths

	return nil
}

func (a *Analyzer) getTOC(ctx context.Context, filePath string) string {
	a.chunk(ctx, filePath)

	toc := fmt.Sprintf("== %s ==\n\n", filePath)

	paths, exists := a.files[filePath]
	if !exists {
		return toc + "<file not found or could not be processed>\n\n"
	}

	for _, chunkPath := range paths {
		id := filePath + "::" + chunkPath
		chunk, err := a.index.GetChunk(ctx, id)
		if err != nil {
			continue
		}

		if chunk.Path != "" {
			toc += "::" + chunk.Path + "\n"
		}

		toc += chunk.Summary + "\n\n"
	}

	return toc
}

func (a *Analyzer) GetTOCs(ctx context.Context, filePaths []string) string {
	tocs := ""
	for _, filePath := range filePaths {
		tocs += a.getTOC(ctx, filePath)
		tocs += "\n"
	}

	return tocs
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

func (a *Analyzer) GetChunkSources(ctx context.Context, ids []string) string {
	chunks := ""
	for _, id := range ids {
		chunks += a.getChunkSource(ctx, id)
	}

	return chunks
}

func (a *Analyzer) SemanticSearch(ctx context.Context, query string) ([]string, error) {
	<-a.initialIndexingDone
	return a.index.Search(ctx, query)
}

func (a *Analyzer) Close() {
	for _, parser := range a.parsers {
		parser.Close()
	}
}
