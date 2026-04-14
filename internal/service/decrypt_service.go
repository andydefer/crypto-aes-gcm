// Package service provides business logic for encryption and decryption operations.
//
// This package orchestrates the encryption and decryption processes by:
//   - Managing progress bars and user feedback
//   - Handling file I/O operations
//   - Coordinating with the cryptolib package for core crypto operations
//   - Providing error handling and cleanup
package service

import (
	"context"
	"encoding/binary"
	"io"
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
// On decryption failure, the function ensures that any partially created output
// file is removed to avoid leaving empty or corrupted files on disk.
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
	return ExecuteDecryptionWithContext(context.Background(), input, output, password, quiet)
}

// ExecuteDecryptionWithContext performs decryption with context support for cancellation.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - input: Path to the encrypted source file
//   - output: Path where decrypted plaintext will be written
//   - password: Passphrase used for encryption (must match original)
//   - quiet: If true, suppresses progress bar output
//
// Returns:
//   - error: Any error encountered during decryption or context cancellation
//
// The function uses a goroutine to perform the actual decryption while listening
// for context cancellation. On cancellation, the operation is aborted and the
// partially created output file is removed.
func ExecuteDecryptionWithContext(ctx context.Context, input, output, password string, quiet bool) (err error) {
	fileInfo, err := os.Stat(input)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	bar := createProgressBar(quiet, fileSize)

	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer closeWithErrorHandling(f, &err)

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

	success := false
	defer func() {
		if closeErr := outFile.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
		if !success {
			_ = os.Remove(output)
		}
	}()

	if err := decryptWithContext(ctx, decryptor, reader, outFile, bar); err != nil {
		return err
	}

	success = true
	_ = bar.Finish()
	ui.PrintSuccess(output, fileSize)
	return nil
}

// createProgressBar initializes a progress bar for tracking decryption progress.
//
// Parameters:
//   - quiet: If true, returns a no-op progress bar that does nothing
//   - fileSize: Total size of the input file in bytes
//
// Returns:
//   - ui.ProgressBar: A progress bar implementation (either real or no-op)
func createProgressBar(quiet bool, fileSize int64) ui.ProgressBar {
	if quiet {
		return &noopProgressBar{}
	}
	return ui.CreateProgressBar(fileSize, "🔓 Decrypting")
}

// decryptWithContext performs streaming decryption with context cancellation support.
//
// Parameters:
//   - ctx: Context for cancellation control
//   - decryptor: Decryptor instance to perform the decryption
//   - reader: Source of encrypted data
//   - writer: Destination for decrypted data
//   - bar: Progress bar for user feedback
//
// Returns:
//   - error: Decryption error or context cancellation error
//
// The function runs decryption in a goroutine and listens for context cancellation.
// If the context is cancelled, the decryption is aborted and the progress bar is cleared.
func decryptWithContext(ctx context.Context, decryptor *cryptolib.Decryptor, reader io.Reader, writer io.Writer, bar ui.ProgressBar) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- decryptor.Decrypt(reader, writer)
	}()

	select {
	case <-ctx.Done():
		_ = bar.Clear()
		return ctx.Err()
	case err := <-errChan:
		if err != nil {
			_ = bar.Clear()
			return err
		}
		return nil
	}
}

// closeWithErrorHandling closes a file and updates the error reference if closing fails.
//
// Parameters:
//   - f: File to close
//   - err: Pointer to error that will be updated if close operation fails and no error exists
//
// This helper function ensures proper error propagation during defer statements.
// If an error already exists, it takes precedence over the close error.
func closeWithErrorHandling(f *os.File, err *error) {
	if closeErr := f.Close(); closeErr != nil && *err == nil {
		*err = closeErr
	}
}
