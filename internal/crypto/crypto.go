// Package crypto provides internal cryptographic utilities for the AES-GCM encryption tool.
//
// This package implements low-level cryptographic operations including:
//   - Argon2id key derivation for secure password-based encryption
//   - HMAC-SHA256 for message authentication
//   - Header serialization for the encrypted file format
//
// The functions in this package are intended for internal use only and should not
// be exposed as public API. The package follows Go best practices for cryptography,
// using constant-time comparisons and secure random number generation where appropriate.
package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"

	"golang.org/x/crypto/argon2"
)

// Argon2Params defines the configuration parameters for the Argon2id key derivation function.
//
// These parameters control the computational cost and memory usage of key derivation:
//   - Time: Number of iterations (higher = slower, more secure)
//   - Memory: Memory usage in KiB (higher = more secure)
//   - Threads: Number of parallel threads (higher = faster on multi-core)
//   - KeyLen: Desired output key length in bytes
//
// The default values are chosen to provide strong security while maintaining
// reasonable performance on modern hardware (≈100ms derivation time).
type Argon2Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultArgon2Params returns secure default parameters for Argon2id key derivation.
//
// The defaults are:
//   - Time: 4 iterations
//   - Memory: 64 MB (65536 KiB)
//   - Threads: 4 parallel threads
//   - KeyLen: 32 bytes (256 bits) for AES-256
//
// Returns:
//   - Argon2Params: Configuration suitable for most use cases
func DefaultArgon2Params() Argon2Params {
	return Argon2Params{
		Time:    4,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
}

// DeriveKey derives a cryptographic key from a passphrase and salt using Argon2id.
//
// Argon2id is a memory-hard key derivation function that provides resistance
// against GPU-based and side-channel attacks. It is the winner of the Password
// Hashing Competition and is recommended for password-based encryption.
//
// Parameters:
//   - passphrase: User-provided password (should be sufficiently complex)
//   - salt: Random salt (16+ bytes, must be unique per encryption operation)
//   - params: Argon2id parameters controlling cost and memory usage
//
// Returns:
//   - []byte: Derived key of length params.KeyLen, suitable for AES encryption
//
// Security considerations:
//   - Salt must be cryptographically random and unique for each encryption
//   - Passphrase should have sufficient entropy (use long, complex passwords)
//   - Default parameters provide ~100ms derivation time on modern hardware
func DeriveKey(passphrase string, salt []byte, params Argon2Params) []byte {
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLen,
	)
}

// ComputeHMAC computes an HMAC-SHA256 message authentication code for the given data.
//
// HMAC (Hash-based Message Authentication Code) provides integrity and authenticity
// verification. It uses a secret key to produce a cryptographic digest that cannot
// be forged without knowledge of the key.
//
// Parameters:
//   - key: Secret key used for HMAC computation (must be kept confidential)
//   - data: Message data to authenticate
//
// Returns:
//   - []byte: 32-byte HMAC-SHA256 digest
//
// Example usage:
//
//	hmac := ComputeHMAC(encryptionKey, headerData)
//	// Store hmac alongside data for later verification
func ComputeHMAC(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// VerifyHMAC securely compares a computed HMAC with a stored HMAC value.
//
// This function uses constant-time comparison to prevent timing attacks
// that could leak information about the HMAC value. Always use this function
// instead of manual comparison when verifying HMACs.
//
// Parameters:
//   - key: Secret key used for HMAC computation
//   - data: Original message data to verify
//   - storedHMAC: Previously computed HMAC value to compare against
//
// Returns:
//   - bool: true if HMACs match (data authentic), false otherwise
//
// Security:
//   - Uses constant-time comparison to prevent timing side-channels
//   - Returns false for invalid inputs without revealing why
func VerifyHMAC(key, data, storedHMAC []byte) bool {
	computed := ComputeHMAC(key, data)
	return hmac.Equal(computed, storedHMAC)
}

// HeaderToBytes serializes a file header to bytes for HMAC calculation.
//
// The header format is fixed and deterministic, allowing the same serialization
// to be used for both header creation and verification. The byte layout ensures
// compatibility across different platforms (big-endian encoding).
//
// Byte layout (25 bytes total):
//   - Bytes 0-3: Magic bytes (4 bytes) - File format identifier
//   - Byte 4: Version (1 byte) - Format version number
//   - Bytes 5-20: Salt (16 bytes) - Argon2id salt
//   - Bytes 21-24: ChunkSize (4 bytes) - Size of each encrypted chunk
//
// Parameters:
//   - magic: 4-byte file magic identifier (typically "CRYP")
//   - version: Format version number (must match current version)
//   - salt: 16-byte random salt used for key derivation
//   - chunkSize: Size of each chunk during streaming encryption
//
// Returns:
//   - []byte: 25-byte serialized header ready for HMAC computation
func HeaderToBytes(magic [4]byte, version byte, salt [16]byte, chunkSize uint32) []byte {
	buf := make([]byte, 4+1+16+4)

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
