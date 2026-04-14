package cli

import (
	"fmt"
	"strings"

	"github.com/andydefer/crypto-aes-gcm/internal/service"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// NewInteractCmd creates the interactive command.
func NewInteractCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "interact",
		Short: "🎮 Interactive mode",
		Run:   runInteractive,
	}
}

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

	// Vérifier si le fichier existe vraiment
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

	// Vérifier si le fichier existe vraiment
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
