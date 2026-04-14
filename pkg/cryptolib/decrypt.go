// Package cryptolib provides AES-256-GCM file encryption and decryption.
//
// This package implements secure file encryption using AES-256-GCM in counter mode
// with Argon2id key derivation. It supports streaming decryption for large files
// and includes integrity verification through HMAC.
//
// The decryption process:
//   - Reads and validates file header with HMAC authentication
//   - Derives encryption key using Argon2id with salt from header
//   - Streams and decrypts chunks using derived nonces
//   - Verifies authenticity of each chunk via GCM
package cryptolib

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/argon2"
	"github.com/andydefer/crypto-aes-gcm/internal/crypto"
	"github.com/andydefer/crypto-aes-gcm/internal/header"
)

// Decryptor handles decryption of data encrypted with Encryptor.
type Decryptor struct {
	key       []byte
	chunkSize int
}

// NewDecryptor creates a Decryptor using the provided passphrase and salt.
func NewDecryptor(passphrase string, salt []byte) (*Decryptor, error) {
	params := argon2.DefaultParams()
	key := argon2.DeriveKey(passphrase, salt, params)

	return &Decryptor{
		key:       key,
		chunkSize: DefaultChunkSize,
	}, nil
}

// DecryptFile decrypts a file at inputPath and writes the result to outputPath.
func (d *Decryptor) DecryptFile(inputPath, outputPath string) error {
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

	return d.Decrypt(input, output)
}

// Decrypt reads encrypted data from r and writes the plaintext to w.
func (d *Decryptor) Decrypt(reader io.Reader, writer io.Writer) error {
	headerData, baseNonce, err := d.readAndVerifyHeader(reader)
	if err != nil {
		return err
	}

	d.chunkSize = int(headerData.ChunkSize)

	block, err := aes.NewCipher(d.key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	return d.processDecryption(reader, writer, gcm, baseNonce)
}

// readAndVerifyHeader extracts and validates the encrypted file header.
func (d *Decryptor) readAndVerifyHeader(reader io.Reader) (FileHeader, []byte, error) {
	var headerData FileHeader
	if err := binary.Read(reader, binary.BigEndian, &headerData); err != nil {
		return FileHeader{}, nil, fmt.Errorf("read header: %w", err)
	}

	if string(headerData.Magic[:]) != Magic {
		return FileHeader{}, nil, ErrInvalidMagic
	}

	if headerData.Version != Version {
		return FileHeader{}, nil, ErrUnsupportedVersion
	}

	storedHMAC := make([]byte, 32)
	if _, err := io.ReadFull(reader, storedHMAC); err != nil {
		return FileHeader{}, nil, fmt.Errorf("read header HMAC: %w", err)
	}

	serialized := header.Serialize(
		headerData.Magic,
		headerData.Version,
		headerData.Salt,
		headerData.ChunkSize,
	)

	if !header.VerifyHMAC(d.key, serialized, storedHMAC) {
		return FileHeader{}, nil, ErrHeaderAuthFailed
	}

	baseNonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(reader, baseNonce); err != nil {
		return FileHeader{}, nil, fmt.Errorf("read nonce: %w", err)
	}

	return headerData, baseNonce, nil
}

// processDecryption streams chunks from the reader, decrypts them, and writes to the writer.
func (d *Decryptor) processDecryption(reader io.Reader, writer io.Writer, gcm cipher.AEAD, baseNonce []byte) error {
	var chunkIndex uint64

	for {
		var chunkLen uint32
		err := binary.Read(reader, binary.BigEndian, &chunkLen)

		if errors.Is(err, io.EOF) {
			return fmt.Errorf("unexpected EOF: missing end marker")
		}
		if err != nil {
			return fmt.Errorf("read chunk length: %w", err)
		}

		if chunkLen == 0 {
			break
		}

		ciphertext := make([]byte, chunkLen)
		if _, err := io.ReadFull(reader, ciphertext); err != nil {
			return fmt.Errorf("read ciphertext chunk %d: %w", chunkIndex, err)
		}

		nonce := crypto.DeriveChunkNonce(baseNonce, chunkIndex)

		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return fmt.Errorf("%w chunk %d: %w", ErrDecryptionFailed, chunkIndex, err)
		}

		if _, err := writer.Write(plaintext); err != nil {
			return fmt.Errorf("write plaintext chunk %d: %w", chunkIndex, err)
		}

		chunkIndex++
	}

	return nil
}
