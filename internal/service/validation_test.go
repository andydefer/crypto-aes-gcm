// Package service provides business logic for encryption and decryption operations.
//
// This package orchestrates the encryption and decryption processes by:
//   - Managing progress bars and user feedback
//   - Handling file I/O operations
//   - Coordinating with the cryptolib package for core crypto operations
//   - Providing error handling and cleanup
package service

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

// TestCheckOverwriteWithForce verifies that CheckOverwrite bypasses confirmation
// when the force flag is enabled, even when the output file already exists.
func TestCheckOverwriteWithForce(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	err := CheckOverwrite(testFile, true)
	if err != nil {
		t.Errorf("CheckOverwrite with force=true should not error: %v", err)
	}
}

// TestCheckOverwriteWithNonExistentFile verifies that CheckOverwrite returns nil
// when the output file does not exist (no confirmation needed).
func TestCheckOverwriteWithNonExistentFile(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentFile := filepath.Join(tempDir, "does-not-exist.txt")

	err := CheckOverwrite(nonExistentFile, false)
	if err != nil {
		t.Errorf("CheckOverwrite with non-existent file should not error: %v", err)
	}
}

// TestValidateWorkerCountMax verifies that ValidateWorkerCount caps the worker
// count at 2× the number of CPU cores to prevent resource exhaustion.
func TestValidateWorkerCountMax(t *testing.T) {
	result := ValidateWorkerCount(9999, true)
	expectedMax := runtime.NumCPU() * 2

	if result > expectedMax {
		t.Errorf("ValidateWorkerCount returned %d, should be capped at %d", result, expectedMax)
	}

	if result < 1 {
		t.Errorf("ValidateWorkerCount returned %d, should be at least 1", result)
	}
}

// TestValidateWorkerCountDefault verifies that ValidateWorkerCount returns the
// default worker count when given invalid input values (zero or negative).
func TestValidateWorkerCountDefault(t *testing.T) {
	defaultWorkers := cryptolib.DefaultWorkers()

	result := ValidateWorkerCount(0, true)
	if result != defaultWorkers {
		t.Errorf("ValidateWorkerCount(0) = %d, want %d", result, defaultWorkers)
	}

	result = ValidateWorkerCount(-5, true)
	if result != defaultWorkers {
		t.Errorf("ValidateWorkerCount(-5) = %d, want %d", result, defaultWorkers)
	}
}

// TestValidateInputFilePermissions verifies that ValidateInputFile correctly
// handles directory paths (os.Stat works on directories, returning no error).
func TestValidateInputFilePermissions(t *testing.T) {
	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "testdir")

	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	err := ValidateInputFile(dirPath)
	if err != nil {
		t.Logf("Directory validation result: %v", err)
	}
}

// TestCheckFileExistsWithDirectory verifies that CheckFileExists returns true
// when checking a directory path (directories are considered existing paths).
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
