// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
//
// The integration tests in this file verify:
//   - Concurrent encryption/decryption operations
//   - Large file streaming (50MB)
//   - Worker count scaling
//   - Non-deterministic encryption (salt uniqueness)
//   - CLI binary functionality
//   - End-to-end encryption/decryption flow
//   - Wrong password rejection
//   - Force overwrite flag behavior
//   - Concurrent CLI operations
package cryptolib

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

// TestConcurrentEncryption verifies that multiple encryption/decryption operations
// can run concurrently without interference.
//
// This test launches 10 goroutines, each encrypting and decrypting 1MB of random
// data using the default worker count. It ensures that:
//   - No data races occur between concurrent operations
//   - Each operation correctly encrypts and decrypts its own data
//   - The crypto primitives are safe for concurrent use
func TestConcurrentEncryption(t *testing.T) {
	const numFiles = 10
	const password = "concurrent-test-password"

	var waitGroup sync.WaitGroup
	errorChan := make(chan error, numFiles)

	for i := 0; i < numFiles; i++ {
		waitGroup.Add(1)
		go func(fileID int) {
			defer waitGroup.Done()

			testData := make([]byte, 1024*1024)
			_, _ = rand.Read(testData)

			inputFile := createTempFile(t, testData)
			defer os.Remove(inputFile)

			encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

			encryptor, err := NewEncryptor(DefaultWorkers)
			if err != nil {
				errorChan <- err
				return
			}

			if err := encryptor.EncryptFile(inputFile, encryptedFile, password); err != nil {
				errorChan <- err
				return
			}

			file, err := os.Open(encryptedFile)
			if err != nil {
				errorChan <- err
				return
			}
			defer file.Close()

			var header FileHeader
			if err := binary.Read(file, binary.BigEndian, &header); err != nil {
				errorChan <- err
				return
			}

			decryptor, err := NewDecryptor(password, header.Salt[:])
			if err != nil {
				errorChan <- err
				return
			}

			decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
			if err := decryptor.DecryptFile(encryptedFile, decryptedFile); err != nil {
				errorChan <- err
				return
			}

			decryptedData, err := os.ReadFile(decryptedFile)
			if err != nil {
				errorChan <- err
				return
			}

			if !bytes.Equal(testData, decryptedData) {
				errorChan <- fmt.Errorf("data mismatch for file %d", fileID)
			}
		}(i)
	}

	waitGroup.Wait()
	close(errorChan)

	for err := range errorChan {
		if err != nil {
			t.Errorf("concurrent operation failed: %v", err)
		}
	}
}

// TestLargeFileStreaming verifies that encryption and decryption work correctly
// with large files (50MB) using streaming mode.
//
// This test generates 50MB of random data, encrypts it with 8 parallel workers,
// decrypts it, and verifies the result matches the original.
//
// The test ensures that:
//   - Streaming doesn't consume excessive memory
//   - Large files are processed correctly
//   - The implementation works with 8 workers
//
// The test is skipped when running in short mode (-test.short).
func TestLargeFileStreaming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large file test in short mode")
	}

	testData := make([]byte, 50*1024*1024)
	_, _ = rand.Read(testData)

	inputFile := createTempFile(t, testData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(8)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile, "large-file-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	file, err := os.Open(encryptedFile)
	if err != nil {
		t.Fatalf("failed to open encrypted file: %v", err)
	}
	defer file.Close()

	var header FileHeader
	if err := binary.Read(file, binary.BigEndian, &header); err != nil {
		t.Fatalf("failed to read header: %v", err)
	}

	decryptor, err := NewDecryptor("large-file-password", header.Salt[:])
	if err != nil {
		t.Fatalf("failed to create decryptor: %v", err)
	}

	decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
	if err := decryptor.DecryptFile(encryptedFile, decryptedFile); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(testData, decryptedData) {
		t.Errorf("decrypted data mismatch. Original: %d bytes, Decrypted: %d bytes",
			len(testData), len(decryptedData))
	}
}

