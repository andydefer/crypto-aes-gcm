// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
package cryptolib

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"

	"github.com/andydefer/crypto-aes-gcm/internal/constants"
)

// TestEncryptor_EncryptDecrypt verifies the complete encryption and decryption cycle.
//
// It tests various scenarios including:
//   - Small and large data sizes
//   - Different worker counts for parallel processing
//   - Different passwords
//   - Empty file handling
func TestEncryptor_EncryptDecrypt(t *testing.T) {
	smallData := []byte("This is secret data that needs to be encrypted. " +
		"It contains multiple chunks to test streaming encryption. " +
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
		"Repeat this several times to ensure chunking works correctly. " +
		"More data to make sure we exceed the chunk size of 1MB. " +
		"Adding even more text to create multiple chunks for testing.")

	largeData := make([]byte, 5*constants.MB)
	_, _ = rand.Read(largeData)

	testCases := []struct {
		name     string
		data     []byte
		workers  int
		password string
	}{
		{
			name:     "small data with single worker",
			data:     smallData,
			workers:  1,
			password: "test-password-123",
		},
		{
			name:     "small data with multiple workers",
			data:     smallData,
			workers:  4,
			password: "another-password",
		},
		{
			name:     "large data 5MB",
			data:     largeData,
			workers:  4,
			password: "strong-password-with-special-chars!@#$%",
		},
		{
			name:     "empty data",
			data:     []byte{},
			workers:  2,
			password: "password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputFile := createTempFile(t, tc.data)
			defer os.Remove(inputFile)

			encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")
			decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")

			encryptor, err := NewEncryptor(tc.workers)
			if err != nil {
				t.Fatalf("failed to create encryptor: %v", err)
			}

			if err := encryptor.EncryptFile(inputFile, encryptedFile, tc.password); err != nil {
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

			decryptor, err := NewDecryptor(tc.password, header.Salt[:])
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

			if !bytes.Equal(tc.data, decryptedData) {
				t.Errorf("decrypted data mismatch: original %d bytes, decrypted %d bytes",
					len(tc.data), len(decryptedData))
			}
		})
	}
}

// TestEncryptor_WrongPassword verifies that decryption fails with an incorrect password.
//
// This test ensures that the authentication mechanism properly rejects
// decryption attempts with wrong credentials.
func TestEncryptor_WrongPassword(t *testing.T) {
	originalData := []byte("secret data")
	inputFile := createTempFile(t, originalData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers())
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile, "correct-password"); err != nil {
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

	decryptor, err := NewDecryptor("wrong-password", header.Salt[:])
	if err != nil {
		t.Fatalf("failed to create decryptor: %v", err)
	}

	decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
	err = decryptor.DecryptFile(encryptedFile, decryptedFile)

	if err == nil {
		t.Error("decryption succeeded with wrong password, expected failure")
	}
}

// TestEncryptor_StreamingInterface verifies the streaming reader/writer API.
//
// This test uses the streaming interfaces directly without temporary files,
// ensuring the encryption and decryption work with any io.Reader/io.Writer.
func TestEncryptor_StreamingInterface(t *testing.T) {
	originalData := []byte("streaming test data for reader/writer interface")
	reader := bytes.NewReader(originalData)
	var encryptedBuf bytes.Buffer
	var decryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers())
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.Encrypt(reader, &encryptedBuf, "stream-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())
	if err := DecryptStream(encryptedReader, &decryptedBuf, "stream-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(originalData, decryptedBuf.Bytes()) {
		t.Errorf("decrypted data mismatch:\nOriginal: %s\nDecrypted: %s",
			originalData, decryptedBuf.Bytes())
	}
}

// TestEncryptor_InvalidWorkers verifies that worker count validation works correctly.
//
// The encryptor should clamp invalid worker values to safe defaults.
func TestEncryptor_InvalidWorkers(t *testing.T) {
	testCases := []struct {
		name    string
		workers int
	}{
		{"zero workers", 0},
		{"negative workers", -5},
		{"excessive workers", 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptor, err := NewEncryptor(tc.workers)
			if err != nil {
				t.Fatalf("NewEncryptor failed: %v", err)
			}

			if encryptor.workers <= 0 {
				t.Errorf("workers should be > 0, got %d", encryptor.workers)
			}
		})
	}
}

