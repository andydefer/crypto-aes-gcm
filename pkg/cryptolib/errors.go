// Package cryptolib provides AES-256-GCM file encryption and decryption.
//
// This package implements secure file encryption using AES-256-GCM in counter mode
// with Argon2id key derivation and parallel streaming capabilities.
package cryptolib

import "errors"

// Sentinel errors for common failure modes during encryption/decryption operations.
var (
	// ErrInvalidMagic indicates the file header does not contain the expected magic bytes.
	//
	// This error occurs when trying to decrypt a file that wasn't created by this
	// encryption tool or when the file is corrupted at the header level.
	ErrInvalidMagic = errors.New("invalid magic bytes: file not encrypted with this tool")

	// ErrUnsupportedVersion indicates the file uses a format version that this library cannot read.
	//
	// This error occurs when trying to decrypt a file created with a newer or
	// incompatible version of the encryption format.
	ErrUnsupportedVersion = errors.New("unsupported file version")

	// ErrHeaderAuthFailed indicates the header HMAC verification failed.
	//
	// This error typically occurs due to:
	//   - Incorrect passphrase provided for decryption
	//   - File corruption affecting the header or HMAC
	//   - Tampering with the encrypted file
	ErrHeaderAuthFailed = errors.New("header authentication failed: wrong passphrase or corrupted file")

	// ErrDecryptionFailed indicates a chunk could not be decrypted.
	//
	// This error occurs when GCM authentication fails for a data chunk,
	// usually due to:
	//   - Corrupted ciphertext data
	//   - Incorrect encryption key
	//   - File tampering or truncation
	ErrDecryptionFailed = errors.New("decryption failed: corrupted data or wrong key")
)
