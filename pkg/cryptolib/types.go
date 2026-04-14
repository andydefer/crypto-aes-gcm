// Package cryptolib provides AES-256-GCM file encryption and decryption.
//
// This package implements secure file encryption using AES-256-GCM in counter mode
// with Argon2id key derivation and parallel streaming capabilities.
package cryptolib

import "runtime"

const (
	// Magic identifies files created by this library.
	//
	// This 4-byte magic number is written at the beginning of every encrypted file
	// to allow the decrypter to quickly verify the file format.
	Magic = "CRYP"

	// Version indicates the current file format version.
	//
	// Incrementing this version allows future format changes while maintaining
	// backward compatibility through version checking.
	Version = 2

	// SaltSize is the length of the Argon2id salt in bytes.
	//
	// The salt is randomly generated for each encryption operation and stored
	// in the file header to prevent rainbow table attacks.
	SaltSize = 16

	// NonceSize is the length of the GCM nonce in bytes.
	//
	// GCM requires a 12-byte nonce for optimal performance and security.
	// Nonces are generated randomly and combined with chunk indices.
	NonceSize = 12

	// KeySize is the length of the derived AES key in bytes (256 bits).
	//
	// AES-256 provides 256-bit encryption keys derived from user passphrases
	// via Argon2id key derivation function.
	KeySize = 32

	// DefaultChunkSize is the size of each encrypted chunk (1MB).
	//
	// Chunks of this size are processed independently, allowing parallel
	// encryption and streaming decryption with bounded memory usage.
	DefaultChunkSize = 1024 * 1024

	// DefaultWorkers is the default number of parallel encryption workers.
	//
	// This value provides a good balance between performance and resource usage
	// on most systems. Workers are automatically clamped to 2×CPU cores.
	// DefaultWorkers = 4

	// DefaultMaxPendingChunks is the default limit of out-of-order chunks buffered in memory
	// during encryption to prevent memory exhaustion attacks.
	//
	// When parallel encryption produces chunks out of order, they are buffered
	// until earlier chunks are written. This limit prevents attackers from
	// causing unbounded memory growth.
	DefaultMaxPendingChunks = 100

	// MaxMaxPendingChunks is the absolute maximum allowed value for pending chunks.
	//
	// This prevents excessive memory allocation even if a user specifies a very large value.
	MaxMaxPendingChunks = 1000
)

// FileHeader represents the header structure at the beginning of every encrypted file.
//
// The header contains:
//   - Magic bytes to identify the file format
//   - Version number for format compatibility
//   - Cryptographic salt for key derivation
//   - Chunk size used during encryption
//
// This structure is serialized in big-endian byte order and verified using HMAC
// to prevent tampering.
type FileHeader struct {
	// Magic identifies the file format (4 bytes: "CRYP")
	Magic [4]byte

	// Version indicates the file format version
	Version byte

	// Salt contains the random salt used for Argon2id key derivation
	Salt [SaltSize]byte

	// ChunkSize specifies the size of each encrypted chunk in bytes
	ChunkSize uint32
}

// EncryptorConfig holds configuration options for the Encryptor.
//
// This allows users to customize the encryption behavior including
// the maximum number of pending chunks allowed in memory.
type EncryptorConfig struct {
	// Workers is the number of parallel encryption workers.
	// Default: DefaultWorkers (4)
	Workers int

	// ChunkSize is the size of each chunk in bytes.
	// Default: DefaultChunkSize (1MB)
	ChunkSize int

	// MaxPendingChunks limits the number of out-of-order chunks buffered.
	// Default: DefaultMaxPendingChunks (100)
	// Maximum: MaxMaxPendingChunks (1000)
	MaxPendingChunks int
}

// DefaultEncryptorConfig returns the default configuration for Encryptor.
func DefaultEncryptorConfig() EncryptorConfig {
	return EncryptorConfig{
		Workers:          DefaultWorkers(),
		ChunkSize:        DefaultChunkSize,
		MaxPendingChunks: DefaultMaxPendingChunks,
	}
}

// DefaultWorkers respects Go's GOMAXPROCS setting
func DefaultWorkers() int {
	maxProcs := runtime.GOMAXPROCS(0)

	return min(max(maxProcs*3/4, 1), 8)
}
