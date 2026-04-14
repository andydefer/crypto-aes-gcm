// Package crypto provides cryptographic utilities for secure data processing.
// This test file verifies the nonce derivation functionality used for chunk-based
// encryption operations.
package crypto

import (
	"bytes"
	"testing"
)

// TestDeriveChunkNonce verifies that DeriveChunkNonce produces correct nonces
// for different chunk indices.
//
// It tests three properties:
// 1. Nonce length matches NonceSize constant
// 2. Different chunk indices produce different nonces
// 3. Same chunk index always produces identical nonce (deterministic)
func TestDeriveChunkNonce(t *testing.T) {
	baseNonce := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}

	nonce1 := DeriveChunkNonce(baseNonce, 0)
	nonce2 := DeriveChunkNonce(baseNonce, 1)

	if len(nonce1) != NonceSize {
		t.Errorf("expected nonce length %d, got %d", NonceSize, len(nonce1))
	}

	if bytes.Equal(nonce1, nonce2) {
		t.Error("different chunk indices should produce different nonces")
	}

	nonce1Again := DeriveChunkNonce(baseNonce, 0)
	if !bytes.Equal(nonce1, nonce1Again) {
		t.Error("same chunk index should produce same nonce")
	}
}

// TestDeriveChunkNonce_DifferentBaseNonce verifies that different base nonces
// produce different derived nonces even for the same chunk index.
func TestDeriveChunkNonce_DifferentBaseNonce(t *testing.T) {
	baseNonce1 := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}
	baseNonce2 := []byte{0xFF, 0xFE, 0xFD, 0xFC, 0xFB, 0xFA, 0xF9, 0xF8, 0xF7, 0xF6, 0xF5, 0xF4}

	nonce1 := DeriveChunkNonce(baseNonce1, 42)
	nonce2 := DeriveChunkNonce(baseNonce2, 42)

	if bytes.Equal(nonce1, nonce2) {
		t.Error("different base nonces should produce different chunk nonces")
	}
}

// TestDeriveChunkNonce_LargeIndex verifies that DeriveChunkNonce handles
// maximum uint64 chunk index values without panicking or truncating incorrectly.
func TestDeriveChunkNonce_LargeIndex(t *testing.T) {
	baseNonce := make([]byte, NonceSize)
	largeIndex := uint64(0xFFFFFFFFFFFFFFFF)

	nonce := DeriveChunkNonce(baseNonce, largeIndex)

	if len(nonce) != NonceSize {
		t.Errorf("expected nonce length %d, got %d", NonceSize, len(nonce))
	}

	// Should not panic
	_ = DeriveChunkNonce(baseNonce, largeIndex+1)
}

// TestDeriveChunkNonceFast verifies that the fast version produces the same
// results as the standard version.
func TestDeriveChunkNonceFast(t *testing.T) {
	baseNonce := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}
	chunkIndex := uint64(12345)

	expected := DeriveChunkNonce(baseNonce, chunkIndex)

	dest := make([]byte, NonceSize)
	err := DeriveChunkNonceFast(dest, baseNonce, chunkIndex)
	if err != nil {
		t.Fatalf("DeriveChunkNonceFast failed: %v", err)
	}

	if !bytes.Equal(expected, dest) {
		t.Errorf("DeriveChunkNonceFast produced different result:\nExpected: %v\nGot: %v", expected, dest)
	}
}

// TestDeriveChunkNonceFast_ShortDest verifies error handling for short destination.
func TestDeriveChunkNonceFast_ShortDest(t *testing.T) {
	baseNonce := make([]byte, NonceSize)
	dest := make([]byte, NonceSize-1) // Trop court

	err := DeriveChunkNonceFast(dest, baseNonce, 0)
	if err == nil {
		t.Error("Expected error for short destination, got nil")
	}
}

// BenchmarkDeriveChunkNonce compares performance of standard vs fast version.
func BenchmarkDeriveChunkNonce(b *testing.B) {
	baseNonce := make([]byte, NonceSize)

	b.Run("Standard", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = DeriveChunkNonce(baseNonce, uint64(i))
		}
	})

	b.Run("Fast", func(b *testing.B) {
		dest := make([]byte, NonceSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = DeriveChunkNonceFast(dest, baseNonce, uint64(i))
		}
	})
}
