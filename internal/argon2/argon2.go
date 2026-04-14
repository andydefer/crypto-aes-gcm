// Package argon2 provides Argon2id key derivation for password-based encryption.
//
// Argon2id is a memory-hard key derivation function (KDF) that provides strong
// protection against GPU-based and side-channel attacks. It is the winner of
// the Password Hashing Competition and is recommended for password-based
// encryption.
//
// This package implements the recommended Argon2id variant with configurable
// parameters for different security and performance trade-offs.
//
// Example:
//
//	params := argon2.DefaultParams()
//	salt := make([]byte, 32)
//	rand.Read(salt)
//	key := argon2.DeriveKey("myPassword", salt, params)
package argon2

import (
	"errors"
	"fmt"
	"runtime"

	"golang.org/x/crypto/argon2"
)

// Params configures the Argon2id key derivation algorithm.
//
// These parameters control the computational cost, memory usage, and parallelism
// of the key derivation process. Higher values increase security but also
// increase computation time.
//
// The parameters follow the Argon2 specification:
//
//	Time:    number of passes over memory (iteration count)
//	Memory:  memory usage in KiB (e.g., 64*1024 = 64 MiB)
//	Threads: degree of parallelism (number of independent computation lanes)
//	KeyLen:  desired output key length in bytes
//
// Recommended values:
//   - Time:    4
//   - Memory:  64*1024 (64 MiB)
//   - Threads: number of CPU cores (max 4)
//   - KeyLen:  32 (256 bits for AES-256)
type Params struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
}

const maxThreads = 32

// Validate checks if the parameters are within safe ranges.
//
// Returns an error if any parameter is outside acceptable bounds:
//   - Memory: between 8 MiB and 1 GiB
//   - Threads: at least 1 and within system capacity
//   - Time: between 1 and 100
//   - KeyLen: between 16 and 64 bytes
func (p Params) Validate() error {
	if p.Memory < 8*1024 {
		return fmt.Errorf("memory too low: %d KiB (minimum 8192 KiB)", p.Memory)
	}
	if p.Memory > 1024*1024 {
		return fmt.Errorf("memory too high: %d KiB (maximum 1,048,576 KiB)", p.Memory)
	}

	if p.Threads < 1 {
		return errors.New("threads must be at least 1")
	}

	if p.Threads > maxThreads {
		return fmt.Errorf("threads too high: %d (maximum %d)", p.Threads, maxThreads)
	}

	if int(p.Threads) > runtime.NumCPU()*2 {
		return fmt.Errorf("threads exceed system capacity: %d (max %d)", p.Threads, runtime.NumCPU()*2)
	}

	if p.Time < 1 {
		return errors.New("time must be at least 1")
	}
	if p.Time > 100 {
		return fmt.Errorf("time too high: %d (maximum 100)", p.Time)
	}

	if p.KeyLen < 16 {
		return fmt.Errorf("key length too short: %d bytes (minimum 16)", p.KeyLen)
	}
	if p.KeyLen > 64 {
		return fmt.Errorf("key length too long: %d bytes (maximum 64)", p.KeyLen)
	}

	return nil
}

// DefaultParams returns secure, production-ready default parameters for Argon2id.
//
// These parameters provide approximately 100ms derivation time on modern hardware
// and are suitable for most applications.
//
// Returns:
//   - Params with Time=4, Memory=64MiB, Threads=4, KeyLen=32
func DefaultParams() Params {
	return Params{
		Time:    4,
		Memory:  64 * 1024,
		Threads: 4,
		KeyLen:  32,
	}
}

// DeriveKey derives a cryptographic key from a passphrase and salt using Argon2id.
//
// The function applies the memory-hard Argon2id KDF to transform a user passphrase
// into a cryptographically strong key suitable for symmetric encryption.
//
// Parameters:
//   - passphrase: user-supplied password or passphrase
//   - salt: cryptographically random byte slice (minimum 16 bytes recommended)
//   - params: Argon2id configuration parameters
//
// Returns:
//   - []byte: derived key of length params.KeyLen bytes
//
// Important security considerations:
//   - Salt MUST be random and unique for each encryption operation
//   - Salt MUST be stored alongside the ciphertext for decryption
//   - Passphrase should have sufficient entropy (use strong passwords)
//   - Changing any parameter produces completely different output
//
// Example:
//
//	salt := make([]byte, 32)
//	rand.Read(salt)
//	key := argon2.DeriveKey("myPassphrase", salt, argon2.DefaultParams())
func DeriveKey(passphrase string, salt []byte, params Params) []byte {
	if err := params.Validate(); err != nil {
		params = DefaultParams()
	}

	return argon2.IDKey(
		[]byte(passphrase),
		salt,
		params.Time,
		params.Memory,
		params.Threads,
		params.KeyLen,
	)
}
