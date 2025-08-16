package main

import (
	"log"
	"os"

	"github.com/st3v3nmw/sourcerer-mcp/internal/mcp"
)

func main() {
	workspaceRoot := os.Getenv("SOURCERER_WORKSPACE_ROOT")
	if workspaceRoot == "" {
		workspaceRoot = "."
	}

	server, err := mcp.NewServer(workspaceRoot)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	if err := server.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
