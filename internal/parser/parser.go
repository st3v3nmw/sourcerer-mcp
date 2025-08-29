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
	chunkSummaryMaxChars = 80
)

// FileType represents the classification of a file within the workspace
type FileType string

const (
	FileTypeSrc    FileType = "src"
	FileTypeTests  FileType = "tests"
	FileTypeDocs   FileType = "docs"
	FileTypeIgnore FileType = "ignore"
)

// File represents a parsed source file with its extracted semantic chunks
type File struct {
	Path   string // path within workspace
	Chunks []*Chunk
	Source []byte

	tree *tree_sitter.Tree
}

// Chunk represents a semantic unit of code extracted from source files
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

// ID returns a unique identifier for this chunk in the format "file::path"
func (c *Chunk) ID() string {
	return c.File + "::" + c.Path
}

// newChunk creates a new Chunk from related tree-sitter nodes
func (p *Parser) newChunk(
	node *tree_sitter.Node,
	source []byte,
	path string,
	usedPaths map[string]bool,
	fileType FileType,
	folded []*tree_sitter.Node,
	extractor *NamedChunkExtractor,
) *Chunk {
	finalPath := resolvePath(path, usedPaths)
	startPos, startByte, endPos, endByte := calculateChunkBounds(node, folded)

	// Determine which node to use for the summary
	summaryNode := node
	if extractor != nil && extractor.SummaryNodeQuery != "" {
		// Use the existing executeQuery method to find the summary node
		matches, err := p.executeQuery(extractor.SummaryNodeQuery, node, source)
		if err == nil && len(matches) > 0 {
			summaryNode = matches[0]
		}
	}

	summaryText := summaryNode.Utf8Text(source)
	fullText := source[startByte:endByte]

	return &Chunk{
		Path:        finalPath,
		Type:        string(fileType),
		Summary:     summarize(summaryText),
		Source:      string(fullText),
		StartLine:   startPos.Row + 1,
		StartColumn: startPos.Column + 1,
		EndLine:     endPos.Row + 1,
		EndColumn:   endPos.Column + 1,
		ParsedAt:    time.Now().Unix(),
	}
}

// resolvePath handles path name conflicts by appending a counter when needed
func resolvePath(path string, usedPaths map[string]bool) string {
	if !usedPaths[path] {
		usedPaths[path] = true
		return path
	}

	counter := 2
	for usedPaths[fmt.Sprintf("%s-%d", path, counter)] {
		counter++
	}

	finalPath := fmt.Sprintf("%s-%d", path, counter)
	usedPaths[finalPath] = true
	return finalPath
}

// calculateChunkBounds determines the start and end positions for a chunk,
// extending to include any preceding folded nodes
func calculateChunkBounds(node *tree_sitter.Node, folded []*tree_sitter.Node) (
	startPos tree_sitter.Point, startByte uint,
	endPos tree_sitter.Point, endByte uint,
) {
	startPos = node.StartPosition()
	startByte = node.StartByte()
	endPos = node.EndPosition()
	endByte = node.EndByte()

	if len(folded) > 0 {
		firstFolded := folded[0]
		startPos = firstFolded.StartPosition()
		startByte = firstFolded.StartByte()
	}

	return startPos, startByte, endPos, endByte
}

// summarize creates a concise summary from source code, truncating at word boundaries
// when the first line exceeds the maximum character limit
func summarize(source string) string {
	source = strings.TrimSpace(source)

	lines := strings.Split(source, "\n")
	if len(lines) == 0 {
		return ""
	}

	firstLine := strings.TrimSpace(lines[0])
	if len(firstLine) <= chunkSummaryMaxChars {
		return firstLine
	}

	nextSpace := strings.Index(firstLine[chunkSummaryMaxChars:], " ")
	if nextSpace >= 0 {
		return firstLine[:chunkSummaryMaxChars+nextSpace] + "..."
	}

	return firstLine[:chunkSummaryMaxChars] + "..."
}

// LanguageSpec defines language-specific parsing behavior for tree-sitter
type LanguageSpec struct {
	NamedChunks       map[string]NamedChunkExtractor // node types that can be extracted by name
	ExtractChildrenIn []string                       // node types whose children should be recursively processed
	FoldIntoNextNode  []string                       // node types to fold into next node, e.g., comments
	SkipTypes         []string                       // node types to completely skip
	FileTypeRules     []FileTypeRule                 // language-specific file type classification rules
}

// NamedChunkExtractor defines tree-sitter queries for extracting named code entities
type NamedChunkExtractor struct {
	NameQuery        string // query to extract the entity name
	ParentNameQuery  string // optional query to extract parent entity name for hierarchical paths
	SummaryNodeQuery string // optional query to extract a specific node for the summary instead of the main node
}

// FileTypeRule defines a pattern-based rule for classifying file types
type FileTypeRule struct {
	Pattern string   // glob pattern to match against file paths
	Type    FileType // the file type to assign when the pattern matches
}

