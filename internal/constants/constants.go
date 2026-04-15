// Package constants provides shared constants used across the project.
package constants

// Size constants for human-readable byte values.
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// Header constants
const (
	// HeaderSize returns the size in bytes of a serialized file header.
	HeaderSize = 4 + 1 + 16 + 4 // magic(4) + version(1) + salt(16) + chunkSize(4)

	// MinSupportedVersion is the earliest protocol version supported.
	MinSupportedVersion = 1
	// MaxSupportedVersion is the latest protocol version supported.
	MaxSupportedVersion = 2

	// MinChunkSize is the smallest allowed encrypted chunk size (1KB).
	MinChunkSize = KB

	// MaxChunkSize is the largest allowed encrypted chunk size (16MB).
	// This limit is used by the header package for validation.
	MaxChunkSize = 16 * MB
)
