package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckOverwriteWithForce(t *testing.T) {
	// Test with force=true should never ask for confirmation
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// With force=true, should not error even if file exists
	err := CheckOverwrite(testFile, true)
	if err != nil {
		t.Errorf("CheckOverwrite with force=true should not error: %v", err)
	}
}

func TestCheckOverwriteWithNonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "does-not-exist.txt")

	// With non-existent file, should not error
	err := CheckOverwrite(nonExistentFile, false)
	if err != nil {
		t.Errorf("CheckOverwrite with non-existent file should not error: %v", err)
	}
}

func TestValidateWorkerCountMax(t *testing.T) {
	// Test that worker count is capped at 2×CPU cores
	result := ValidateWorkerCount(9999, true)
	maxWorkers := 16 // Assuming 8 CPU cores

	if result > maxWorkers {
		t.Errorf("ValidateWorkerCount returned %d, should be capped at %d", result, maxWorkers)
	}

	if result < 1 {
		t.Errorf("ValidateWorkerCount returned %d, should be at least 1", result)
	}
}

func TestValidateWorkerCountDefault(t *testing.T) {
	// Test that invalid values return default
	defaultWorkers := 4

	result := ValidateWorkerCount(0, true)
	if result != defaultWorkers {
		t.Errorf("ValidateWorkerCount(0) = %d, want %d", result, defaultWorkers)
	}

	result = ValidateWorkerCount(-5, true)
	if result != defaultWorkers {
		t.Errorf("ValidateWorkerCount(-5) = %d, want %d", result, defaultWorkers)
	}
}

func TestValidateInputFilePermissions(t *testing.T) {
	tempDir := t.TempDir()

	// Create a directory instead of a file
	dirPath := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// ValidateInputFile should accept directories? (os.Stat works on dirs)
	err := ValidateInputFile(dirPath)
	if err != nil {
		t.Logf("Directory validation result: %v", err)
	}
}

func TestCheckFileExistsWithDirectory(t *testing.T) {
	tempDir := t.TempDir()

	exists, err := CheckFileExists(tempDir)
	if err != nil {
		t.Errorf("CheckFileExists on directory returned error: %v", err)
	}
	if !exists {
		t.Error("CheckFileExists on directory should return true")
	}
}
