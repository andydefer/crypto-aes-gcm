// Package cryptolib provides AES-256-GCM file encryption and decryption.
//
// This package implements secure file encryption using AES-256-GCM in counter mode
// with Argon2id key derivation and parallel streaming capabilities.
package cryptolib

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"github.com/andydefer/crypto-aes-gcm/internal/argon2"
	"github.com/andydefer/crypto-aes-gcm/internal/crypto"
	"github.com/andydefer/crypto-aes-gcm/internal/header"
)

// DecryptStream decrypts data from reader and writes plaintext to writer.
//
// This convenience function handles the entire decryption process including
// header parsing, HMAC verification, and streaming chunk decryption without
// loading the entire file into memory.
//
// Parameters:
//   - reader: Reader containing encrypted data in the expected format
//   - writer: Writer where decrypted plaintext will be written
//   - passphrase: Password used for encryption
//
// Returns:
//   - error: If header validation fails, decryption fails, or IO operations fail
func DecryptStream(reader io.Reader, writer io.Writer, passphrase string) error {
	_, key, baseNonce, err := readAndValidateHeader(reader, passphrase)
	if err != nil {
		return err
	}

	gcm, err := createGCMCipher(key)
	if err != nil {
		return err
	}

	return decryptChunks(reader, writer, gcm, baseNonce)
}

// readAndValidateHeader reads the file header, validates magic bytes and version,
// verifies HMAC, and derives the encryption key.
//
// Parameters:
//   - reader: Reader positioned at the start of an encrypted file
//   - passphrase: Password for key derivation
//
// Returns:
//   - FileHeader: Validated header data
//   - []byte: Derived encryption key
//   - []byte: Base nonce for chunk decryption
//   - error: If header reading or validation fails
func readAndValidateHeader(reader io.Reader, passphrase string) (FileHeader, []byte, []byte, error) {
	var headerData FileHeader
	if err := binary.Read(reader, binary.BigEndian, &headerData); err != nil {
		return FileHeader{}, nil, nil, fmt.Errorf("read header: %w", err)
	}

	if string(headerData.Magic[:]) != Magic {
		return FileHeader{}, nil, nil, ErrInvalidMagic
	}

	if headerData.Version != Version {
		return FileHeader{}, nil, nil, ErrUnsupportedVersion
	}

	storedHMAC := make([]byte, 32)
	if _, err := io.ReadFull(reader, storedHMAC); err != nil {
		return FileHeader{}, nil, nil, fmt.Errorf("read header HMAC: %w", err)
	}

	params := argon2.DefaultParams()
	key := argon2.DeriveKey(passphrase, headerData.Salt[:], params)

	serialized := header.Serialize(
		headerData.Magic,
		headerData.Version,
		headerData.Salt,
		headerData.ChunkSize,
	)

	if !header.VerifyHMAC(key, serialized, storedHMAC) {
		return FileHeader{}, nil, nil, ErrHeaderAuthFailed
	}

	baseNonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(reader, baseNonce); err != nil {
		return FileHeader{}, nil, nil, fmt.Errorf("read nonce: %w", err)
	}

	return headerData, key, baseNonce, nil
}

// createGCMCipher creates an AES-GCM cipher from the derived key.
//
// Parameters:
//   - key: Cryptographic key derived from passphrase and salt
//
// Returns:
//   - cipher.AEAD: GCM cipher for authenticated decryption
//   - error: If cipher creation fails
func createGCMCipher(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	return gcm, nil
}

// decryptChunks reads and decrypts chunks from the reader, writing plaintext to the writer.
//
// The decryption process:
//  1. Reads chunk length (uint32) - zero indicates end of stream
//  2. Reads ciphertext of specified length
//  3. Derives chunk-specific nonce by XORing base nonce with chunk index
//  4. Decrypts using GCM which also verifies authentication
//  5. Writes plaintext immediately to minimize memory usage
//
// Parameters:
//   - reader: Reader positioned after the header (at first chunk length)
//   - writer: Writer for decrypted plaintext
//   - gcm: GCM cipher for authenticated decryption
//   - baseNonce: Base nonce from file header
//
// Returns:
//   - error: If chunk reading, decryption, or writing fails
func decryptChunks(reader io.Reader, writer io.Writer, gcm cipher.AEAD, baseNonce []byte) error {
	var chunkIndex uint64
	var nonceBuf [crypto.NonceSize]byte

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

		if err := crypto.DeriveChunkNonceFast(nonceBuf[:], baseNonce, chunkIndex); err != nil {
			return fmt.Errorf("derive nonce for chunk %d: %w", chunkIndex, err)
		}

		plaintext, err := gcm.Open(nil, nonceBuf[:], ciphertext, nil)
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
