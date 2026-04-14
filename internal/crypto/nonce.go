// Package crypto provides cryptographic utilities for the encryption system.
//
// This package contains shared cryptographic helper functions used across
// multiple packages to avoid code duplication.
package crypto

import "encoding/binary"

const (
	// NonceSize is the length of the GCM nonce in bytes.
	NonceSize = 12
)

// DeriveChunkNonce creates a chunk-specific nonce by XORing the base nonce
// with the chunk index bytes to prevent nonce reuse.
//
// This function is used by both encryption and decryption to generate
// deterministic, unique nonces for each chunk based on the base nonce
// and the chunk's sequential index.
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
	nonce := make([]byte, NonceSize)
	copy(nonce, baseNonce)

	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, chunkIndex)

	for i := 0; i < 8 && i < NonceSize-4; i++ {
		nonce[4+i] ^= indexBytes[i]
	}

	return nonce
}
