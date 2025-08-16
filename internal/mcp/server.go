package mcp

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/st3v3nmw/sourcerer-mcp/internal/editor"
	"github.com/st3v3nmw/sourcerer-mcp/internal/fs"
)

type Server struct {
	workspaceRoot string
	mcp           *server.MCPServer
	editor        *editor.Editor
}

func NewServer(workspaceRoot string) (*Server, error) {
	e, err := editor.New(context.Background(), workspaceRoot)
	if err != nil {
		return nil, err
	}

	s := &Server{
		workspaceRoot: workspaceRoot,
		editor:        e,
	}

	s.mcp = server.NewMCPServer(
		"Sourcerer",
		"1.0.0",
		server.WithInstructions(`
You have access to Sourcerer MCP tools for efficient codebase navigation and editing.
Sourcerer provides surgical precision - you can jump directly to specific functions,
classes, and code chunks without reading entire files or burning tokens on broad exploration.

SEARCH STRATEGY:
Sourcerer's semantic search understands concepts and relationships.
Describe what you're looking for conceptually and functionally:

Good queries:
- "user authentication and session management logic"
- "error handling and exception processing code"
- "file parsing and syntax analysis functionality"
- "database operations and data persistence"
- "HTTP request routing and API endpoints"
- "configuration loading and environment setup"

Effective approaches:
- Describe the purpose/behavior you're seeking
- Use natural language to explain the concept
- Include context about what the code should accomplish
- Mention related functionality or typical patterns

CHUNK IDs:
Chunks use stable addressing: path/to/file.ext::ClassName::methodName
- Classes: src/auth.js::AuthService
- Functions: src/auth.js::AuthService::login
- Top-level: src/utils.js::validateEmail
- Unnamed chunks, like imports: src/utils.js::af81a7ff

This addressing is persistent and won't break with minor code changes.
Use get_source_code with these precise ids to get exactly the code you need.

BATCHING:
Prefer batched operations - get_table_of_contents for multiple files, get_source_code for multiple chunks.
When you need multiple related chunks, collect the chunk ids first then batch them in
a single get_source_code call.
This is better than making separate requests which waste tokens and time (round-trips).
`),
	)

	s.mcp.AddTool(
		mcp.NewTool("list_files",
			mcp.WithDescription("Understand project structure and find entry points"),
			mcp.WithString("in",
				mcp.Description("Directory path to explore e.g., 'pkg/fs', '.' for root"),
			),
			mcp.WithNumber("depth",
				mcp.Description("How many levels deep to traverse (default: 3)"),
			),
		),
		s.listFiles,
	)

	s.mcp.AddTool(
		mcp.NewTool("semantic_search",
			mcp.WithDescription("Find relevant code using semantic understanding"),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Your search, returns chunk ids and a small summary of that chunk"),
			),
		),
		s.semanticSearch,
	)

	s.mcp.AddTool(
		mcp.NewTool("get_table_of_contents",
			mcp.WithDescription("Get table of contents showing file structure and chunk IDs"),
			mcp.WithArray("files",
				mcp.WithStringItems(),
				mcp.MinItems(1),
				mcp.Required(),
				mcp.Description("File paths to analyze"),
			),
		),
		s.getTableOfContents,
	)

	s.mcp.AddTool(
		mcp.NewTool("get_source_code",
			mcp.WithDescription("Get the actual code you need to examine/modify"),
			mcp.WithArray("ids",
				mcp.WithStringItems(),
				mcp.MinItems(1),
				mcp.Required(),
				mcp.Description(`
					IDs to chunks to get source code for e.g.
					['pkg/fs/files.go::File::IsDir', 'src/auth/login.js::generateJWT', 'src/utils.js::af81a7ff']
				`),
			),
		),
		s.getSourceCode,
	)

	return s, nil
}

func (s *Server) Serve() error {
	return server.ServeStdio(s.mcp)
}

func (s *Server) listFiles(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	in := request.GetString("in", ".")
	depth := request.GetInt("depth", 3)

	dirPath := path.Join(s.workspaceRoot, in)
	info, err := os.Stat(dirPath)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to access directory: %v", err)), nil
	}

	if !info.IsDir() {
		return mcp.NewToolResultError(fmt.Sprintf("Path is not a directory: %s", dirPath)), nil
	}

	result, err := fs.BuildDirectoryTree(dirPath, depth)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build directory tree: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

func (s *Server) semanticSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")

	results, err := s.editor.SemanticSearch(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No matching chunks found."), nil
	}

	content := strings.Join(results, "\n")
	return mcp.NewToolResultText(content), nil
}

func (s *Server) getTableOfContents(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePaths := request.GetStringSlice("files", []string{})
	tocs := s.editor.GetTOCs(ctx, filePaths)
	return mcp.NewToolResultText(tocs), nil
}

func (s *Server) getSourceCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ids := request.GetStringSlice("ids", []string{})
	chunks := s.editor.GetChunkSources(ctx, ids)
	return mcp.NewToolResultText(chunks), nil
}

func (s *Server) Close() error {
	if s.editor != nil {
		s.editor.Close()
	}

	return nil
}
