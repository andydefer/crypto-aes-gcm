// Package service provides business logic for encryption and decryption operations.
//
// This package orchestrates the encryption and decryption processes by:
//   - Managing progress bars and user feedback
//   - Handling file I/O operations
//   - Coordinating with the cryptolib package for core crypto operations
//   - Providing error handling and cleanup
package service

import (
	"fmt"
	"io"
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

// ExecuteEncryption performs the encryption operation with progress feedback.
//
// This function orchestrates the complete encryption workflow:
//  1. Reads the source file size for progress tracking
//  2. Initializes a progress bar (unless quiet mode is enabled)
//  3. Creates an encryptor with the specified number of parallel workers
//  4. Streams the source file through the encryptor to the output file
//  5. Displays success information upon completion
//
// Parameters:
//   - input: Path to the plaintext source file
//   - output: Path where encrypted data will be written
//   - password: Passphrase used for encryption (will be derived via Argon2id)
//   - workerCount: Number of parallel workers for chunk encryption (clamped to 1-2×CPU)
//   - quiet: If true, suppresses progress bar output
//
// Returns:
//   - error: Any error encountered during file operations, encryptor creation,
//     or the encryption process itself
func ExecuteEncryption(input, output, password string, workerCount int, quiet bool) error {
	fileInfo, err := os.Stat(input)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	var bar ui.ProgressBar
	if !quiet {
		bar = ui.CreateProgressBar(fileSize, "🔒 Encrypting")
	} else {
		bar = &noopProgressBar{}
	}

	// Create encryptor with default configuration
	encryptor, err := cryptolib.NewEncryptor(workerCount)
	if err != nil {
		return err
	}

	inputFile, err := openFile(input)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader := &progressReader{
		r:     inputFile,
		bar:   bar,
		total: fileSize,
	}

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := encryptor.Encrypt(reader, outFile, password); err != nil {
		_ = bar.Clear()
		return err
	}

	_ = bar.Finish()
	ui.PrintSuccess(output, fileSize)
	return nil
}

// ExecuteEncryptionWithConfig performs encryption with custom encryptor configuration.
//
// This function allows fine-grained control over the encryption parameters including
// chunk size and maximum pending chunks limit.
//
// Parameters:
//   - input: Path to the plaintext source file
//   - output: Path where encrypted data will be written
//   - password: Passphrase used for encryption
//   - config: Encryptor configuration (workers, chunk size, max pending chunks)
//   - quiet: If true, suppresses progress bar output
//
// Returns:
//   - error: Any error encountered during encryption
func ExecuteEncryptionWithConfig(input, output, password string, config cryptolib.EncryptorConfig, quiet bool) error {
	fileInfo, err := os.Stat(input)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	var bar ui.ProgressBar
	if !quiet {
		bar = ui.CreateProgressBar(fileSize, "🔒 Encrypting")
	} else {
		bar = &noopProgressBar{}
	}

	encryptor, err := cryptolib.NewEncryptorWithConfig(config)
	if err != nil {
		return err
	}

	inputFile, err := openFile(input)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	reader := &progressReader{
		r:     inputFile,
		bar:   bar,
		total: fileSize,
	}

	outFile, err := os.Create(output)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if err := encryptor.Encrypt(reader, outFile, password); err != nil {
		_ = bar.Clear()
		return err
	}

	_ = bar.Finish()
	ui.PrintSuccess(output, fileSize)
	return nil
}

// progressReader wraps an io.ReadCloser to track reading progress.
//
// This struct intercepts Read calls to update a progress bar, providing
// visual feedback during long-running file operations.
type progressReader struct {
	r     io.ReadCloser  // Underlying reader being wrapped
	bar   ui.ProgressBar // Progress bar to update
	total int64          // Total bytes expected to read
	read  int64          // Bytes read so far
}

// Read reads data from the underlying reader and updates the progress bar.
//
// Parameters:
//   - p: Byte slice to read data into
//
// Returns:
//   - n: Number of bytes read
//   - err: Error from underlying reader (typically io.EOF at end of file)
func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)
	pr.read += int64(n)
	_ = pr.bar.Set64(pr.read)
	return n, err
}

// Close closes the underlying reader.
//
// Returns:
//   - error: Any error from closing the underlying reader
func (pr *progressReader) Close() error {
	return pr.r.Close()
}

// noopProgressBar is a progress bar implementation that does nothing.
//
// This is used when quiet mode is enabled to avoid nil pointer checks.
type noopProgressBar struct{}

// Set64 implements ProgressBar.Set64 but does nothing.
func (n *noopProgressBar) Set64(int64) error { return nil }

// Finish implements ProgressBar.Finish but does nothing.
func (n *noopProgressBar) Finish() error { return nil }

// Clear implements ProgressBar.Clear but does nothing.
func (n *noopProgressBar) Clear() error { return nil }

// openFile opens a file and returns an error on failure.
//
// Parameters:
//   - path: Path to the file to open
//
// Returns:
//   - io.ReadCloser: Opened file handle (caller must close)
//   - error: Any error encountered during opening
func openFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ouverture du fichier '%s': %w", path, err)
	}
	return f, nil
}
