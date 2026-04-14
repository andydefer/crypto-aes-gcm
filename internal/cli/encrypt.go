// Package cli provides the command-line interface for cryptool.
//
// It implements the encrypt command for encrypting files using AES-256-GCM.
// The command handles argument parsing, validation, and delegates the actual
// encryption work to the service layer.
package cli

import (
	"os"

	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/spf13/cobra"
)

// NewEncryptCmd creates and configures the encrypt command.
//
// The command expects two positional arguments:
//   - input: Path to the plaintext source file
//   - output: Path where encrypted data will be written
//
// Flags:
//   - --pass, -p: Passphrase for encryption (required)
//   - --workers, -w: Number of parallel workers (default: cryptolib.DefaultWorkers)
//   - --force, -f: Overwrite output file without confirmation
//   - --quiet, -q: Suppress progress output
//
// Returns:
//   - *cobra.Command: Configured Cobra command ready for registration
func NewEncryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt [input] [output]",
		Short: "🔒 Encrypt a file",
		Long: `Encrypt a file using AES-256-GCM with Argon2id key derivation.

The encryption process:
  1. Generates a random salt and nonce
  2. Derives a 256-bit key using Argon2id
  3. Splits the input into chunks (default 1MB)
  4. Encrypts chunks in parallel using the specified number of workers
  5. Writes header, HMAC, nonce, and encrypted chunks to the output file

Examples:
  cryptool encrypt secret.txt secret.enc --pass myPassword
  cryptool encrypt data.txt output.enc --pass secure123 --force
  cryptool encrypt large.bin result.enc --pass pass123 --workers 8 --quiet`,
		Args: cobra.ExactArgs(2),
		Run:  runEncrypt,
	}

	cmd.Flags().StringVarP(&pass, "pass", "p", "", "Passphrase for encryption (required)")
	cmd.Flags().IntVarP(&workers, "workers", "w", cryptolib.DefaultWorkers, "Number of parallel workers for chunk encryption")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing output file without confirmation")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress progress bar output")
	_ = cmd.MarkFlagRequired("pass")

	return cmd
}

// runEncrypt executes the encryption operation.
//
// It validates the input file, checks for output file conflicts, validates
// the worker count, and delegates the actual encryption to the service layer.
// On any error, it prints an error message and exits with code 1.
//
// Parameters:
//   - cmd: The Cobra command (provides stderr output)
//   - args: Command arguments containing input and output file paths
func runEncrypt(cmd *cobra.Command, args []string) {
	input := args[0]
	output := args[1]

	workerCount := service.ValidateWorkerCount(workers, quiet)

	// Validate input file exists
	if err := service.ValidateInputFile(input); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Check for existing output file
	if err := service.CheckOverwrite(output, force); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Execute encryption
	if err := service.ExecuteEncryption(input, output, pass, workerCount, quiet); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		os.Exit(1)
	}
}
