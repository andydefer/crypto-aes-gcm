// Package cli provides the interactive command for aescryptool.
//
// The interactive mode guides users through encryption and decryption operations
// with step-by-step prompts, real-time validation, and visual feedback.
package cli

import (
	"fmt"
	"strings"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// NewInteractCmd creates the interactive command.
//
// Returns:
//   - *cobra.Command: configured interactive command
func NewInteractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interact",
		Short: "🎮 Interactive mode",
		Long:  lang.T(lang.CmdEncryptLong),
		RunE:  runInteractive,
	}
}

// runInteractive is the main entry point for interactive mode.
//
// It displays the welcome header and enters a loop that repeatedly presents
// the operation menu until the user chooses to exit.
//
// Parameters:
//   - cmd: the Cobra command (unused, kept for interface compliance)
//   - args: command arguments (unused)
//
// Returns:
//   - error: always nil (interactive mode exits via normal flow)
func runInteractive(cmd *cobra.Command, args []string) error {
	applyLanguage(GlobalConfig.Lang)

	ui.PrintInteractiveHeader()

	for {
		switch ui.PromptOperation() {
		case "encrypt":
			fmt.Println()
			runInteractiveEncrypt()
		case "decrypt":
			fmt.Println()
			runInteractiveDecrypt()
		case "exit":
			ui.PrintInteractiveGoodbye()
			return nil
		}
	}
}

// runInteractiveEncrypt guides the user through the encryption process.
//
// The flow:
//  1. Prompt for source file path
//  2. Prompt for output file path (defaults to input + ".enc")
//  3. Prompt for password with strength validation and confirmation
//  4. Prompt for worker count
//  5. Check if output file exists and ask for overwrite confirmation
//  6. Execute encryption with progress bar
//  7. Wait for user to press Enter before returning to menu
func runInteractiveEncrypt() {
	ui.PrintEncryptHeader()

	input := promptInputFile()
	if input == "" {
		return
	}

	output := promptOutputFile(input, ".enc")
	if output == "" {
		return
	}

	password := promptAndConfirmPassword()
	if password == "" {
		return
	}

	workerCount := ui.PromptWorkers()
	fmt.Println()

	if err := checkAndConfirmOverwrite(output); err != nil {
		ui.ErrorColor.Printf(lang.T(lang.CliError), err)
		fmt.Println()
		waitForUser()
		return
	}

	if err := service.ExecuteEncryption(input, output, password, workerCount, false); err != nil {
		ui.ErrorColor.Printf(lang.T(lang.CliError), err)
	}

	fmt.Println()
	waitForUser()
}

// runInteractiveDecrypt guides the user through the decryption process.
//
// The flow:
//  1. Prompt for encrypted source file path
//  2. Prompt for output file path (defaults to input without ".enc" or input + ".dec")
//  3. Prompt for password (no confirmation needed)
//  4. Check if output file exists and ask for overwrite confirmation
//  5. Execute decryption with progress bar
//  6. Wait for user to press Enter before returning to menu
func runInteractiveDecrypt() {
	ui.PrintDecryptHeader()

	input := ui.PromptFilePath(lang.T(lang.InteractiveEncryptedFile), true, "")
	if input == "" {
		return
	}
	fmt.Println()

	defaultOutput := strings.TrimSuffix(input, ".enc")
	if defaultOutput == input {
		defaultOutput = input + ".dec"
	}
	output := ui.PromptFilePath(lang.T(lang.InteractiveOutputFile), false, defaultOutput)
	if output == "" {
		output = defaultOutput
	}
	fmt.Println()

	password := ui.PromptPassword(lang.T(lang.InteractivePassword), false)
	if password == "" {
		return
	}
	fmt.Println()

	if err := checkAndConfirmOverwrite(output); err != nil {
		ui.ErrorColor.Printf(lang.T(lang.CliError), err)
		fmt.Println()
		waitForUser()
		return
	}

	if err := service.ExecuteDecryption(input, output, password, false); err != nil {
		ui.ErrorColor.Printf(lang.T(lang.CliError), err)
	}

	fmt.Println()
	waitForUser()
}

// promptInputFile prompts the user for the source file path.
//
// Returns:
//   - string: the input file path, or empty string if cancelled
func promptInputFile() string {
	input := ui.PromptFilePath(lang.T(lang.InteractiveFileToEncrypt), true, "")
	if input == "" {
		return ""
	}
	fmt.Println()
	return input
}

// promptOutputFile prompts for output file path with a default value.
//
// Parameters:
//   - input: source file path used to generate default output name
//   - extension: file extension to append for default output (e.g., ".enc")
//
// Returns:
//   - string: the output file path
func promptOutputFile(input, extension string) string {
	defaultOutput := input + extension
	output := ui.PromptFilePath(lang.T(lang.InteractiveOutputFile), false, defaultOutput)
	if output == "" {
		output = defaultOutput
	}
	fmt.Println()
	return output
}

// promptAndConfirmPassword handles password entry with confirmation.
//
// Prompts for password, validates strength, asks for confirmation,
// and verifies that both entries match.
//
// Returns:
//   - string: the confirmed password, or empty string if validation fails
func promptAndConfirmPassword() string {
	password := ui.PromptPassword(lang.T(lang.InteractivePassword), true)
	if password == "" {
		return ""
	}
	fmt.Println()

	confirm := ui.PromptPassword(lang.T(lang.InteractiveConfirm), true)
	if confirm == "" {
		return ""
	}
	fmt.Println()

	if password != confirm {
		ui.ErrorColor.Println(lang.T(lang.InteractivePasswordsNotMatch))
		fmt.Println()
		return ""
	}

	return password
}

// checkAndConfirmOverwrite checks if a file exists and prompts for overwrite confirmation.
//
// Parameters:
//   - output: path to the file to check
//
// Returns:
//   - error: if the file exists and user declines overwrite, or if the check fails
func checkAndConfirmOverwrite(output string) error {
	exists, err := service.CheckFileExists(output)
	if err != nil {
		return fmt.Errorf("check file existence: %w", err)
	}

	if exists {
		if !ui.PromptConfirm(lang.T(lang.InteractiveOverwrite), true) {
			ui.InfoColor.Println(lang.T(lang.InteractiveCancel))
			return fmt.Errorf("user cancelled overwrite")
		}
		fmt.Println()
	}

	return nil
}

// waitForUser pauses execution until the user presses Enter.
func waitForUser() {
	ui.InfoColor.Println(lang.T(lang.InteractivePressEnter))
	fmt.Scanln()
}
