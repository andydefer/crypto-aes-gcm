// Package argon2 provides Argon2id key derivation for password-based encryption.
//
// This package implements the memory-hard Argon2id KDF with configurable parameters
// for different security and performance trade-offs.
package argon2

import (
	"bytes"
	"crypto/rand"
	"testing"
)

// TestDeriveKey verifies that key derivation produces consistent outputs
// with identical inputs and different outputs with different inputs.
func TestDeriveKey(t *testing.T) {
	passphrase := "test-passphrase-123"
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultParams()

	// Act: derive key
	key := DeriveKey(passphrase, salt, params)

	// Assert: key length matches expected
	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDeriveKeyDeterminism verifies that the same inputs produce the same output.
func TestDeriveKeyDeterminism(t *testing.T) {
	passphrase := "deterministic-test-passphrase"
	salt := []byte("fixed-salt-for-testing-1234567890")
	params := DefaultParams()

	// Act: derive key twice
	key1 := DeriveKey(passphrase, salt, params)
	key2 := DeriveKey(passphrase, salt, params)

	// Assert: both derivations produce identical results
	if !bytes.Equal(key1, key2) {
		t.Error("key derivation is not deterministic with same inputs")
	}
}

// TestDeriveKeySaltUniqueness verifies that different salts produce different keys.
func TestDeriveKeySaltUniqueness(t *testing.T) {
	passphrase := "test-passphrase"
	salt1 := []byte("salt-1-for-testing-purpose-123456")
	salt2 := []byte("salt-2-for-testing-purpose-123456")
	params := DefaultParams()

	// Act: derive keys with different salts
	key1 := DeriveKey(passphrase, salt1, params)
	key2 := DeriveKey(passphrase, salt2, params)

	// Assert: keys should be different
	if bytes.Equal(key1, key2) {
		t.Error("different salts produced identical keys")
	}
}

// TestDeriveKeyPassphraseUniqueness verifies that different passphrases produce different keys.
func TestDeriveKeyPassphraseUniqueness(t *testing.T) {
	passphrase1 := "first-passphrase"
	passphrase2 := "second-passphrase"
	salt := []byte("common-salt-for-testing-1234567890")
	params := DefaultParams()

	// Act: derive keys with different passphrases
	key1 := DeriveKey(passphrase1, salt, params)
	key2 := DeriveKey(passphrase2, salt, params)

	// Assert: keys should be different
	if bytes.Equal(key1, key2) {
		t.Error("different passphrases produced identical keys")
	}
}

// TestDeriveKeyDifferentParams verifies that different parameters produce different keys.
func TestDeriveKeyDifferentParams(t *testing.T) {
	passphrase := "test-passphrase"
	salt := []byte("fixed-salt-for-param-testing-12345678")
	params1 := DefaultParams()
	params2 := Params{
		Time:    params1.Time + 1,
		Memory:  params1.Memory,
		Threads: params1.Threads,
		KeyLen:  params1.KeyLen,
	}

	// Act: derive keys with different parameters
	key1 := DeriveKey(passphrase, salt, params1)
	key2 := DeriveKey(passphrase, salt, params2)

	// Assert: keys should be different
	if bytes.Equal(key1, key2) {
		t.Error("different parameters produced identical keys")
	}
}

// TestDeriveKeyEmptyPassphrase verifies that empty passphrases are handled gracefully.
func TestDeriveKeyEmptyPassphrase(t *testing.T) {
	passphrase := ""
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultParams()

	// Act: derive key with empty passphrase
	key := DeriveKey(passphrase, salt, params)

	// Assert: key should be generated (not panic) and have correct length
	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDeriveKeyEmptySalt verifies that empty salts are handled (though not recommended).
func TestDeriveKeyEmptySalt(t *testing.T) {
	passphrase := "test-passphrase"
	salt := []byte{}
	params := DefaultParams()

	// Act: derive key with empty salt
	key := DeriveKey(passphrase, salt, params)

	// Assert: key should be generated (not panic) and have correct length
	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDefaultParams verifies that default parameters return expected values.
func TestDefaultParams(t *testing.T) {
	// Act: get default parameters
	params := DefaultParams()

	// Assert: all fields have expected values
	if params.Time != 4 {
		t.Errorf("expected Time=4, got %d", params.Time)
	}
	if params.Memory != 64*1024 {
		t.Errorf("expected Memory=%d, got %d", 64*1024, params.Memory)
	}
	if params.Threads != 4 {
		t.Errorf("expected Threads=4, got %d", params.Threads)
	}
	if params.KeyLen != 32 {
		t.Errorf("expected KeyLen=32, got %d", params.KeyLen)
	}
}

// TestDeriveKeyKeyLength verifies that key length parameter is respected.
func TestDeriveKeyKeyLength(t *testing.T) {
	passphrase := "test-passphrase"
	salt := []byte("fixed-salt-for-length-testing-1234567890")
	testCases := []struct {
		name    string
		keyLen  uint32
		params  Params
		wantLen int
	}{
		{
			name:    "16 bytes (128 bits)",
			keyLen:  16,
			params:  Params{Time: 4, Memory: 64 * 1024, Threads: 4, KeyLen: 16},
			wantLen: 16,
		},
		{
			name:    "32 bytes (256 bits)",
			keyLen:  32,
			params:  Params{Time: 4, Memory: 64 * 1024, Threads: 4, KeyLen: 32},
			wantLen: 32,
		},
		{
			name:    "64 bytes (512 bits)",
			keyLen:  64,
			params:  Params{Time: 4, Memory: 64 * 1024, Threads: 4, KeyLen: 64},
			wantLen: 64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act: derive key with specified length
			key := DeriveKey(passphrase, salt, tc.params)

			// Assert: key length matches expected
			if len(key) != tc.wantLen {
				t.Errorf("expected key length %d, got %d", tc.wantLen, len(key))
			}
		})
	}
}

// BenchmarkDeriveKey measures key derivation performance with default parameters.
func BenchmarkDeriveKey(b *testing.B) {
	passphrase := "benchmark-passphrase"
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultParams()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DeriveKey(passphrase, salt, params)
	}
}

// BenchmarkDeriveKeyWithParams measures performance with different parameter sets.
func BenchmarkDeriveKeyWithParams(b *testing.B) {
	passphrase := "benchmark-passphrase"
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)

	benchmarks := []struct {
		name   string
		params Params
	}{
		{"default", DefaultParams()},
		{"time_2", Params{Time: 2, Memory: 64 * 1024, Threads: 4, KeyLen: 32}},
		{"time_8", Params{Time: 8, Memory: 64 * 1024, Threads: 4, KeyLen: 32}},
		{"memory_32MB", Params{Time: 4, Memory: 32 * 1024, Threads: 4, KeyLen: 32}},
		{"memory_128MB", Params{Time: 4, Memory: 128 * 1024, Threads: 4, KeyLen: 32}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = DeriveKey(passphrase, salt, bm.params)
			}
		})
	}
}
