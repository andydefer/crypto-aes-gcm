// Package crypto provides internal cryptographic utilities for the AES-GCM encryption tool.
//
// This package implements low-level cryptographic operations including:
//   - Argon2id key derivation for secure password-based encryption
//   - HMAC-SHA256 for message authentication
//   - Header serialization for the encrypted file format
package crypto

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// TestDeriveKey verifies that key derivation produces consistent outputs
// with identical inputs.
func TestDeriveKey(t *testing.T) {
	passphrase := "test-passphrase-123"
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultArgon2Params()

	// Act: derive key
	key := DeriveKey(passphrase, salt, params)

	// Assert: key length matches expected
	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDeriveKeyDeterminism verifies that the same inputs produce the same output.
func TestDeriveKeyDeterminism(t *testing.T) {
	passphrase := "deterministic-test-passphrase"
	salt := []byte("fixed-salt-for-testing-1234567890")
	params := DefaultArgon2Params()

	// Act: derive key twice
	key1 := DeriveKey(passphrase, salt, params)
	key2 := DeriveKey(passphrase, salt, params)

	// Assert: both derivations produce identical results
	if !bytes.Equal(key1, key2) {
		t.Error("key derivation is not deterministic with same inputs")
	}
}

// TestDeriveKeySaltUniqueness verifies that different salts produce different keys.
func TestDeriveKeySaltUniqueness(t *testing.T) {
	passphrase := "test-passphrase"
	salt1 := []byte("salt-1-for-testing-purpose-123456")
	salt2 := []byte("salt-2-for-testing-purpose-123456")
	params := DefaultArgon2Params()

	// Act: derive keys with different salts
	key1 := DeriveKey(passphrase, salt1, params)
	key2 := DeriveKey(passphrase, salt2, params)

	// Assert: keys should be different
	if bytes.Equal(key1, key2) {
		t.Error("different salts produced identical keys")
	}
}

// TestDeriveKeyPassphraseUniqueness verifies that different passphrases produce different keys.
func TestDeriveKeyPassphraseUniqueness(t *testing.T) {
	passphrase1 := "first-passphrase"
	passphrase2 := "second-passphrase"
	salt := []byte("common-salt-for-testing-1234567890")
	params := DefaultArgon2Params()

	// Act: derive keys with different passphrases
	key1 := DeriveKey(passphrase1, salt, params)
	key2 := DeriveKey(passphrase2, salt, params)

	// Assert: keys should be different
	if bytes.Equal(key1, key2) {
		t.Error("different passphrases produced identical keys")
	}
}

// TestDeriveKeyEmptyPassphrase verifies that empty passphrases are handled gracefully.
func TestDeriveKeyEmptyPassphrase(t *testing.T) {
	passphrase := ""
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultArgon2Params()

	// Act: derive key with empty passphrase
	key := DeriveKey(passphrase, salt, params)

	// Assert: key should be generated and have correct length
	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDefaultArgon2Params verifies that default parameters return expected values.
func TestDefaultArgon2Params(t *testing.T) {
	// Act: get default parameters
	params := DefaultArgon2Params()

	// Assert: all fields have expected values
	if params.Time != 4 {
		t.Errorf("expected Time=4, got %d", params.Time)
	}
	if params.Memory != 64*1024 {
		t.Errorf("expected Memory=%d, got %d", 64*1024, params.Memory)
	}
	if params.Threads != 4 {
		t.Errorf("expected Threads=4, got %d", params.Threads)
	}
	if params.KeyLen != 32 {
		t.Errorf("expected KeyLen=32, got %d", params.KeyLen)
	}
}

// TestComputeHMAC verifies that HMAC computation produces consistent outputs.
func TestComputeHMAC(t *testing.T) {
	key := []byte("test-key-1234567890123456")
	data := []byte("test data to authenticate")

	// Act: compute HMAC
	hmac1 := ComputeHMAC(key, data)
	hmac2 := ComputeHMAC(key, data)

	// Assert: same inputs produce same HMAC
	if !bytes.Equal(hmac1, hmac2) {
		t.Error("ComputeHMAC not deterministic with same inputs")
	}

	// Assert: HMAC length is 32 bytes (SHA-256)
	if len(hmac1) != 32 {
		t.Errorf("expected HMAC length 32, got %d", len(hmac1))
	}
}

// TestComputeHMACKeySensitivity verifies that different keys produce different HMACs.
func TestComputeHMACKeySensitivity(t *testing.T) {
	key1 := []byte("key-1-for-hmac-testing-12345")
	key2 := []byte("key-2-for-hmac-testing-12345")
	data := []byte("same data for both keys")

	// Act: compute HMACs with different keys
	hmac1 := ComputeHMAC(key1, data)
	hmac2 := ComputeHMAC(key2, data)

	// Assert: HMACs should be different
	if bytes.Equal(hmac1, hmac2) {
		t.Error("different keys produced identical HMAC")
	}
}

// TestComputeHMACDataSensitivity verifies that different data produces different HMACs.
func TestComputeHMACDataSensitivity(t *testing.T) {
	key := []byte("fixed-key-for-testing-12345678")
	data1 := []byte("first message to authenticate")
	data2 := []byte("second message to authenticate")

	// Act: compute HMACs with different data
	hmac1 := ComputeHMAC(key, data1)
	hmac2 := ComputeHMAC(key, data2)

	// Assert: HMACs should be different
	if bytes.Equal(hmac1, hmac2) {
		t.Error("different data produced identical HMAC")
	}
}

// TestVerifyHMAC verifies that HMAC verification works correctly.
func TestVerifyHMAC(t *testing.T) {
	key := []byte("verification-test-key-12345678")
	data := []byte("data to authenticate")

	// Arrange: compute correct HMAC
	correctHMAC := ComputeHMAC(key, data)

	// Act & Assert: verification with correct HMAC should succeed
	if !VerifyHMAC(key, data, correctHMAC) {
		t.Error("VerifyHMAC failed with correct HMAC")
	}

	// Act & Assert: verification with incorrect HMAC should fail
	incorrectHMAC := make([]byte, 32)
	_, _ = rand.Read(incorrectHMAC)
	if VerifyHMAC(key, data, incorrectHMAC) {
		t.Error("VerifyHMAC succeeded with incorrect HMAC")
	}

	// Act & Assert: verification with wrong key should fail
	wrongKey := []byte("wrong-key-for-verification-1234567")
	if VerifyHMAC(wrongKey, data, correctHMAC) {
		t.Error("VerifyHMAC succeeded with wrong key")
	}
}

// TestVerifyHMACTamperedData verifies that tampered data fails verification.
func TestVerifyHMACTamperedData(t *testing.T) {
	key := []byte("tamper-test-key-123456789012")
	originalData := []byte("original message to authenticate")
	tamperedData := []byte("tampered message to authenticate")

	// Arrange: compute HMAC for original data
	originalHMAC := ComputeHMAC(key, originalData)

	// Act & Assert: verification with tampered data should fail
	if VerifyHMAC(key, tamperedData, originalHMAC) {
		t.Error("VerifyHMAC succeeded with tampered data")
	}
}

// TestHeaderToBytes verifies header serialization format and length.
func TestHeaderToBytes(t *testing.T) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	chunkSize := uint32(1024 * 1024)

	// Act: serialize header
	serialized := HeaderToBytes(magic, version, salt, chunkSize)

	// Assert: correct total length (4+1+16+4 = 25 bytes)
	expectedLen := 25
	if len(serialized) != expectedLen {
		t.Errorf("expected header length %d, got %d", expectedLen, len(serialized))
	}

	// Assert: magic bytes are correctly encoded (big-endian)
	// Magic "CRYP" should be encoded as 0x43525950
	if serialized[0] != 'C' || serialized[1] != 'R' || serialized[2] != 'Y' || serialized[3] != 'P' {
		t.Errorf("magic bytes incorrectly serialized: got %v", serialized[0:4])
	}

	// Assert: version byte is at correct position
	if serialized[4] != version {
		t.Errorf("expected version at position 4, got %d", serialized[4])
	}

	// Assert: salt is at correct position
	for i := 0; i < 16; i++ {
		if serialized[5+i] != salt[i] {
			t.Errorf("salt byte %d mismatch: expected %d, got %d", i, salt[i], serialized[5+i])
		}
	}
}

// TestHeaderToBytesDeterminism verifies that header serialization is deterministic.
func TestHeaderToBytesDeterminism(t *testing.T) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	chunkSize := uint32(1024 * 1024)

	// Act: serialize header twice
	serialized1 := HeaderToBytes(magic, version, salt, chunkSize)
	serialized2 := HeaderToBytes(magic, version, salt, chunkSize)

	// Assert: both serializations are identical
	if !bytes.Equal(serialized1, serialized2) {
		t.Error("HeaderToBytes is not deterministic")
	}
}

