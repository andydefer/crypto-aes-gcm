// Package cryptolib provides secure, streaming AES-256-GCM encryption with Argon2id key derivation.
//
// The package implements authenticated encryption with parallel chunk processing for large files,
// and HMAC-SHA256 integrity verification for both headers and content.
package cryptolib

import "errors"

// Sentinel errors for common failure modes during encryption/decryption operations.
var (
	// ErrInvalidMagic indicates the file header does not contain the expected magic bytes.
	ErrInvalidMagic = errors.New("invalid magic bytes: file not encrypted with this tool")

	// ErrUnsupportedVersion indicates the file uses a format version that this library cannot read.
	ErrUnsupportedVersion = errors.New("unsupported file version")

	// ErrHeaderAuthFailed indicates the header HMAC verification failed,
	// typically due to an incorrect passphrase or file corruption.
	ErrHeaderAuthFailed = errors.New("header authentication failed: wrong passphrase or corrupted file")

	// ErrGlobalHMACFailed indicates the global HMAC verification failed,
	// meaning the file content has been corrupted.
	ErrGlobalHMACFailed = errors.New("global HMAC verification failed: file may be corrupted")

	// ErrDecryptionFailed indicates a chunk could not be decrypted,
	// usually due to corrupted data or an incorrect key.
	ErrDecryptionFailed = errors.New("decryption failed: corrupted data or wrong key")
)
