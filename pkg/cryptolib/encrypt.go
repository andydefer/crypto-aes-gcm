// Package cryptolib provides AES-256-GCM file encryption and decryption.
//
// This package implements secure file encryption using AES-256-GCM in counter mode
// with Argon2id key derivation. It supports parallel streaming encryption for large
// files and includes integrity verification through HMAC.
package cryptolib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/andydefer/crypto-aes-gcm/internal/argon2"
	"github.com/andydefer/crypto-aes-gcm/internal/header"
)

// Encryptor handles parallel streaming encryption of data.
type Encryptor struct {
	workers          int
	chunkSize        int
	maxPendingChunks int
	bufferPool       sync.Pool
}

// chunkJob represents a single chunk of plaintext to be encrypted.
type chunkJob struct {
	index uint64
	data  []byte
}

// chunkResult represents an encrypted chunk with its original index.
type chunkResult struct {
	index      uint64
	ciphertext []byte
}

// NewEncryptor creates an Encryptor with the specified number of workers.
// The worker count is clamped between 1 and 2×CPU cores.
// Deprecated: Use NewEncryptorWithConfig instead for full configuration.
func NewEncryptor(workers int) (*Encryptor, error) {
	return NewEncryptorWithConfig(EncryptorConfig{
		Workers:          workers,
		ChunkSize:        DefaultChunkSize,
		MaxPendingChunks: DefaultMaxPendingChunks,
	})
}

// NewEncryptorWithConfig creates an Encryptor with the provided configuration.
//
// The configuration values are validated and clamped to safe ranges:
//   - Workers: clamped between 1 and 2×CPU cores
//   - ChunkSize: clamped between 1KB and 1GB
//   - MaxPendingChunks: clamped between 1 and MaxMaxPendingChunks (1000)
//
// Returns an error if the configuration is invalid after clamping.
func NewEncryptorWithConfig(config EncryptorConfig) (*Encryptor, error) {
	// Validate and clamp workers
	workers := config.Workers
	if workers <= 0 {
		workers = DefaultWorkers
	}
	maxWorkers := runtime.NumCPU() * 2
	if workers > maxWorkers {
		workers = maxWorkers
	}

	// Validate and clamp chunk size
	chunkSize := config.ChunkSize
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if chunkSize < header.MinChunkSize {
		chunkSize = header.MinChunkSize
	}
	if chunkSize > header.MaxChunkSize {
		chunkSize = header.MaxChunkSize
	}

	// Validate and clamp max pending chunks
	maxPending := config.MaxPendingChunks
	if maxPending <= 0 {
		maxPending = DefaultMaxPendingChunks
	}
	if maxPending > MaxMaxPendingChunks {
		maxPending = MaxMaxPendingChunks
	}

	return &Encryptor{
		workers:          workers,
		chunkSize:        chunkSize,
		maxPendingChunks: maxPending,
		bufferPool: sync.Pool{
			New: func() interface{} {
				buffer := make([]byte, chunkSize)
				return &buffer
			},
		},
	}, nil
}

// EncryptFile encrypts a file at inputPath and writes the result to outputPath.
func (e *Encryptor) EncryptFile(inputPath, outputPath, passphrase string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer output.Close()

	return e.Encrypt(input, output, passphrase)
}

// Encrypt reads from r, encrypts the data, and writes to w.
func (e *Encryptor) Encrypt(reader io.Reader, writer io.Writer, passphrase string) error {
	var salt [SaltSize]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return fmt.Errorf("generate salt: %w", err)
	}

	key := argon2.DeriveKey(passphrase, salt[:], argon2.DefaultParams())

	if err := e.writeHeader(writer, key, salt); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	baseNonce := make([]byte, NonceSize)
	if _, err := rand.Read(baseNonce); err != nil {
		return fmt.Errorf("generate nonce: %w", err)
	}

	if _, err := writer.Write(baseNonce); err != nil {
		return fmt.Errorf("write nonce: %w", err)
	}

	return e.processEncryption(reader, writer, gcm, baseNonce)
}

// writeHeader writes the file header and its HMAC to the writer.
func (e *Encryptor) writeHeader(writer io.Writer, key []byte, salt [SaltSize]byte) error {
	headerData := FileHeader{
		Magic:     [4]byte{Magic[0], Magic[1], Magic[2], Magic[3]},
		Version:   Version,
		Salt:      salt,
		ChunkSize: uint32(e.chunkSize),
	}

	if err := binary.Write(writer, binary.BigEndian, &headerData); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	headerHMAC := header.ComputeHMAC(key, header.Serialize(
		headerData.Magic,
		headerData.Version,
		headerData.Salt,
		headerData.ChunkSize,
	))

	if _, err := writer.Write(headerHMAC); err != nil {
		return fmt.Errorf("write header HMAC: %w", err)
	}

	return nil
}

