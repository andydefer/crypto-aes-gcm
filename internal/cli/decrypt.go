// Package cli provides the command-line interface for cryptool.
//
// It implements Cobra commands for encryption, decryption, interactive mode,
// and version display. The package delegates business logic to the service layer
// and UI rendering to the ui package.
package cli

import (
	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/spf13/cobra"
)

// NewDecryptCmd creates the decrypt command.
//
// The command expects two arguments: input file (encrypted) and output file (plaintext).
// Flags:
//   - --pass, -p: Passphrase (required)
//   - --workers, -w: Number of parallel workers (default: DefaultWorkers)
//   - --force, -f: Force overwrite existing output file
//   - --quiet, -q: Suppress progress output
//
// Returns:
//   - *cobra.Command: Configured Cobra command
func NewDecryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt [input] [output]",
		Short: "🔓 Decrypt a file",
		Args:  cobra.ExactArgs(2),
		Run:   runDecrypt,
	}

	cmd.Flags().StringVarP(&pass, "pass", "p", "", "Passphrase (required)")
	cmd.Flags().IntVarP(&workers, "workers", "w", cryptolib.DefaultWorkers, "Parallel workers")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
	_ = cmd.MarkFlagRequired("pass")

	return cmd
}

// runDecrypt executes the decryption operation.
//
// It validates the input file, checks for output file conflicts, and delegates
// the actual decryption to the service layer.
//
// Parameters:
//   - cmd: The Cobra command (provides stderr output)
//   - args: Command arguments containing input and output file paths
func runDecrypt(cmd *cobra.Command, args []string) {
	input := args[0]
	output := args[1]

	if err := service.ValidateInputFile(input); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return
	}

	if err := service.CheckOverwrite(output, force); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return
	}

	if err := service.ExecuteDecryption(input, output, pass, quiet); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
	}
}
