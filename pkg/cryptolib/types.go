package cryptolib

const (
	// Magic identifies files created by this library.
	Magic = "CRYP"

	// Version indicates the current file format version.
	Version = 2

	// SaltSize is the length of the Argon2id salt in bytes.
	SaltSize = 16

	// NonceSize is the length of the GCM nonce in bytes.
	NonceSize = 12

	// KeySize is the length of the derived AES key in bytes (256 bits).
	KeySize = 32

	// DefaultChunkSize is the size of each encrypted chunk (1MB).
	DefaultChunkSize = 1024 * 1024

	// DefaultWorkers is the default number of parallel encryption workers.
	DefaultWorkers = 4
)

// FileHeader represents the header structure at the beginning of every encrypted file.
type FileHeader struct {
	Magic     [4]byte
	Version   byte
	Salt      [SaltSize]byte
	ChunkSize uint32
}
