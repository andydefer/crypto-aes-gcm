// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
package cryptolib

import (
	"bytes"
	"runtime"
	"testing"

	"github.com/andydefer/crypto-aes-gcm/internal/constants"
)

// TestDefaultEncryptorConfig verifies that default configuration values are correct.
func TestDefaultEncryptorConfig(t *testing.T) {
	config := DefaultEncryptorConfig()

	if config.Workers != DefaultWorkers() {
		t.Errorf("Expected Workers=%d, got %d", DefaultWorkers(), config.Workers)
	}
	if config.ChunkSize != DefaultChunkSize {
		t.Errorf("Expected ChunkSize=%d, got %d", DefaultChunkSize, config.ChunkSize)
	}
	if config.MaxPendingChunks != DefaultMaxPendingChunks {
		t.Errorf("Expected MaxPendingChunks=%d, got %d", DefaultMaxPendingChunks, config.MaxPendingChunks)
	}
}

// TestMaxChunkSizeConstant verifies that MaxChunkSize is properly defined
// and has a reasonable value for security.
func TestMaxChunkSizeConstant(t *testing.T) {
	// Verify MaxChunkSize is positive
	if MaxChunkSize <= 0 {
		t.Errorf("MaxChunkSize should be positive, got %d", MaxChunkSize)
	}

	// Verify MaxChunkSize is at least DefaultChunkSize
	if MaxChunkSize < DefaultChunkSize {
		t.Errorf("MaxChunkSize (%d) should be at least DefaultChunkSize (%d)", MaxChunkSize, DefaultChunkSize)
	}

	// Verify MaxChunkSize is not unreasonably large (should be <= 100MB for safety)
	reasonableMax := 100 * constants.MB
	if MaxChunkSize > reasonableMax {
		t.Errorf("MaxChunkSize (%d) exceeds reasonable maximum of %d bytes", MaxChunkSize, reasonableMax)
	}

	// Verify MaxChunkSize is a multiple of KB (reasonable alignment)
	if MaxChunkSize%constants.KB != 0 {
		t.Logf("Warning: MaxChunkSize (%d) is not a multiple of %d bytes", MaxChunkSize, constants.KB)
	}

	// Log the current value for documentation
	t.Logf("MaxChunkSize = %d bytes (%.2f MB)", MaxChunkSize, float64(MaxChunkSize)/float64(constants.MB))
}

// TestNewEncryptorWithConfig_Clamping verifies that configuration values are properly clamped.
func TestNewEncryptorWithConfig_Clamping(t *testing.T) {
	maxWorkers := runtime.NumCPU() * 2

	tests := []struct {
		name            string
		config          EncryptorConfig
		expectedWorker  int
		expectedChunk   int
		expectedPending int
	}{
		{
			name: "zero workers uses default",
			config: EncryptorConfig{
				Workers:          0,
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "negative workers uses default",
			config: EncryptorConfig{
				Workers:          -5,
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "excessive workers capped",
			config: EncryptorConfig{
				Workers:          9999,
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  maxWorkers,
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "zero chunk size uses default",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        0,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "negative chunk size uses default",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        -constants.KB,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "too small chunk size clamped to min",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        512,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   constants.KB,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "too large chunk size clamped to max",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        2 * constants.GB,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   10 * constants.MB, // Valeur réelle de cryptolib.MaxChunkSize
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "chunk size exactly at MaxChunkSize",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        MaxChunkSize,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   MaxChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "zero pending chunks uses default",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: 0,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "negative pending chunks uses default",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: -10,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "excessive pending chunks capped",
			config: EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: 5000,
			},
			expectedWorker:  DefaultWorkers(),
			expectedChunk:   DefaultChunkSize,
			expectedPending: MaxMaxPendingChunks,
		},
		{
			name: "custom valid values",
			config: EncryptorConfig{
				Workers:          8,
				ChunkSize:        2 * constants.MB,
				MaxPendingChunks: 50,
			},
			expectedWorker:  8,
			expectedChunk:   2 * constants.MB,
			expectedPending: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptor, err := NewEncryptorWithConfig(tt.config)
			if err != nil {
				t.Fatalf("NewEncryptorWithConfig failed: %v", err)
			}

			if encryptor.workers != tt.expectedWorker {
				t.Errorf("workers: expected %d, got %d", tt.expectedWorker, encryptor.workers)
			}
			if encryptor.chunkSize != tt.expectedChunk {
				t.Errorf("chunkSize: expected %d, got %d", tt.expectedChunk, encryptor.chunkSize)
			}
			if encryptor.maxPendingChunks != tt.expectedPending {
				t.Errorf("maxPendingChunks: expected %d, got %d", tt.expectedPending, encryptor.maxPendingChunks)
			}
		})
	}
}

// TestNewEncryptorWithConfig_EncryptionDecryption verifies that custom config works correctly.
func TestNewEncryptorWithConfig_EncryptionDecryption(t *testing.T) {
	originalData := []byte("test data for custom config encryption")
	password := "custom-config-password"

	config := EncryptorConfig{
		Workers:          4,
		ChunkSize:        64 * constants.KB,
		MaxPendingChunks: 25,
	}

	encryptor, err := NewEncryptorWithConfig(config)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	var encryptedBuf bytes.Buffer
	reader := bytes.NewReader(originalData)

	if err := encryptor.Encrypt(reader, &encryptedBuf, password); err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, password); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(originalData, decryptedBuf.Bytes()) {
		t.Errorf("data mismatch: original %d bytes, decrypted %d bytes",
			len(originalData), len(decryptedBuf.Bytes()))
	}
}

