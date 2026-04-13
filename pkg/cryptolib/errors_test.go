// ==== ./pkg/cryptolib/errors_test.go ===

package cryptolib

import (
	"errors"
	"testing"
)

func TestErrorTypes(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{"ErrInvalidMagic", ErrInvalidMagic},
		{"ErrUnsupportedVersion", ErrUnsupportedVersion},
		{"ErrHeaderAuthFailed", ErrHeaderAuthFailed},
		{"ErrGlobalHMACFailed", ErrGlobalHMACFailed},
		{"ErrDecryptionFailed", ErrDecryptionFailed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.err == nil {
				t.Errorf("expected non-nil error for %s", tc.name)
			}

			// Verify errors are distinct
			for _, other := range testCases {
				if tc.name != other.name && errors.Is(tc.err, other.err) {
					t.Errorf("error %s should not be equal to %s", tc.name, other.name)
				}
			}
		})
	}
}

func TestErrorWrapping(t *testing.T) {
	finalErr := errors.Join(ErrDecryptionFailed, errors.New("additional context"))

	if !errors.Is(finalErr, ErrDecryptionFailed) {
		t.Errorf("expected error to wrap ErrDecryptionFailed")
	}
}
