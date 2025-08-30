# Contributing to Sourcerer MCP

First off, thanks for wanting to contribute to Sourcerer!

## Getting Started

### Prerequisites

- **Go 1.24.5+**: Required for building and running the project
- **OpenAI API Key**: Needed for embedding generation during development & testing

### Development Setup

1. **Clone/fork the repository:**

```bash
git clone git@github.com:st3v3nmw/sourcerer-mcp.git
cd sourcerer-mcp
```

2. **Install dependencies:**

```bash
go mod download
```

3. **Set up environment:**

```bash
export OPENAI_API_KEY=your-openai-api-key
export SOURCERER_WORKSPACE_ROOT=$(pwd)
```

4. **Run the project:**

```bash
go run ./cmd/sourcerer
```

or:

```bash
claude mcp add sourcerer -e OPENAI_API_KEY=your-openai-api-key -e SOURCERER_WORKSPACE_ROOT=$(pwd) -- go run ./cmd/sourcerer
```

5. **Run tests:**

```bash
go test -v -coverprofile coverage.out  ./... && go tool cover -html=coverage.out -o coverage.html
```

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker to report bugs or suggest features
- Include relevant details: Steps to reproduce, error messages, OS, MCP version, etc
- Check existing issues to avoid duplicates

### Code Contributions

#### Types of Contributions Welcome

1. **Language Support**: Add Tree-sitter parsers for new programming languages
2. **Bug Fixes**: Fix existing functionality issues
3. **Performance Improvements**: Optimize indexing, search, or memory usage
4. **Feature Enhancements**: Improve existing MCP tools or add new ones
5. **Documentation**: Improve existing docs or add new documentation

#### Development Workflow

1. **Fork the repository** and create a feature branch from `main`
2. **Make your changes** following the coding standards below
3. **Add tests** for new parser-related functionality
4. **Ensure all tests pass**: `go test ./...`
5. **Submit a pull request** with a clear description

### Adding Language Support

To add support for a new programming language:

1. **Add Tree-sitter dependency** to `go.mod`. Look for grammars in [gh/tree-sitter](https://github.com/tree-sitter) or [gh/tree-sitter-grammars](https://github.com/tree-sitter-grammars) first
2. **Create parser file** in `internal/parser/` (e.g., `rust.go`)
3. **Write Tree-sitter queries** to extract functions, classes, structs, interfaces, methods, etc
4. **Add comprehensive tests** in `internal/parser/<language>_test.go` & `testdata/<language>/`
5. **Update README.md** to list the new supported language

### Coding Standards

#### Go Style Guidelines

- Follow standard Go conventions and `go fmt`
- Use meaningful variable and function names
- Write clear, concise comments for exported functions
- Keep functions focused and reasonably sized

#### Testing Requirements

- Write unit tests for all parser-related functionality i.e., in `internal/parser/`
- Achieve reasonable test coverage for new parser-related code
- Test edge cases and error conditions
- Use table-driven tests where appropriate

#### Commit Message Format

Use clear, descriptive commit messages:

```
<type>: brief description

Optional longer explanation of the change, including:
- Why the change was needed
- How it addresses the issue
- Any breaking changes or migration notes
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`, `ci`

## Code Review Process

### Pull Request Requirements

- Clear description of changes and motivation
- All tests must pass
- Code follows project conventions
- Documentation updated if needed

## Getting Help

- **[Issues](https://github.com/st3v3nmw/sourcerer-mcp/issues)**: Use GitHub issues for questions about contributing
- **[Discussions](https://github.com/st3v3nmw/sourcerer-mcp/discussions)**: For broader design discussions or ideas

## License

By contributing to Sourcerer MCP, you agree that your contributions will be licensed under the same license as the project.
