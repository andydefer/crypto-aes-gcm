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
//   - Automatic trimming of leading/trailing spaces
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

		// Trim spaces from the result
		result = strings.TrimSpace(result)

		if result == "" && defaultValue != "" {
			result = strings.TrimSpace(defaultValue)
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

		// Trim spaces from password as well (though unlikely)
		result = strings.TrimSpace(result)

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

		result = strings.TrimSpace(result)

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
