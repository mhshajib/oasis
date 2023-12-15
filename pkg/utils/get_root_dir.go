package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetGoModFile(path string) (string, string, error) {
	// Check if go.mod exists in the current directory
	goModPath := filepath.Join(path, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return path, goModPath, nil
	}

	// If reached the root of the file system, stop
	if path == "/" || path == "." {
		return "", "", fmt.Errorf("no go.mod file found")
	}

	// Move up a directory and try again
	return GetGoModFile(filepath.Dir(path))
}
