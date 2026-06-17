package scanner

import (
	"fmt"
	"os"
	"path/filepath"
)

// ScanResult holds the results of a directory scan
type ScanResult struct {
	Path            string
	FileCount       int64
	DirectoryCount  int64
	IsGitRepository bool
	Error           error
}

// Scan performs a recursive scan of the given path
func Scan(path string) *ScanResult {
	result := &ScanResult{
		Path: path,
	}

	// Validate path exists
	info, err := os.Stat(path)
	if err != nil {
		result.Error = fmt.Errorf("failed to access path: %w", err)
		return result
	}

	if !info.IsDir() {
		result.Error = fmt.Errorf("path is not a directory")
		return result
	}

	// Count files and directories
	err = filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() {
			// Don't count the root directory itself, only subdirectories
			if filePath != path {
				result.DirectoryCount++
			}
		} else {
			result.FileCount++
		}

		return nil
	})

	if err != nil {
		result.Error = fmt.Errorf("failed to scan directory: %w", err)
		return result
	}

	// Check if it's a Git repository
	gitPath := filepath.Join(path, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		result.IsGitRepository = true
	}

	return result
}
