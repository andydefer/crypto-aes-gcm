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
	"github.com/andydefer/crypto-aes-gcm/internal/lang"
)

// Decryptor handles decryption of data encrypted with Encryptor.
// It manages the cryptographic state and streaming decryption process.
type Decryptor struct {
	key       []byte
	chunkSize int
}

// NewDecryptor creates a Decryptor using the provided passphrase and salt.
//
// Parameters:
//   - passphrase: The user's secret passphrase for key derivation
//   - salt: Cryptographic salt from the encrypted file header
//
// Returns:
//   - *Decryptor: Configured decryptor instance
//   - error: If key derivation fails or parameters are invalid
//
// The function derives an encryption key using Argon2id with default parameters.
func NewDecryptor(passphrase string, salt []byte) (*Decryptor, error) {
	params := argon2.DefaultParams()
	key := argon2.DeriveKey(passphrase, salt, params)

	return &Decryptor{
		key:       key,
		chunkSize: DefaultChunkSize,
	}, nil
}

// DecryptFile decrypts a file at inputPath and writes the result to outputPath.
//
// Parameters:
//   - inputPath: Path to the encrypted file
//   - outputPath: Path where decrypted content will be written
//
// Returns:
//   - error: Any error encountered during file operations or decryption
//
// This is a convenience wrapper around Decrypt that handles file opening/closing.
func (d *Decryptor) DecryptFile(inputPath, outputPath string) error {
	input, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrOpenInput), err)
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrCreateOutput), err)
	}
	defer output.Close()

	return d.Decrypt(input, output)
}

// Decrypt reads encrypted data from r and writes the plaintext to w.
//
// Parameters:
//   - reader: Source of encrypted data (typically a file or network stream)
//   - writer: Destination for decrypted plaintext
//
// Returns:
//   - error: If header validation, key setup, or decryption fails
//
// The function validates the file header, sets up the AES-GCM cipher,
// and processes the encrypted stream chunk by chunk.
func (d *Decryptor) Decrypt(reader io.Reader, writer io.Writer) error {
	headerData, baseNonce, err := d.readAndVerifyHeader(reader)
	if err != nil {
		return err
	}

	d.chunkSize = int(headerData.ChunkSize)

	block, err := aes.NewCipher(d.key)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrCreateCipher), err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf(lang.T(lang.CryptolibErrCreateGCM), err)
	}

	return d.processDecryption(reader, writer, gcm, baseNonce)
}

// readAndVerifyHeader extracts and validates the encrypted file header.
//
// Parameters:
//   - reader: Source containing the encrypted file header
//
// Returns:
//   - FileHeader: Parsed header information
//   - []byte: Base nonce for chunk counter derivation
//   - error: If header validation fails
//
// The function verifies magic bytes, version compatibility, and HMAC integrity.
func (d *Decryptor) readAndVerifyHeader(reader io.Reader) (FileHeader, []byte, error) {
	var headerData FileHeader
	if err := binary.Read(reader, binary.BigEndian, &headerData); err != nil {
		return FileHeader{}, nil, fmt.Errorf(lang.T(lang.CryptolibErrReadHeader), err)
	}

	if string(headerData.Magic[:]) != Magic {
		return FileHeader{}, nil, ErrInvalidMagic
	}

	if headerData.Version != Version {
		return FileHeader{}, nil, ErrUnsupportedVersion
	}

	storedHMAC := make([]byte, 32)
	if _, err := io.ReadFull(reader, storedHMAC); err != nil {
		return FileHeader{}, nil, fmt.Errorf(lang.T(lang.CryptolibErrReadHeaderHMAC), err)
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
		return FileHeader{}, nil, fmt.Errorf(lang.T(lang.CryptolibErrReadNonce), err)
	}

	return headerData, baseNonce, nil
}

// processDecryption streams chunks from the reader, decrypts them, and writes to the writer.
//
// Parameters:
//   - reader: Source of encrypted chunks
//   - writer: Destination for decrypted plaintext
//   - gcm: GCM cipher for authenticated decryption
//   - baseNonce: Base nonce for counter-based nonce derivation
//
// Returns:
//   - error: If chunk reading, decryption, or writing fails
//
// The function reads chunk length prefixes, decrypts authenticated ciphertext,
// and writes plaintext sequentially. A zero-length chunk marks the end of stream.
// Nonces are derived efficiently using a pre-allocated buffer to minimize allocations.
func (d *Decryptor) processDecryption(reader io.Reader, writer io.Writer, gcm cipher.AEAD, baseNonce []byte) error {
	var chunkIndex uint64

	var nonceBuf [crypto.NonceSize]byte

	for {
		var chunkLen uint32
		err := binary.Read(reader, binary.BigEndian, &chunkLen)

		if errors.Is(err, io.EOF) {
			return errors.New(lang.T(lang.CryptolibErrUnexpectedEOF))
		}
		if err != nil {
			return fmt.Errorf(lang.T(lang.CryptolibErrReadChunkLen), err)
		}

		if chunkLen == 0 {
			break
		}

		ciphertext := make([]byte, chunkLen)
		if _, err := io.ReadFull(reader, ciphertext); err != nil {
			return fmt.Errorf(lang.T(lang.CryptolibErrReadCiphertext), chunkIndex, err)
		}

		if err := crypto.DeriveChunkNonceFast(nonceBuf[:], baseNonce, chunkIndex); err != nil {
			return fmt.Errorf(lang.T(lang.CryptolibErrDeriveNonce), chunkIndex, err)
		}

		plaintext, err := gcm.Open(nil, nonceBuf[:], ciphertext, nil)
		if err != nil {
			return fmt.Errorf("%w chunk %d: %w", ErrDecryptionFailed, chunkIndex, err)
		}

		if _, err := writer.Write(plaintext); err != nil {
			return fmt.Errorf(lang.T(lang.CryptolibErrWritePlaintext), chunkIndex, err)
		}

		chunkIndex++
	}

	return nil
}
