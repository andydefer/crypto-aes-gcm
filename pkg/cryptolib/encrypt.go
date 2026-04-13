package cryptolib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/andydefer/crypto-aes-gcm/internal/argon2"
	"github.com/andydefer/crypto-aes-gcm/internal/header"
)

// Encryptor handles parallel streaming encryption of data.
type Encryptor struct {
	workers    int
	chunkSize  int
	bufferPool sync.Pool
}

type chunkJob struct {
	index uint64
	data  []byte
}

type chunkResult struct {
	index      uint64
	ciphertext []byte
}

// NewEncryptor creates an Encryptor with the specified number of workers.
// The worker count is clamped between 1 and 2×CPU cores.
func NewEncryptor(workers int) (*Encryptor, error) {
	if workers <= 0 {
		workers = DefaultWorkers
	}
	if workers > runtime.NumCPU()*2 {
		workers = runtime.NumCPU() * 2
	}

	return &Encryptor{
		workers:   workers,
		chunkSize: DefaultChunkSize,
		bufferPool: sync.Pool{
			New: func() interface{} {
				b := make([]byte, DefaultChunkSize)
				return &b
			},
		},
	}, nil
}

// EncryptFile encrypts a file at inputPath and writes the result to outputPath.
func (e *Encryptor) EncryptFile(inputPath, outputPath, passphrase string) error {
	in, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}
	defer out.Close()

	return e.Encrypt(in, out, passphrase)
}

// Encrypt reads from r, encrypts the data, and writes to w.
func (e *Encryptor) Encrypt(r io.Reader, w io.Writer, passphrase string) error {
	var salt [SaltSize]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return fmt.Errorf("generate salt: %w", err)
	}

	key := argon2.DeriveKey(passphrase, salt[:], argon2.DefaultParams())

	if err := e.writeHeader(w, key, salt); err != nil {
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
	if _, err := w.Write(baseNonce); err != nil {
		return fmt.Errorf("write nonce: %w", err)
	}

	return e.processEncryption(r, w, gcm, key, baseNonce)
}

func (e *Encryptor) writeHeader(w io.Writer, key []byte, salt [SaltSize]byte) error {
	headerData := FileHeader{
		Magic:     [4]byte{Magic[0], Magic[1], Magic[2], Magic[3]},
		Version:   Version,
		Salt:      salt,
		ChunkSize: uint32(e.chunkSize),
	}

	if err := binary.Write(w, binary.BigEndian, &headerData); err != nil {
		return fmt.Errorf("write header: %w", err)
	}

	headerHMAC := header.ComputeHMAC(key, header.ToBytes(
		headerData.Magic,
		headerData.Version,
		headerData.Salt,
		headerData.ChunkSize,
	))
	if _, err := w.Write(headerHMAC); err != nil {
		return fmt.Errorf("write header HMAC: %w", err)
	}

	return nil
}

func (e *Encryptor) processEncryption(r io.Reader, w io.Writer, gcm cipher.AEAD, key []byte, baseNonce []byte) error {
	jobs := make(chan chunkJob, e.workers*2)
	results := make(chan chunkResult, e.workers*2)
	var wg sync.WaitGroup

	for i := 0; i < e.workers; i++ {
		wg.Add(1)
		go e.encryptionWorker(gcm, baseNonce, jobs, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go e.readChunks(r, jobs)

	globalHMAC := hmac.New(sha256.New, key)
	if err := e.writeResults(results, w, globalHMAC); err != nil {
		return err
	}

	if err := binary.Write(w, binary.BigEndian, uint32(0)); err != nil {
		return fmt.Errorf("write end marker: %w", err)
	}

	if _, err := w.Write(globalHMAC.Sum(nil)); err != nil {
		return fmt.Errorf("write global HMAC: %w", err)
	}

	return nil
}

func (e *Encryptor) encryptionWorker(gcm cipher.AEAD, baseNonce []byte, jobs <-chan chunkJob, results chan<- chunkResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		nonce := make([]byte, NonceSize)
		copy(nonce, baseNonce)
		binary.BigEndian.PutUint64(nonce[4:], job.index)

		ciphertext := gcm.Seal(nil, nonce, job.data, nil)
		results <- chunkResult{
			index:      job.index,
			ciphertext: ciphertext,
		}
		e.bufferPool.Put(&job.data)
	}
}

func (e *Encryptor) readChunks(r io.Reader, jobs chan<- chunkJob) {
	defer close(jobs)
	var index uint64

	for {
		bufPtr := e.bufferPool.Get().(*[]byte)
		buf := *bufPtr
		n, err := io.ReadFull(r, buf[:e.chunkSize])

		if err == io.EOF {
			e.bufferPool.Put(bufPtr)
			break
		}

		if err == io.ErrUnexpectedEOF {
			if n > 0 {
				jobs <- chunkJob{index: index, data: buf[:n]}
				index++
			}
			e.bufferPool.Put(bufPtr)
			break
		}

		if err != nil {
			e.bufferPool.Put(bufPtr)
			break
		}

		jobs <- chunkJob{index: index, data: buf[:n]}
		index++
	}
}

func (e *Encryptor) writeResults(results <-chan chunkResult, w io.Writer, globalHMAC hash.Hash) error {
	expectedIndex := uint64(0)
	pending := make(map[uint64][]byte)

	for result := range results {
		pending[result.index] = result.ciphertext

		for {
			ciphertext, ok := pending[expectedIndex]
			if !ok {
				break
			}

			chunkLen := uint32(len(ciphertext))
			if err := binary.Write(w, binary.BigEndian, chunkLen); err != nil {
				return fmt.Errorf("write chunk length: %w", err)
			}

			if _, err := w.Write(ciphertext); err != nil {
				return fmt.Errorf("write ciphertext: %w", err)
			}

			globalHMAC.Write(ciphertext)
			delete(pending, expectedIndex)
			expectedIndex++
		}
	}

	return nil
}