// TestHeaderToBytesDifferentInputs verifies that different inputs produce different outputs.
func TestHeaderToBytesDifferentInputs(t *testing.T) {
	baseMagic := [4]byte{'C', 'R', 'Y', 'P'}
	baseVersion := byte(2)
	baseSalt := [16]byte{}
	baseChunkSize := uint32(1024 * 1024)

	// Test different magic
	magic2 := [4]byte{'T', 'E', 'S', 'T'}
	serialized1 := HeaderToBytes(baseMagic, baseVersion, baseSalt, baseChunkSize)
	serialized2 := HeaderToBytes(magic2, baseVersion, baseSalt, baseChunkSize)

	if bytes.Equal(serialized1, serialized2) {
		t.Error("different magic produced same serialization")
	}

	// Test different version
	version2 := byte(3)
	serialized1 = HeaderToBytes(baseMagic, baseVersion, baseSalt, baseChunkSize)
	serialized2 = HeaderToBytes(baseMagic, version2, baseSalt, baseChunkSize)

	if bytes.Equal(serialized1, serialized2) {
		t.Error("different version produced same serialization")
	}

	// Test different chunk size
	chunkSize2 := uint32(2 * 1024 * 1024)
	serialized1 = HeaderToBytes(baseMagic, baseVersion, baseSalt, baseChunkSize)
	serialized2 = HeaderToBytes(baseMagic, baseVersion, baseSalt, chunkSize2)

	if bytes.Equal(serialized1, serialized2) {
		t.Error("different chunk size produced same serialization")
	}
}

// BenchmarkDeriveKey measures key derivation performance with default parameters.
func BenchmarkDeriveKey(b *testing.B) {
	passphrase := "benchmark-passphrase"
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultArgon2Params()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DeriveKey(passphrase, salt, params)
	}
}

// BenchmarkComputeHMAC measures HMAC computation performance.
func BenchmarkComputeHMAC(b *testing.B) {
	key := make([]byte, 32)
	data := make([]byte, 1024)
	_, _ = rand.Read(key)
	_, _ = rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComputeHMAC(key, data)
	}
}

// BenchmarkHeaderToBytes measures header serialization performance.
func BenchmarkHeaderToBytes(b *testing.B) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{}
	chunkSize := uint32(1024 * 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = HeaderToBytes(magic, version, salt, chunkSize)
	}
}
