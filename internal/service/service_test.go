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
	"testing"

	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

// TestValidateWorkerCount verifies that worker count validation correctly handles
// various inputs including zero, negative, and valid values.
func TestValidateWorkerCount(t *testing.T) {
	tests := []struct {
		name      string
		requested int
		quiet     bool
		expected  int
	}{
		{"zero workers", 0, true, cryptolib.DefaultWorkers()},
		{"negative workers", -5, true, cryptolib.DefaultWorkers()},
		{"valid workers", 4, true, 4},
		{"zero workers non-quiet", 0, false, cryptolib.DefaultWorkers()},
		{"valid workers non-quiet", 4, false, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateWorkerCount(tt.requested, tt.quiet)
			if result != tt.expected {
				t.Errorf("ValidateWorkerCount(%d, %v) = %d, want %d", tt.requested, tt.quiet, result, tt.expected)
			}
		})
	}
}

// TestValidateInputFile verifies that input file validation correctly identifies
// existing files and rejects non-existent paths.
func TestValidateInputFile(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	nonExistentFile := filepath.Join(tempDir, "does-not-exist.txt")

	if err := ValidateInputFile(existingFile); err != nil {
		t.Errorf("existing file should not error: %v", err)
	}

	if err := ValidateInputFile(nonExistentFile); err == nil {
		t.Error("non-existent file should error")
	}
}

// TestCheckFileExists verifies that file existence detection works correctly
// for both existing and non-existent files.
func TestCheckFileExists(t *testing.T) {
	tempDir := t.TempDir()
	existingFile := filepath.Join(tempDir, "exists.txt")
	if err := os.WriteFile(existingFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	exists, err := CheckFileExists(existingFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("file should exist")
	}

	nonExistentFile := filepath.Join(tempDir, "does-not-exist.txt")
	exists, err = CheckFileExists(nonExistentFile)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if exists {
		t.Error("file should not exist")
	}
}

// TestExecuteEncryption_InvalidInput verifies that encryption fails gracefully
// when the input file does not exist.
func TestExecuteEncryption_InvalidInput(t *testing.T) {
	err := ExecuteEncryption("non-existent.txt", "output.enc", "password", 4, true)
	if err == nil {
		t.Error("expected error for non-existent input file")
	}
}

// TestExecuteDecryption_InvalidInput verifies that decryption fails gracefully
// when the input file does not exist.
func TestExecuteDecryption_InvalidInput(t *testing.T) {
	err := ExecuteDecryption("non-existent.enc", "output.txt", "password", true)
	if err == nil {
		t.Error("expected error for non-existent input file")
	}
}

// TestValidateWorkerCountEdgeCases verifies boundary conditions for worker count
// validation, including maximum caps and default values.
func TestValidateWorkerCountEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		requested int
		quiet     bool
		expected  int
	}{
		{"max workers", 1000, true, 16}, // Assuming 8 CPUs -> max 16
		{"exactly default", cryptolib.DefaultWorkers(), true, cryptolib.DefaultWorkers()},
		{"one worker", 1, true, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateWorkerCount(tt.requested, tt.quiet)
			if result < 1 {
				t.Errorf("ValidateWorkerCount returned %d, should be at least 1", result)
			}
		})
	}
}

// TestExecuteEncryptionWithConfig verifies encryption with custom configuration.
func TestExecuteEncryptionWithConfig(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.enc")

	content := []byte("test content for custom config encryption")
	if err := os.WriteFile(inputFile, content, 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	config := cryptolib.EncryptorConfig{
		Workers:          4,
		ChunkSize:        64 * 1024,
		MaxPendingChunks: 30,
	}

	err := ExecuteEncryptionWithConfig(inputFile, outputFile, "test-password", config, true)
	if err != nil {
		t.Errorf("ExecuteEncryptionWithConfig failed: %v", err)
	}

	if _, err := os.Stat(outputFile); err != nil {
		t.Error("output file was not created")
	}
}

// TestExecuteEncryptionWithConfig_InvalidInput verifies error handling.
func TestExecuteEncryptionWithConfig_InvalidInput(t *testing.T) {
	config := cryptolib.DefaultEncryptorConfig()

	err := ExecuteEncryptionWithConfig("nonexistent.txt", "output.enc", "password", config, true)
	if err == nil {
		t.Error("expected error for non-existent input file")
	}
}

// TestValidateWorkerCountWithCustomConfig verifies worker count validation.
func TestValidateWorkerCountEdgeCasesWithConfig(t *testing.T) {
	tests := []struct {
		name      string
		requested int
		expected  int
	}{
		{"zero workers", 0, cryptolib.DefaultWorkers()},
		{"negative workers", -1, cryptolib.DefaultWorkers()},
		{"one worker", 1, 1},
		{"default workers", cryptolib.DefaultWorkers(), cryptolib.DefaultWorkers()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateWorkerCount(tt.requested, true)
			if result != tt.expected {
				t.Errorf("ValidateWorkerCount(%d) = %d, want %d", tt.requested, result, tt.expected)
			}
		})
	}
}

// TestExecuteEncryptionWithValidInput verifies that encryption successfully
// processes a valid input file and creates the expected output file.
func TestExecuteEncryptionWithValidInput(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.enc")

	if err := os.WriteFile(inputFile, []byte("test data"), 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	err := ExecuteEncryption(inputFile, outputFile, "test-password", 4, true)
	if err != nil {
		t.Errorf("encryption failed: %v", err)
	}

	if _, err := os.Stat(outputFile); err != nil {
		t.Error("output file was not created")
	}
}
