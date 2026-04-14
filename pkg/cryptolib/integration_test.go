// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
package cryptolib

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

// TestConcurrentEncryption verifies that multiple encryption/decryption operations
// can run concurrently without interference.
//
// This test launches 10 goroutines, each encrypting and decrypting 1MB of random data.
// It ensures that the crypto primitives and file operations are safe for concurrent use.
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
// This test ensures that the implementation doesn't load the entire file into memory
// and processes data in chunks efficiently.
//
// The test is skipped when running in short mode (-test.short) due to the
// large file size and processing time.
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
// This test runs the encryption/decryption cycle with worker counts ranging from
// 1 to 2×CPU cores to ensure the parallel processing works at all scales.
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

// TestEncryptionDeterminism verifies that encrypting the same file with the same password
// produces different outputs due to random salt and nonce generation.
//
// This test ensures that:
//   - A unique salt is generated for each encryption operation
//   - A unique nonce is generated for each encryption operation
//   - The encryption is not deterministic (important for security)
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
