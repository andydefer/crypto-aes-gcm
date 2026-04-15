// Package cli provides the command-line interface for aescryptool.
//
// It implements the decrypt command for decrypting files encrypted with aescryptool.
// The command handles argument parsing, validation, and delegates the actual
// decryption work to the service layer.
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

// NewDecryptCmd creates and configures the decrypt command.
//
// The command expects two positional arguments:
//   - input: path to the encrypted file
//   - output: path where decrypted content will be written
//
// Flags:
//   - --pass, -p: passphrase for decryption (optional, prompts if omitted)
//   - --workers, -w: number of parallel workers (defaults to cryptolib.DefaultWorkers)
//   - --force, -f: overwrite output file without confirmation
//   - --quiet, -q: suppress progress output
//   - --lang: language for UI messages (en or fr)
//
// Returns:
//   - *cobra.Command: configured decrypt command ready for registration
func NewDecryptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "decrypt [input] [output]",
		Short: lang.T(lang.CmdDecryptShort),
		Long:  lang.T(lang.CmdDecryptLong),
		Args:  cobra.ExactArgs(2),
		RunE:  runDecrypt,
	}

	cmd.Flags().StringVarP(&GlobalConfig.Pass, "pass", "p", "", lang.T(lang.FlagPassDesc))
	cmd.Flags().IntVarP(&GlobalConfig.Workers, "workers", "w", cryptolib.DefaultWorkers(), lang.T(lang.FlagWorkersDesc))
	cmd.Flags().BoolVarP(&GlobalConfig.Force, "force", "f", false, lang.T(lang.FlagForceDesc))
	cmd.Flags().BoolVarP(&GlobalConfig.Quiet, "quiet", "q", false, lang.T(lang.FlagQuietDesc))

	return cmd
}

// runDecrypt executes the decryption operation.
//
// Parameters:
//   - cmd: cobra command providing error output streams
//   - args: positional arguments (input file path, output file path)
//
// Returns:
//   - error: nil on success, or an error describing the failure
//
// The function performs the following steps:
//  1. Applies language settings for UI messages
//  2. Validates input file existence
//  3. Checks output file overwrite conditions
//  4. Resolves the passphrase (flag or interactive prompt)
//  5. Delegates decryption to the service layer
//
// Error conditions:
//   - Input file does not exist
//   - Output file exists and user declines overwrite
//   - Password resolution fails
//   - Decryption operation fails (wrong password, corrupted file, etc.)
func runDecrypt(cmd *cobra.Command, args []string) error {
	applyLanguage(GlobalConfig.Lang)

	input := args[0]
	output := args[1]

	// Validate input file exists
	if err := service.ValidateInputFile(input); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err
	}

	// Check for existing output file with interactive confirmation
	err := service.CheckOverwrite(output, GlobalConfig.Force)
	if err == service.ErrFileExists && !GlobalConfig.Force {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf(lang.T(lang.CliFileExists), output),
			IsConfirm: true,
			Default:   "n",
		}
		result, errPrompt := prompt.Run()
		if errPrompt != nil || (result != "y" && result != "Y") {
			ui.InfoColor.Println(lang.T(lang.CliOperationCancelled))
			return nil
		}
	} else if err != nil && err != service.ErrFileExists {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err
	}

	// Resolve password (prompt if not provided via flag)
	password, err := ResolvePassword(GlobalConfig.Pass, false)
	if err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err
	}

	// Execute decryption
	if err := service.ExecuteDecryption(input, output, password, GlobalConfig.Quiet); err != nil {
		ui.ErrorColor.Fprintf(cmd.ErrOrStderr(), lang.T(lang.CliError), err)
		return err
	}

	return nil
}
