package cryptolib

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/argon2"
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

	return d.Decrypt(in, out)
}

// Decrypt reads encrypted data from r and writes the plaintext to w.
func (d *Decryptor) Decrypt(r io.Reader, w io.Writer) error {
	headerData, baseNonce, err := d.readAndVerifyHeader(r)
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

	return d.processDecryption(r, w, gcm, baseNonce)
}

func (d *Decryptor) readAndVerifyHeader(r io.Reader) (FileHeader, []byte, error) {
	var headerData FileHeader
	if err := binary.Read(r, binary.BigEndian, &headerData); err != nil {
		return FileHeader{}, nil, fmt.Errorf("read header: %w", err)
	}

	if string(headerData.Magic[:]) != Magic {
		return FileHeader{}, nil, ErrInvalidMagic
	}

	if headerData.Version != Version {
		return FileHeader{}, nil, ErrUnsupportedVersion
	}

	storedHMAC := make([]byte, 32)
	if _, err := io.ReadFull(r, storedHMAC); err != nil {
		return FileHeader{}, nil, fmt.Errorf("read header HMAC: %w", err)
	}

	if !header.VerifyHMAC(d.key, header.ToBytes(
		headerData.Magic,
		headerData.Version,
		headerData.Salt,
		headerData.ChunkSize,
	), storedHMAC) {
		return FileHeader{}, nil, ErrHeaderAuthFailed
	}

	baseNonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(r, baseNonce); err != nil {
		return FileHeader{}, nil, fmt.Errorf("read nonce: %w", err)
	}

	return headerData, baseNonce, nil
}

func (d *Decryptor) processDecryption(r io.Reader, w io.Writer, gcm cipher.AEAD, baseNonce []byte) error {
	ciphertexts, globalHMAC, err := d.readCiphertexts(r)
	if err != nil {
		return err
	}

	if err := d.verifyGlobalHMAC(r, globalHMAC); err != nil {
		return err
	}

	return d.decryptAndWrite(ciphertexts, w, gcm, baseNonce)
}

func (d *Decryptor) readCiphertexts(r io.Reader) ([][]byte, hash.Hash, error) {
	var ciphertexts [][]byte
	globalHMAC := hmac.New(sha256.New, d.key)

	for {
		var chunkLen uint32
		if err := binary.Read(r, binary.BigEndian, &chunkLen); err != nil {
			if err == io.EOF {
				return nil, nil, fmt.Errorf("unexpected EOF: missing end marker")
			}
			return nil, nil, fmt.Errorf("read chunk length: %w", err)
		}

		if chunkLen == 0 {
			break
		}

		ciphertext := make([]byte, chunkLen)
		if _, err := io.ReadFull(r, ciphertext); err != nil {
			return nil, nil, fmt.Errorf("read ciphertext: %w", err)
		}
		ciphertexts = append(ciphertexts, ciphertext)
		globalHMAC.Write(ciphertext)
	}

	return ciphertexts, globalHMAC, nil
}

func (d *Decryptor) verifyGlobalHMAC(r io.Reader, computedHMAC hash.Hash) error {
	stored := make([]byte, 32)
	if _, err := io.ReadFull(r, stored); err != nil {
		return fmt.Errorf("read global HMAC: %w", err)
	}

	if !hmac.Equal(stored, computedHMAC.Sum(nil)) {
		return ErrGlobalHMACFailed
	}

	return nil
}

func (d *Decryptor) decryptAndWrite(ciphertexts [][]byte, w io.Writer, gcm cipher.AEAD, baseNonce []byte) error {
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
