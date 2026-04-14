// Package cryptolib provides AES-256-GCM file encryption and decryption.
//
// This package implements secure file encryption using AES-256-GCM in counter mode
// with Argon2id key derivation. It supports parallel streaming encryption for large
// files and includes integrity verification through HMAC.
package cryptolib

import (
	"context"
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
	"github.com/andydefer/crypto-aes-gcm/internal/crypto"
	"github.com/andydefer/crypto-aes-gcm/internal/header"
	"github.com/andydefer/crypto-aes-gcm/internal/lang"
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
//
// Parameters:
//   - workers: Number of parallel encryption workers (clamped between 1 and 2×CPU cores)
//
// Returns:
//   - *Encryptor: Configured encryptor instance
//   - error: If configuration is invalid after clamping
//
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
// Parameters:
//   - config: Configuration parameters for the encryptor
//
// Returns:
//   - *Encryptor: Configured encryptor instance
//   - error: If configuration is invalid after clamping
//
// The configuration values are validated and clamped to safe ranges:
//   - Workers: clamped between 1 and 2×CPU cores
//   - ChunkSize: clamped between 1KB and 1GB
//   - MaxPendingChunks: clamped between 1 and MaxMaxPendingChunks (1000)
func NewEncryptorWithConfig(config EncryptorConfig) (*Encryptor, error) {
	workers := clampWorkers(config.Workers)
	chunkSize := clampChunkSize(config.ChunkSize)
	maxPending := clampMaxPending(config.MaxPendingChunks)

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
//
// Parameters:
//   - inputPath: Path to the plaintext file to encrypt
//   - outputPath: Path where encrypted output will be written
//   - passphrase: User passphrase for key derivation
//
// Returns:
//   - error: Any error encountered during file operations or encryption
func (e *Encryptor) EncryptFile(inputPath, outputPath, passphrase string) error {
	return e.EncryptFileWithContext(context.Background(), inputPath, outputPath, passphrase)
}

// EncryptFileWithContext encrypts a file with context support for cancellation.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - inputPath: Path to the plaintext file to encrypt
//   - outputPath: Path where encrypted output will be written
//   - passphrase: User passphrase for key derivation
//
// Returns:
//   - error: Any error encountered during file operations, encryption, or context cancellation
func (e *Encryptor) EncryptFileWithContext(ctx context.Context, inputPath, outputPath, passphrase string) (err error) {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrOpenInputEnc), err)
	}
	defer closeWithErrorHandler(input, &err, lang.T(lang.CryptolibErrCloseInput))

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrCreateOutputEnc), err)
	}
	defer closeWithErrorHandler(output, &err, lang.T(lang.CryptolibErrCloseOutput))

	return e.EncryptWithContext(ctx, input, output, passphrase)
}

// Encrypt reads from r, encrypts the data, and writes to w.
//
// Parameters:
//   - reader: Source of plaintext data
//   - writer: Destination for encrypted data
//   - passphrase: User passphrase for key derivation
//
// Returns:
//   - error: Any error encountered during encryption
func (e *Encryptor) Encrypt(reader io.Reader, writer io.Writer, passphrase string) error {
	return e.EncryptWithContext(context.Background(), reader, writer, passphrase)
}

// EncryptWithContext reads from r, encrypts the data, and writes to w with context support.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - reader: Source of plaintext data
//   - writer: Destination for encrypted data
//   - passphrase: User passphrase for key derivation
//
// Returns:
//   - error: Any error encountered during encryption or context cancellation
func (e *Encryptor) EncryptWithContext(ctx context.Context, reader io.Reader, writer io.Writer, passphrase string) error {
	var salt [SaltSize]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrGenerateSalt), err)
	}

	key := argon2.DeriveKey(passphrase, salt[:], argon2.DefaultParams())

	if err := e.writeHeader(writer, key, salt); err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrCreateCipherEnc), err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrCreateGCMEnc), err)
	}

	baseNonce := make([]byte, NonceSize)
	if _, err := rand.Read(baseNonce); err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrGenerateNonce), err)
	}

	if _, err := writer.Write(baseNonce); err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrWriteNonce), err)
	}

	return e.processEncryptionWithContext(ctx, reader, writer, gcm, baseNonce)
}

