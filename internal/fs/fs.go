package fs

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func shouldIgnorePath(workspaceRoot, path string) bool {
	name := filepath.Base(path)
	if name == ".git" {
		return true
	}

	cmd := exec.Command("git", "check-ignore", path)
	cmd.Dir = workspaceRoot
	return cmd.Run() == nil
}

func WalkSourceFiles(workspaceRoot string, supportedExts []string, callback func(filePath string) error) error {
	extMap := make(map[string]bool, len(supportedExts))
	for _, ext := range supportedExts {
		extMap[ext] = true
	}

	return filepath.Walk(workspaceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if shouldIgnorePath(workspaceRoot, path) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !extMap[ext] {
			return nil
		}

		relPath, err := filepath.Rel(workspaceRoot, path)
		if err != nil {
			relPath = path
		}

		return callback(relPath)
	})
}
