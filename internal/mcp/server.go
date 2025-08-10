package mcp

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/st3v3nmw/sourcerer-mcp/internal/editor"
	"github.com/st3v3nmw/sourcerer-mcp/internal/lsp"
)

type Server struct {
	workspaceRoot string
	mcp           *server.MCPServer
	editor        *editor.Editor
	lsp           *lsp.LSP
}

func NewServer(workspaceRoot string) (*Server, error) {
	s := &Server{
		workspaceRoot: workspaceRoot,
		editor:        editor.New(workspaceRoot),
	}

	s.mcp = server.NewMCPServer(
		"Sourcerer",
		"1.0.0",
		server.WithInstructions(`
			An MCP server that helps AI agents work with large codebases
			efficiently without burning through costly tokens.
		`,
		),
	)

	s.mcp.AddTool(
		mcp.NewTool("list_files",
			mcp.WithDescription("Get directory tree up to a specified depth"),
			mcp.WithString("in",
				mcp.Description("Directory path to explore (e.g., 'pkg/fs', '.' for root)"),
			),
			mcp.WithNumber("depth",
				mcp.Description("How many levels deep to traverse (default: 3)"),
			),
		),
		s.listFiles,
	)

	s.mcp.AddTool(
		mcp.NewTool("get_file_overviews",
			mcp.WithDescription("Get lay-of-the-land summaries of files. Use list_files to find the files."),
			mcp.WithArray("paths",
				mcp.WithStringItems(),
				mcp.MinItems(1),
				mcp.Required(),
				mcp.Description("List of file paths to analyze"),
			),
		),
		s.getFileOverviews,
	)

	s.mcp.AddTool(
		mcp.NewTool("get_implementations",
			mcp.WithDescription("Retrieve full implementation of specific functions, classes, or other chunks. Use get_file_overviews to get the chunks."),
			mcp.WithArray("paths",
				mcp.WithStringItems(),
				mcp.MinItems(1),
				mcp.Required(),
				mcp.Description("Paths to chunks to get implementations for (e.g., ['pkg/fs/files.go::File::IsDir', 'src/auth/login.js::generateJWT'])"),
			),
		),
		s.getImplementations,
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

	result, err := buildDirectoryTree(dirPath, depth, 0, "")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to build directory tree: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

func (s *Server) getFileOverviews(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	filePaths := request.GetStringSlice("paths", []string{})
	overviews := s.editor.GetOverviews(filePaths)
	return mcp.NewToolResultText(overviews), nil
}

func (s *Server) getImplementations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	paths := request.GetStringSlice("paths", []string{})
	chunks := s.editor.GetChunks(paths)
	return mcp.NewToolResultText(chunks), nil
}

func (s *Server) Close() error {
	if s.editor != nil {
		s.editor.Close()
	}

	return nil
}
