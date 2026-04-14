// Package service provides business logic for encryption and decryption operations.
//
// This package orchestrates the encryption and decryption processes by:
//   - Managing progress bars and user feedback
//   - Handling file I/O operations
//   - Coordinating with the cryptolib package for core crypto operations
//   - Providing error handling and cleanup
package service

import (
	"encoding/binary"
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

// ExecuteDecryption performs the decryption operation with progress feedback.
//
// This function orchestrates the complete decryption workflow:
//  1. Reads the source file size for progress tracking
//  2. Initializes a progress bar (unless quiet mode is enabled)
//  3. Opens the encrypted file and reads its header
//  4. Creates a decryptor with the provided password and extracted salt
//  5. Streams the decrypted data to the output file
//  6. Displays success information upon completion
//
// Parameters:
//   - input: Path to the encrypted source file
//   - output: Path where decrypted plaintext will be written
//   - password: Passphrase used for encryption (must match original)
//   - quiet: If true, suppresses progress bar output
//
// Returns:
//   - error: Any error encountered during file operations, header parsing,
//     decryptor creation, or the decryption process itself
func ExecuteDecryption(input, output, password string, quiet bool) error {
	fileInfo, err := os.Stat(input)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	var bar ui.ProgressBar
	if !quiet {
		bar = ui.CreateProgressBar(fileSize, "🔓 Decrypting")
	} else {
		bar = &noopProgressBar{}
	}

	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()

	var header cryptolib.FileHeader
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		return err
	}

	decryptor, err := cryptolib.NewDecryptor(password, header.Salt[:])
	if err != nil {
		return err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}

	reader := &progressReader{
		r:     f,
		bar:   bar,
		total: fileSize,
	}

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := decryptor.Decrypt(reader, outFile); err != nil {
		_ = bar.Clear()
		return err
	}

	_ = bar.Finish()
	ui.PrintSuccess(output, fileSize)
	return nil
}
