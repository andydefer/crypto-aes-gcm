package cryptolib

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"testing"
)

func TestDecryptor_CorruptedFile(t *testing.T) {
	originalData := []byte("test data for corruption testing")
	inputFile := createTempFile(t, originalData)
	defer os.Remove(inputFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	err = encryptor.EncryptFile(inputFile, encryptedFile, "corruption-test")
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	testCases := []struct {
		name        string
		corruptFunc func([]byte) []byte
	}{
		{
			name: "corrupt magic bytes",
			corruptFunc: func(data []byte) []byte {
				if len(data) > 0 {
					data[0] = 0xFF
				}
				return data
			},
		},
		{
			name: "corrupt header HMAC",
			corruptFunc: func(data []byte) []byte {
				if len(data) > 25+32 {
					data[25+16] = 0xFF
				}
				return data
			},
		},
		{
			name: "corrupt ciphertext",
			corruptFunc: func(data []byte) []byte {
				if len(data) > 100 {
					data[100] = 0xFF
				}
				return data
			},
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
			err = os.WriteFile(corruptedFile, corruptedData, 0644)
			if err != nil {
				t.Fatalf("failed to write corrupted file: %v", err)
			}

			f, err := os.Open(corruptedFile)
			if err != nil {
				t.Fatalf("failed to open corrupted file: %v", err)
			}
			defer f.Close()

			var header FileHeader
			if err := binary.Read(f, binary.BigEndian, &header); err != nil {
				return
			}

			decryptor, err := NewDecryptor("corruption-test", header.Salt[:])
			if err != nil {
				t.Fatalf("failed to create decryptor: %v", err)
			}

			decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
			err = decryptor.DecryptFile(corruptedFile, decryptedFile)

			if err == nil {
				t.Error("expected decryption to fail with corrupted file, but it succeeded")
			}
		})
	}
}

func TestDecryptor_EmptyFile(t *testing.T) {
	emptyFile := createTempFile(t, []byte{})
	defer os.Remove(emptyFile)

	encryptedFile := filepath.Join(t.TempDir(), "encrypted.bin")

	encryptor, err := NewEncryptor(DefaultWorkers)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	err = encryptor.EncryptFile(emptyFile, encryptedFile, "password")
	if err != nil {
		t.Fatalf("encryption of empty file failed: %v", err)
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

	decryptor, err := NewDecryptor("password", header.Salt[:])
	if err != nil {
		t.Fatalf("failed to create decryptor: %v", err)
	}

	decryptedFile := filepath.Join(t.TempDir(), "decrypted.txt")
	err = decryptor.DecryptFile(encryptedFile, decryptedFile)

	if err != nil {
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
