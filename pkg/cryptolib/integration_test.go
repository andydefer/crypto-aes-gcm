package cryptolib

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestConcurrentEncryption(t *testing.T) {
	const numFiles = 10
	const password = "concurrent-test-password"

	var wg sync.WaitGroup
	errors := make(chan error, numFiles)

	for i := 0; i < numFiles; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			data := make([]byte, 1024*1024) // 1MB each
			rand.Read(data)

			inputFile := createTempFile(t, data)
			defer os.Remove(inputFile)

			encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

			encryptor, err := NewEncryptor(DefaultWorkers)
			if err != nil {
				errors <- err
				return
			}

			err = encryptor.EncryptFile(inputFile, encryptedFile, password)
			if err != nil {
				errors <- err
				return
			}

			f, err := os.Open(encryptedFile)
			if err != nil {
				errors <- err
				return
			}
			defer f.Close()

			var header FileHeader
			if err := binary.Read(f, binary.BigEndian, &header); err != nil {
				errors <- err
				return
			}

			decryptor, err := NewDecryptor(password, header.Salt[:])
			if err != nil {
				errors <- err
				return
			}

			decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
			err = decryptor.DecryptFile(encryptedFile, decryptedFile)
			if err != nil {
				errors <- err
				return
			}

			decryptedData, err := os.ReadFile(decryptedFile)
			if err != nil {
				errors <- err
				return
			}

			if !bytes.Equal(data, decryptedData) {
				errors <- fmt.Errorf("data mismatch")
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Errorf("concurrent operation failed: %v", err)
		}
	}
}

func TestLargeFileStreaming(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping large file test in short mode")
	}

	data := make([]byte, 50*1024*1024)
	rand.Read(data)

	inputFile := createTempFile(t, data)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(8)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	err = encryptor.EncryptFile(inputFile, encryptedFile, "large-file-password")
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

	decryptor, err := NewDecryptor("large-file-password", header.Salt[:])
	if err != nil {
		t.Fatalf("failed to create decryptor: %v", err)
	}

	decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
	err = decryptor.DecryptFile(encryptedFile, decryptedFile)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	decryptedData, err := os.ReadFile(decryptedFile)
	if err != nil {
		t.Fatalf("failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(data, decryptedData) {
		t.Errorf("decrypted data doesn't match original. Original size: %d, Decrypted size: %d",
			len(data), len(decryptedData))
	}
}

func createTempFile(t *testing.T, data []byte) string {
	t.Helper()
	tmpFile := filepath.Join(t.TempDir(), "input.txt")
	err := os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return tmpFile
}
