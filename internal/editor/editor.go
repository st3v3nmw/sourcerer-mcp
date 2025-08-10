package editor

import (
	"fmt"
	"strings"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
)

type Editor struct {
	workspaceRoot string
	parsers       map[Language]parser.Parser
	files         map[string]parser.File
}

func New(workspaceRoot string) *Editor {
	return &Editor{
		workspaceRoot: workspaceRoot,
		parsers:       map[Language]parser.Parser{},
		files:         map[string]parser.File{},
	}
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

func (e *Editor) chunk(filePath string) {
	parser, _ := e.getParser(filePath) // TODO: handle error

	file, _ := parser.Chunk(filePath) // TODO: handle error

	e.files[filePath] = file
}

func (e *Editor) getOverview(filePath string) string { // TODO: handle errors
	e.chunk(filePath)

	overview := fmt.Sprintf("== %s ==\n\n", filePath)
	file := e.files[filePath]
	for _, chunk := range file.Chunks {
		if chunk.Path != "" {
			overview += "::" + chunk.Path + "\n"
		}

		overview += chunk.Summary + "\n\n"
	}

	return overview
}

func (e *Editor) GetOverviews(filePaths []string) string {
	overviews := ""
	for _, filePath := range filePaths {
		overviews += e.getOverview(filePath)
		overviews += "\n"
	}

	return overviews
}

func (e *Editor) getChunk(path string) (string, error) { // TODO: handle errors
	parts := strings.SplitN(path, "::", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid path %s", path)
	}

	filePath, pathInFile := parts[0], parts[1]
	e.chunk(filePath)

	file := e.files[filePath]
	for _, chunk := range file.Chunks {
		if chunk.Path == pathInFile {
			return fmt.Sprintf("== %s ==\n\n%s\n\n", path, chunk.Source), nil
		}
	}

	return "", fmt.Errorf("chunk not found at path %s", path)
}

func (e *Editor) GetChunks(paths []string) string {
	chunks := ""
	for _, path := range paths {
		chunk, _ := e.getChunk(path) // TODO: handle errors

		chunks += chunk
		chunks += "\n"
	}

	return chunks
}

func (e *Editor) Close() {
	for _, parser := range e.parsers {
		parser.Close()
	}
}
