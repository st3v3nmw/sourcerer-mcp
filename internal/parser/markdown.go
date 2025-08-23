package parser

import (
	tree_sitter_markdown "github.com/tree-sitter-grammars/tree-sitter-markdown/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

var MarkdownSpec = &LanguageSpec{
	ExtractChildrenIn: []string{"section"},
	IgnoreTypes: []string{
		// Headings are organizational markers, not containers.
		"atx_heading", "setext_heading",
		// We're chunking by section so lower level nodes don't get their own chunks
		//  since this would lead to a lot of noise as sections overlap.
		"block_quote",
		"block_continuation",
		"fenced_code_block", "indented_code_block",
		"html_block",
		"link_reference_definition",
		"list",
		"paragraph",
		"pipe_table",
		"thematic_break",
	},
	FileTypeRules: []FileTypeRule{
		{Pattern: "**/*.md", Type: FileTypeDocs},
	},
}

func NewMarkdownParser(workspaceRoot string) (*Parser, error) {
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_markdown.Language()))

	return &Parser{
		workspaceRoot: workspaceRoot,
		parser:        parser,
		spec:          MarkdownSpec,
	}, nil
}
