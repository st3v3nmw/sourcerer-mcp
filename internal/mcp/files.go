package mcp

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func shouldIgnore(name string, isDir bool) bool {
	ignoredDirs := []string{
		".git", ".hg", ".svn", ".bzr", // Version control
		"node_modules", "vendor", "Godeps", // Dependencies
		".next", ".nuxt", "dist", "build", "out", "target", // Build outputs
		".vscode", ".idea", ".vim", ".emacs", // IDE/editor
		"__pycache__", ".pytest_cache", ".mypy_cache", // Python
		".cache", ".tmp", "tmp", "temp", // Cache/temp
		"coverage", ".nyc_output", // Coverage
		".DS_Store", "Thumbs.db", // OS files
	}

	if isDir {
		for _, ignored := range ignoredDirs {
			if name == ignored {
				return true
			}
		}
		// Ignore hidden directories (except .github, .vscode for some projects)
		if strings.HasPrefix(name, ".") && name != ".github" && name != ".claude" {
			return true
		}
	} else {
		// Ignore some common files
		if name == ".DS_Store" || name == "Thumbs.db" {
			return true
		}
	}

	return false
}

func buildDirectoryTree(dirPath string, maxDepth, currentDepth int, prefix string) (string, error) {
	if currentDepth >= maxDepth {
		return "", nil
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// Filter out ignored entries
	var filteredEntries []os.DirEntry
	for _, entry := range entries {
		if !shouldIgnore(entry.Name(), entry.IsDir()) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	entries = filteredEntries

	// Sort entries: directories first, then files, alphabetically within each group
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return entries[i].Name() < entries[j].Name()
	})

	var result strings.Builder

	// If this is the root call, add the directory name
	if currentDepth == 0 {
		result.WriteString(filepath.Base(dirPath))
		if dirPath == "." {
			wd, _ := os.Getwd()
			result.WriteString(filepath.Base(wd))
		}
		result.WriteString("/\n")
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1

		// Build the tree characters
		var connector, nextPrefix string
		if isLast {
			connector = "└── "
			nextPrefix = prefix + "    "
		} else {
			connector = "├── "
			nextPrefix = prefix + "│   "
		}

		result.WriteString(prefix + connector + entry.Name())

		entryPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			result.WriteString("/\n")
			// Recursively build subtree if we haven't reached max depth
			if currentDepth+1 < maxDepth {
				subtree, err := buildDirectoryTree(entryPath, maxDepth, currentDepth+1, nextPrefix)
				if err != nil {
					return "", err
				}
				result.WriteString(subtree)
			}
		} else {
			result.WriteString("\n")
		}
	}

	return result.String(), nil
}
