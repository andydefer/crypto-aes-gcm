// Package crypto provides cryptographic utilities for the encryption system.
//
// This package contains shared cryptographic helper functions used across
// multiple packages to avoid code duplication.
package crypto

import (
	"encoding/binary"
	"fmt"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
)

const (
	// NonceSize is the length of the GCM nonce in bytes.
	// GCM standard requires 12 bytes for optimal performance and security.
	NonceSize = 12

	// nonceXOROffset defines where to start XORing the chunk index.
	// We preserve the first 4 bytes of the base nonce to maintain entropy
	// and XOR the chunk index into the remaining 8 bytes.
	nonceXOROffset = 4

	// nonceXORBytes is the number of bytes available for XOR operations.
	// With NonceSize=12 and offset=4, we have 8 bytes for the chunk counter.
	// This provides 2^64 unique nonces before wrap-around.
	nonceXORBytes = NonceSize - nonceXOROffset // = 8
)

// DeriveChunkNonce creates a chunk-specific nonce by XORing the base nonce
// with the chunk index bytes to prevent nonce reuse.
//
// This function is used by both encryption and decryption to generate
// deterministic, unique nonces for each chunk based on the base nonce
// and the chunk's sequential index.
//
// The XOR operation is applied to the last 8 bytes of the nonce (bytes 4-11)
// to preserve the first 4 bytes which may contain fixed header information.
// This provides 2^64 unique nonces per base nonce, which with 1MB chunks
// allows files up to approximately 16 exabytes before nonce repetition.
//
// Parameters:
//   - baseNonce: Base nonce from the file header (must be NonceSize bytes)
//   - chunkIndex: Sequential index of the current chunk (0-based)
//
// Returns:
//   - []byte: A new nonce slice of length NonceSize
//
// Example:
//
//	baseNonce := make([]byte, NonceSize)
//	nonce := DeriveChunkNonce(baseNonce, 42)
func DeriveChunkNonce(baseNonce []byte, chunkIndex uint64) []byte {
	// Use a fixed-size array to avoid allocation
	var nonce [NonceSize]byte
	copy(nonce[:], baseNonce)

	// XOR the chunk index into the last 8 bytes of the nonce
	var indexBytes [nonceXORBytes]byte
	binary.BigEndian.PutUint64(indexBytes[:], chunkIndex)

	// XOR starting at offset to preserve first 4 bytes
	for i := 0; i < nonceXORBytes; i++ {
		nonce[nonceXOROffset+i] ^= indexBytes[i]
	}

	// Return a copy to prevent accidental modification of the array
	return append([]byte(nil), nonce[:]...)
}

// DeriveChunkNonceFast is an optimized version that reuses a buffer.
// It writes the result directly into the provided destination slice.
//
// This function is faster than DeriveChunkNonce when called repeatedly
// because it avoids allocations. Use this in performance-critical code
// where you can guarantee the destination slice length.
//
// Theoretical limits:
//   - Maximum file size with 1MB chunks: 2^64 * 1MB ≈ 16 exabytes
//   - This exceeds practical filesystem limits on all current systems
//
// Parameters:
//   - dest: Destination slice (must be length >= NonceSize)
//   - baseNonce: Base nonce from the file header
//   - chunkIndex: Sequential index of the current chunk
//
// Returns:
//   - error: nil on success, or an error if dest is too short
func DeriveChunkNonceFast(dest []byte, baseNonce []byte, chunkIndex uint64) error {
	if len(dest) < NonceSize {
		return fmt.Errorf(lang.T(lang.ErrDestSliceTooShort), NonceSize, len(dest))
	}
	copy(dest, baseNonce)

	// XOR chunk index into nonce
	var indexBytes [nonceXORBytes]byte
	binary.BigEndian.PutUint64(indexBytes[:], chunkIndex)

	for i := 0; i < nonceXORBytes; i++ {
		dest[nonceXOROffset+i] ^= indexBytes[i]
	}
	return nil
}
