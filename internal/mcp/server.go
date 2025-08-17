package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/st3v3nmw/sourcerer-mcp/internal/analyzer"
)

type Server struct {
	workspaceRoot string
	mcp           *server.MCPServer
	analyzer      *analyzer.Analyzer
}

func NewServer(workspaceRoot, version string) (*Server, error) {
	a, err := analyzer.New(context.Background(), workspaceRoot)
	if err != nil {
		return nil, err
	}

	s := &Server{
		workspaceRoot: workspaceRoot,
		analyzer:      a,
	}

	s.mcp = server.NewMCPServer(
		"Sourcerer",
		version,
		server.WithInstructions(`
You have access to Sourcerer MCP tools for efficient codebase navigation.
Sourcerer provides surgical precision - you can jump directly to specific functions,
classes, and code chunks without reading entire files, reducing token waste & cognitive load.

SEARCH STRATEGY:
Sourcerer's semantic search understands concepts and relationships.
Describe what you're looking for conceptually and functionally:

Good queries:
- "user authentication and session management logic"
- "database operations and data persistence"
- "HTTP request routing and API endpoints"
- "configuration loading and environment setup"

Effective approaches:
- Describe the purpose/behavior you're seeking
- Use natural language to explain the concept
- Include context about what the code should accomplish
- Mention related functionality or typical patterns

AVOID SEMANTIC SEARCH FOR STRUCTURAL QUERIES:
Before using semantic search, ask: "Am I looking for a specific named thing?"
If yes, use pattern-based tools instead:

DON'T semantic search for:
- "function definition of X" → use grep with function patterns
- "interface implementation" → use grep for "type.*interface"
- "struct definition" → use glob for "*.go" & grep for "type.*struct"
- "method calls to X" → use grep for "X(" patterns

Semantic search is for CONCEPTS and RELATIONSHIPS, not NAMES and STRUCTURES.
Use it for "authentication logic" or "error handling patterns",
not "Parser interface" or "ExtractReferences function".

CHUNK IDs:
Chunks use stable addressing: path/to/file.ext::TypeName::methodName
- Classes: src/auth.js::AuthService
- Methods: src/auth.js::AuthService::login
- Top-level: src/utils.js::validateEmail
- Unnamed chunks, like imports: src/utils.js::af81a7ff

Chunk IDs are stable across minor edits but update when code structure changes
(renames, moves, deletions). Use get_source_code with these precise ids to get
exactly the code you need.

If you already know the specific function/class/method/struct/etc
and file location from previous context, construct the chunk ID yourself
and use get_source_code directly rather than semantic searching again.

BATCHING:
When you need multiple related chunks, collect the chunk ids first then batch them in
a single get_source_code call.
This is better than making separate requests which waste tokens and time (round-trips).
`),
	)

	s.mcp.AddTool(
		mcp.NewTool("semantic_search",
			mcp.WithDescription("Find relevant code using semantic understanding"),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Your search, returns chunk ids, a chunk summary, and line numbers"),
			),
		),
		s.semanticSearch,
	)

	s.mcp.AddTool(
		mcp.NewTool("get_source_code",
			mcp.WithDescription("Get the actual code you need to examine/modify"),
			mcp.WithArray("ids",
				mcp.WithStringItems(),
				mcp.MinItems(1),
				mcp.Required(),
				mcp.Description(`
					FULL chunk IDs to get source code for e.g.
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

func (s *Server) semanticSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")

	results, err := s.analyzer.SemanticSearch(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	if len(results) == 0 {
		return mcp.NewToolResultText("No matching chunks found."), nil
	}

	content := strings.Join(results, "\n")
	return mcp.NewToolResultText(content), nil
}

func (s *Server) getSourceCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ids := request.GetStringSlice("ids", []string{})
	chunks := s.analyzer.GetChunkSources(ctx, ids)
	return mcp.NewToolResultText(chunks), nil
}

func (s *Server) Close() error {
	if s.analyzer != nil {
		s.analyzer.Close()
	}

	return nil
}
