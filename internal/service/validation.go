// Package service provides business logic for encryption and decryption operations.
//
// This package orchestrates the encryption and decryption processes by:
//   - Managing progress bars and user feedback
//   - Handling file I/O operations
//   - Coordinating with the cryptolib package for core crypto operations
//   - Providing error handling and cleanup
package service

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

// ErrFileExists indicates that the output file already exists and force is false.
var ErrFileExists = errors.New(lang.T(lang.ErrFileAlreadyExists))

// ValidateWorkerCount ensures the worker count is within reasonable bounds.
//
// This function clamps the requested worker count to safe values:
//   - Values <= 0 return the default worker count
//   - Values exceeding 2×CPU cores are capped at that maximum
//
// A warning is printed to stderr when the value is capped, unless quiet mode is enabled.
//
// Parameters:
//   - requested: The desired number of parallel workers
//   - quiet: If true, suppresses warning messages when capping worker count
//
// Returns:
//   - int: A valid worker count between 1 and 2×runtime.NumCPU()
func ValidateWorkerCount(requested int, quiet bool) int {
	if requested <= 0 {
		return cryptolib.DefaultWorkers()
	}
	maxWorkers := runtime.NumCPU() * 2
	if requested > maxWorkers {
		if !quiet {
			ui.WarningColor.Printf(lang.T(lang.WarnWorkersReduced), maxWorkers)
		}
		return maxWorkers
	}
	return requested
}

// ValidateInputFile checks if the input file exists and is accessible.
//
// Parameters:
//   - path: Filesystem path to the input file
//
// Returns:
//   - error: nil if the file exists, otherwise an error describing the issue
func ValidateInputFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("%s", fmt.Sprintf(lang.T(lang.ErrFileNotFound), path))
	}
	return nil
}

// CheckFileExists determines whether a file exists at the given path.
//
// This function distinguishes between a file that doesn't exist (returns false, nil)
// and other errors like permission denied (returns false, err).
//
// Parameters:
//   - path: Filesystem path to check
//
// Returns:
//   - bool: true if the file exists, false otherwise
//   - error: Any error encountered during stat (except IsNotExist)
func CheckFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckOverwrite checks if output file exists and returns ErrFileExists if it does
// and force is false. This function is non-interactive and returns an error that
// the caller (CLI layer) can handle appropriately.
//
// Parameters:
//   - output: Path to the output file that may already exist
//   - force: If true, overwrites without checking
//
// Returns:
//   - error: nil if overwrite is allowed, ErrFileExists if file exists and force is false
func CheckOverwrite(output string, force bool) error {
	if force {
		return nil
	}
	if _, err := os.Stat(output); err == nil {
		return ErrFileExists
	}
	return nil
}
