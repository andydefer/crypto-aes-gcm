// Package cryptolib provides cryptographic file encryption and decryption.
//
// This package implements AES-256-GCM encryption with Argon2id key derivation
// and parallel streaming capabilities for large files.
//
// The benchmarks in this file measure performance of:
//   - Encryption with various data sizes
//   - Decryption with various data sizes
//   - Parallel encryption with different worker counts
//   - Argon2id key derivation
//   - HMAC-SHA256 computation
package cryptolib

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/andydefer/crypto-aes-gcm/internal/argon2"
	"github.com/andydefer/crypto-aes-gcm/internal/constants"
	"github.com/andydefer/crypto-aes-gcm/internal/header"
)

// BenchmarkEncrypt measures encryption performance for different data sizes.
//
// It tests sizes: 1KB, 1MB, and 10MB to evaluate how encryption scales
// with input size. The benchmark creates random data and encrypts it
// repeatedly, measuring throughput and memory allocation.
//
// Results are useful for understanding performance characteristics
// and detecting regressions.
func BenchmarkEncrypt(b *testing.B) {
	sizes := []int{constants.KB, constants.MB, 10 * constants.MB}
	passphrase := "benchmark-password"

	for _, size := range sizes {
		b.Run(formatSize(size), func(b *testing.B) {
			data := make([]byte, size)
			_, _ = rand.Read(data)
			reader := bytes.NewReader(data)

			encryptor, err := NewEncryptor(DefaultWorkers())
			if err != nil {
				b.Fatalf("failed to create encryptor: %v", err)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var buf bytes.Buffer
				reader.Seek(0, 0)
				if err := encryptor.Encrypt(reader, &buf, passphrase); err != nil {
					b.Fatalf("encryption failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkDecrypt measures decryption performance for different data sizes.
//
// It tests sizes: 1KB, 1MB, and 10MB. The benchmark pre-encrypts the data
// once, then repeatedly decrypts it to measure throughput.
func BenchmarkDecrypt(b *testing.B) {
	sizes := []int{constants.KB, constants.MB, 10 * constants.MB}
	passphrase := "benchmark-password"

	for _, size := range sizes {
		b.Run(formatSize(size), func(b *testing.B) {
			originalData := make([]byte, size)
			_, _ = rand.Read(originalData)

			// Pre-encrypt data once
			var encryptedBuf bytes.Buffer
			encryptor, err := NewEncryptor(DefaultWorkers())
			if err != nil {
				b.Fatalf("failed to create encryptor: %v", err)
			}
			reader := bytes.NewReader(originalData)
			if err := encryptor.Encrypt(reader, &encryptedBuf, passphrase); err != nil {
				b.Fatalf("encryption failed: %v", err)
			}
			encryptedData := encryptedBuf.Bytes()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var decryptedBuf bytes.Buffer
				encryptedReader := bytes.NewReader(encryptedData)
				if err := DecryptStream(encryptedReader, &decryptedBuf, passphrase); err != nil {
					b.Fatalf("decryption failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkEncryptParallel measures parallel encryption performance with
// different worker counts.
//
// It tests sizes: 1MB and 10MB with worker counts: 1, 2, 4, 8.
// This benchmark helps determine optimal worker count for different
// hardware configurations and data sizes.
func BenchmarkEncryptParallel(b *testing.B) {
	sizes := []int{constants.MB, 10 * constants.MB}
	workerCounts := []int{1, 2, 4, 8}
	passphrase := "benchmark-password"

	for _, size := range sizes {
		for _, workers := range workerCounts {
			b.Run(formatSize(size)+"/workers_"+itoa(workers), func(b *testing.B) {
				data := make([]byte, size)
				_, _ = rand.Read(data)

				encryptor, err := NewEncryptor(workers)
				if err != nil {
					b.Fatalf("failed to create encryptor: %v", err)
				}

				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						var buf bytes.Buffer
						reader := bytes.NewReader(data)
						if err := encryptor.Encrypt(reader, &buf, passphrase); err != nil {
							b.Fatalf("encryption failed: %v", err)
						}
					}
				})
			})
		}
	}
}

// BenchmarkArgon2KeyDerivation measures key derivation performance.
//
// This benchmark evaluates the Argon2id KDF with default parameters.
// Results help tune the Time, Memory, and Threads parameters for
// the target hardware.
func BenchmarkArgon2KeyDerivation(b *testing.B) {
	passphrase := "benchmark-passphrase"
	salt := make([]byte, 32)
	_, _ = rand.Read(salt)
	params := argon2.DefaultParams()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = argon2.DeriveKey(passphrase, salt, params)
	}
}

// BenchmarkHMACComputation measures HMAC-SHA256 performance.
//
// This benchmark evaluates HMAC computation on 1KB of data.
// It helps understand the overhead of header authentication.
func BenchmarkHMACComputation(b *testing.B) {
	key := make([]byte, 32)
	data := make([]byte, constants.KB)
	_, _ = rand.Read(key)
	_, _ = rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = header.ComputeHMAC(key, data)
	}
}

// formatSize converts a byte count to a human-readable string.
//
// It returns strings like "512", "1KB", "5MB", "2GB" for use in
// benchmark subtest names, making them readable and unique.
//
// Parameters:
//   - bytes: Number of bytes to format
//
// Returns:
//   - Human-readable size string
func formatSize(bytes int) string {
	if bytes < constants.KB {
		return itoa(bytes)
	}
	if bytes < constants.MB {
		return itoa(bytes/constants.KB) + "KB"
	}
	if bytes < constants.GB {
		return itoa(bytes/constants.MB) + "MB"
	}
	return itoa(bytes/constants.GB) + "GB"
}

// itoa converts an integer to a string without using strconv.Itoa.
//
// This is a simple helper for benchmark subtest names to avoid
// importing strconv, which would add overhead to benchmarks.
//
// Parameters:
//   - n: Integer to convert
//
// Returns:
//   - String representation of the integer
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	negative := n < 0
	if negative {
		n = -n
	}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
