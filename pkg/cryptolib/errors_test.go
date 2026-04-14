// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
package cryptolib

import (
	"errors"
	"testing"
)

// TestErrorTypes verifies that all error variables are properly initialized
// and distinct from each other.
//
// This test ensures that:
//   - Each error constant is non-nil
//   - Errors are unique and not equal to each other
//   - Error comparison works correctly with errors.Is
func TestErrorTypes(t *testing.T) {
	errorConstants := []struct {
		name string
		err  error
	}{
		{"ErrInvalidMagic", ErrInvalidMagic},
		{"ErrUnsupportedVersion", ErrUnsupportedVersion},
		{"ErrHeaderAuthFailed", ErrHeaderAuthFailed},
		{"ErrDecryptionFailed", ErrDecryptionFailed},
	}

	for _, current := range errorConstants {
		t.Run(current.name, func(t *testing.T) {
			if current.err == nil {
				t.Errorf("error %s is nil", current.name)
			}

			for _, other := range errorConstants {
				if current.name != other.name && errors.Is(current.err, other.err) {
					t.Errorf("error %s should not be equal to %s", current.name, other.name)
				}
			}
		})
	}
}

// TestErrorWrapping verifies that errors can be properly wrapped and unwrapped.
//
// This test uses errors.Join (Go 1.20+) to combine errors and ensures that
// errors.Is correctly identifies wrapped errors.
func TestErrorWrapping(t *testing.T) {
	additionalErr := errors.New("additional context")
	wrappedErr := errors.Join(ErrDecryptionFailed, additionalErr)

	if !errors.Is(wrappedErr, ErrDecryptionFailed) {
		t.Error("wrapped error does not contain ErrDecryptionFailed")
	}

	if !errors.Is(wrappedErr, additionalErr) {
		t.Error("wrapped error does not contain additional context")
	}
}

// TestErrorMessageContent verifies that error messages are user-friendly
// and non-empty.
//
// Error messages should be descriptive enough for users to understand
// what went wrong without exposing internal implementation details.
func TestErrorMessageContent(t *testing.T) {
	errorConstants := []struct {
		name string
		err  error
	}{
		{"ErrInvalidMagic", ErrInvalidMagic},
		{"ErrUnsupportedVersion", ErrUnsupportedVersion},
		{"ErrHeaderAuthFailed", ErrHeaderAuthFailed},
		{"ErrDecryptionFailed", ErrDecryptionFailed},
	}

	for _, current := range errorConstants {
		t.Run(current.name, func(t *testing.T) {
			if current.err.Error() == "" {
				t.Errorf("error %s has empty message", current.name)
			}
		})
	}
}
