package utils

import (
	"os"
	"path/filepath"
)

func CreateDirectory(basePath, newDir string) error {
	newPath := filepath.Join(basePath, newDir)
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return os.Mkdir(newPath, 0755)
	}
	return nil
}
