// Package header provides file header serialization and HMAC authentication utilities.
//
// This package handles the deterministic binary serialization of encryption file headers
// and provides HMAC-SHA256 functions for integrity verification. All operations are
// designed to be cross-platform compatible and resistant to timing attacks.
package header

import (
	"bytes"
	"testing"
)

// TestSerialize verifies that header serialization produces the correct byte length.
func TestSerialize(t *testing.T) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	chunkSize := uint32(1024 * 1024)

	serialized := Serialize(magic, version, salt, chunkSize)

	if len(serialized) != HeaderSize {
		t.Errorf("expected length %d, got %d", HeaderSize, len(serialized))
	}
}

// TestSerializeDeterminism verifies that serialization produces identical output
// for identical inputs.
func TestSerializeDeterminism(t *testing.T) {
	magic := [4]byte{'C', 'R', 'Y', 'P'}
	version := byte(2)
	salt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	chunkSize := uint32(1024 * 1024)

	serialized1 := Serialize(magic, version, salt, chunkSize)
	serialized2 := Serialize(magic, version, salt, chunkSize)

	if !bytes.Equal(serialized1, serialized2) {
		t.Error("Serialize is not deterministic")
	}
}

// TestComputeHMAC verifies that HMAC computation is deterministic and produces
// the correct output length (32 bytes for SHA-256).
func TestComputeHMAC(t *testing.T) {
	key := []byte("test-key-1234567890123456")
	data := []byte("test data")

	hmac1 := ComputeHMAC(key, data)
	hmac2 := ComputeHMAC(key, data)

	if !bytes.Equal(hmac1, hmac2) {
		t.Error("ComputeHMAC not deterministic")
	}
	if len(hmac1) != 32 {
		t.Errorf("expected HMAC length 32, got %d", len(hmac1))
	}
}

// TestVerifyHMAC verifies that HMAC verification correctly validates
// authentic tags and rejects invalid ones.
func TestVerifyHMAC(t *testing.T) {
	key := []byte("test-key")
	data := []byte("test data")
	correctHMAC := ComputeHMAC(key, data)

	if !VerifyHMAC(key, data, correctHMAC) {
		t.Error("VerifyHMAC failed with correct HMAC")
	}

	wrongKey := []byte("wrong-key")
	if VerifyHMAC(wrongKey, data, correctHMAC) {
		t.Error("VerifyHMAC succeeded with wrong key")
	}
}

// TestDeserializeMagic verifies magic byte extraction from serialized headers.
func TestDeserializeMagic(t *testing.T) {
	expected := [4]byte{'C', 'R', 'Y', 'P'}
	serialized := Serialize(expected, 2, [16]byte{}, 1024)

	magic := DeserializeMagic(serialized)
	if magic != expected {
		t.Errorf("expected %v, got %v", expected, magic)
	}

	shortData := []byte{1, 2, 3}
	magic = DeserializeMagic(shortData)
	if magic != [4]byte{} {
		t.Error("short input should return zero magic")
	}
}

// TestDeserializeVersion verifies version byte extraction from serialized headers.
func TestDeserializeVersion(t *testing.T) {
	expected := byte(2)
	serialized := Serialize([4]byte{}, expected, [16]byte{}, 1024)

	version := DeserializeVersion(serialized)
	if version != expected {
		t.Errorf("expected %d, got %d", expected, version)
	}

	shortData := []byte{1, 2, 3, 4}
	version = DeserializeVersion(shortData)
	if version != 0 {
		t.Error("short input should return 0")
	}
}

// TestDeserializeSalt verifies salt extraction from serialized headers.
func TestDeserializeSalt(t *testing.T) {
	expected := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	serialized := Serialize([4]byte{}, 2, expected, 1024)

	salt := DeserializeSalt(serialized)
	if salt != expected {
		t.Errorf("expected %v, got %v", expected, salt)
	}
}

// TestDeserializeChunkSize verifies chunk size extraction from serialized headers.
func TestDeserializeChunkSize(t *testing.T) {
	expected := uint32(1024 * 1024)
	serialized := Serialize([4]byte{}, 2, [16]byte{}, expected)

	chunkSize := DeserializeChunkSize(serialized)
	if chunkSize != expected {
		t.Errorf("expected %d, got %d", expected, chunkSize)
	}
}

// TestParseHeader verifies complete header deserialization into all components.
func TestParseHeader(t *testing.T) {
	expectedMagic := [4]byte{'C', 'R', 'Y', 'P'}
	expectedVersion := byte(2)
	expectedSalt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	expectedChunkSize := uint32(1024 * 1024)

	serialized := Serialize(expectedMagic, expectedVersion, expectedSalt, expectedChunkSize)
	magic, version, salt, chunkSize, ok := ParseHeader(serialized)

	if !ok {
		t.Error("ParseHeader returned false")
	}
	if magic != expectedMagic {
		t.Errorf("magic mismatch")
	}
	if version != expectedVersion {
		t.Errorf("version mismatch")
	}
	if salt != expectedSalt {
		t.Errorf("salt mismatch")
	}
	if chunkSize != expectedChunkSize {
		t.Errorf("chunkSize mismatch")
	}
}

// TestValidateHeader verifies that header validation correctly identifies
// valid headers and rejects invalid ones.
func TestValidateHeader(t *testing.T) {
	validMagic := [4]byte{'C', 'R', 'Y', 'P'}
	validVersion := byte(2)
	validSalt := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	validChunkSize := uint32(1024 * 1024)

	if !ValidateHeader(validMagic, validVersion, validSalt, validChunkSize) {
		t.Error("valid header should pass")
	}

	if ValidateHeader([4]byte{}, validVersion, validSalt, validChunkSize) {
		t.Error("invalid magic should fail")
	}
}
