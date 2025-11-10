package parser

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

var JavaScriptSpec = &LanguageSpec{
	NamedChunks: map[string]NamedChunkExtractor{
		"function_declaration": {
			NameQuery: `(function_declaration name: (identifier) @name)`,
		},
		"generator_function_declaration": {
			NameQuery: `(generator_function_declaration name: (identifier) @name)`,
		},
		"class_declaration": {
			NameQuery: `(class_declaration name: (identifier) @name)`,
		},
		"lexical_declaration": {
			NameQuery: `(lexical_declaration (variable_declarator name: (identifier) @name))`,
		},
		"variable_declaration": {
			NameQuery: `(variable_declaration (variable_declarator name: (identifier) @name))`,
		},
		"method_definition": {
			NameQuery: `(method_definition name: (property_identifier) @name)`,
		},
		"field_definition": {
			NameQuery: `(field_definition property: (property_identifier) @name)`,
		},
	},
	ExtractChildrenIn: []string{
		"class_declaration",
		"class_body",
		"export_statement",
	},
	FoldIntoNextNode: []string{"comment", "export", "default"},
	SkipTypes: []string{
		// Imports pollute search results
		"import_statement",
		// Skip punctuation and keyword tokens
		"{", "}", ";",
		"class", "extends", "implements",
		// Skip identifier tokens (they're part of declarations)
		"identifier",
		// Skip class heritage clauses
		"class_heritage",
		// Skip container nodes (but still extract their children)
		"class_body",
		"export_statement",
	},
	FileTypeRules: []FileTypeRule{
		{Pattern: "**/*.test.js", Type: FileTypeTests},
		{Pattern: "**/*.test.jsx", Type: FileTypeTests},
		{Pattern: "**/*.spec.js", Type: FileTypeTests},
		{Pattern: "**/*.spec.jsx", Type: FileTypeTests},
		{Pattern: "**/node_modules/**", Type: FileTypeIgnore},
		{Pattern: "**/dist/**", Type: FileTypeIgnore},
		{Pattern: "**/build/**", Type: FileTypeIgnore},
	},
}

func NewJavaScriptParser(workspaceRoot string) (*Parser, error) {
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_javascript.Language()))

	return &Parser{
		workspaceRoot: workspaceRoot,
		parser:        parser,
		spec:          JavaScriptSpec,
	}, nil
}
