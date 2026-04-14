// Package cli provides the command-line interface for aescryptool.
//
// It implements the decrypt command for decrypting files encrypted with aescryptool.
// The command handles argument parsing, validation, and delegates the actual
// decryption work to the service layer.
package cli

import (
	"fmt"

	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// NewDecryptCmd creates and configures the decrypt command.
//
// The command expects two positional arguments:
//   - input: Path to the encrypted source file
//   - output: Path where decrypted plaintext will be written
//
// Flags:
//   - --pass, -p: Passphrase used for encryption (optional - will prompt if omitted)
//   - --workers, -w: Number of parallel workers (default: cryptolib.DefaultWorkers)
//   - --force, -f: Overwrite output file without confirmation
//   - --quiet, -q: Suppress progress output
//
// Returns:
//   - *cobra.Command: Configured Cobra command ready for registration
func NewDecryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt [input] [output]",
		Short: "🔓 Decrypt a file",
		Long: `Decrypt a file that was encrypted with the encrypt command.

The decryption process:
  1. Validates the input file exists
  2. Reads and verifies the file header
  3. Derives the encryption key using Argon2id with the salt from header
  4. Streams and decrypts the data to the output file
  5. Verifies integrity of each chunk via GCM authentication

Password can be provided via:
  - --pass flag (visible in process list, not recommended for shared environments)
  - Interactive prompt (recommended for manual use)

Examples:
  aescryptool decrypt secret.enc secret.txt              # Prompts for password
  aescryptool decrypt secret.enc secret.txt --pass myPassword
  aescryptool decrypt data.enc output.txt --pass secure123 --force
  aescryptool decrypt large.enc result.bin --workers 8 --quiet`,
		Args: cobra.ExactArgs(2),
		RunE: runDecrypt,
	}

	cmd.Flags().StringVarP(&pass, "pass", "p", "", "Passphrase used for encryption (optional - will prompt if omitted)")
	cmd.Flags().IntVarP(&workers, "workers", "w", cryptolib.DefaultWorkers, "Number of parallel workers")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing output file without confirmation")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress progress bar output")

	return cmd
}

// runDecrypt executes the decryption operation.
//
// It performs validation steps and delegates the actual work to the service layer.
// On any error, it prints an error message and returns the error.
//
// Parameters:
//   - cmd: The Cobra command (provides stderr output)
//   - args: Command arguments containing input and output file paths
//
// Returns:
//   - error: Any error encountered during decryption, or nil on success
func runDecrypt(cmd *cobra.Command, args []string) error {
	input := args[0]
	output := args[1]

	// Validate input file exists
	if err := service.ValidateInputFile(input); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return err
	}

	// Check for existing output file with interactive confirmation
	err := service.CheckOverwrite(output, force)
	if err == service.ErrFileExists && !force {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Fichier '%s' existe déjà. Écraser ?", output),
			IsConfirm: true,
			Default:   "n",
		}
		result, errPrompt := prompt.Run()
		if errPrompt != nil || (result != "y" && result != "Y") {
			ui.InfoColor.Println("❌ Opération annulée")
			return nil
		}
	} else if err != nil && err != service.ErrFileExists {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return err
	}

	// Resolve password (prompt if not provided via flag)
	// For decryption, needConfirmation=false (single prompt, no confirmation)
	password, err := resolvePassword(pass, false)
	if err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return err
	}

	// Execute decryption
	if err := service.ExecuteDecryption(input, output, password, quiet); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return err
	}

	return nil
}
