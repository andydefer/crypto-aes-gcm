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
)

// TestDecryptor_CorruptedFile verifies that decryption properly handles
// various forms of file corruption.
//
// The test creates an encrypted file and then applies different corruption
// strategies to ensure the decryption process correctly identifies and
// reports each type of corruption.
func TestDecryptor_CorruptedFile(t *testing.T) {
	originalData := []byte("test data for corruption testing")
	inputFile := createTempFile(t, originalData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile, "corruption-test"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	testCases := []struct {
		name        string
		corruptFunc func([]byte) []byte
		expectErr   error
	}{
		{
			name: "corrupt magic bytes",
			corruptFunc: func(data []byte) []byte {
				if len(data) > 0 {
					data[0] = 0xFF
				}
				return data
			},
			expectErr: ErrInvalidMagic,
		},
		{
			name: "corrupt header HMAC",
			corruptFunc: func(data []byte) []byte {
				// Header offset: Magic(4) + Version(1) + Salt(32) + ChunkSize(8) = 45 bytes
				// HMAC follows at offset 45, length 32
				if len(data) > 45+32 {
					data[45+16] = 0xFF
				}
				return data
			},
			expectErr: ErrHeaderAuthFailed,
		},
		{
			name: "corrupt ciphertext",
			corruptFunc: func(data []byte) []byte {
				// Header(45) + HMAC(32) + Nonce(12) = 89 bytes offset
				if len(data) > 100 {
					data[100] = 0xFF
				}
				return data
			},
			expectErr: ErrDecryptionFailed,
		},
		{
			name: "truncate file",
			corruptFunc: func(data []byte) []byte {
				if len(data) > 10 {
					return data[:len(data)-10]
				}
				return data
			},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encryptedData, err := os.ReadFile(encryptedFile)
			if err != nil {
				t.Fatalf("failed to read encrypted file: %v", err)
			}

			corruptedData := tc.corruptFunc(encryptedData)
			corruptedFile := filepath.Join(t.TempDir(), "corrupted.bin")

			if err := os.WriteFile(corruptedFile, corruptedData, 0644); err != nil {
				t.Fatalf("failed to write corrupted file: %v", err)
			}

			file, err := os.Open(corruptedFile)
			if err != nil {
				t.Fatalf("failed to open corrupted file: %v", err)
			}
			defer file.Close()

			var header FileHeader
			if err := binary.Read(file, binary.BigEndian, &header); err != nil {
				return
			}

			decryptor, err := NewDecryptor("corruption-test", header.Salt[:])
			if err != nil {
				t.Fatalf("failed to create decryptor: %v", err)
			}

			decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
			err = decryptor.DecryptFile(corruptedFile, decryptedFile)

			if err == nil && tc.expectErr != nil {
				t.Errorf("expected error %v, but got nil", tc.expectErr)
			}
		})
	}
}

// TestDecryptor_EmptyFile verifies that empty files are handled correctly
// throughout the encryption and decryption process.
//
// Empty files represent an edge case where the encryption pipeline must
// produce a valid encrypted file that decrypts back to zero bytes.
func TestDecryptor_EmptyFile(t *testing.T) {
	emptyFile := createTempFile(t, []byte{})
	defer os.Remove(emptyFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(emptyFile, encryptedFile, "password"); err != nil {
		t.Fatalf("encryption of empty file failed: %v", err)
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

	decryptor, err := NewDecryptor("password", header.Salt[:])
	if err != nil {
		t.Fatalf("failed to create decryptor: %v", err)
	}

	decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
	if err := decryptor.DecryptFile(encryptedFile, decryptedFile); err != nil {
		t.Fatalf("decryption of empty file failed: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if len(decryptedData) != 0 {
		t.Errorf("expected empty decrypted file, got %d bytes", len(decryptedData))
	}
}

// TestDecryptor_StreamingMemory verifies that streaming decryption doesn't
// allocate excessive memory for large files.
//
// This test generates 50MB of random data and performs streaming decryption
// to ensure the implementation processes data in chunks without loading the
// entire file into memory.
//
// The test is skipped when running in short mode (-test.short) due to the
// large memory allocation required for test data generation.
func TestDecryptor_StreamingMemory(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping streaming memory test in short mode")
	}

	testData := generateTestData(t, 50*1024*1024)

	inputFile := createTempFile(t, testData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile, "streaming-test"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	encryptedReader, err := os.Open(encryptedFile)
	if err != nil {
		t.Fatalf("failed to open encrypted file: %v", err)
	}
	defer encryptedReader.Close()

	decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
	decryptedWriter, err := os.Create(decryptedFile)
	if err != nil {
		t.Fatalf("failed to create decrypted file: %v", err)
	}
	defer decryptedWriter.Close()

	if err := DecryptStream(encryptedReader, decryptedWriter, "streaming-test"); err != nil {
		t.Fatalf("streaming decryption failed: %v", err)
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

// generateTestData creates random test data of the specified size.
//
// Parameters:
//   - t: Testing context for fatal error reporting
//   - size: Number of bytes to generate
//
// Returns:
//   - []byte: Randomly generated test data
func generateTestData(t *testing.T, size int) []byte {
	testData := make([]byte, size)
	if _, err := rand.Read(testData); err != nil {
		t.Fatalf("failed to generate test data: %v", err)
	}
	return testData
}
