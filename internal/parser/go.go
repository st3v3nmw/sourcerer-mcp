package parser

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

var GoSpec = &LanguageSpec{
	NamedChunks: map[string]NamedChunkExtractor{
		"function_declaration": {
			NameQuery: `(function_declaration name: (identifier) @name)`,
		},
		"method_declaration": {
			NameQuery: `(method_declaration name: (field_identifier) @name)`,
			ParentNameQuery: `
				(method_declaration
					receiver: (parameter_list
						(parameter_declaration
							type: [
								(type_identifier) @name
								(pointer_type
									(type_identifier) @name)
								(generic_type
									(type_identifier) @name)
								(pointer_type
									(generic_type
										(type_identifier) @name))])))`,
		},
		"type_declaration": {
			NameQuery: `
				(type_declaration [
					(type_spec name: (type_identifier) @name)
					(type_alias name: (type_identifier) @name)])`,
		},
		"var_declaration": {
			NameQuery: `(var_declaration (var_spec name: (identifier) @name))`,
		},
		"const_declaration": {
			NameQuery: `(const_declaration (const_spec name: (identifier) @name))`,
		},
	},
	FoldIntoNextNode: []string{"comment"},
	SkipTypes: []string{
		// These pollute search results
		"package_clause",
		"import_declaration",
	},
	FileTypeRules: []FileTypeRule{
		{Pattern: "**/*_test.go", Type: FileTypeTests},
		{Pattern: "vendor/**", Type: FileTypeIgnore},
		{Pattern: "third_party/**", Type: FileTypeIgnore},
	},
}

func NewGoParser(workspaceRoot string) (*Parser, error) {
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_go.Language()))

	return &Parser{
		workspaceRoot: workspaceRoot,
		parser:        parser,
		spec:          GoSpec,
	}, nil
}
