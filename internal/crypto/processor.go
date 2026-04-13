package crypto

const (
	Magic            = "CRYP"
	Version          = 2
	SaltSize         = 16
	NonceSize        = 12
	KeySize          = 32
	DefaultChunkSize = 1024 * 1024
	DefaultWorkers   = 4
)