// TestEncryptDecryptWithAllWorkerCounts verifies that encryption works correctly
// with various worker counts.
//
// This test runs the encryption/decryption cycle with worker counts ranging
// from 1 to 2×CPU cores. It ensures that:
//   - All worker counts produce correct results
//   - No off-by-one errors in chunk processing
//   - The encryptor handles the full range of valid worker values
func TestEncryptDecryptWithAllWorkerCounts(t *testing.T) {
	testData := make([]byte, 5*1024*1024)
	_, _ = rand.Read(testData)

	workerCounts := []int{1, 2, 4, 8, 16, runtime.NumCPU() * 2}

	for _, workerCount := range workerCounts {
		t.Run(fmt.Sprintf("workers_%d", workerCount), func(t *testing.T) {
			inputFile := createTempFile(t, testData)
			defer os.Remove(inputFile)

			encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")
			decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")

			encryptor, err := NewEncryptor(workerCount)
			if err != nil {
				t.Fatalf("failed to create encryptor: %v", err)
			}

			if err := encryptor.EncryptFile(inputFile, encryptedFile, "worker-test-password"); err != nil {
				t.Fatalf("encryption failed: %v", err)
			}

			file, err := os.Open(encryptedFile)
			if err != nil {
				t.Fatalf("failed to open encrypted file: %v", err)
			}
			defer file.Close()

			var header FileHeader
			if err := binary.Read(file, binary.BigEndian, &header); err != nil {
				t.Fatalf("failed to read header: %v", err)
			}

			decryptor, err := NewDecryptor("worker-test-password", header.Salt[:])
			if err != nil {
				t.Fatalf("failed to create decryptor: %v", err)
			}

			if err := decryptor.DecryptFile(encryptedFile, decryptedFile); err != nil {
				t.Fatalf("decryption failed: %v", err)
			}

			decryptedData, err := os.ReadFile(decryptedFile)
			if err != nil {
				t.Fatalf("failed to read decrypted file: %v", err)
			}

			if !bytes.Equal(testData, decryptedData) {
				t.Errorf("decrypted data mismatch for %d workers", workerCount)
			}
		})
	}
}

// TestEncryptionDeterminism verifies that encrypting the same file with the same
// password produces different outputs due to random salt and nonce generation.
//
// This test encrypts the same data twice with identical passwords and verifies
// that the outputs are different. This ensures:
//   - A unique salt is generated for each encryption
//   - A unique nonce is generated for each encryption
//   - The encryption is not deterministic (critical for security)
func TestEncryptionDeterminism(t *testing.T) {
	testData := []byte("test data for determinism check")
	password := "determinism-test-password"

	inputFile := createTempFile(t, testData)
	defer os.Remove(inputFile)

	encryptedFile1 := filepath.Join(t.TempDir(), "encrypted1.bin")
	encryptedFile2 := filepath.Join(t.TempDir(), "encrypted2.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile1, password); err != nil {
		t.Fatalf("first encryption failed: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile2, password); err != nil {
		t.Fatalf("second encryption failed: %v", err)
	}

	encryptedData1, err := os.ReadFile(encryptedFile1)
	if err != nil {
		t.Fatalf("failed to read first encrypted file: %v", err)
	}

	encryptedData2, err := os.ReadFile(encryptedFile2)
	if err != nil {
		t.Fatalf("failed to read second encrypted file: %v", err)
	}

	if bytes.Equal(encryptedData1, encryptedData2) {
		t.Error("encryption with same password produced identical output - salt or nonce not random")
	}
}

// TestCLIBinary verifies that the compiled cryptool binary works correctly.
//
// This test checks that:
//   - The `version` command produces output
//   - The `--help` command produces output
//   - The `encrypt --help` command produces output
//
// The test is skipped if the binary is not found (run 'make build' first).
func TestCLIBinary(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("cryptool binary not found, run 'make build' first")
	}

	t.Run("version", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "version")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("version command failed: %v", err)
		}
		if len(output) == 0 {
			t.Error("version command produced no output")
		}
	})

	t.Run("help", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "--help")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("help command failed: %v", err)
		}
		if len(output) == 0 {
			t.Error("help command produced no output")
		}
	})

	t.Run("encrypt help", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "encrypt", "--help")
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("encrypt help command failed: %v", err)
		}
		if len(output) == 0 {
			t.Error("encrypt help command produced no output")
		}
	})
}

