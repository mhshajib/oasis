package utils

import (
	"path"
	"strings"
)

func NormalizePath(input string) string {
	// Use path.Clean to normalize the slashes
	cleanedPath := path.Clean(input)

	trimmedPath := strings.TrimPrefix(cleanedPath, "/")

	// Remove leading slash if it exists
	return "/" + trimmedPath
}
