// Package header provides file header serialization and HMAC utilities.
package header

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
)

// ToBytes serializes a file header to bytes for HMAC calculation.
// The byte layout is fixed and deterministic across platforms.
func ToBytes(magic [4]byte, version byte, salt [16]byte, chunkSize uint32) []byte {
	buf := make([]byte, 4+1+16+4)

	magicValue := uint32(magic[0])<<24 |
		uint32(magic[1])<<16 |
		uint32(magic[2])<<8 |
		uint32(magic[3])
	binary.BigEndian.PutUint32(buf[0:4], magicValue)

	buf[4] = version

	copy(buf[5:21], salt[:])

	binary.BigEndian.PutUint32(buf[21:25], chunkSize)

	return buf
}

// ComputeHMAC computes an HMAC-SHA256 digest for the given data.
func ComputeHMAC(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// VerifyHMAC securely compares a computed HMAC with a stored value.
// Uses constant-time comparison to prevent timing attacks.
func VerifyHMAC(key, data, stored []byte) bool {
	computed := ComputeHMAC(key, data)
	return hmac.Equal(computed, stored)
}
