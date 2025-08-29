package parser

import (
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_python "github.com/tree-sitter/tree-sitter-python/bindings/go"
)

var PythonSpec = &LanguageSpec{
	NamedChunks: map[string]NamedChunkExtractor{
		"function_definition": {
			NameQuery: `(function_definition name: (identifier) @name)`,
		},
		"class_definition": {
			NameQuery: `(class_definition name: (identifier) @name)`,
		},
		"decorated_definition": {
			NameQuery: `(decorated_definition definition: [
				(function_definition name: (identifier) @name)
				(class_definition name: (identifier) @name)
			])`,
			SummaryNodeQuery: `(decorated_definition definition: [
				(function_definition) @summary
				(class_definition) @summary
			])`,
		},
	},
	FoldIntoNextNode: []string{"comment"},
	SkipTypes: []string{
		// These pollute search results
		"import_statement",
	},
	FileTypeRules: []FileTypeRule{
		{Pattern: "**/test*.py", Type: FileTypeTests},
		{Pattern: "**/*_test.py", Type: FileTypeTests},
		{Pattern: "**/__pycache__/**", Type: FileTypeIgnore},
		{Pattern: "**/venv/**", Type: FileTypeIgnore},
		{Pattern: "**/.venv/**", Type: FileTypeIgnore},
		{Pattern: "**/env/**", Type: FileTypeIgnore},
		{Pattern: "**/.env/**", Type: FileTypeIgnore},
		{Pattern: "**/site-packages/**", Type: FileTypeIgnore},
	},
}

func NewPythonParser(workspaceRoot string) (*Parser, error) {
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_python.Language()))

	return &Parser{
		workspaceRoot: workspaceRoot,
		parser:        parser,
		spec:          PythonSpec,
	}, nil
}
