// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
package cryptolib

import (
	"bytes"
	"runtime"
	"testing"
)

// Size constants for human-readable byte values.
const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

// TestDefaultEncryptorConfig verifies that default configuration values are correct.
func TestDefaultEncryptorConfig(t *testing.T) {
	config := DefaultEncryptorConfig()

	if config.Workers != DefaultWorkers {
		t.Errorf("Expected Workers=%d, got %d", DefaultWorkers, config.Workers)
	}
	if config.ChunkSize != DefaultChunkSize {
		t.Errorf("Expected ChunkSize=%d, got %d", DefaultChunkSize, config.ChunkSize)
	}
	if config.MaxPendingChunks != DefaultMaxPendingChunks {
		t.Errorf("Expected MaxPendingChunks=%d, got %d", DefaultMaxPendingChunks, config.MaxPendingChunks)
	}
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
			expectedWorker:  DefaultWorkers,
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
			expectedWorker:  DefaultWorkers,
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
				Workers:          DefaultWorkers,
				ChunkSize:        0,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers,
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "negative chunk size uses default",
			config: EncryptorConfig{
				Workers:          DefaultWorkers,
				ChunkSize:        -KB,
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers,
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "too small chunk size clamped to min",
			config: EncryptorConfig{
				Workers:          DefaultWorkers,
				ChunkSize:        512, // MinChunkSize is KB (1024)
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers,
			expectedChunk:   KB,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "too large chunk size clamped to max",
			config: EncryptorConfig{
				Workers:          DefaultWorkers,
				ChunkSize:        2 * GB, // 2GB > MaxChunkSize (1GB)
				MaxPendingChunks: DefaultMaxPendingChunks,
			},
			expectedWorker:  DefaultWorkers,
			expectedChunk:   GB,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "zero pending chunks uses default",
			config: EncryptorConfig{
				Workers:          DefaultWorkers,
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: 0,
			},
			expectedWorker:  DefaultWorkers,
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "negative pending chunks uses default",
			config: EncryptorConfig{
				Workers:          DefaultWorkers,
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: -10,
			},
			expectedWorker:  DefaultWorkers,
			expectedChunk:   DefaultChunkSize,
			expectedPending: DefaultMaxPendingChunks,
		},
		{
			name: "excessive pending chunks capped",
			config: EncryptorConfig{
				Workers:          DefaultWorkers,
				ChunkSize:        DefaultChunkSize,
				MaxPendingChunks: 5000,
			},
			expectedWorker:  DefaultWorkers,
			expectedChunk:   DefaultChunkSize,
			expectedPending: MaxMaxPendingChunks,
		},
		{
			name: "custom valid values",
			config: EncryptorConfig{
				Workers:          8,
				ChunkSize:        2 * MB, // 2MB
				MaxPendingChunks: 50,
			},
			expectedWorker:  8,
			expectedChunk:   2 * MB,
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
		ChunkSize:        64 * KB, // 64KB chunks
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
	dataSize := 2 * MB   // 2MB
	chunkSize := 64 * KB // 64KB chunks -> about 32 chunks total
	smallPendingLimit := 5

	data := make([]byte, dataSize)
	for i := range data {
		data[i] = byte(i % 256)
	}

	config := EncryptorConfig{
		Workers:          1, // Single worker = sequential = no out-of-order
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
	dataSize := 5 * MB   // 5MB
	chunkSize := 64 * KB // 64KB chunks -> about 80 chunks total
	unreasonableLimit := 3

	data := make([]byte, dataSize)
	for i := range data {
		data[i] = byte(i % 256)
	}

	config := EncryptorConfig{
		Workers:          4, // Multiple workers = potential out-of-order
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

	if err != nil && !bytes.Contains([]byte(err.Error()), []byte("too many pending chunks")) {
		t.Errorf("expected 'too many pending chunks' error, got: %v", err)
	}
}

// TestNewEncryptorWithConfig_ZeroPendingLimit verifies that zero/negative values
// are properly clamped to default.
func TestNewEncryptorWithConfig_ZeroPendingLimit(t *testing.T) {
	testData := []byte("test data for zero pending limit")

	config := EncryptorConfig{
		Workers:          1,
		ChunkSize:        64 * KB,
		MaxPendingChunks: 0, // Should be clamped to DefaultMaxPendingChunks
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

// BenchmarkEncryptWithConfig measures performance with different pending chunk limits.
func BenchmarkEncryptWithConfig(b *testing.B) {
	dataSize := 10 * MB // 10MB
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
				Workers:          DefaultWorkers,
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