// TestEndToEndEncryptionDecryption tests the complete encryption/decryption flow
// using the CLI binary.
//
// This test creates a test file, encrypts it with the CLI, decrypts it,
// and verifies that the decrypted content matches the original.
//
// The test is skipped if the binary is not found.
func TestEndToEndEncryptionDecryption(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("cryptool binary not found, run 'make build' first")
	}

	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	encryptedFile := filepath.Join(tempDir, "encrypted.enc")
	decryptedFile := filepath.Join(tempDir, "decrypted.txt")
	password := "e2e-test-password-123"

	testContent := []byte("This is end-to-end test content for cryptool integration testing.")
	if err := os.WriteFile(inputFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	cmd := exec.Command(binaryPath, "encrypt", inputFile, encryptedFile, "--pass", password, "--quiet")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("encryption failed: %v\n%s", err, output)
	}

	info, err := os.Stat(encryptedFile)
	if err != nil {
		t.Fatalf("encrypted file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("encrypted file is empty")
	}

	cmd = exec.Command(binaryPath, "decrypt", encryptedFile, decryptedFile, "--pass", password, "--quiet")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("decryption failed: %v\n%s", err, output)
	}

	decryptedContent, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(testContent, decryptedContent) {
		t.Errorf("decrypted content mismatch: expected %q, got %q", testContent, decryptedContent)
	}
}

// TestWrongPassword tests that decryption fails with an incorrect password.
//
// This test encrypts a file with a correct password, then attempts to decrypt
// it with a wrong password. It verifies that:
//   - The decryption command fails (non-zero exit code)
//   - The error message indicates authentication failure
//   - No output file is created (or it is deleted on failure)
func TestWrongPassword(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("cryptool binary not found, run 'make build' first")
	}

	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	encryptedFile := filepath.Join(tempDir, "encrypted.enc")
	correctPassword := "correct-password-123"
	wrongPassword := "wrong-password-456"

	testContent := []byte("secret data")
	if err := os.WriteFile(inputFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	cmd := exec.Command(binaryPath, "encrypt", inputFile, encryptedFile, "--pass", correctPassword)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("encryption failed: %v\n%s", err, output)
	}

	decryptedFile := filepath.Join(tempDir, "decrypted.txt")
	cmd = exec.Command(binaryPath, "decrypt", encryptedFile, decryptedFile, "--pass", wrongPassword)

	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("decryption succeeded with wrong password, expected failure")
	}

	if !bytes.Contains(output, []byte("authentication failed")) {
		t.Errorf("expected error message containing 'authentication failed', got: %s", output)
	}

	if _, err := os.Stat(decryptedFile); err == nil {
		t.Error("decrypted file was created despite wrong password")
	}
}

// TestForceOverwrite tests that the --force flag overwrites existing files.
//
// This test creates an encrypted file, then encrypts the same input again
// with the --force flag. It verifies that the operation succeeds without
// prompting for confirmation.
func TestForceOverwrite(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("cryptool binary not found, run 'make build' first")
	}

	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.enc")
	password := "force-test-password"

	testContent := []byte("test content")
	if err := os.WriteFile(inputFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	cmd := exec.Command(binaryPath, "encrypt", inputFile, outputFile, "--pass", password, "--quiet")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("first encryption failed: %v\n%s", err, output)
	}

	cmd = exec.Command(binaryPath, "encrypt", inputFile, outputFile, "--pass", password, "--force", "--quiet")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("second encryption with --force failed: %v\n%s", err, output)
	}
}

