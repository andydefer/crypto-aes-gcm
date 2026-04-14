// Package ui provides terminal user interface utilities for cryptool.
//
// This package handles all user interaction including:
//   - Colored output for different message types (info, success, error, warning)
//   - Progress bars for long-running operations
//   - Interactive prompts for file paths, passwords, and confirmations
//   - Banner displays for interactive mode
//
// All UI functions are designed to work consistently across different terminals
// and operating systems.
package ui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/manifoldco/promptui"
)

// PromptOperation displays a selection menu and returns the user's choice.
//
// The menu presents three options:
//   - Encrypt a file
//   - Decrypt a file
//   - Exit the application
//
// If the user cancels (Ctrl+C or Ctrl+D), the application exits gracefully.
//
// Returns:
//   - "encrypt" if the user selects encryption
//   - "decrypt" if the user selects decryption
//   - "exit" if the user selects exit
func PromptOperation() string {
	prompt := promptui.Select{
		Label: "Que souhaitez-vous faire",
		Items: []string{
			"🔒  Chiffrer un fichier",
			"🔓  Déchiffrer un fichier",
			"🚪  Quitter",
		},
		Size: 5,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		fmt.Println()
		SuccessColor.Println("👋 Merci d'avoir utilisé CRYPTOOL !")
		fmt.Println()
		os.Exit(0)
	}

	switch idx {
	case 0:
		return "encrypt"
	case 1:
		return "decrypt"
	default:
		return "exit"
	}
}

// PromptFilePath asks the user for a file path with optional validation.
//
// The prompt includes:
//   - Input validation to ensure the path is not empty
//   - Optional existence check (if mustExist is true)
//   - Default value support for common paths
//
// Parameters:
//   - label: The prompt text displayed to the user
//   - mustExist: If true, validates that the file exists on disk
//   - defaultValue: Default path if user presses Enter without typing
//
// Returns:
//   - The validated file path, or empty string if user cancels (Ctrl+C)
func PromptFilePath(label string, mustExist bool, defaultValue string) string {
	for {
		prompt := promptui.Prompt{
			Label: label,
		}
		if defaultValue != "" {
			prompt.Default = defaultValue
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				return ""
			}
			fmt.Println()
			SuccessColor.Println("👋 Merci d'avoir utilisé CRYPTOOL !")
			fmt.Println()
			os.Exit(0)
		}

		if result == "" && defaultValue != "" {
			result = defaultValue
		}

		if result == "" {
			ErrorColor.Println("❌ Le chemin ne peut pas être vide")
			continue
		}

		if mustExist {
			if _, err := os.Stat(result); os.IsNotExist(err) {
				ErrorColor.Printf("❌ Le fichier '%s' n'existe pas\n", result)
				continue
			}
		}

		SuccessColor.Printf("   ✓ %s\n", result)
		return result
	}
}

// PromptPassword asks the user for a password with masked input.
//
// When needValidation is true, the password must meet security requirements:
//   - Minimum 8 characters
//   - At least one uppercase letter (A-Z)
//   - At least one lowercase letter (a-z)
//   - At least one digit (0-9)
//
// Parameters:
//   - label: The prompt text displayed to the user
//   - needValidation: If true, enforces password strength requirements
//
// Returns:
//   - The entered password, or empty string if user cancels (Ctrl+C)
func PromptPassword(label string, needValidation bool) string {
	for {
		prompt := promptui.Prompt{
			Label: label,
			Mask:  '*',
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				return ""
			}
			fmt.Println()
			SuccessColor.Println("👋 Merci d'avoir utilisé CRYPTOOL !")
			fmt.Println()
			os.Exit(0)
		}

		if needValidation {
			if len(result) < 8 {
				ErrorColor.Println("❌ 8 caractères minimum")
				continue
			}
			if !regexp.MustCompile(`[A-Z]`).MatchString(result) {
				ErrorColor.Println("❌ Une majuscule requise")
				continue
			}
			if !regexp.MustCompile(`[a-z]`).MatchString(result) {
				ErrorColor.Println("❌ Une minuscule requise")
				continue
			}
			if !regexp.MustCompile(`[0-9]`).MatchString(result) {
				ErrorColor.Println("❌ Un chiffre requis")
				continue
			}
		}

		SuccessColor.Printf("   ✓ %s\n", strings.Repeat("*", len(result)))
		return result
	}
}

// PromptWorkers asks the user for the number of parallel encryption workers.
//
// The worker count is limited to between 1 and 2×CPU cores for optimal performance.
// The default value is cryptolib.DefaultWorkers (typically 4).
//
// Returns:
//   - The selected worker count, or the default value if user cancels (Ctrl+C)
func PromptWorkers() int {
	maxWorkers := runtime.NumCPU() * 2

	for {
		prompt := promptui.Prompt{
			Label:   fmt.Sprintf("⚙️  Workers (défaut: %d, max: %d)", cryptolib.DefaultWorkers, maxWorkers),
			Default: fmt.Sprintf("%d", cryptolib.DefaultWorkers),
		}

		result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				return cryptolib.DefaultWorkers
			}
			fmt.Println()
			SuccessColor.Println("👋 Merci d'avoir utilisé CRYPTOOL !")
			fmt.Println()
			os.Exit(0)
		}

		if result == "" {
			SuccessColor.Printf("   ✓ %d workers\n", cryptolib.DefaultWorkers)
			return cryptolib.DefaultWorkers
		}

		var w int
		_, err = fmt.Sscanf(result, "%d", &w)
		if err != nil || w < 1 {
			ErrorColor.Println("❌ Nombre valide requis (>=1)")
			continue
		}
		if w > maxWorkers {
			ErrorColor.Printf("❌ Maximum %d workers\n", maxWorkers)
			continue
		}

		SuccessColor.Printf("   ✓ %d workers\n", w)
		return w
	}
}

// PromptConfirm asks the user for a yes/no confirmation.
//
// The function accepts multiple affirmative responses:
//   - "y", "Y", "yes", "o", "O", "oui"
//
// And negative responses:
//   - "n", "N", "no", "non"
//
// Pressing Enter (empty input) returns the default value.
//
// Parameters:
//   - label: The prompt text displayed to the user
//   - defaultValue: Value returned when user presses Enter (true=Yes, false=No)
//
// Returns:
//   - true for affirmative responses, false for negative responses
func PromptConfirm(label string, defaultValue bool) bool {
	defaultDisplay := "Y/n"
	if !defaultValue {
		defaultDisplay = "y/N"
	}

	for {
		fmt.Printf("❓ %s [%s]: ", label, defaultDisplay)

		reader := bufio.NewReader(os.Stdin)
		result, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			SuccessColor.Println("👋 Merci d'avoir utilisé CRYPTOOL !")
			fmt.Println()
			os.Exit(0)
		}

		result = strings.TrimSpace(strings.ToLower(result))

		if result == "" {
			return defaultValue
		}

		if result == "y" || result == "yes" || result == "o" || result == "oui" {
			return true
		}
		if result == "n" || result == "no" || result == "non" {
			return false
		}

		ErrorColor.Println("❌ Répondez par y/n")
	}
}
