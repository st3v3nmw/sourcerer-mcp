package main

import (
	"flag"
	"log"

	"github.com/st3v3nmw/sourcerer-mcp/internal/mcp"
)

func main() {
	var workspaceRoot string
	flag.StringVar(&workspaceRoot, "workspace", ".", "Project workspace directory")
	flag.Parse()

	server, err := mcp.NewServer(workspaceRoot)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	if err := server.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
