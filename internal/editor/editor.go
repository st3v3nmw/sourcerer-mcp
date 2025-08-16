package editor

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/st3v3nmw/sourcerer-mcp/internal/fs"
	"github.com/st3v3nmw/sourcerer-mcp/internal/index"
	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

type Editor struct {
	workspaceRoot string
	parsers       map[Language]parser.Parser
	files         map[string][]string // file path -> chunk paths
	index         *index.Index
}

func New(ctx context.Context, workspaceRoot string) (*Editor, error) {
	index, err := index.New(workspaceRoot)
	if err != nil {
		return nil, err
	}

	editor := &Editor{
		workspaceRoot: workspaceRoot,
		parsers:       map[Language]parser.Parser{},
		files:         map[string][]string{},
		index:         index,
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
	fs.WalkSourceFiles(e.workspaceRoot, func(filePath string) error {
		e.chunk(ctx, filePath)
		return nil
	})
}

func (e *Editor) getParser(filePath string) (parser.Parser, error) {
	language := detectLanguage(filePath)
	if language == UnknownLang {
		return nil, fmt.Errorf("unsupported file %s", filePath)
	}

	var err error
	parser, exists := e.parsers[language]
	if !exists {
		parser, err = newParserForLanguage(language, e.workspaceRoot)
		if err != nil {
			return nil, err
		}
	}

	return parser, nil
}

func (e *Editor) chunk(ctx context.Context, filePath string) {
	info, err := os.Stat(filePath)
	if err != nil {
		return
	}

	if paths, exists := e.files[filePath]; exists {
		allFresh := true
		for _, chunkPath := range paths {
			id := filePath + "::" + chunkPath
			isStale, err := e.index.IsChunkStale(ctx, id, info.ModTime())
			if err != nil || isStale {
				allFresh = false
				break
			}
		}

		if allFresh {
			return
		}
	}

	parser, err := e.getParser(filePath) // TODO: handle error
	if err != nil {
		return
	}

	file, err := parser.Chunk(filePath) // TODO: handle error
	if err != nil {
		return
	}

	err = e.index.Upsert(ctx, &file) // TODO: handle error
	if err != nil {
		panic(err)
	}

	paths := make([]string, len(file.Chunks))
	for i, chunk := range file.Chunks {
		paths[i] = chunk.Path
	}

	e.files[filePath] = paths
}

func (e *Editor) getOverview(ctx context.Context, filePath string) string {
	e.chunk(ctx, filePath)

	overview := fmt.Sprintf("== %s ==\n\n", filePath)

	paths, exists := e.files[filePath]
	if !exists {
		return overview + "<file not found or could not be processed>\n\n"
	}

	for _, chunkPath := range paths {
		id := filePath + "::" + chunkPath
		chunk, err := e.index.GetChunk(ctx, id)
		if err != nil {
			continue
		}

		if chunk.Path != "" {
			overview += "::" + chunk.Path + "\n"
		}

		overview += chunk.Summary + "\n\n"
	}

	return overview
}

func (e *Editor) GetOverviews(ctx context.Context, filePaths []string) string {
	overviews := ""
	for _, filePath := range filePaths {
		overviews += e.getOverview(ctx, filePath)
		overviews += "\n"
	}

	return overviews
}

func (e *Editor) getChunkSource(ctx context.Context, id string) string {
	parts := strings.SplitN(id, "::", 2)
	if len(parts) != 2 {
		return fmt.Sprintf("== %s ==\n\n<invalid chunk id>\n\n", id)
	}

	e.chunk(ctx, parts[0])

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
	return e.index.Search(ctx, query)
}

func (e *Editor) Close() {
	for _, parser := range e.parsers {
		parser.Close()
	}
}
