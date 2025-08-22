package parser

import (
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/cespare/xxhash"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

const (
	summaryMaxChars = 80
)

type FileType string

const (
	FileTypeSrc    FileType = "src"
	FileTypeTests  FileType = "tests"
	FileTypeDocs   FileType = "docs"
	FileTypeIgnore FileType = "ignore"
)

type File struct {
	Path   string // path within workspace
	Chunks []*Chunk
	Source []byte

	tree *tree_sitter.Tree
}

type Chunk struct {
	File        string // file path within workspace
	Type        string
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

func newChunk(node *tree_sitter.Node, source []byte, path string, usedPaths map[string]bool, fileType FileType) *Chunk {
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
		Type:        string(fileType),
		Summary:     summarize(sourceText),
		Source:      sourceText,
		StartLine:   start.Row + 1,
		StartColumn: start.Column + 1,
		EndLine:     end.Row + 1,
		EndColumn:   end.Column + 1,
		ParsedAt:    time.Now().Unix(),
	}
}

func summarize(source string) string {
	source = strings.TrimSpace(source)

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

type LanguageSpec struct {
	NamedChunks       map[string]NamedChunkExtractor
	ExtractChildrenIn []string
	Ignore            []string
	FileTypeRules     []FileTypeRule
}

type NamedChunkExtractor struct {
	NameQuery       string
	ParentNameQuery string
}

type FileTypeRule struct {
	Pattern string
	Type    FileType
}

var globalFileTyleRules = []FileTypeRule{
	{Pattern: "tests/**", Type: FileTypeTests},
	{Pattern: "test/**", Type: FileTypeTests},
	{Pattern: "**/testdata/**", Type: FileTypeTests},

	{Pattern: "docs/**", Type: FileTypeDocs},
	{Pattern: "doc/**", Type: FileTypeDocs},

	{Pattern: ".git/**", Type: FileTypeIgnore},
}

type Parser struct {
	workspaceRoot string
	parser        *tree_sitter.Parser
	spec          *LanguageSpec
}

func (p *Parser) parse(filePath string) (*File, error) {
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

func (p *Parser) Chunk(filePath string) (*File, error) {
	fileType := p.classifyFileType(filePath)
	if fileType == FileTypeIgnore {
		return nil, fmt.Errorf("file %s is marked as ignore", filePath)
	}

	file, err := p.parse(filePath)
	if err != nil {
		return nil, err
	}

	file.Chunks = p.extractChunks(file.tree.RootNode(), file.Source, "", fileType)
	for i := range len(file.Chunks) {
		file.Chunks[i].File = file.Path
	}

	return file, nil
}

func (p *Parser) classifyFileType(filePath string) FileType {
	for _, rule := range globalFileTyleRules {
		matched, _ := doublestar.PathMatch(rule.Pattern, filePath)
		if matched {
			return rule.Type
		}
	}

	for _, rule := range p.spec.FileTypeRules {
		matched, _ := doublestar.PathMatch(rule.Pattern, filePath)
		if matched {
			return rule.Type
		}
	}

	return FileTypeSrc
}

func (p *Parser) extractChunks(node *tree_sitter.Node, source []byte, parentPath string, fileType FileType) []*Chunk {
	var chunks []*Chunk
	usedPaths := map[string]bool{}
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		kind := child.Kind()
		if slices.Contains(p.spec.Ignore, kind) {
			continue
		}

		path := parentPath
		extractor, exists := p.spec.NamedChunks[kind]
		if exists {
			chunkPath, err := p.buildChunkPath(extractor, child, source, parentPath)
			if err != nil {
				// Query failed, fall back to content-hash extraction
				chunks = append(chunks, p.extractNode(child, source, usedPaths, fileType))
			} else {
				chunk := newChunk(child, source, chunkPath, usedPaths, fileType)
				chunks = append(chunks, chunk)
				path = chunkPath
			}
		} else {
			chunks = append(chunks, p.extractNode(child, source, usedPaths, fileType))
		}

		if slices.Contains(p.spec.ExtractChildrenIn, kind) {
			childChunks := p.extractChunks(child, source, path, fileType)
			chunks = append(chunks, childChunks...)
		}
	}

	return chunks
}

func (p *Parser) extractNode(node *tree_sitter.Node, source []byte, usedPaths map[string]bool, fileType FileType) *Chunk {
	nodeSource := node.Utf8Text(source)
	hash := fmt.Sprintf("%x", xxhash.Sum64String(nodeSource))
	return newChunk(node, source, hash, usedPaths, fileType)
}

func (p *Parser) buildChunkPath(extractor NamedChunkExtractor, child *tree_sitter.Node, source []byte, parentPath string) (string, error) {
	path, err := p.getNamedNodePath(extractor.NameQuery, child, source)
	if err != nil {
		return "", err
	}

	if extractor.ParentNameQuery != "" {
		parentName, err := p.getNamedNodePath(extractor.ParentNameQuery, child, source)
		if err != nil {
			return "", err
		}
		parentPath = parentName
	}

	if parentPath != "" {
		path = parentPath + "::" + path
	}

	return path, nil
}

func (p *Parser) getNamedNodePath(query string, node *tree_sitter.Node, source []byte) (string, error) {
	nodes, err := p.executeQuery(query, node, source)
	if err != nil {
		return "", err
	}

	if len(nodes) == 1 {
		return nodes[0].Utf8Text(source), nil
	}

	if len(nodes) > 1 {
		return "", errors.New("too many matches")
	}

	return "", errors.New("no matches found")
}

func (p *Parser) executeQuery(rawQuery string, node *tree_sitter.Node, source []byte) ([]*tree_sitter.Node, error) {
	query, err := tree_sitter.NewQuery(p.parser.Language(), rawQuery)
	if err != nil {
		return nil, fmt.Errorf("invalid tree-sitter query: %s\nquery: %s", err, rawQuery)
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

	return results, nil
}

func (p *Parser) Close() {
	p.parser.Close()
}
