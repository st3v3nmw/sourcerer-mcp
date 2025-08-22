package parser_test

import (
	"testing"

	"github.com/st3v3nmw/sourcerer-mcp/internal/parser"
	"github.com/stretchr/testify/suite"
)

type MarkdownParserTestSuite struct {
	ParserBaseTestSuite
}

func (s *MarkdownParserTestSuite) SetupSuite() {
	s.ParserBaseTestSuite.SetupSuite()

	var err error
	s.parser, err = parser.NewMarkdownParser(s.workspaceRoot)
	s.Require().NoError(err)
}

func (s *MarkdownParserTestSuite) TestSectionParsing() {
	chunks := s.getChunks("markdown/comprehensive.md")

	tests := []struct {
		name      string
		path      string
		summary   string
		source    string
		startLine uint
		endLine   uint
	}{
		{
			name:    "Document Root Section",
			path:    "eb1f15e27e8bc1a8",
			summary: "---",
			source: `---
type: test
---
`,
			startLine: 1,
			endLine:   4,
		},
		{
			name:    "Paragraph Before Headings",
			path:    "226eabafabc87fec",
			summary: "This document contains various markdown elements for testing the parser.",
			source: `
This document contains various markdown elements for testing the parser.
It starts with content before any headings to test parsing behavior.

Note: This test file primarily uses ATX headings (# ## ###) because setext headings
(underlined with === or ---) don't create section boundaries in the Tree-sitter grammar.

`,
			startLine: 4,
			endLine:   11,
		},
		{
			name:    "ATX Heading Level 1 Section",
			path:    "c77b4b821e14029",
			summary: "# ATX Heading Level 1",
			source: `# ATX Heading Level 1

This is content under the first level heading. It demonstrates basic paragraph text that should be parsed correctly.

Here's a fenced code block:

` + "```python\ndef hello_world():\n    print(\"Hello, World!\")\n    return True\n```" + `

` + "```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

Here's an indented code block:

    func indentedCode() {
        return "This is indented code"
    }

    var x = 42

## ATX Heading Level 2

More content with different markdown elements.

Here's some HTML content:

<div class="example">
    <p>This is an HTML block</p>
    <span>With multiple elements</span>
</div>

<script>
console.log("JavaScript in HTML block");
</script>

### ATX Heading Level 3

Content under level 3 heading.

> This is a block quote.
> It can span multiple lines and contain various content.
>
> Block quotes can have multiple paragraphs.

`,
			startLine: 11,
			endLine:   65,
		},
		{
			name:    "ATX Heading Level 2 Section",
			path:    "11a1130cf96965f5",
			summary: "## ATX Heading Level 2",
			source: `## ATX Heading Level 2

More content with different markdown elements.

Here's some HTML content:

<div class="example">
    <p>This is an HTML block</p>
    <span>With multiple elements</span>
</div>

<script>
console.log("JavaScript in HTML block");
</script>

### ATX Heading Level 3

Content under level 3 heading.

> This is a block quote.
> It can span multiple lines and contain various content.
>
> Block quotes can have multiple paragraphs.

`,
			startLine: 41,
			endLine:   65,
		},
		{
			name:    "ATX Heading Level 3 Section",
			path:    "ac3c606644fb5932",
			summary: "### ATX Heading Level 3",
			source: `### ATX Heading Level 3

Content under level 3 heading.

> This is a block quote.
> It can span multiple lines and contain various content.
>
> Block quotes can have multiple paragraphs.

`,
			startLine: 56,
			endLine:   65,
		},
		{
			name:    "Lists and Content Section",
			path:    "685a9a31c2605cb3",
			summary: "# Lists and Content",
			source: `# Lists and Content

This section demonstrates various list types and content organization.

Lists come in different varieties:

- Unordered list item 1
- Unordered list item 2
  - Nested item 2.1
  - Nested item 2.2
- Unordered list item 3

1. Ordered list item 1
2. Ordered list item 2
   1. Nested ordered item 2.1
   2. Nested ordered item 2.2
3. Ordered list item 3

## Tables and Tasks

This section contains tables and task lists.

Here's a pipe table:

| Name | Age | City | Country |
|------|-----|------|---------|
| Alice | 30 | NYC | USA |
| Bob | 25 | LA | USA |
| Carol | 35 | London | UK |
| David | 28 | Tokyo | Japan |

Task lists are also supported:

- [x] Completed task
- [ ] Pending task
- [ ] Another pending task

Here are some thematic breaks:

---

Content after first thematic break.

***

Content after second thematic break.

___

Content after third thematic break.

Here are various link types:

[Inline link](https://example.com)
[Link with title](https://example.com "Example Title")

[Reference link][ref1]
[Another reference][ref2]

[ref1]: https://reference1.com "Reference 1"
[ref2]: https://reference2.com "Reference 2 Title"

Some inline code: ` + "`fmt.Println(\"Hello\")`" + ` and more text.

Final paragraph with **bold text** and *italic text* and even ***bold italic***.
`,
			startLine: 65,
			endLine:   130,
		},
		{
			name:    "Tables and Tasks Section",
			path:    "6a3ad991e647bd76",
			summary: "## Tables and Tasks",
			source: `## Tables and Tasks

This section contains tables and task lists.

Here's a pipe table:

| Name | Age | City | Country |
|------|-----|------|---------|
| Alice | 30 | NYC | USA |
| Bob | 25 | LA | USA |
| Carol | 35 | London | UK |
| David | 28 | Tokyo | Japan |

Task lists are also supported:

- [x] Completed task
- [ ] Pending task
- [ ] Another pending task

Here are some thematic breaks:

---

Content after first thematic break.

***

Content after second thematic break.

___

Content after third thematic break.

Here are various link types:

[Inline link](https://example.com)
[Link with title](https://example.com "Example Title")

[Reference link][ref1]
[Another reference][ref2]

[ref1]: https://reference1.com "Reference 1"
[ref2]: https://reference2.com "Reference 2 Title"

Some inline code: ` + "`fmt.Println(\"Hello\")`" + ` and more text.

Final paragraph with **bold text** and *italic text* and even ***bold italic***.
`,
			startLine: 83,
			endLine:   130,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			chunk, exists := chunks[test.path]
			s.Require().True(exists, "chunk %s not found", test.path)
			s.Require().NotNil(chunk)

			s.Equal("docs", chunk.Type)
			s.Equal(test.path, chunk.Path)
			s.Equal(test.summary, chunk.Summary)
			s.Equal(test.source, chunk.Source)
			s.Equal(test.startLine, chunk.StartLine)
			s.Equal(test.endLine, chunk.EndLine)
			s.Equal("markdown/comprehensive.md::"+test.path, chunk.ID())
		})
	}
}

func (s *MarkdownParserTestSuite) TearDownSuite() {
	if s.parser != nil {
		s.parser.Close()
	}
}

func TestMarkdownParserTestSuite(t *testing.T) {
	suite.Run(t, new(MarkdownParserTestSuite))
}
