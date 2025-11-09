package main

import (
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/st3v3nmw/sourcerer-mcp/internal/mcp"
)

//go:embed VERSION
var Version string

func main() {
	Version = strings.TrimSpace(Version)

	workspaceRoot := os.Getenv("SOURCERER_WORKSPACE_ROOT")
	if workspaceRoot == "" {
		workspaceRoot = "."
	}

	server, err := mcp.NewServer(workspaceRoot, Version)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	err = server.Serve()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
