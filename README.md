# Sourcerer MCP 🧙

An MCP server that helps AI agents work with large codebases efficiently without burning through costly tokens.
It breaks down your codebase into semantically meaningful chunks and where possible,
provides stable path addressing for precise code navigation and editing.

Your AI agent becomes a sourcerer 🧙, wielding these spells:

- [x] `list_files(in: path, depth?: int = 3)`: Get directory tree up to a specified depth
- [x] `get_file_overviews(paths: path[])`: Get lay-of-the-land summaries of files
- [x] `get_implementations(paths: path[])`: Retrieve full implementation of specific functions, classes, or other chunks
- [ ] `find(pattern: glob, include?: path[], exclude?: path[])`: Find patterns across files and return the full chunks containing matches
- [ ] `edit(ops: operation[])`: Execute batch editing operations
  - [ ] `{"op": "replace", "what": ("<path>" OR "<content>"), "with": "<new>"}`
  - [ ] `{"op": "insert-before", "where": ("<path>" OR "<content>"), "content": "<new>"}`
  - [ ] `{"op": "insert-after", "where": ("<path>" OR "<content>"), "content": "<new>"}`
  - [ ] `{"op": "delete", "what": ("<path>" OR "<content>")}`
  - [ ] `{"op": "rename", "what": "<path>", "to": "<new-name>"}`
  - [ ] `{"op": "move-before", "what": ("<path>" OR "<content>"), "where": ("<path>" OR "<content>")}`
  - [ ] `{"op": "move-after", "what": ("<path>" OR "<content>"), "where": ("<path>" OR "<content>")}`

Examples of paths: `pkg/fs/files.go::File::IsDir`, `"src/auth.js::generateJWT`
