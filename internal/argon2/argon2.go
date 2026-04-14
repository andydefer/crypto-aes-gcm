// Package argon2 provides Argon2id key derivation for password-based encryption.
//
// Argon2id is a memory-hard key derivation function (KDF) that provides strong
// protection against GPU-based and side-channel attacks. This package implements
// the recommended Argon2id variant with configurable parameters for different
// security and performance trade-offs.
package argon2

import "golang.org/x/crypto/argon2"

// Params configures the Argon2id key derivation algorithm.
//
// These parameters control the computational cost, memory usage, and parallelism
// of the key derivation process. Higher values increase security but also
// increase computation time.
//
// The parameters follow Argon2 specification:
//   - Time: number of passes over memory (iteration count)
//   - Memory: memory usage in KiB (e.g., 64*1024 = 64 MiB)
//   - Threads: degree of parallelism (number of independent computation lanes)
//   - KeyLen: desired output key length in bytes
type Params struct {
	Time    uint32 // Number of iterations (recommended: 4)
	Memory  uint32 // Memory usage in KiB (recommended: 64*1024)
	Threads uint8  // Degree of parallelism (recommended: number of CPU cores)
	KeyLen  uint32 // Output key length in bytes (recommended: 32 for AES-256)
}

// DefaultParams returns secure, production-ready default parameters for Argon2id.
//
// These parameters provide approximately 100ms derivation time on modern hardware
// and are suitable for most applications:
//   - Time: 4 iterations for balanced security/performance
//   - Memory: 64 MiB (64*1024 KiB) for strong memory-hardness
//   - Threads: 4 parallel lanes for multi-core efficiency
//   - KeyLen: 32 bytes (256 bits) for AES-256 compatibility
//
// Returns:
//   - Params with secure defaults suitable for general use cases
func DefaultParams() Params {
	return Params{
		Time:    4,
		Memory:  64 * 1024, // 64 MiB in KiB
		Threads: 4,
		KeyLen:  32, // 256 bits for AES-256
	}
}

// DeriveKey derives a cryptographic key from a passphrase and salt using Argon2id.
//
// The function applies the memory-hard Argon2id KDF to transform a user passphrase
// into a cryptographically strong key suitable for symmetric encryption.
//
// Parameters:
//   - passphrase: user-supplied password or passphrase (should be strong)
//   - salt: cryptographically random byte slice (minimum 16 bytes recommended)
//   - params: Argon2id configuration parameters (use DefaultParams for most cases)
//
// Returns:
//   - []byte: derived key of length params.KeyLen bytes
//
// Important:
//   - Salt MUST be random, unique per encryption operation
//   - Salt MUST be stored alongside the ciphertext for decryption
//   - Passphrase should have sufficient entropy (use strong passwords)
//   - Changing any parameter produces completely different output
//
// Example:
//
//	salt := make([]byte, 32)
//	rand.Read(salt)
//	key := DeriveKey("myPassphrase", salt, DefaultParams())
func DeriveKey(passphrase string, salt []byte, params Params) []byte {
	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLen,
	)
}
