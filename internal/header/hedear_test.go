// Package header provides file header serialization and HMAC authentication utilities.
//
// This package handles the deterministic binary serialization of encryption file headers
// and provides HMAC-SHA256 functions for integrity verification.
package header

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// TestSerialize verifies that header serialization produces correct output format.
func TestSerialize(t *testing.T) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	chunkSize := uint32(1024 * 1024)

	// Act: serialize header
	serialized := Serialize(magic, version, salt, chunkSize)

	// Assert: correct total length
	if len(serialized) != HeaderSize {
		t.Errorf("expected header length %d, got %d", HeaderSize, len(serialized))
	}

	// Assert: magic bytes are correctly encoded
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

	// Assert: chunk size is correctly encoded (big-endian)
	decodedChunkSize := uint32(serialized[21])<<24 |
		uint32(serialized[22])<<16 |
		uint32(serialized[23])<<8 |
		uint32(serialized[24])
	if decodedChunkSize != chunkSize {
		t.Errorf("expected chunk size %d, got %d", chunkSize, decodedChunkSize)
	}
}

// TestSerializeDeterminism verifies that serialization is deterministic.
func TestSerializeDeterminism(t *testing.T) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	chunkSize := uint32(1024 * 1024)

	// Act: serialize header twice
	serialized1 := Serialize(magic, version, salt, chunkSize)
	serialized2 := Serialize(magic, version, salt, chunkSize)

	// Assert: both serializations are identical
	if !bytes.Equal(serialized1, serialized2) {
		t.Error("Serialize is not deterministic")
	}
}