// globalFileTyleRules contains universal file type classification patterns
// that apply across all programming languages
var globalFileTyleRules = []FileTypeRule{
	{Pattern: "**/tests/**", Type: FileTypeTests},
	{Pattern: "**/test/**", Type: FileTypeTests},
	{Pattern: "**/testdata/**", Type: FileTypeTests},

	{Pattern: "docs/**", Type: FileTypeDocs},
	{Pattern: "doc/**", Type: FileTypeDocs},

	{Pattern: ".git/**", Type: FileTypeIgnore},
	{Pattern: "coverage/**", Type: FileTypeIgnore},
	{Pattern: ".coverage/**", Type: FileTypeIgnore},
}

// Parser handles parsing and semantic chunk extraction from source files
// using tree-sitter for language-aware AST processing
type Parser struct {
	workspaceRoot string              // absolute path to the workspace root
	parser        *tree_sitter.Parser // tree-sitter parser instance
	spec          *LanguageSpec       // language-specific parsing configuration
}

// parse reads and parses a file using tree-sitter, returning the AST and source
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

// Chunk parses a file and extracts semantic chunks from its AST
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

// classifyFileType determines the file type based on path patterns,
// checking global rules first, then language-specific rules
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

// extractChunks recursively extracts semantic chunks from an AST node.
func (p *Parser) extractChunks(
	node *tree_sitter.Node,
	source []byte,
	parentPath string,
	fileType FileType,
) []*Chunk {
	var chunks []*Chunk
	usedPaths := map[string]bool{}
	var folded []*tree_sitter.Node

	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		kind := child.Kind()

		if slices.Contains(p.spec.SkipTypes, kind) {
			// Process any remaining folded nodes as standalone chunks
			for _, foldedNode := range folded {
				chunks = append(chunks, p.extractNode(foldedNode, source, usedPaths, fileType, nil))
			}
			folded = nil

			continue
		}

		if slices.Contains(p.spec.FoldIntoNextNode, kind) {
			folded = append(folded, child)
			continue
		}

		// Process code nodes & folded nodes, if any
		chunk, path := p.createChunkFromNode(child, source, parentPath, fileType, usedPaths, folded)
		chunks = append(chunks, chunk)
		folded = nil

		// Recursively process children if specified
		if slices.Contains(p.spec.ExtractChildrenIn, kind) {
			childChunks := p.extractChunks(child, source, path, fileType)
			chunks = append(chunks, childChunks...)
		}
	}

	// Process any remaining folded nodes as standalone chunks
	for _, foldedNode := range folded {
		chunks = append(chunks, p.extractNode(foldedNode, source, usedPaths, fileType, nil))
	}

	return chunks
}

// createChunkFromNode creates a chunk from a code node, attempting named extraction first
func (p *Parser) createChunkFromNode(
	node *tree_sitter.Node,
	source []byte,
	parentPath string,
	fileType FileType,
	usedPaths map[string]bool,
	folded []*tree_sitter.Node,
) (*Chunk, string) {
	kind := node.Kind()
	extractor, exists := p.spec.NamedChunks[kind]

	if exists {
		chunkPath, err := p.buildChunkPath(extractor, node, source, parentPath)
		if err == nil {
			chunk := p.newChunk(node, source, chunkPath, usedPaths, fileType, folded, &extractor)
			return chunk, chunkPath
		}
	}

	// No named extractor or building chunk path failed, use content-hash
	return p.extractNode(node, source, usedPaths, fileType, folded), parentPath
}

// extractNode creates a chunk from a node using content-based hashing for the path
func (p *Parser) extractNode(
	node *tree_sitter.Node,
	source []byte,
	usedPaths map[string]bool,
	fileType FileType,
	folded []*tree_sitter.Node,
) *Chunk {
	nodeSource := node.Utf8Text(source)
	hash := fmt.Sprintf("%x", xxhash.Sum64String(nodeSource))

	return p.newChunk(node, source, hash, usedPaths, fileType, folded, nil)
}

// buildChunkPath constructs a hierarchical path for a named chunk using tree-sitter queries
func (p *Parser) buildChunkPath(
	extractor NamedChunkExtractor,
	child *tree_sitter.Node,
	source []byte,
	parentPath string,
) (string, error) {
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

// getNamedNodePath extracts a name from a node using a tree-sitter query
func (p *Parser) getNamedNodePath(
	query string,
	node *tree_sitter.Node,
	source []byte,
) (string, error) {
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

// executeQuery runs a tree-sitter query against a node and returns all matching nodes
func (p *Parser) executeQuery(
	rawQuery string,
	node *tree_sitter.Node,
	source []byte,
) ([]*tree_sitter.Node, error) {
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

// Close releases resources used by the tree-sitter parser
func (p *Parser) Close() {
	p.parser.Close()
}