// writeHeader writes the file header and its HMAC to the writer.
//
// Parameters:
//   - writer: Destination for header data
//   - key: Encryption key for HMAC computation
//   - salt: Cryptographic salt to store in header
//
// Returns:
//   - error: If header writing or HMAC computation fails
func (e *Encryptor) writeHeader(writer io.Writer, key []byte, salt [SaltSize]byte) error {
	headerData := FileHeader{
		Magic:     [4]byte{Magic[0], Magic[1], Magic[2], Magic[3]},
		Version:   Version,
		Salt:      salt,
		ChunkSize: uint32(e.chunkSize),
	}

	if err := binary.Write(writer, binary.BigEndian, &headerData); err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrWriteHeader), err)
	}

	headerHMAC := header.ComputeHMAC(key, header.Serialize(
		headerData.Magic,
		headerData.Version,
		headerData.Salt,
		headerData.ChunkSize,
	))

	if _, err := writer.Write(headerHMAC); err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrWriteHeaderHMAC), err)
	}

	return nil
}

// processEncryptionWithContext orchestrates the parallel encryption pipeline with context support.
//
// Parameters:
//   - ctx: Context for cancellation control
//   - reader: Source of plaintext data
//   - writer: Destination for encrypted data
//   - gcm: GCM cipher for authenticated encryption
//   - baseNonce: Base nonce for counter-based nonce derivation
//
// Returns:
//   - error: Any error encountered during the encryption pipeline
//
// The function sets up a worker pool for parallel encryption, reads chunks from
// the input, distributes them to workers, and writes results in order.
func (e *Encryptor) processEncryptionWithContext(ctx context.Context, reader io.Reader, writer io.Writer, gcm cipher.AEAD, baseNonce []byte) error {
	jobs := make(chan chunkJob, e.workers*2)
	results := make(chan chunkResult, e.workers*2)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go e.encryptionWorker(ctx, gcm, baseNonce, jobs, results, &wg, errChan)
	}

	go func() {
		if err := e.readChunks(ctx, reader, jobs); err != nil {
			errChan <- err
			cancel()
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	if err := e.writeResultsWithContext(ctx, results, writer); err != nil {
		cancel()
		return err
	}

	select {
	case err := <-errChan:
		return err
	default:
	}

	if err := binary.Write(writer, binary.BigEndian, uint32(0)); err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrWriteEndMarker), err)
	}

	return nil
}

// encryptionWorker encrypts chunks from the jobs channel and sends results.
//
// Parameters:
//   - ctx: Context for cancellation
//   - gcm: GCM cipher for authenticated encryption
//   - baseNonce: Base nonce for nonce derivation
//   - jobs: Channel receiving plaintext chunks
//   - results: Channel sending encrypted chunks
//   - wg: WaitGroup for worker coordination
//   - errChan: Channel for error reporting
func (e *Encryptor) encryptionWorker(ctx context.Context, gcm cipher.AEAD, baseNonce []byte, jobs <-chan chunkJob, results chan<- chunkResult, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	var nonceBuf [crypto.NonceSize]byte

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			if err := crypto.DeriveChunkNonceFast(nonceBuf[:], baseNonce, job.index); err != nil {
				select {
				case <-ctx.Done():
				case errChan <- fmt.Errorf(lang.T(lang.CryptolibErrNonceDerivation), err):
				}
				return
			}

			ciphertext := gcm.Seal(nil, nonceBuf[:], job.data, nil)

			select {
			case <-ctx.Done():
				return
			case results <- chunkResult{
				index:      job.index,
				ciphertext: ciphertext,
			}:
			}

			e.bufferPool.Put(&job.data)
		}
	}
}

