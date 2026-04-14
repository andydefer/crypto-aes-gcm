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

// TestDeriveChunkNonce_PreservesFirstBytes verifies that the first 4 bytes
// of the base nonce remain unchanged after derivation.
func TestDeriveChunkNonce_PreservesFirstBytes(t *testing.T) {
	baseNonce := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}

	nonce := DeriveChunkNonce(baseNonce, 12345)

	// First 4 bytes should be unchanged
	for i := 0; i < nonceXOROffset; i++ {
		if nonce[i] != baseNonce[i] {
			t.Errorf("byte %d should be preserved: expected %d, got %d",
				i, baseNonce[i], nonce[i])
		}
	}
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
	dest := make([]byte, NonceSize-1) // Too short

	err := DeriveChunkNonceFast(dest, baseNonce, 0)
	if err == nil {
		t.Error("Expected error for short destination, got nil")
	}
}

// TestConstantsConsistency verifies that the constant values are consistent
// with the NonceSize constant.
func TestConstantsConsistency(t *testing.T) {
	if nonceXOROffset+nonceXORBytes != NonceSize {
		t.Errorf("offset+bytes should equal NonceSize: %d+%d=%d, want %d",
			nonceXOROffset, nonceXORBytes, nonceXOROffset+nonceXORBytes, NonceSize)
	}

	if nonceXORBytes != 8 {
		t.Errorf("nonceXORBytes should be 8, got %d", nonceXORBytes)
	}
}

// TestDeriveChunkNonceFast_SequentialUniqueness verifies that sequential
// chunk indices produce different nonces.
func TestDeriveChunkNonceFast_SequentialUniqueness(t *testing.T) {
	baseNonce := make([]byte, NonceSize)
	// Initialize base nonce with known values for predictability
	for i := range baseNonce {
		baseNonce[i] = byte(i)
	}

	dest0 := make([]byte, NonceSize)
	dest1 := make([]byte, NonceSize)
	dest2 := make([]byte, NonceSize)

	err := DeriveChunkNonceFast(dest0, baseNonce, 0)
	if err != nil {
		t.Fatal(err)
	}

	err = DeriveChunkNonceFast(dest1, baseNonce, 1)
	if err != nil {
		t.Fatal(err)
	}

	err = DeriveChunkNonceFast(dest2, baseNonce, 256)
	if err != nil {
		t.Fatal(err)
	}

	// All nonces should be different
	if bytes.Equal(dest0, dest1) {
		t.Error("chunk 0 and chunk 1 should produce different nonces")
	}
	if bytes.Equal(dest0, dest2) {
		t.Error("chunk 0 and chunk 256 should produce different nonces")
	}
	if bytes.Equal(dest1, dest2) {
		t.Error("chunk 1 and chunk 256 should produce different nonces")
	}

	// Verify that XORing the same index twice produces the same result
	dest0Again := make([]byte, NonceSize)
	err = DeriveChunkNonceFast(dest0Again, baseNonce, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(dest0, dest0Again) {
		t.Error("same chunk index should produce same nonce (deterministic)")
	}

	// Verify that the first nonceXOROffset bytes remain unchanged
	for i := 0; i < nonceXOROffset; i++ {
		if dest0[i] != baseNonce[i] {
			t.Errorf("first %d bytes should be preserved at position %d: expected %d, got %d",
				nonceXOROffset, i, baseNonce[i], dest0[i])
		}
	}
}

// TestDeriveChunkNonceFast_XORPattern verifies the XOR pattern with explicit values.
func TestDeriveChunkNonceFast_XORPattern(t *testing.T) {
	// Use a zero base nonce for predictable XOR results
	baseNonce := make([]byte, NonceSize)

	testCases := []struct {
		chunkIndex uint64
		expected   byte // Expected value at LSB position (last byte of XOR region)
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{255, 255},
		{256, 0}, // 256 % 256 = 0, but higher bytes will be set
	}

	for _, tc := range testCases {
		dest := make([]byte, NonceSize)
		err := DeriveChunkNonceFast(dest, baseNonce, tc.chunkIndex)
		if err != nil {
			t.Fatalf("failed for index %d: %v", tc.chunkIndex, err)
		}

		// The LSB of chunkIndex should appear at the last byte of XOR region
		lsbPos := nonceXOROffset + nonceXORBytes - 1
		if dest[lsbPos] != tc.expected {
			t.Errorf("chunkIndex %d: expected LSB byte %d at position %d, got %d",
				tc.chunkIndex, tc.expected, lsbPos, dest[lsbPos])
		}
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

// BenchmarkDeriveChunkNonceFast_NoAlloc verifies the fast version doesn't allocate.
func BenchmarkDeriveChunkNonceFast_NoAlloc(b *testing.B) {
	baseNonce := make([]byte, NonceSize)
	dest := make([]byte, NonceSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DeriveChunkNonceFast(dest, baseNonce, uint64(i))
	}
}