// TestNewEncryptorWithConfig_PendingChunksLimit verifies that the pending chunks
// limit is respected. Uses single worker to guarantee in-order processing
// so the limit is never exceeded.
func TestNewEncryptorWithConfig_PendingChunksLimit(t *testing.T) {
	dataSize := 2 * constants.MB
	chunkSize := 64 * constants.KB
	smallPendingLimit := 5

	data := make([]byte, dataSize)
	for i := range data {
		data[i] = byte(i % 256)
	}

	config := EncryptorConfig{
		Workers:          1,
		ChunkSize:        chunkSize,
		MaxPendingChunks: smallPendingLimit,
	}

	encryptor, err := NewEncryptorWithConfig(config)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	var encryptedBuf bytes.Buffer
	reader := bytes.NewReader(data)

	if err := encryptor.Encrypt(reader, &encryptedBuf, "test-password"); err != nil {
		t.Fatalf("encryption with single worker and pending limit %d failed: %v", smallPendingLimit, err)
	}

	var decryptedBuf bytes.Buffer
	encryptedReader := bytes.NewReader(encryptedBuf.Bytes())

	if err := DecryptStream(encryptedReader, &decryptedBuf, "test-password"); err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if !bytes.Equal(data, decryptedBuf.Bytes()) {
		t.Error("decrypted data mismatch")
	}
}

// TestNewEncryptorWithConfig_PendingChunksLimitExceeded verifies that the
// encryptor fails when the pending chunks limit is set too low for parallel workers.
func TestNewEncryptorWithConfig_PendingChunksLimitExceeded(t *testing.T) {
	dataSize := 5 * constants.MB
	chunkSize := 64 * constants.KB
	unreasonableLimit := 3

	data := make([]byte, dataSize)
	for i := range data {
		data[i] = byte(i % 256)
	}

	config := EncryptorConfig{
		Workers:          4,
		ChunkSize:        chunkSize,
		MaxPendingChunks: unreasonableLimit,
	}

	encryptor, err := NewEncryptorWithConfig(config)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	var encryptedBuf bytes.Buffer
	reader := bytes.NewReader(data)

	err = encryptor.Encrypt(reader, &encryptedBuf, "test-password")

	if err == nil {
		t.Error("expected encryption to fail with unreasonably low pending limit, but it succeeded")
	}
}

// TestNewEncryptorWithConfig_ZeroPendingLimit verifies that zero/negative values
// are properly clamped to default.
func TestNewEncryptorWithConfig_ZeroPendingLimit(t *testing.T) {
	testData := []byte("test data for zero pending limit")

	config := EncryptorConfig{
		Workers:          1,
		ChunkSize:        64 * constants.KB,
		MaxPendingChunks: 0,
	}

	encryptor, err := NewEncryptorWithConfig(config)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	if encryptor.maxPendingChunks != DefaultMaxPendingChunks {
		t.Errorf("maxPendingChunks should be clamped to %d, got %d",
			DefaultMaxPendingChunks, encryptor.maxPendingChunks)
	}

	var encryptedBuf bytes.Buffer
	reader := bytes.NewReader(testData)

	if err := encryptor.Encrypt(reader, &encryptedBuf, "test-password"); err != nil {
		t.Fatalf("encryption with default pending limit failed: %v", err)
	}
}

// TestMemoryLeak verifies that encryption doesn't leak memory.
func TestMemoryLeak(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping memory leak test in short mode")
	}

	data := make([]byte, 10*constants.MB)
	password := "test-password"

	var memStats1, memStats2 runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&memStats1)

	const iterations = 10
	for i := 0; i < iterations; i++ {
		encryptor, err := NewEncryptor(DefaultWorkers())
		if err != nil {
			t.Fatalf("failed to create encryptor: %v", err)
		}

		var buf bytes.Buffer
		reader := bytes.NewReader(data)
		if err := encryptor.Encrypt(reader, &buf, password); err != nil {
			t.Fatalf("encryption failed: %v", err)
		}
	}

	runtime.GC()
	runtime.ReadMemStats(&memStats2)

	allocDiff := int64(memStats2.Alloc) - int64(memStats1.Alloc)
	growthPercent := float64(allocDiff) / float64(memStats1.Alloc) * 100

	maxAllowedGrowthPercent := 10.0

	if growthPercent > maxAllowedGrowthPercent {
		t.Errorf("Memory growth %.1f%% exceeds allowed limit of %.0f%% (growth: %d bytes)",
			growthPercent, maxAllowedGrowthPercent, allocDiff)
	}

	t.Logf("Memory growth: %d bytes (%.1f%%)", allocDiff, growthPercent)
}

// BenchmarkEncryptWithConfig measures performance with different pending chunk limits.
func BenchmarkEncryptWithConfig(b *testing.B) {
	dataSize := 10 * constants.MB
	data := make([]byte, dataSize)

	configs := []struct {
		name  string
		limit int
	}{
		{"pending_10", 10},
		{"pending_50", 50},
		{"pending_100", 100},
		{"pending_500", 500},
		{"pending_1000", 1000},
	}

	for _, cfg := range configs {
		b.Run(cfg.name, func(b *testing.B) {
			encryptor, err := NewEncryptorWithConfig(EncryptorConfig{
				Workers:          DefaultWorkers(),
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: cfg.limit,
			})
			if err != nil {
				b.Fatalf("failed to create encryptor: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var buf bytes.Buffer
				reader := bytes.NewReader(data)
				if err := encryptor.Encrypt(reader, &buf, "benchmark"); err != nil {
					b.Fatalf("encryption failed: %v", err)
				}
			}
		})
	}
}
