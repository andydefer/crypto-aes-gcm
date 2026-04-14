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

// PromptOperation asks the user which operation to perform.
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

// PromptFilePath asks the user for a file path with validation.
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

// PromptPassword asks the user for a password (masked input).
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

// PromptWorkers asks the user for the number of parallel workers.
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
// Enter key defaults to YES (true) for better UX.
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
