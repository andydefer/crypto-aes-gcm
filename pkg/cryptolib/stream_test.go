// Package cryptolib provides secure, streaming AES-256-GCM encryption with Argon2id key derivation.
//
// The package implements authenticated encryption with parallel chunk processing for large files,
// and HMAC-SHA256 integrity verification for headers.
package cryptolib

import (
	"bytes"
	"crypto/rand"
	"os"
	"path/filepath"
	"testing"
)

// TestDecryptStream verifies the streaming decryption convenience function.
//
// This test ensures that DecryptStream correctly decrypts data that was
// encrypted with the standard Encryptor.
func TestDecryptStream(t *testing.T) {
	originalData := []byte("test data for streaming decryption")
	reader := bytes.NewReader(originalData)
	var encryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.Encrypt(reader, &encryptedBuf, "stream-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, "stream-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(originalData, decryptedBuf.Bytes()) {
		t.Errorf("decrypted data mismatch.\nOriginal: %s\nDecrypted: %s",
			originalData, decryptedBuf.Bytes())
	}
}

// TestDecryptStream_EmptyFile verifies that DecryptStream handles empty files correctly.
func TestDecryptStream_EmptyFile(t *testing.T) {
	originalData := []byte{}
	reader := bytes.NewReader(originalData)
	var encryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.Encrypt(reader, &encryptedBuf, "empty-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, "empty-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if len(decryptedBuf.Bytes()) != 0 {
		t.Errorf("expected empty decrypted data, got %d bytes", len(decryptedBuf.Bytes()))
	}
}

// TestDecryptStream_WrongPassword verifies that DecryptStream fails with incorrect password.
func TestDecryptStream_WrongPassword(t *testing.T) {
	originalData := []byte("secret data")
	reader := bytes.NewReader(originalData)
	var encryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.Encrypt(reader, &encryptedBuf, "correct-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	err = DecryptStream(encryptedReader, &decryptedBuf, "wrong-password")
	if err == nil {
		t.Error("decryption succeeded with wrong password, expected failure")
	}
}

// TestDecryptStream_CorruptedData verifies that DecryptStream detects corrupted ciphertext.
func TestDecryptStream_CorruptedData(t *testing.T) {
	originalData := []byte("test data for corruption test")
	reader := bytes.NewReader(originalData)
	var encryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.Encrypt(reader, &encryptedBuf, "corruption-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	encryptedData := encryptedBuf.Bytes()

	// ✅ Vérification que la donnée est assez longue pour être corrompue
	if len(encryptedData) < 100 {
		t.Skip("encrypted data too short for corruption test")
	}

	encryptedData[100] ^= 0xFF

	var decryptedBuf bytes.Buffer
	corruptedReader := bytes.NewReader(encryptedData)

	err = DecryptStream(corruptedReader, &decryptedBuf, "corruption-password")
	if err == nil {
		t.Error("decryption succeeded with corrupted data, expected failure")
	}
}

// TestDecryptStream_LargeData verifies streaming decryption works with large data.
//
// This test is skipped in short mode due to memory requirements.
func TestDecryptStream_LargeData(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large data test in short mode")
	}

	testData := make([]byte, 10*1024*1024)
	_, _ = rand.Read(testData)

	reader := bytes.NewReader(testData)
	var encryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.Encrypt(reader, &encryptedBuf, "large-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, "large-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Errorf("decrypted data mismatch. Original: %d bytes, Decrypted: %d bytes",
			len(testData), len(decryptedBuf.Bytes()))
	}
}

// TestDecryptStream_FileInterface verifies DecryptStream works with actual files.
func TestDecryptStream_FileInterface(t *testing.T) {
	testData := []byte("file-based streaming test")
	inputFile := createTempFile(t, testData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.EncryptFile(inputFile, encryptedFile, "file-password"); err != nil {
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

	if err := DecryptStream(encryptedReader, decryptedWriter, "file-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(testData, decryptedData) {
		t.Errorf("decrypted data mismatch.\nOriginal: %s\nDecrypted: %s",
			testData, decryptedData)
	}
}

// TestDecryptStream_ChunkBoundary verifies decryption works across chunk boundaries.
func TestDecryptStream_ChunkBoundary(t *testing.T) {
	// Create data that spans multiple chunks
	testData := make([]byte, DefaultChunkSize*3+1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	reader := bytes.NewReader(testData)
	var encryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if err := encryptor.Encrypt(reader, &encryptedBuf, "boundary-password"); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, "boundary-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(testData, decryptedBuf.Bytes()) {
		t.Errorf("decrypted data mismatch across chunk boundaries")
	}
}
