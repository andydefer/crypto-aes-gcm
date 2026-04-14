// Package cli provides the command-line interface for aescryptool.
//
// It implements the encrypt command for encrypting files using AES-256-GCM.
// The command handles argument parsing, validation, and delegates the actual
// encryption work to the service layer.
package cli

import (
	"fmt"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// NewEncryptCmd creates and configures the encrypt command.
//
// The command expects two positional arguments:
//   - input: path to the plaintext source file
//   - output: path where encrypted data will be written
//
// Flags:
//   - --pass, -p: passphrase for encryption (optional, prompts if omitted)
//   - --workers, -w: number of parallel workers for chunk encryption
//   - --force, -f: overwrite existing output file without confirmation
//   - --quiet, -q: suppress progress bar output
//
// Returns:
//   - *cobra.Command: configured Cobra command ready for registration
func NewEncryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "encrypt [input] [output]",
		Short: lang.T(lang.CmdEncryptShort),
		Long:  lang.T(lang.CmdEncryptLong),
		Args:  cobra.ExactArgs(2),
		RunE:  runEncrypt,
	}

	cmd.Flags().StringVarP(&GlobalConfig.Pass, "pass", "p", "", lang.T(lang.FlagPassDesc))
	cmd.Flags().IntVarP(&GlobalConfig.Workers, "workers", "w", cryptolib.DefaultWorkers(), lang.T(lang.FlagWorkersDesc))
	cmd.Flags().BoolVarP(&GlobalConfig.Force, "force", "f", false, lang.T(lang.FlagForceDesc))
	cmd.Flags().BoolVarP(&GlobalConfig.Quiet, "quiet", "q", false, lang.T(lang.FlagQuietDesc))

	return cmd
}

// runEncrypt executes the encryption operation.
//
// It validates the input file, checks for output file conflicts, validates
// the worker count, resolves the password, and delegates encryption to the service layer.
//
// Parameters:
//   - cmd: the Cobra command (provides stderr output)
//   - args: command arguments containing input and output file paths
//
// Returns:
//   - error: any error encountered during encryption, or nil on success
func runEncrypt(cmd *cobra.Command, args []string) error {
	applyLanguage(GlobalConfig.Lang)

	input := args[0]
	output := args[1]

	workerCount := service.ValidateWorkerCount(GlobalConfig.Workers, GlobalConfig.Quiet)

	if err := service.ValidateInputFile(input); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err
	}

	if err := handleOutputOverwrite(output, cmd); err != nil {
		return err
	}

	password, err := ResolvePassword(GlobalConfig.Pass, true)
	if err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err
	}

	if err := service.ExecuteEncryption(input, output, password, workerCount, GlobalConfig.Quiet); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err
	}

	return nil
}

// handleOutputOverwrite checks if the output file exists and handles overwrite confirmation.
//
// Parameters:
//   - output: path to the output file
//   - cmd: the Cobra command for error output
//
// Returns:
//   - error: if the operation is cancelled or an unexpected error occurs
func handleOutputOverwrite(output string, cmd *cobra.Command) error {
	err := service.CheckOverwrite(output, GlobalConfig.Force)

	switch {
	case err == service.ErrFileExists && !GlobalConfig.Force:
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf(lang.T(lang.CliFileExists), output),
			IsConfirm: true,
			Default:   "n",
		}
		result, promptErr := prompt.Run()
		if promptErr != nil || (result != "y" && result != "Y") {
			ui.InfoColor.Println(lang.T(lang.CliOperationCancelled))
			return nil // User cancelled, no error
		}
		return nil

	case err != nil && err != service.ErrFileExists:
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err

	default:
		return nil
	}
}
