package parser

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

var GoSpec = &LanguageSpec{
	ChunkExtractors: map[string]ChunkExtractor{
		"function_declaration": {
			NameQuery: `(function_declaration name: (identifier) @name)`,
		},
		"method_declaration": {
			NameQuery: `(method_declaration name: (field_identifier) @name)`,
			ParentNameQuery: `
			  (method_declaration
			    receiver: (parameter_list
			      (parameter_declaration
			        type: [(type_identifier) @receiver_type
			               (pointer_type
			                 (type_identifier) @receiver_type)])))`,
			ParentNameInParent: false,
		},
		"type_declaration": {
			NameQuery: `(type_declaration (type_spec name: (type_identifier) @name))`,
		},
	},
	RefQueries: map[string]string{
		"function_calls": `(call_expression function: (identifier) @call)`,
		"method_calls":   `(call_expression function: (selector_expression field: (field_identifier) @method))`,
		"imports":        `(import_spec path: (interpreted_string_literal) @import)`,
	},
}

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
			spec:          GoSpec,
		},
	}, nil
}
