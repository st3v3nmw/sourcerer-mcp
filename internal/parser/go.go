package parser

import (
	"fmt"
	"os"
	"path"
	"time"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

type GoParser struct {
	ParserBase
}

func NewGoParser(workspaceRoot string) (*GoParser, error) {
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language()))

	return &GoParser{
		ParserBase: ParserBase{
			workspaceRoot: workspaceRoot,
			parser:        parser,
			cache:         map[string]*File{},
		},
	}, nil
}

func (p *GoParser) Chunk(filePath string) (File, error) {
	file := p.getFileFromCache(filePath)
	if file != nil && !file.isStale(p.workspaceRoot) {
		return file.Copy(), nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	fullPath := path.Join(p.workspaceRoot, filePath)
	source, err := os.ReadFile(fullPath)
	if err != nil {
		return File{}, err
	}

	if file == nil {
		file = &File{Path: filePath}
	}

	tree := p.parser.Parse(source, file.tree)
	if tree == nil {
		return File{}, fmt.Errorf("couldn't parse %s", filePath)
	}

	file.Source = source
	file.ParsedAt = time.Now()
	file.Chunks = p.extractChunks(tree.RootNode(), source)

	p.cache[filePath] = file

	return file.Copy(), nil
}

func (p *GoParser) extractChunks(node *tree_sitter.Node, source []byte) []*Chunk {
	var chunks []*Chunk
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		childKind := child.Kind()

		switch childKind {
		case "function_declaration":
			chunks = append(chunks, p.extractFunction(child, source))
		case "type_declaration":
			chunks = append(chunks, p.extractTypeDeclaration(child, source))
		case "method_declaration":
			chunks = append(chunks, p.extractMethod(child, source))
		default:
			chunks = append(chunks, p.extractNode(child, source))
		}
	}

	return chunks
}

func (p *GoParser) extractFunction(node *tree_sitter.Node, source []byte) *Chunk {
	nameQuery := `(function_declaration name: (identifier) @function_name)`
	name := p.getTextWithQuery(nameQuery, node, source)

	summary := p.getFilteredNodeSource(node, source, []string{"body"})

	return p.createChunk(node, source, name, summary)
}

func (p *GoParser) extractMethod(node *tree_sitter.Node, source []byte) *Chunk {
	nameQuery := `(method_declaration name: (field_identifier) @method_name)`
	name := p.getTextWithQuery(nameQuery, node, source)

	receiverQuery := `
  (method_declaration
    receiver: (parameter_list
      (parameter_declaration
        type: [(type_identifier) @receiver_type
               (pointer_type
                 (type_identifier) @receiver_type)])))
  `
	receiver := p.getTextWithQuery(receiverQuery, node, source)

	path := receiver + "::" + name
	summary := p.getFilteredNodeSource(node, source, []string{"body"})

	return p.createChunk(node, source, path, summary)
}

func (p *GoParser) extractTypeDeclaration(node *tree_sitter.Node, source []byte) *Chunk {
	typeQuery := `(type_declaration (type_spec name: (type_identifier) @type_name))`
	name := p.getTextWithQuery(typeQuery, node, source)

	return p.createChunk(node, source, name, node.Utf8Text(source))
}
