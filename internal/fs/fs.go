package fs

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileFilter struct {
	workspaceRoot string
	supportedExts map[string]bool
}

func NewFileFilter(workspaceRoot string, supportedExts []string) *FileFilter {
	extMap := make(map[string]bool, len(supportedExts))
	for _, ext := range supportedExts {
		extMap[ext] = true
	}

	return &FileFilter{
		workspaceRoot: workspaceRoot,
		supportedExts: extMap,
	}
}

func (f *FileFilter) ShouldIgnore(path string) bool {
	if f.isGitIgnored(path) {
		return true
	}

	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return false
	}

	return !f.supportedExts[ext]
}

func (f *FileFilter) isGitIgnored(path string) bool {
	name := filepath.Base(path)
	if name == ".git" {
		return true
	}

	cmd := exec.Command("git", "check-ignore", path)
	cmd.Dir = f.workspaceRoot
	return cmd.Run() == nil
}

func WalkSourceFiles(workspaceRoot string, supportedExts []string, callback func(filePath string) error) error {
	filter := NewFileFilter(workspaceRoot, supportedExts)

	return filepath.Walk(workspaceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filter.ShouldIgnore(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		relPath, err := filepath.Rel(workspaceRoot, path)
		if err != nil {
			relPath = path
		}

		return callback(relPath)
	})
}
