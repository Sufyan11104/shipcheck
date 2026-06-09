package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScan_ValidDirectory(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	
	// Create some test files and directories
	os.Mkdir(filepath.Join(tmpDir, "subdir1"), 0755)
	os.Mkdir(filepath.Join(tmpDir, "subdir2"), 0755)
	os.Create(filepath.Join(tmpDir, "file1.txt"))
	os.Create(filepath.Join(tmpDir, "file2.txt"))
	os.Create(filepath.Join(tmpDir, "subdir1", "file3.txt"))

	result := Scan(tmpDir)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}

	if result.FileCount != 3 {
		t.Errorf("expected 3 files, got %d", result.FileCount)
	}

	if result.DirectoryCount != 2 {
		t.Errorf("expected 2 directories, got %d", result.DirectoryCount)
	}

	if result.IsGitRepository {
		t.Errorf("expected IsGitRepository to be false")
	}
}

func TestScan_GitRepository(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.Mkdir(gitDir, 0755)

	result := Scan(tmpDir)

	if result.Error != nil {
		t.Fatalf("expected no error, got: %v", result.Error)
	}

	if !result.IsGitRepository {
		t.Errorf("expected IsGitRepository to be true")
	}
}

func TestScan_InvalidPath(t *testing.T) {
	result := Scan("/nonexistent/path/that/does/not/exist")

	if result.Error == nil {
		t.Errorf("expected error for invalid path, got nil")
	}
}

func TestScan_PathIsNotDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "file.txt")
	os.Create(filePath)

	result := Scan(filePath)

	if result.Error == nil {
		t.Errorf("expected error when path is a file, got nil")
	}
}
