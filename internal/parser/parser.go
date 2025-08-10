package parser

import (
	"os"
	"path"
	"slices"
	"strings"
	"sync"
	"time"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

type Parser interface {
	Chunk(filePath string) (File, error)
	Close()
}

var _ Parser = (*GoParser)(nil)

type ParserBase struct {
	workspaceRoot string
	parser        *tree_sitter.Parser
	cache         map[string]*File
	mu            sync.RWMutex
}

func (p *ParserBase) executeQuery(rawQuery string, node *tree_sitter.Node, source []byte) []*tree_sitter.Node {
	query, err := tree_sitter.NewQuery(p.parser.Language(), rawQuery)
	if err != nil {
		panic(err)
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

func (p *ParserBase) createChunk(node *tree_sitter.Node, source []byte, path, summary string) *Chunk {
	return &Chunk{
		Path:      path, // TODO: handle ambiguity when more then one symbol exists at the same level with the same name
		Summary:   summary,
		Source:    node.Utf8Text(source),
		StartByte: node.StartByte(),
		EndByte:   node.EndByte(),
		node:      node,
	}
}

func (p *ParserBase) getFileFromCache(filePath string) *File {
	p.mu.RLock()
	defer p.mu.RUnlock()

	file, exists := p.cache[filePath]
	if !exists {
		return nil
	}

	return file
}

func (p *ParserBase) extractNode(node *tree_sitter.Node, source []byte) *Chunk {
	nodeSource := node.Utf8Text(source)
	return &Chunk{
		Summary:   nodeSource,
		Source:    nodeSource,
		StartByte: node.StartByte(),
		EndByte:   node.EndByte(),
		node:      node,
	}
}

func (p *ParserBase) getFilteredNodeSource(node *tree_sitter.Node, source []byte, exclude []string) string {
	slices.Sort(exclude)

	var parts []string
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)

		fieldName := node.FieldNameForChild(uint32(i))
		_, found := slices.BinarySearch(exclude, fieldName)
		if found {
			continue
		}

		parts = append(parts, child.Utf8Text(source))
	}

	return strings.Join(parts, " ")
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

func (f *File) isStale(workspaceRoot string) bool {
	fullPath := path.Join(workspaceRoot, f.Path)
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return false
	}

	return f.ParsedAt.Before(fileInfo.ModTime())
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

	source := []byte{}
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