// TestEncryptor_MemoryUsage verifies that encryption handles large files correctly.
//
// This test generates a 10MB file and ensures the encrypted output has
// reasonable size expectations (larger than original due to crypto overhead).
//
// The test is skipped when running in short mode (-test.short) due to the
// large file generation.
func TestEncryptor_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping memory test in short mode")
	}

	testData := make([]byte, 10*constants.MB)
	_, _ = rand.Read(testData)

	inputFile := createTempFile(t, testData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers())
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile, "memory-test-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	fileInfo, err := os.Stat(encryptedFile)
	if err != nil {
		t.Fatalf("failed to stat encrypted file: %v", err)
	}

	if fileInfo.Size() < int64(len(testData)) {
		t.Errorf("encrypted file size %d is smaller than original %d",
			fileInfo.Size(), len(testData))
	}
}

// TestEncryptorWithCustomConfig verifies the new configuration API works.
func TestEncryptorWithCustomConfig(t *testing.T) {
	testData := []byte("test data for custom config")

	config := EncryptorConfig{
		Workers:          2,
		ChunkSize:        32 * constants.KB,
		MaxPendingChunks: 20,
	}

	encryptor, err := NewEncryptorWithConfig(config)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	var encryptedBuf bytes.Buffer
	reader := bytes.NewReader(testData)

	if err := encryptor.Encrypt(reader, &encryptedBuf, "custom-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, "custom-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Errorf("data mismatch")
	}
}

// TestNewEncryptorBackwardCompatibility verifies the old API still works.
func TestNewEncryptorBackwardCompatibility(t *testing.T) {
	testData := []byte("backward compatibility test")

	encryptor, err := NewEncryptor(DefaultWorkers())
	if err != nil {
		t.Fatalf("NewEncryptor failed: %v", err)
	}

	var encryptedBuf bytes.Buffer
	reader := bytes.NewReader(testData)

	if err := encryptor.Encrypt(reader, &encryptedBuf, "password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, "password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Errorf("data mismatch")
	}
}

// TestEncryptor_ChunkSizeClamping verifies that chunk size validation works correctly.
//
// This test ensures that the encryptor clamps invalid chunk sizes to safe ranges
// and that encryption still works with clamped values.
func TestEncryptor_ChunkSizeClamping(t *testing.T) {
	testData := []byte("test data for chunk size clamping")

	testCases := []struct {
		name          string
		requestedSize int
		expectedMin   int
		expectedMax   int
	}{
		{
			name:          "negative chunk size",
			requestedSize: -100,
			expectedMin:   constants.KB,
			expectedMax:   DefaultChunkSize,
		},
		{
			name:          "zero chunk size",
			requestedSize: 0,
			expectedMin:   constants.KB,
			expectedMax:   DefaultChunkSize,
		},
		{
			name:          "very small chunk size",
			requestedSize: 512,
			expectedMin:   constants.KB,
			expectedMax:   constants.KB,
		},
		{
			name:          "very large chunk size",
			requestedSize: 100 * constants.MB,
			expectedMin:   DefaultChunkSize,
			expectedMax:   16 * constants.MB,
		},
		{
			name:          "normal chunk size",
			requestedSize: 64 * constants.KB,
			expectedMin:   64 * constants.KB,
			expectedMax:   64 * constants.KB,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        tc.requestedSize,
				MaxPendingChunks: DefaultMaxPendingChunks,
			}

			encryptor, err := NewEncryptorWithConfig(config)
			if err != nil {
				t.Fatalf("failed to create encryptor: %v", err)
			}

			if encryptor.chunkSize < tc.expectedMin {
				t.Errorf("chunk size %d is below minimum expected %d", encryptor.chunkSize, tc.expectedMin)
			}
			if encryptor.chunkSize > tc.expectedMax {
				t.Errorf("chunk size %d is above maximum expected %d", encryptor.chunkSize, tc.expectedMax)
			}

			var encryptedBuf bytes.Buffer
			reader := bytes.NewReader(testData)

			if err := encryptor.Encrypt(reader, &encryptedBuf, "clamp-test"); err != nil {
				t.Fatalf("encryption failed with clamped chunk size %d: %v", encryptor.chunkSize, err)
			}

			var decryptedBuf bytes.Buffer
			encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

			if err := DecryptStream(encryptedReader, &decryptedBuf, "clamp-test"); err != nil {
				t.Fatalf("decryption failed: %v", err)
			}

			if !bytes.Equal(testData, decryptedBuf.Bytes()) {
				t.Errorf("data mismatch after chunk size clamping")
			}
		})
	}
}

// createTempFile creates a temporary file with the provided data.
//
// Parameters:
//   - t: Testing context for fatal error reporting and cleanup
//   - data: Bytes to write to the temporary file
//
// Returns:
//   - string: Path to the created temporary file
func createTempFile(t *testing.T, data []byte) string {
	t.Helper()

	tmpFile := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return tmpFile
}
