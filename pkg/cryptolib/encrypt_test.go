package cryptolib

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptor_EncryptDecrypt(t *testing.T) {
	originalData := []byte("This is secret data that needs to be encrypted. " +
		"It contains multiple chunks to test streaming encryption. " +
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
		"Repeat this several times to ensure chunking works correctly. " +
		"More data to make sure we exceed the chunk size of 1MB. " +
		"Adding even more text to create multiple chunks for testing.")

	largeData := make([]byte, 5*1024*1024)
	rand.Read(largeData)

	testCases := []struct {
		name     string
		data     []byte
		workers  int
		password string
	}{
		{
			name:     "small data",
			data:     originalData,
			workers:  1,
			password: "test-password-123",
		},
		{
			name:     "small data with multiple workers",
			data:     originalData,
			workers:  4,
			password: "another-password",
		},
		{
			name:     "large data (5MB)",
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

			err = encryptor.EncryptFile(inputFile, encryptedFile, tc.password)
			if err != nil {
				t.Fatalf("encryption failed: %v", err)
			}

			f, err := os.Open(encryptedFile)
			if err != nil {
				t.Fatalf("failed to open encrypted file: %v", err)
			}
			defer f.Close()

			var header FileHeader
			if err := binary.Read(f, binary.BigEndian, &header); err != nil {
				t.Fatalf("failed to read header: %v", err)
			}

			decryptor, err := NewDecryptor(tc.password, header.Salt[:])
			if err != nil {
				t.Fatalf("failed to create decryptor: %v", err)
			}

			err = decryptor.DecryptFile(encryptedFile, decryptedFile)
			if err != nil {
				t.Fatalf("decryption failed: %v", err)
			}

			decryptedData, err := os.ReadFile(decryptedFile)
			if err != nil {
				t.Fatalf("failed to read decrypted file: %v", err)
			}

			if !bytes.Equal(tc.data, decryptedData) {
				t.Errorf("decrypted data doesn't match original. Original length: %d, Decrypted length: %d",
					len(tc.data), len(decryptedData))
			}
		})
	}
}

func TestEncryptor_WrongPassword(t *testing.T) {
	originalData := []byte("secret data")
	inputFile := createTempFile(t, originalData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	err = encryptor.EncryptFile(inputFile, encryptedFile, "correct-password")
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	f, err := os.Open(encryptedFile)
	if err != nil {
		t.Fatalf("failed to open encrypted file: %v", err)
	}
	defer f.Close()

	var header FileHeader
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		t.Fatalf("failed to read header: %v", err)
	}

	decryptor, err := NewDecryptor("wrong-password", header.Salt[:])
	if err != nil {
		t.Fatalf("failed to create decryptor: %v", err)
	}

	decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
	err = decryptor.DecryptFile(encryptedFile, decryptedFile)

	if err == nil {
		t.Error("expected decryption to fail with wrong password, but it succeeded")
	}
}

func TestEncryptor_StreamingInterface(t *testing.T) {
	originalData := []byte("streaming test data for reader/writer interface")
	reader := bytes.NewReader(originalData)
	var encryptedBuf bytes.Buffer
	var decryptedBuf bytes.Buffer

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	err = encryptor.Encrypt(reader, &encryptedBuf, "stream-password")
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	// Use DecryptStream for simpler decryption
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())
	err = DecryptStream(encryptedReader, &decryptedBuf, "stream-password")
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(originalData, decryptedBuf.Bytes()) {
		t.Errorf("decrypted data doesn't match original.\nOriginal: %s\nDecrypted: %s",
			originalData, decryptedBuf.Bytes())
	}
}

func TestEncryptor_InvalidWorkers(t *testing.T) {
	testCases := []struct {
		name    string
		workers int
	}{
		{"zero workers", 0},
		{"negative workers", -5},
		{"too many workers", 1000},
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