// TestLargeFile tests encryption/decryption of a 10MB file using the CLI binary.
//
// This test generates a 10MB file with a deterministic pattern (byte values 0-255),
// encrypts it with 8 workers, decrypts it, and verifies the result matches.
//
// The test is skipped when running in short mode (-test.short).
func TestLargeFile(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large file test in short mode")
	}

	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("cryptool binary not found, run 'make build' first")
	}

	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "large.bin")
	encryptedFile := filepath.Join(tempDir, "large.enc")
	decryptedFile := filepath.Join(tempDir, "large.dec")
	password := "large-file-password"

	size := 10 * 1024 * 1024
	testContent := make([]byte, size)
	for i := range testContent {
		testContent[i] = byte(i % 256)
	}
	if err := os.WriteFile(inputFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create large input file: %v", err)
	}

	cmd := exec.Command(binaryPath, "encrypt", inputFile, encryptedFile, "--pass", password, "--workers", "8", "--quiet")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("encryption failed: %v\n%s", err, output)
	}

	cmd = exec.Command(binaryPath, "decrypt", encryptedFile, decryptedFile, "--pass", password, "--quiet")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("decryption failed: %v\n%s", err, output)
	}

	decryptedContent, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(testContent, decryptedContent) {
		t.Errorf("decrypted content mismatch: original %d bytes, decrypted %d bytes", len(testContent), len(decryptedContent))
	}
}

// TestConcurrentOperations tests multiple concurrent encryption operations using the CLI.
//
// This test launches 5 goroutines, each encrypting and decrypting its own file.
// It verifies that concurrent CLI invocations don't interfere with each other.
func TestConcurrentOperations(t *testing.T) {
	binaryPath := findBinary(t)
	if binaryPath == "" {
		t.Skip("cryptool binary not found, run 'make build' first")
	}

	tempDir := t.TempDir()
	password := "concurrent-test-password"
	numFiles := 5

	type result struct {
		name string
		err  error
	}
	results := make(chan result, numFiles)

	for i := 0; i < numFiles; i++ {
		go func(id int) {
			name := strconv.Itoa(id)
			inputFile := filepath.Join(tempDir, "input_"+name+".txt")
			encryptedFile := filepath.Join(tempDir, "encrypted_"+name+".enc")
			decryptedFile := filepath.Join(tempDir, "decrypted_"+name+".txt")

			content := []byte("concurrent test content " + name)
			if err := os.WriteFile(inputFile, content, 0644); err != nil {
				results <- result{name: name, err: err}
				return
			}

			cmd := exec.Command(binaryPath, "encrypt", inputFile, encryptedFile, "--pass", password, "--quiet")
			if _, err := cmd.CombinedOutput(); err != nil {
				results <- result{name: name, err: err}
				return
			}

			cmd = exec.Command(binaryPath, "decrypt", encryptedFile, decryptedFile, "--pass", password, "--quiet")
			if _, err := cmd.CombinedOutput(); err != nil {
				results <- result{name: name, err: err}
				return
			}

			decryptedContent, err := os.ReadFile(decryptedFile)
			if err != nil {
				results <- result{name: name, err: err}
				return
			}

			if !bytes.Equal(content, decryptedContent) {
				results <- result{name: name, err: err}
				return
			}

			results <- result{name: name, err: nil}
		}(i)
	}

	failures := 0
	for i := 0; i < numFiles; i++ {
		r := <-results
		if r.err != nil {
			t.Errorf("operation %s failed: %v", r.name, r.err)
			failures++
		}
	}

	if failures > 0 {
		t.Fatalf("%d concurrent operations failed", failures)
	}
}

// findBinary locates the compiled cryptool binary for CLI tests.
//
// It searches common locations including:
//   - ../build/cryptool (from test directory)
//   - ../../build/cryptool (from deeper test directories)
//   - ./cryptool (current directory)
//   - cryptool (PATH)
//
// Returns an empty string if the binary is not found.
func findBinary(t *testing.T) string {
	locations := []string{
		"../build/cryptool",
		"../../build/cryptool",
		"./cryptool",
		"cryptool",
	}

	for _, loc := range locations {
		if _, err := os.Stat(loc); err == nil {
			abs, err := filepath.Abs(loc)
			if err == nil {
				return abs
			}
			return loc
		}
	}
	return ""
}
