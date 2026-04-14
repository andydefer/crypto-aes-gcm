package cli

import (
	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/spf13/cobra"
)

// NewDecryptCmd creates the decrypt command.
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