// TestComputeHMAC verifies HMAC computation produces consistent outputs.
func TestComputeHMAC(t *testing.T) {
	key := []byte("test-key-1234567890123456")
	data := []byte("test data to authenticate")

	// Act: compute HMAC twice
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

// TestDeserializeMagic verifies magic byte extraction from serialized header.
func TestDeserializeMagic(t *testing.T) {
	expectedMagic := [4]byte{'C', 'R', 'Y', 'P'}
	serialized := Serialize(expectedMagic, 2, [16]byte{}, 1024)

	// Act: deserialize magic
	magic := DeserializeMagic(serialized)

	// Assert: magic matches expected
	if magic != expectedMagic {
		t.Errorf("expected magic %v, got %v", expectedMagic, magic)
	}

	// Assert: short input returns zero magic
	shortData := []byte{1, 2, 3}
	magic = DeserializeMagic(shortData)
	if magic != [4]byte{} {
		t.Error("short input should return zero magic")
	}
}

// TestDeserializeVersion verifies version extraction from serialized header.
func TestDeserializeVersion(t *testing.T) {
	expectedVersion := byte(2)
	serialized := Serialize([4]byte{}, expectedVersion, [16]byte{}, 1024)

	// Act: deserialize version
	version := DeserializeVersion(serialized)

	// Assert: version matches expected
	if version != expectedVersion {
		t.Errorf("expected version %d, got %d", expectedVersion, version)
	}

	// Assert: short input returns 0
	shortData := []byte{1, 2, 3, 4}
	version = DeserializeVersion(shortData)
	if version != 0 {
		t.Error("short input should return 0")
	}
}

// TestDeserializeSalt verifies salt extraction from serialized header.
func TestDeserializeSalt(t *testing.T) {
	expectedSalt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	serialized := Serialize([4]byte{}, 2, expectedSalt, 1024)

	// Act: deserialize salt
	salt := DeserializeSalt(serialized)

	// Assert: salt matches expected
	if salt != expectedSalt {
		t.Errorf("expected salt %v, got %v", expectedSalt, salt)
	}

	// Assert: short input returns zero salt
	shortData := []byte{1, 2, 3, 4, 5}
	salt = DeserializeSalt(shortData)
	if salt != [16]byte{} {
		t.Error("short input should return zero salt")
	}
}

// TestDeserializeChunkSize verifies chunk size extraction from serialized header.
func TestDeserializeChunkSize(t *testing.T) {
	expectedChunkSize := uint32(1024 * 1024)
	serialized := Serialize([4]byte{}, 2, [16]byte{}, expectedChunkSize)

	// Act: deserialize chunk size
	chunkSize := DeserializeChunkSize(serialized)

	// Assert: chunk size matches expected
	if chunkSize != expectedChunkSize {
		t.Errorf("expected chunk size %d, got %d", expectedChunkSize, chunkSize)
	}

	// Assert: short input returns 0
	shortData := make([]byte, 24)
	chunkSize = DeserializeChunkSize(shortData)
	if chunkSize != 0 {
		t.Error("short input should return 0")
	}
}

// TestParseHeader verifies complete header deserialization.
func TestParseHeader(t *testing.T) {
	expectedMagic := [4]byte{'C', 'R', 'Y', 'P'}
	expectedVersion := byte(2)
	expectedSalt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	expectedChunkSize := uint32(1024 * 1024)

	serialized := Serialize(expectedMagic, expectedVersion, expectedSalt, expectedChunkSize)

	// Act: parse header
	magic, version, salt, chunkSize, ok := ParseHeader(serialized)

	// Assert: all components match expected
	if !ok {
		t.Error("ParseHeader returned false for valid header")
	}
	if magic != expectedMagic {
		t.Errorf("expected magic %v, got %v", expectedMagic, magic)
	}
	if version != expectedVersion {
		t.Errorf("expected version %d, got %d", expectedVersion, version)
	}
	if salt != expectedSalt {
		t.Errorf("expected salt %v, got %v", expectedSalt, salt)
	}
	if chunkSize != expectedChunkSize {
		t.Errorf("expected chunk size %d, got %d", expectedChunkSize, chunkSize)
	}

	// Assert: invalid length returns false
	invalidData := make([]byte, 10)
	_, _, _, _, ok = ParseHeader(invalidData)
	if ok {
		t.Error("ParseHeader should return false for invalid length")
	}
}

// TestValidateMagic verifies magic byte validation.
func TestValidateMagic(t *testing.T) {
	testCases := []struct {
		name     string
		magic    [4]byte
		expected bool
	}{
		{"valid magic", [4]byte{'C', 'R', 'Y', 'P'}, true},
		{"zero magic", [4]byte{0, 0, 0, 0}, false},
		{"partial zero", [4]byte{'C', 0, 'Y', 'P'}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateMagic(tc.magic)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

// TestValidateVersion verifies protocol version validation.
func TestValidateVersion(t *testing.T) {
	testCases := []struct {
		name     string
		version  byte
		expected bool
	}{
		{"version too low", 0, false},
		{"minimum supported", MinSupportedVersion, true},
		{"supported version 1", 1, true},
		{"supported version 2", 2, true},
		{"version too high", 3, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateVersion(tc.version)
			if result != tc.expected {
				t.Errorf("version %d: expected %v, got %v", tc.version, tc.expected, result)
			}
		})
	}
}

// TestValidateSalt verifies salt validation.
func TestValidateSalt(t *testing.T) {
	validSalt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	zeroSalt := [16]byte{}

	if !ValidateSalt(validSalt) {
		t.Error("valid salt should return true")
	}
	if ValidateSalt(zeroSalt) {
		t.Error("zero salt should return false")
	}
}

// TestValidateChunkSize verifies chunk size validation.
func TestValidateChunkSize(t *testing.T) {
	testCases := []struct {
		name      string
		chunkSize uint32
		expected  bool
	}{
		{"too small", 512, false},
		{"minimum", MinChunkSize, true},
		{"valid", 1024 * 1024, true},
		{"maximum", MaxChunkSize, true},
		{"too large", MaxChunkSize + 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ValidateChunkSize(tc.chunkSize)
			if result != tc.expected {
				t.Errorf("chunk size %d: expected %v, got %v", tc.chunkSize, tc.expected, result)
			}
		})
	}
}

// TestValidateHeader verifies comprehensive header validation.
func TestValidateHeader(t *testing.T) {
	validMagic := [4]byte{'C', 'R', 'Y', 'P'}
	validVersion := byte(2)
	validSalt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	validChunkSize := uint32(1024 * 1024)

	// Assert: valid header passes
	if !ValidateHeader(validMagic, validVersion, validSalt, validChunkSize) {
		t.Error("valid header should pass validation")
	}

	// Assert: invalid magic fails
	if ValidateHeader([4]byte{0, 0, 0, 0}, validVersion, validSalt, validChunkSize) {
		t.Error("invalid magic should fail validation")
	}

	// Assert: invalid version fails
	if ValidateHeader(validMagic, 99, validSalt, validChunkSize) {
		t.Error("invalid version should fail validation")
	}

	// Assert: invalid salt fails
	if ValidateHeader(validMagic, validVersion, [16]byte{}, validChunkSize) {
		t.Error("invalid salt should fail validation")
	}

	// Assert: invalid chunk size fails
	if ValidateHeader(validMagic, validVersion, validSalt, 0) {
		t.Error("invalid chunk size should fail validation")
	}
}

// BenchmarkSerialize measures header serialization performance.
func BenchmarkSerialize(b *testing.B) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{}
	chunkSize := uint32(1024 * 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Serialize(magic, version, salt, chunkSize)
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

// BenchmarkVerifyHMAC measures HMAC verification performance.
func BenchmarkVerifyHMAC(b *testing.B) {
	key := make([]byte, 32)
	data := make([]byte, 1024)
	_, _ = rand.Read(key)
	_, _ = rand.Read(data)
	hmac := ComputeHMAC(key, data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = VerifyHMAC(key, data, hmac)
	}
}

// BenchmarkParseHeader measures complete header parsing performance.
func BenchmarkParseHeader(b *testing.B) {
	serialized := Serialize([4]byte{'C', 'R', 'Y', 'P'}, 2, [16]byte{}, 1024*1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, _, _ = ParseHeader(serialized)
	}
}
