package editor

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/st3v3nmw/sourcerer-mcp/internal/fs"
	"github.com/st3v3nmw/sourcerer-mcp/internal/index"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

type Editor struct {
	workspaceRoot string
	parsers       map[Language]parser.Parser
	files         map[string][]string // file path -> chunk paths

	index               *index.Index
	initialIndexingDone chan struct{}
}

func New(ctx context.Context, workspaceRoot string) (*Editor, error) {
	index, err := index.New(workspaceRoot)
	if err != nil {
		return nil, err
	}

	editor := &Editor{
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
		editor.files[filePath] = chunkPaths
	}

	editor.indexWorkspace(ctx)

	return editor, nil
}

func (e *Editor) indexWorkspace(ctx context.Context) {
	fs.WalkSourceFiles(e.workspaceRoot, languages.supportedExts(), func(filePath string) error {
		e.chunk(ctx, filePath)
		return nil
	})

	close(e.initialIndexingDone)
}

func (e *Editor) getParser(filePath string) (parser.Parser, error) {
	lang := languages.detect(filepath.Ext(filePath))
	parser, exists := e.parsers[lang]
	if exists {
		return parser, nil
	}

	return languages.createParser(e.workspaceRoot, lang)
}

func (e *Editor) chunk(ctx context.Context, filePath string) error {
	if paths, exists := e.files[filePath]; exists {
		allFresh := true
		for _, chunkPath := range paths {
			id := filePath + "::" + chunkPath
			isStale, err := e.index.IsChunkStale(ctx, id)
			if err != nil || isStale {
				allFresh = false
				break
			}
		}

		if allFresh {
			return nil
		}
	}

	parser, err := e.getParser(filePath)
	if err != nil {
		return err
	}

	file, err := parser.Chunk(filePath)
	if err != nil {
		return err
	}

	err = e.index.Upsert(ctx, file)
	if err != nil {
		return err
	}

	paths := make([]string, len(file.Chunks))
	for i, chunk := range file.Chunks {
		paths[i] = chunk.Path
	}

	e.files[filePath] = paths

	return nil
}

func (e *Editor) getTOC(ctx context.Context, filePath string) string {
	e.chunk(ctx, filePath)

	toc := fmt.Sprintf("== %s ==\n\n", filePath)

	paths, exists := e.files[filePath]
	if !exists {
		return toc + "<file not found or could not be processed>\n\n"
	}

	for _, chunkPath := range paths {
		id := filePath + "::" + chunkPath
		chunk, err := e.index.GetChunk(ctx, id)
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

func (e *Editor) GetTOCs(ctx context.Context, filePaths []string) string {
	tocs := ""
	for _, filePath := range filePaths {
		tocs += e.getTOC(ctx, filePath)
		tocs += "\n"
	}

	return tocs
}

func (e *Editor) getChunkSource(ctx context.Context, id string) string {
	parts := strings.SplitN(id, "::", 2)
	if len(parts) != 2 {
		return fmt.Sprintf("== %s ==\n\n<invalid chunk id>\n\n", id)
	}

	err := e.chunk(ctx, parts[0])
	if err != nil {
		return fmt.Sprintf("== %s ==\n\n<processing error: %v>\n\n", id, err)
	}

	chunk, err := e.index.GetChunk(ctx, id)
	if err != nil {
		return fmt.Sprintf("== %s ==\n\n<source not found for chunk>\n\n", id)
	}

	return fmt.Sprintf("== %s ==\n\n%s\n\n", id, chunk.Source)
}

func (e *Editor) GetChunkSources(ctx context.Context, ids []string) string {
	chunks := ""
	for _, id := range ids {
		chunks += e.getChunkSource(ctx, id)
	}

	return chunks
}

func (e *Editor) SemanticSearch(ctx context.Context, query string) ([]string, error) {
	<-e.initialIndexingDone
	return e.index.Search(ctx, query)
}

func (e *Editor) Close() {
	for _, parser := range e.parsers {
		parser.Close()
	}
}
