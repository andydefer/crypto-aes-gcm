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
	"bytes"
	"crypto/rand"
	"testing"
)

// TestDeriveKey verifies that key derivation produces a key of the expected length.
func TestDeriveKey(t *testing.T) {
	passphrase := "test-passphrase-123"
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultParams()

	key := DeriveKey(passphrase, salt, params)

	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDeriveKeyDeterminism verifies that the same inputs produce identical keys.
func TestDeriveKeyDeterminism(t *testing.T) {
	passphrase := "deterministic-test-passphrase"
	salt := []byte("fixed-salt-for-testing-1234567890")
	params := DefaultParams()

	key1 := DeriveKey(passphrase, salt, params)
	key2 := DeriveKey(passphrase, salt, params)

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

	key1 := DeriveKey(passphrase, salt1, params)
	key2 := DeriveKey(passphrase, salt2, params)

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

	key1 := DeriveKey(passphrase1, salt, params)
	key2 := DeriveKey(passphrase2, salt, params)

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

	key1 := DeriveKey(passphrase, salt, params1)
	key2 := DeriveKey(passphrase, salt, params2)

	if bytes.Equal(key1, key2) {
		t.Error("different parameters produced identical keys")
	}
}

// TestDeriveKeyEmptyPassphrase verifies that empty passphrases are handled without panic.
func TestDeriveKeyEmptyPassphrase(t *testing.T) {
	passphrase := ""
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := DefaultParams()

	key := DeriveKey(passphrase, salt, params)

	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDeriveKeyEmptySalt verifies that empty salts are handled (though not recommended).
func TestDeriveKeyEmptySalt(t *testing.T) {
	passphrase := "test-passphrase"
	salt := []byte{}
	params := DefaultParams()

	key := DeriveKey(passphrase, salt, params)

	if len(key) != int(params.KeyLen) {
		t.Errorf("expected key length %d, got %d", params.KeyLen, len(key))
	}
}

// TestDefaultParams verifies that default parameters return expected values.
func TestDefaultParams(t *testing.T) {
	params := DefaultParams()

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

// TestDeriveKeyKeyLength verifies that the key length parameter is respected.
func TestDeriveKeyKeyLength(t *testing.T) {
	passphrase := "test-passphrase"
	salt := []byte("fixed-salt-for-length-testing-1234567890")

	tests := []struct {
		name    string
		keyLen  uint32
		wantLen int
	}{
		{"16 bytes (128 bits)", 16, 16},
		{"32 bytes (256 bits)", 32, 32},
		{"64 bytes (512 bits)", 64, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := Params{Time: 4, Memory: 64 * 1024, Threads: 4, KeyLen: tt.keyLen}
			key := DeriveKey(passphrase, salt, params)

			if len(key) != tt.wantLen {
				t.Errorf("expected key length %d, got %d", tt.wantLen, len(key))
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