// processEncryption orchestrates the parallel encryption pipeline.
func (e *Encryptor) processEncryption(reader io.Reader, writer io.Writer, gcm cipher.AEAD, baseNonce []byte) error {
	jobs := make(chan chunkJob, e.workers*2)
	results := make(chan chunkResult, e.workers*2)
	var workerWaitGroup sync.WaitGroup

	for i := 0; i < e.workers; i++ {
		workerWaitGroup.Add(1)
		go e.encryptionWorker(gcm, baseNonce, jobs, results, &workerWaitGroup)
	}

	go func() {
		workerWaitGroup.Wait()
		close(results)
	}()

	go e.readChunks(reader, jobs)

	if err := e.writeResults(results, writer); err != nil {
		return err
	}

	if err := binary.Write(writer, binary.BigEndian, uint32(0)); err != nil {
		return fmt.Errorf("write end marker: %w", err)
	}

	return nil
}

// encryptionWorker encrypts chunks from the jobs channel and sends results.
func (e *Encryptor) encryptionWorker(gcm cipher.AEAD, baseNonce []byte, jobs <-chan chunkJob, results chan<- chunkResult, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for job := range jobs {
		nonce := e.deriveChunkNonce(baseNonce, job.index)
		ciphertext := gcm.Seal(nil, nonce, job.data, nil)

		results <- chunkResult{
			index:      job.index,
			ciphertext: ciphertext,
		}

		e.bufferPool.Put(&job.data)
	}
}

// deriveChunkNonce creates a chunk-specific nonce by XORing the base nonce
// with the chunk index bytes to prevent nonce reuse.
func (e *Encryptor) deriveChunkNonce(baseNonce []byte, chunkIndex uint64) []byte {
	nonce := make([]byte, NonceSize)
	copy(nonce, baseNonce)

	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, chunkIndex)

	for i := 0; i < 8 && i < NonceSize-4; i++ {
		nonce[4+i] ^= indexBytes[i]
	}

	return nonce
}

// readChunks reads plaintext chunks from the reader and sends them to the jobs channel.
func (e *Encryptor) readChunks(reader io.Reader, jobs chan<- chunkJob) {
	defer close(jobs)

	var chunkIndex uint64

	for {
		bufferPtr := e.bufferPool.Get().(*[]byte)
		buffer := *bufferPtr

		bytesRead, err := io.ReadFull(reader, buffer[:e.chunkSize])

		if err == io.EOF {
			e.bufferPool.Put(bufferPtr)
			break
		}

		if err == io.ErrUnexpectedEOF {
			if bytesRead > 0 {
				jobs <- chunkJob{index: chunkIndex, data: buffer[:bytesRead]}
				chunkIndex++
			}
			e.bufferPool.Put(bufferPtr)
			break
		}

		if err != nil {
			e.bufferPool.Put(bufferPtr)
			break
		}

		jobs <- chunkJob{index: chunkIndex, data: buffer[:bytesRead]}
		chunkIndex++
	}
}

// writeResults writes ciphertext chunks in order with bounded memory usage.
//
// This function maintains a pending map for out-of-order chunks and limits
// the maximum pending size using e.maxPendingChunks to prevent memory exhaustion attacks.
func (e *Encryptor) writeResults(results <-chan chunkResult, writer io.Writer) error {
	expectedIndex := uint64(0)
	pending := make(map[uint64][]byte)

	for result := range results {
		if len(pending) > e.maxPendingChunks {
			return fmt.Errorf("too many pending chunks (limit %d) - possible reordering attack", e.maxPendingChunks)
		}

		pending[result.index] = result.ciphertext

		for {
			ciphertext, exists := pending[expectedIndex]
			if !exists {
				break
			}

			chunkLen := uint32(len(ciphertext))
			if err := binary.Write(writer, binary.BigEndian, chunkLen); err != nil {
				return fmt.Errorf("write chunk length: %w", err)
			}

			if _, err := writer.Write(ciphertext); err != nil {
				return fmt.Errorf("write ciphertext: %w", err)
			}

			delete(pending, expectedIndex)
			expectedIndex++
		}
	}

	if len(pending) > 0 {
		return fmt.Errorf("missing chunks: expected index %d, have %d pending", expectedIndex, len(pending))
	}

	return nil
}