// readChunks reads plaintext chunks from the reader and sends them to the jobs channel.
//
// Parameters:
//   - ctx: Context for cancellation
//   - reader: Source of plaintext data
//   - jobs: Channel for sending plaintext chunks
//
// Returns:
//   - error: If reading fails or context is cancelled
func (e *Encryptor) readChunks(ctx context.Context, reader io.Reader, jobs chan<- chunkJob) error {
	var chunkIndex uint64

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		bufferPtr := e.bufferPool.Get().(*[]byte)
		buffer := *bufferPtr

		bytesRead, err := io.ReadFull(reader, buffer[:e.chunkSize])

		if err == io.EOF {
			e.bufferPool.Put(bufferPtr)
			break
		}

		if err == io.ErrUnexpectedEOF {
			if bytesRead > 0 {
				select {
				case <-ctx.Done():
					e.bufferPool.Put(bufferPtr)
					return ctx.Err()
				case jobs <- chunkJob{index: chunkIndex, data: buffer[:bytesRead]}:
					chunkIndex++
				}
			}
			e.bufferPool.Put(bufferPtr)
			break
		}

		if err != nil {
			e.bufferPool.Put(bufferPtr)
			return fmt.Errorf(lang.T(lang.CryptolibErrReadChunk), err)
		}

		select {
		case <-ctx.Done():
			e.bufferPool.Put(bufferPtr)
			return ctx.Err()
		case jobs <- chunkJob{index: chunkIndex, data: buffer[:bytesRead]}:
			chunkIndex++
		}
	}

	return nil
}

// writeResultsWithContext writes ciphertext chunks in order with context support.
//
// Parameters:
//   - ctx: Context for cancellation
//   - results: Channel receiving encrypted chunks
//   - writer: Destination for encrypted data
//
// Returns:
//   - error: If chunk reordering is detected or writing fails
//
// The function buffers out-of-order chunks and writes them sequentially.
// A pending chunk limit prevents memory exhaustion attacks.
func (e *Encryptor) writeResultsWithContext(ctx context.Context, results <-chan chunkResult, writer io.Writer) error {
	expectedIndex := uint64(0)
	pending := make(map[uint64][]byte)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case result, ok := <-results:
			if !ok {
				if len(pending) > 0 {
					return fmt.Errorf(lang.T(lang.CryptolibErrMissingChunks), expectedIndex, len(pending))
				}
				return nil
			}

			if len(pending) > e.maxPendingChunks {
				return fmt.Errorf(lang.T(lang.CryptolibErrTooManyPending), e.maxPendingChunks)
			}

			pending[result.index] = result.ciphertext

			for {
				ciphertext, exists := pending[expectedIndex]
				if !exists {
					break
				}

				chunkLen := uint32(len(ciphertext))
				if err := binary.Write(writer, binary.BigEndian, chunkLen); err != nil {
					return fmt.Errorf(lang.T(lang.CryptolibErrWriteChunkLen), err)
				}

				if _, err := writer.Write(ciphertext); err != nil {
					return fmt.Errorf(lang.T(lang.CryptolibErrWriteCiphertext), err)
				}

				delete(pending, expectedIndex)
				expectedIndex++
			}
		}
	}
}

// clampWorkers ensures the worker count is within acceptable bounds.
func clampWorkers(workers int) int {
	if workers <= 0 {
		workers = DefaultWorkers()
	}
	maxWorkers := runtime.NumCPU() * 2
	if workers > maxWorkers {
		workers = maxWorkers
	}
	return workers
}

// clampChunkSize ensures the chunk size is within acceptable bounds.
func clampChunkSize(chunkSize int) int {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if chunkSize < header.MinChunkSize {
		chunkSize = header.MinChunkSize
	}
	if chunkSize > header.MaxChunkSize {
		chunkSize = header.MaxChunkSize
	}
	return chunkSize
}

// clampMaxPending ensures the max pending chunks count is within acceptable bounds.
func clampMaxPending(maxPending int) int {
	if maxPending <= 0 {
		maxPending = DefaultMaxPendingChunks
	}
	if maxPending > MaxMaxPendingChunks {
		maxPending = MaxMaxPendingChunks
	}
	return maxPending
}

// closeWithErrorHandler closes a file and updates the error reference if closing fails.
func closeWithErrorHandler(f *os.File, err *error, context string) {
	if closeErr := f.Close(); closeErr != nil && *err == nil {
		*err = fmt.Errorf("%s: %w", context, closeErr)
	}
}
