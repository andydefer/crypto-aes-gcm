// Package cli provides the interactive command for aescryptool.
//
// The interactive mode guides users through encryption and decryption operations
// with step-by-step prompts, real-time validation, and visual feedback.
package cli

import (
	"fmt"
	"strings"

	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// NewInteractCmd creates the interactive command.
//
// This command launches a user-friendly interactive shell that prompts for
// all necessary inputs (file paths, passwords, options) with real-time
// validation and visual feedback.
//
// Returns:
//   - *cobra.Command: Configured Cobra command for interactive mode
func NewInteractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interact",
		Short: "🎮 Interactive mode",
		Long:  "Run aescryptool in interactive mode with guided prompts for all inputs",
		RunE:  runInteractive,
	}
}

// runInteractive is the main entry point for interactive mode.
//
// It displays the welcome header and enters a loop that repeatedly presents
// the operation menu until the user chooses to exit or sends Ctrl+D.
//
// Returns:
//   - error: nil always (interactive mode exits via normal flow)
func runInteractive(cmd *cobra.Command, args []string) error {
	ui.PrintInteractiveHeader()

	for {
		choice := ui.PromptOperation()
		switch choice {
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
//  3. Prompt for password with strength validation
//  4. Prompt for password confirmation
//  5. Prompt for worker count (parallel processing)
//  6. Check if output file exists and ask for overwrite confirmation
//  7. Execute encryption with progress bar
//  8. Wait for user to press Enter before returning to menu
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
		ui.ErrorColor.Printf("❌ Error: %v\n", err)
		fmt.Println()
		waitForUser()
		return
	}

	if err := service.ExecuteEncryption(input, output, password, workerCount, false); err != nil {
		ui.ErrorColor.Printf("❌ Error: %v\n", err)
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

	input := ui.PromptFilePath("📁 Encrypted file", true, "")
	if input == "" {
		return
	}
	fmt.Println()

	defaultOutput := strings.TrimSuffix(input, ".enc")
	if defaultOutput == input {
		defaultOutput = input + ".dec"
	}
	output := ui.PromptFilePath("📂 Output file", false, defaultOutput)
	if output == "" {
		output = defaultOutput
	}
	fmt.Println()

	password := ui.PromptPassword("🔑 Password", false)
	if password == "" {
		return
	}
	fmt.Println()

	if err := checkAndConfirmOverwrite(output); err != nil {
		ui.ErrorColor.Printf("❌ Error: %v\n", err)
		fmt.Println()
		waitForUser()
		return
	}

	if err := service.ExecuteDecryption(input, output, password, false); err != nil {
		ui.ErrorColor.Printf("❌ Error: %v\n", err)
	}

	fmt.Println()
	waitForUser()
}

// promptInputFile prompts the user for the source file path.
//
// Returns:
//   - string: The input file path, or empty string if cancelled
func promptInputFile() string {
	input := ui.PromptFilePath("📁 File to encrypt", true, "")
	if input == "" {
		return ""
	}
	fmt.Println()
	return input
}

// promptOutputFile prompts for output file path with a default value.
//
// Parameters:
//   - input: Source file path used to generate default output name
//   - extension: File extension to append for default output (e.g., ".enc")
//
// Returns:
//   - string: The output file path
func promptOutputFile(input, extension string) string {
	defaultOutput := input + extension
	output := ui.PromptFilePath("📂 Output file", false, defaultOutput)
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
//   - string: The confirmed password, or empty string if validation fails
func promptAndConfirmPassword() string {
	password := ui.PromptPassword("🔑 Password", true)
	if password == "" {
		return ""
	}
	fmt.Println()

	confirm := ui.PromptPassword("✅ Confirmation", true)
	if confirm == "" {
		return ""
	}
	fmt.Println()

	if password != confirm {
		ui.ErrorColor.Println("❌ Passwords do not match")
		fmt.Println()
		return ""
	}

	return password
}

// checkAndConfirmOverwrite checks if a file exists and prompts for overwrite confirmation.
//
// Parameters:
//   - output: Path to the file to check
//
// Returns:
//   - error: If the file exists and user declines overwrite, or if check fails
func checkAndConfirmOverwrite(output string) error {
	exists, err := service.CheckFileExists(output)
	if err != nil {
		return fmt.Errorf("check file existence: %w", err)
	}

	if exists {
		if !ui.PromptConfirm("⚠️  File already exists. Overwrite?", true) {
			ui.InfoColor.Println("❌ Operation cancelled")
			return fmt.Errorf("user cancelled overwrite")
		}
		fmt.Println()
	}

	return nil
}

// waitForUser pauses execution until the user presses Enter.
func waitForUser() {
	ui.InfoColor.Println("🔁 Press Enter to continue...")
	fmt.Scanln()
}
