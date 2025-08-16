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
		},
	}, nil
}

func (p *GoParser) Chunk(filePath string) (File, error) {
	fullPath := path.Join(p.workspaceRoot, filePath)
	source, err := os.ReadFile(fullPath)
	if err != nil {
		return File{}, err
	}

	tree := p.parser.Parse(source, nil)
	if tree == nil {
		return File{}, fmt.Errorf("couldn't parse %s", filePath)
	}

	file := File{
		Path:     filePath,
		Source:   source,
		ParsedAt: time.Now(),
		Chunks:   p.extractChunks(tree.RootNode(), source),
		tree:     tree,
	}

	return file, nil
}

func (p *GoParser) extractChunks(node *tree_sitter.Node, source []byte) []*Chunk {
	var chunks []*Chunk
	usedPaths := map[string]bool{}

	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		childKind := child.Kind()

		switch childKind {
		case "function_declaration":
			chunks = append(chunks, p.extractFunction(child, source, usedPaths))
		case "type_declaration":
			chunks = append(chunks, p.extractTypeDeclaration(child, source, usedPaths))
		case "method_declaration":
			chunks = append(chunks, p.extractMethod(child, source, usedPaths))
		default:
			chunks = append(chunks, p.extractNode(child, source, usedPaths))
		}
	}

	return chunks
}

func (p *GoParser) extractFunction(node *tree_sitter.Node, source []byte, usedPaths map[string]bool) *Chunk {
	nameQuery := `(function_declaration name: (identifier) @function_name)`
	name := p.getTextWithQuery(nameQuery, node, source)

	return p.createChunk(node, source, name, usedPaths)
}

func (p *GoParser) extractMethod(node *tree_sitter.Node, source []byte, usedPaths map[string]bool) *Chunk {
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

	return p.createChunk(node, source, path, usedPaths)
}

func (p *GoParser) extractTypeDeclaration(node *tree_sitter.Node, source []byte, usedPaths map[string]bool) *Chunk {
	typeQuery := `(type_declaration (type_spec name: (type_identifier) @type_name))`
	name := p.getTextWithQuery(typeQuery, node, source)

	return p.createChunk(node, source, name, usedPaths)
}
