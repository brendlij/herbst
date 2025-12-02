package util

import (
	"os"
	"path/filepath"
)

// ResolveDir determines which directory to use based on:
// 1. Environment variable (if set)
// 2. Default dev path (if exists)
// 3. Container fallback path
func ResolveDir(envVar, devDefault, containerFallback string) string {
	// Check environment variable first
	if dir := os.Getenv(envVar); dir != "" {
		return dir
	}

	// Check if dev default exists
	if _, err := os.Stat(devDefault); err == nil {
		abs, err := filepath.Abs(devDefault)
		if err == nil {
			return abs
		}
		return devDefault
	}

	// Fall back to container path
	return containerFallback
}

// EnsureDir creates the directory if it doesn't exist
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}
