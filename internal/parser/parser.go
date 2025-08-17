package parser

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cespare/xxhash"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

const (
	summaryMaxChars = 80
)

type File struct {
	Path   string // path within workspace
	Chunks []*Chunk
	Source []byte

	tree *tree_sitter.Tree
}

type Chunk struct {
	File        string // file path within workspace
	Path        string // path within file
	Summary     string
	Source      string
	StartLine   uint
	StartColumn uint
	EndLine     uint
	EndColumn   uint
	ParsedAt    int64
}

func (c *Chunk) ID() string {
	return c.File + "::" + c.Path
}

type Parser interface {
	Chunk(filePath string) (*File, error)
	Close()
}

var _ Parser = (*GoParser)(nil)

type ParserBase struct {
	workspaceRoot string
	parser        *tree_sitter.Parser
}

func (p *ParserBase) parse(filePath string) (*File, error) {
	fullPath := path.Join(p.workspaceRoot, filePath)
	source, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	tree := p.parser.Parse(source, nil)
	if tree == nil {
		return nil, fmt.Errorf("couldn't parse %s", filePath)
	}

	return &File{
		Path:   filePath,
		Source: source,
		tree:   tree,
	}, nil
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
	start := node.StartPosition()
	end := node.EndPosition()

	return &Chunk{
		Path:        finalPath,
		Summary:     summarize(sourceText),
		Source:      sourceText,
		StartLine:   start.Row + 1,
		StartColumn: start.Column + 1,
		EndLine:     end.Row + 1,
		EndColumn:   end.Column + 1,
		ParsedAt:    time.Now().UnixMicro(),
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
