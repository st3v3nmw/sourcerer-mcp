---
type: test
---

This document contains various markdown elements for testing the parser.
It starts with content before any headings to test parsing behavior.

Note: This test file primarily uses ATX headings (# ## ###) because setext headings
(underlined with === or ---) don't create section boundaries in the Tree-sitter grammar.

# ATX Heading Level 1

This is content under the first level heading. It demonstrates basic paragraph text that should be parsed correctly.

Here's a fenced code block:

```python
def hello_world():
    print("Hello, World!")
    return True
```

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

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

# Lists and Content

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

Some inline code: `fmt.Println("Hello")` and more text.

Final paragraph with **bold text** and *italic text* and even ***bold italic***.
