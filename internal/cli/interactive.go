// Package cli provides the interactive command for cryptool.
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
		Long:  "Run cryptool in interactive mode with guided prompts for all inputs",
		Run:   runInteractive,
	}
}

// runInteractive is the main entry point for interactive mode.
//
// It displays the welcome header and enters a loop that repeatedly presents
// the operation menu until the user chooses to exit or sends Ctrl+D.
func runInteractive(cmd *cobra.Command, args []string) {
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
			return
		}
	}
}

// runInteractiveEncrypt guides the user through the encryption process.
//
// The flow:
//  1. Prompt for source file path
//  2. Prompt for output file path (defaults to input + ".enc")
//  3. Prompt for password (with strength validation)
//  4. Prompt for password confirmation
//  5. Prompt for worker count (parallel processing)
//  6. Check if output file exists and ask for overwrite confirmation
//  7. Execute encryption with progress bar
//  8. Wait for user to press Enter before returning to menu
func runInteractiveEncrypt() {
	ui.PrintEncryptHeader()

	input := ui.PromptFilePath("📁 Fichier à chiffrer", true, "")
	if input == "" {
		return
	}
	fmt.Println()

	defaultOutput := input + ".enc"
	output := ui.PromptFilePath("📂 Fichier de sortie", false, defaultOutput)
	if output == "" {
		output = defaultOutput
	}
	fmt.Println()

	password := ui.PromptPassword("🔑 Mot de passe", true)
	if password == "" {
		return
	}
	fmt.Println()

	confirm := ui.PromptPassword("✅ Confirmation", true)
	if confirm == "" {
		return
	}
	if password != confirm {
		ui.ErrorColor.Println("❌ Les mots de passe ne correspondent pas")
		fmt.Println()
		return
	}
	fmt.Println()

	workerCount := ui.PromptWorkers()
	fmt.Println()

	exists, err := service.CheckFileExists(output)
	if err != nil {
		ui.ErrorColor.Printf("❌ Erreur lors de la vérification: %v\n", err)
		fmt.Println()
		return
	}
	if exists {
		if !ui.PromptConfirm("⚠️  Le fichier existe déjà. Écraser ?", true) {
			ui.InfoColor.Println("❌ Opération annulée")
			fmt.Println()
			return
		}
		fmt.Println()
	}

	if err := service.ExecuteEncryption(input, output, password, workerCount, false); err != nil {
		ui.ErrorColor.Printf("❌ Erreur: %v\n", err)
	}

	fmt.Println()
	ui.InfoColor.Println("🔁 Appuyez sur Entrée pour continuer...")
	fmt.Scanln()
}

// runInteractiveDecrypt guides the user through the decryption process.
//
// The flow:
//  1. Prompt for encrypted source file path
//  2. Prompt for output file path (defaults to input without ".enc" or input + ".dec")
//  3. Prompt for password (no confirmation needed, just validation)
//  4. Check if output file exists and ask for overwrite confirmation
//  5. Execute decryption with progress bar
//  6. Wait for user to press Enter before returning to menu
func runInteractiveDecrypt() {
	ui.PrintDecryptHeader()

	input := ui.PromptFilePath("📁 Fichier chiffré", true, "")
	if input == "" {
		return
	}
	fmt.Println()

	defaultOutput := strings.TrimSuffix(input, ".enc")
	if defaultOutput == input {
		defaultOutput = input + ".dec"
	}
	output := ui.PromptFilePath("📂 Fichier de sortie", false, defaultOutput)
	if output == "" {
		output = defaultOutput
	}
	fmt.Println()

	password := ui.PromptPassword("🔑 Mot de passe", false)
	if password == "" {
		return
	}
	fmt.Println()

	exists, err := service.CheckFileExists(output)
	if err != nil {
		ui.ErrorColor.Printf("❌ Erreur lors de la vérification: %v\n", err)
		fmt.Println()
		return
	}
	if exists {
		if !ui.PromptConfirm("⚠️  Le fichier existe déjà. Écraser ?", true) {
			ui.InfoColor.Println("❌ Opération annulée")
			fmt.Println()
			return
		}
		fmt.Println()
	}

	if err := service.ExecuteDecryption(input, output, password, false); err != nil {
		ui.ErrorColor.Printf("❌ Erreur: %v\n", err)
	}

	fmt.Println()
	ui.InfoColor.Println("🔁 Appuyez sur Entrée pour continuer...")
	fmt.Scanln()
}
