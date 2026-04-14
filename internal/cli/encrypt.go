// Package cli provides the command-line interface for cryptool.
//
// It implements the encrypt, decrypt, interact, and version commands using
// the Cobra CLI framework. The package delegates business logic to the
// service layer and UI rendering to the ui package.
package cli

import (
	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/spf13/cobra"
)

// NewEncryptCmd creates the encrypt command.
//
// The command accepts two positional arguments: input file path and output file path.
// It requires a passphrase via the --pass flag and supports optional flags for
// parallel workers, force overwrite, and quiet mode.
//
// Flags:
//   - --pass, -p: Passphrase for encryption (required)
//   - --workers, -w: Number of parallel workers (default: cryptolib.DefaultWorkers)
//   - --force, -f: Overwrite output file without confirmation
//   - --quiet, -q: Suppress progress output
//
// Returns:
//   - *cobra.Command: Configured Cobra command ready for registration.
func NewEncryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt [input] [output]",
		Short: "🔒 Encrypt a file",
		Args:  cobra.ExactArgs(2),
		Run:   runEncrypt,
	}

	cmd.Flags().StringVarP(&pass, "pass", "p", "", "Passphrase (required)")
	cmd.Flags().IntVarP(&workers, "workers", "w", cryptolib.DefaultWorkers, "Parallel workers")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite")
	cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress output")
	_ = cmd.MarkFlagRequired("pass")

	return cmd
}

// runEncrypt executes the encryption command.
//
// It validates the input file, checks for output file conflicts, validates
// the worker count, and delegates the actual encryption to the service layer.
//
// Parameters:
//   - cmd: The Cobra command being executed.
//   - args: Command-line arguments containing input and output paths.
func runEncrypt(cmd *cobra.Command, args []string) {
	input := args[0]
	output := args[1]

	workerCount := service.ValidateWorkerCount(workers, quiet)

	if err := service.ValidateInputFile(input); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return
	}

	if err := service.CheckOverwrite(output, force); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
		return
	}

	if err := service.ExecuteEncryption(input, output, pass, workerCount, quiet); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), "❌ Error: %v\n", err)
	}
}
