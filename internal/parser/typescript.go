package parser

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_typescript "github.com/tree-sitter/tree-sitter-typescript/bindings/go"
)

var TypeScriptSpec = &LanguageSpec{
	NamedChunks: map[string]NamedChunkExtractor{
		"function_declaration": {
			NameQuery: `(function_declaration name: (identifier) @name)`,
		},
		"function_signature": {
			NameQuery: `(function_signature name: (identifier) @name)`,
		},
		"generator_function_declaration": {
			NameQuery: `(generator_function_declaration name: (identifier) @name)`,
		},
		"class_declaration": {
			NameQuery: `(class_declaration name: (type_identifier) @name)`,
		},
		"abstract_class_declaration": {
			NameQuery: `(abstract_class_declaration name: (type_identifier) @name)`,
		},
		"interface_declaration": {
			NameQuery: `(interface_declaration name: (type_identifier) @name)`,
		},
		"type_alias_declaration": {
			NameQuery: `(type_alias_declaration name: (type_identifier) @name)`,
		},
		"lexical_declaration": {
			NameQuery: `(lexical_declaration (variable_declarator name: (identifier) @name))`,
		},
		"variable_declaration": {
			NameQuery: `(variable_declaration (variable_declarator name: (identifier) @name))`,
		},
		"ambient_declaration": {
			NameQuery: `(ambient_declaration (variable_declaration (variable_declarator name: (identifier) @name)))`,
		},
		"export_statement": {
			NameQuery: `(export_statement (declaration name: (identifier) @name))`,
		},
		"enum_declaration": {
			NameQuery: `(enum_declaration name: (identifier) @name)`,
		},
		"module": {
			NameQuery: `(module name: (identifier) @name)`,
		},
	},
	FoldIntoNextNode: []string{"comment"},
	SkipTypes: []string{
		// These pollute search results
		"import_statement",
		"import_alias",
	},
	FileTypeRules: []FileTypeRule{
		{Pattern: "**/*.test.ts", Type: FileTypeTests},
		{Pattern: "**/*.test.tsx", Type: FileTypeTests},
		{Pattern: "**/*.spec.ts", Type: FileTypeTests},
		{Pattern: "**/*.spec.tsx", Type: FileTypeTests},
		{Pattern: "**/node_modules/**", Type: FileTypeIgnore},
		{Pattern: "**/dist/**", Type: FileTypeIgnore},
		{Pattern: "**/build/**", Type: FileTypeIgnore},
	},
}

func NewTypeScriptParser(workspaceRoot string) (*Parser, error) {
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_typescript.LanguageTypescript()))

	return &Parser{
		workspaceRoot: workspaceRoot,
		parser:        parser,
		spec:          TypeScriptSpec,
	}, nil
}
