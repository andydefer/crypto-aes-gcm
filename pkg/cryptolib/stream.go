package cryptolib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/andydefer/crypto-aes-gcm/internal/argon2"
	"github.com/andydefer/crypto-aes-gcm/internal/header"
)

// DecryptStream decrypts data from r and writes plaintext to w.
// This convenience function handles the entire decryption process including
// header parsing, HMAC verification, and chunk decryption.
func DecryptStream(r io.Reader, w io.Writer, passphrase string) error {
	var headerData FileHeader
	if err := binary.Read(r, binary.BigEndian, &headerData); err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	if string(headerData.Magic[:]) != Magic {
		return ErrInvalidMagic
	}

	if headerData.Version != Version {
		return ErrUnsupportedVersion
	}

	storedHMAC := make([]byte, 32)
	if _, err := io.ReadFull(r, storedHMAC); err != nil {
		return fmt.Errorf("read header HMAC: %w", err)
	}

	params := argon2.DefaultParams()
	key := argon2.DeriveKey(passphrase, headerData.Salt[:], params)

	if !header.VerifyHMAC(key, header.ToBytes(
		headerData.Magic,
		headerData.Version,
		headerData.Salt,
		headerData.ChunkSize,
	), storedHMAC) {
		return ErrHeaderAuthFailed
	}

	baseNonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(r, baseNonce); err != nil {
		return fmt.Errorf("read nonce: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("create GCM: %w", err)
	}

	var ciphertexts [][]byte
	globalHMAC := hmac.New(sha256.New, key)

	for {
		var chunkLen uint32
		if err := binary.Read(r, binary.BigEndian, &chunkLen); err != nil {
			if err == io.EOF {
				return fmt.Errorf("unexpected EOF: missing end marker")
			}
			return fmt.Errorf("read chunk length: %w", err)
		}

		if chunkLen == 0 {
			break
		}

		ciphertext := make([]byte, chunkLen)
		if _, err := io.ReadFull(r, ciphertext); err != nil {
			return fmt.Errorf("read ciphertext: %w", err)
		}
		ciphertexts = append(ciphertexts, ciphertext)
		globalHMAC.Write(ciphertext)
	}

	storedGlobalHMAC := make([]byte, 32)
	if _, err := io.ReadFull(r, storedGlobalHMAC); err != nil {
		return fmt.Errorf("read global HMAC: %w", err)
	}
	if !hmac.Equal(storedGlobalHMAC, globalHMAC.Sum(nil)) {
		return ErrGlobalHMACFailed
	}

	for idx, ciphertext := range ciphertexts {
		nonce := make([]byte, NonceSize)
		copy(nonce, baseNonce)
		binary.BigEndian.PutUint64(nonce[4:], uint64(idx))

		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return fmt.Errorf("%w (chunk %d): %v", ErrDecryptionFailed, idx, err)
		}

		if _, err := w.Write(plaintext); err != nil {
			return fmt.Errorf("write plaintext: %w", err)
		}
	}

	return nil
}
