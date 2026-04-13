// Package argon2 provides Argon2id key derivation for password-based encryption.
package argon2

import "golang.org/x/crypto/argon2"

// Params defines the configuration for Argon2id key derivation.
type Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

// DefaultParams returns secure default parameters for Argon2id.
// These provide approximately 100ms derivation time on modern hardware.
func DefaultParams() Params {
	return Params{
		Time:    4,
		Memory:  64 * 1024, // 64 MB
		Threads: 4,
		KeyLen:  32, // 256 bits for AES-256
	}
}

// DeriveKey derives a cryptographic key from a passphrase and salt using Argon2id.
// The salt must be cryptographically random and unique for each encryption operation.
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
