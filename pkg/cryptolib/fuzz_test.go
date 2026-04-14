// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
//
// The package includes fuzz tests to verify:
//   - Round-trip encryption/decryption with random data
//   - Graceful handling of corrupted ciphertext
//   - Header serialization consistency
package cryptolib

import (
	"bytes"
	"testing"

	"github.com/andydefer/crypto-aes-gcm/internal/header"
)

// FuzzEncryptDecrypt verifies that encryption and decryption are inverse operations.
//
// This fuzz test generates random byte slices and ensures that:
//   - Encryption always succeeds
//   - Decryption of the ciphertext recovers the original plaintext
//   - No panics occur under any input
//
// Seed corpus includes:
//   - Simple strings ("hello world", "a")
//   - Empty data
//   - Binary data with null bytes
//   - Repeated patterns (1KB of 'A', repeated alphabet)
//
// The test runs with the default worker count (4) and a fixed password.
func FuzzEncryptDecrypt(f *testing.F) {
	f.Add([]byte("hello world"))
	f.Add([]byte(""))
	f.Add([]byte("a"))
	f.Add([]byte{0x00, 0x01, 0x02, 0x03, 0xFF})
	f.Add(bytes.Repeat([]byte("A"), 1024))
	f.Add(bytes.Repeat([]byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), 100))

	f.Fuzz(func(t *testing.T, original []byte) {
		password := "fuzz-test-password-123"

		encryptor, err := NewEncryptor(DefaultWorkers())
		if err != nil {
			t.Fatalf("failed to create encryptor: %v", err)
		}

		var encryptedBuf bytes.Buffer
		reader := bytes.NewReader(original)
		if err := encryptor.Encrypt(reader, &encryptedBuf, password); err != nil {
			t.Fatalf("encryption failed: %v", err)
		}

		var decryptedBuf bytes.Buffer
		encryptedReader := bytes.NewReader(encryptedBuf.Bytes())
		if err := DecryptStream(encryptedReader, &decryptedBuf, password); err != nil {
			t.Fatalf("decryption failed: %v", err)
		}

		if !bytes.Equal(original, decryptedBuf.Bytes()) {
			t.Errorf("round-trip failed: original %d bytes, decrypted %d bytes",
				len(original), len(decryptedBuf.Bytes()))
		}
	})
}

// FuzzDecryptCorrupted verifies that decryption fails gracefully on corrupted data.
//
// This fuzz test takes a valid encrypted file as a seed and applies random
// mutations to the ciphertext. It then verifies that:
//   - Decryption either fails (returns error) or
//   - If it succeeds, it doesn't panic (output may be garbage)
//
// The test ensures that no input causes a crash or unexpected behavior.
// The seed corpus contains a properly formatted encrypted file of the string
// "test data for corruption fuzzing".
func FuzzDecryptCorrupted(f *testing.F) {
	original := []byte("test data for corruption fuzzing")
	password := "corruption-test-password"

	encryptor, err := NewEncryptor(DefaultWorkers())
	if err != nil {
		f.Fatalf("failed to create encryptor: %v", err)
	}

	var encryptedBuf bytes.Buffer
	reader := bytes.NewReader(original)
	if err := encryptor.Encrypt(reader, &encryptedBuf, password); err != nil {
		f.Fatalf("encryption failed: %v", err)
	}

	f.Add(encryptedBuf.Bytes())

	f.Fuzz(func(t *testing.T, corrupted []byte) {
		var decryptedBuf bytes.Buffer
		corruptedReader := bytes.NewReader(corrupted)
		err := DecryptStream(corruptedReader, &decryptedBuf, password)

		if err == nil {
			_ = decryptedBuf.Bytes()
		}
	})
}

// FuzzHeaderSerialization verifies that header serialization is consistent.
//
// This fuzz test generates random magic bytes, version numbers, and chunk sizes,
// then verifies that:
//   - Serialization produces a valid header
//   - Deserialization recovers the original values
//   - The parse operation never fails on validly serialized data
//
// The test ensures that the header format is stable and that any valid input
// to Serialize produces output that ParseHeader can correctly decode.
func FuzzHeaderSerialization(f *testing.F) {
	f.Add([]byte("CRYP"), byte(2), uint32(1024*1024))

	f.Fuzz(func(t *testing.T, magicBytes []byte, version byte, chunkSize uint32) {
		var magic [4]byte
		copy(magic[:], magicBytes)

		salt := [16]byte{}
		for i := 0; i < 16 && i < len(magicBytes); i++ {
			salt[i] = magicBytes[i%len(magicBytes)]
		}

		serialized := header.Serialize(magic, version, salt, chunkSize)

		deserializedMagic, deserializedVersion, _, deserializedChunkSize, ok := header.ParseHeader(serialized)

		if !ok {
			t.Errorf("failed to parse header that was just serialized")
		}
		if deserializedMagic != magic {
			t.Errorf("magic mismatch: got %v, want %v", deserializedMagic, magic)
		}
		if deserializedVersion != version {
			t.Errorf("version mismatch: got %d, want %d", deserializedVersion, version)
		}
		if deserializedChunkSize != chunkSize {
			t.Errorf("chunkSize mismatch: got %d, want %d", deserializedChunkSize, chunkSize)
		}
	})
}
