// Package header provides file header serialization and HMAC authentication utilities.
//
// This package handles the deterministic binary serialization of encryption file headers
// and provides HMAC-SHA256 functions for integrity verification. All operations are
// designed to be cross-platform compatible and resistant to timing attacks.
package header

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
)

const (
	// HeaderSize returns the size in bytes of a serialized file header.
	HeaderSize = 4 + 1 + 16 + 4 // magic(4) + version(1) + salt(16) + chunkSize(4)

	// MinSupportedVersion is the earliest protocol version supported.
	MinSupportedVersion = 1
	// MaxSupportedVersion is the latest protocol version supported.
	MaxSupportedVersion = 2

	// MinChunkSize is the smallest allowed encrypted chunk size (1KB).
	MinChunkSize = 1024
	// MaxChunkSize is the largest allowed encrypted chunk size (1GB).
	MaxChunkSize = 1024 * 1024 * 1024
)

// Serialize converts file header components into a fixed-size byte slice.
//
// The byte layout is:
//   - Bytes 0-3:   Magic bytes (4 bytes, big-endian uint32)
//   - Byte 4:      Version (1 byte)
//   - Bytes 5-20:  Salt (16 bytes)
//   - Bytes 21-24: Chunk size (4 bytes, big-endian uint32)
//
// This layout is deterministic and consistent across all platforms.
//
// Parameters:
//   - magic: 4-byte identifier for file type validation
//   - version: protocol version number
//   - salt: 16-byte cryptographic salt for key derivation
//   - chunkSize: size of each encrypted chunk in bytes
//
// Returns:
//   - A 25-byte slice containing the serialized header
func Serialize(magic [4]byte, version byte, salt [16]byte, chunkSize uint32) []byte {
	buf := make([]byte, HeaderSize)

	magicValue := uint32(magic[0])<<24 |
		uint32(magic[1])<<16 |
		uint32(magic[2])<<8 |
		uint32(magic[3])
	binary.BigEndian.PutUint32(buf[0:4], magicValue)

	buf[4] = version
	copy(buf[5:21], salt[:])
	binary.BigEndian.PutUint32(buf[21:25], chunkSize)

	return buf
}

// ComputeHMAC calculates an HMAC-SHA256 digest for data integrity verification.
//
// HMAC (Hash-based Message Authentication Code) provides both integrity and
// authenticity verification using a secret key.
//
// Parameters:
//   - key: secret key used for HMAC generation (typically derived from password)
//   - data: the message to authenticate
//
// Returns:
//   - 32-byte HMAC-SHA256 digest
func ComputeHMAC(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// VerifyHMAC checks if a stored HMAC matches the computed value for given data.
//
// This function uses constant-time comparison to prevent timing side-channel attacks.
// It is the recommended way to validate HMAC tags.
//
// Parameters:
//   - key: secret key used for HMAC generation
//   - data: the message to authenticate
//   - stored: previously computed HMAC value to verify against
//
// Returns:
//   - true if the HMAC is valid, false otherwise
func VerifyHMAC(key, data, stored []byte) bool {
	computed := ComputeHMAC(key, data)
	return hmac.Equal(computed, stored)
}

// DeserializeMagic extracts magic bytes from a serialized header.
//
// This is a convenience function for reading the magic identifier without
// parsing the entire header structure.
//
// Parameters:
//   - data: serialized header bytes (must be at least 4 bytes)
//
// Returns:
//   - 4-byte magic identifier as an array
func DeserializeMagic(data []byte) [4]byte {
	var magic [4]byte
	if len(data) < 4 {
		return magic
	}
	magicValue := binary.BigEndian.Uint32(data[0:4])
	magic[0] = byte(magicValue >> 24)
	magic[1] = byte(magicValue >> 16)
	magic[2] = byte(magicValue >> 8)
	magic[3] = byte(magicValue)
	return magic
}

// DeserializeVersion extracts version byte from a serialized header.
//
// Parameters:
//   - data: serialized header bytes (must be at least 5 bytes)
//
// Returns:
//   - version byte, or 0 if data is too short
func DeserializeVersion(data []byte) byte {
	if len(data) < 5 {
		return 0
	}
	return data[4]
}

// DeserializeSalt extracts salt bytes from a serialized header.
//
// Parameters:
//   - data: serialized header bytes (must be at least 21 bytes)
//
// Returns:
//   - 16-byte salt array, zero-initialized if data is too short
func DeserializeSalt(data []byte) [16]byte {
	var salt [16]byte
	if len(data) < 21 {
		return salt
	}
	copy(salt[:], data[5:21])
	return salt
}

// DeserializeChunkSize extracts chunk size from a serialized header.
//
// Parameters:
//   - data: serialized header bytes (must be at least 25 bytes)
//
// Returns:
//   - chunk size in bytes, or 0 if data is too short
func DeserializeChunkSize(data []byte) uint32 {
	if len(data) < HeaderSize {
		return 0
	}
	return binary.BigEndian.Uint32(data[21:25])
}

// ParseHeader deserializes a complete header from bytes into its components.
//
// This is a convenience function that combines all deserialization operations.
//
// Parameters:
//   - data: complete serialized header (must be exactly HeaderSize bytes)
//
// Returns:
//   - magic: 4-byte file identifier
//   - version: protocol version
//   - salt: 16-byte cryptographic salt
//   - chunkSize: chunk size in bytes
//   - ok: true if parsing succeeded, false if data length is invalid
func ParseHeader(data []byte) (magic [4]byte, version byte, salt [16]byte, chunkSize uint32, ok bool) {
	if len(data) < HeaderSize {
		return magic, version, salt, chunkSize, false
	}

	magic = DeserializeMagic(data)
	version = DeserializeVersion(data)
	salt = DeserializeSalt(data)
	chunkSize = DeserializeChunkSize(data)

	return magic, version, salt, chunkSize, true
}

// ValidateMagic checks if magic bytes are valid (non-zero).
func ValidateMagic(magic [4]byte) bool {
	return !(magic[0] == 0 && magic[1] == 0 && magic[2] == 0 && magic[3] == 0)
}

// ValidateVersion checks if the protocol version is supported.
func ValidateVersion(version byte) bool {
	return version >= MinSupportedVersion && version <= MaxSupportedVersion
}

// ValidateSalt checks if salt is not all zeros.
func ValidateSalt(salt [16]byte) bool {
	var zeroSalt [16]byte
	return salt != zeroSalt
}

// ValidateChunkSize checks if chunk size is within reasonable bounds.
func ValidateChunkSize(chunkSize uint32) bool {
	return chunkSize >= MinChunkSize && chunkSize <= MaxChunkSize
}

// ValidateHeader performs comprehensive validation of all header components.
//
// Parameters:
//   - magic: expected magic identifier
//   - version: protocol version
//   - salt: cryptographic salt
//   - chunkSize: chunk size in bytes
//
// Returns:
//   - true if all components pass their respective validations
func ValidateHeader(magic [4]byte, version byte, salt [16]byte, chunkSize uint32) bool {
	return ValidateMagic(magic) &&
		ValidateVersion(version) &&
		ValidateSalt(salt) &&
		ValidateChunkSize(chunkSize)
}
