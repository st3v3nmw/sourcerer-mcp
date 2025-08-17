package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/cespare/xxhash"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

const (
	summaryMaxChars = 80
)

type Parser interface {
	Chunk(filePath string) (File, error)
	Close()
}

var _ Parser = (*GoParser)(nil)

type ParserBase struct {
	workspaceRoot string
	parser        *tree_sitter.Parser
}

func (p *ParserBase) executeQuery(rawQuery string, node *tree_sitter.Node, source []byte) []*tree_sitter.Node {
	query, err := tree_sitter.NewQuery(p.parser.Language(), rawQuery)
	if err != nil {
		panic(fmt.Sprintf("invalid tree-sitter query: %s\nquery: %s", err, rawQuery))
	}

	cursor := tree_sitter.NewQueryCursor()
	defer cursor.Close()

	var results []*tree_sitter.Node
	matches := cursor.Matches(query, node, source)
	for match := matches.Next(); match != nil; match = matches.Next() {
		for _, capture := range match.Captures {
			results = append(results, &capture.Node)
		}
	}

	return results
}

func (p *ParserBase) getTextWithQuery(query string, node *tree_sitter.Node, source []byte) string {
	nodes := p.executeQuery(query, node, source)
	if len(nodes) > 0 {
		return nodes[0].Utf8Text(source)
	}

	return ""
}

func (p *ParserBase) createChunk(node *tree_sitter.Node, source []byte, path string, usedPaths map[string]bool) *Chunk {
	finalPath := path
	if usedPaths[path] {
		counter := 2
		for usedPaths[fmt.Sprintf("%s-%d", path, counter)] {
			counter++
		}

		finalPath = fmt.Sprintf("%s-%d", path, counter)
	}
	usedPaths[finalPath] = true

	sourceText := node.Utf8Text(source)

	return &Chunk{
		Path:      finalPath,
		Summary:   summarize(sourceText),
		Source:    sourceText,
		StartByte: node.StartByte(),
		EndByte:   node.EndByte(),
		node:      node,
	}
}

func (p *ParserBase) extractNode(node *tree_sitter.Node, source []byte, usedPaths map[string]bool) *Chunk {
	nodeSource := node.Utf8Text(source)
	hash := fmt.Sprintf("%x", xxhash.Sum64String(nodeSource))
	return p.createChunk(node, source, hash, usedPaths)
}

func (p *ParserBase) Close() {
	p.parser.Close()
}

type File struct {
	Path     string // path within namespace
	Chunks   []*Chunk
	Source   []byte
	ParsedAt time.Time

	tree *tree_sitter.Tree
}

func (f *File) Copy() File {
	chunks := make([]*Chunk, len(f.Chunks))
	for i, chunk := range f.Chunks {
		chunks[i] = &Chunk{
			Path:      chunk.Path,
			Summary:   chunk.Summary,
			Source:    chunk.Source,
			StartByte: chunk.StartByte,
			EndByte:   chunk.EndByte,
		}
	}

	source := make([]byte, len(f.Source))
	copy(f.Source, source)

	return File{
		Path:     f.Path,
		Chunks:   chunks,
		Source:   source,
		ParsedAt: f.ParsedAt,
	}
}

type Chunk struct {
	Path      string // path within file
	Summary   string
	Source    string
	StartByte uint
	EndByte   uint

	Children []*Chunk

	node *tree_sitter.Node
}

func summarize(source string) string {
	lines := strings.Split(source, "\n")
	if len(lines) == 0 {
		return ""
	}

	firstLine := strings.TrimSpace(lines[0])
	if len(firstLine) <= summaryMaxChars {
		return firstLine
	}

	nextSpace := strings.Index(firstLine[summaryMaxChars:], " ")
	if nextSpace >= 0 {
		return firstLine[:summaryMaxChars+nextSpace] + "..."
	}

	return firstLine[:summaryMaxChars] + "..."
}
