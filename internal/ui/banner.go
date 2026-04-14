// Package ui provides terminal user interface components for cryptool.
//
// This package handles all user-facing output including:
//   - Colored banners and headers for different operation modes
//   - Interactive mode welcome and goodbye messages
//   - Progress bars for file operations
//   - Color-coded informational, success, error, and warning messages
//
// The package uses the fatih/color library to ensure consistent and
// visually appealing terminal output across different platforms.
package ui

import "fmt"

// PrintInteractiveHeader displays the interactive mode welcome banner.
//
// This header is shown when the user starts interactive mode. It includes:
//   - Application name and mode indicator
//   - Brief instructions on how to use the interactive prompts
//   - Keyboard shortcuts (Ctrl+C to return to menu, Ctrl+D to quit)
//
// The banner is formatted with a decorative ASCII art border and uses
// the HeaderColor for visual emphasis.
func PrintInteractiveHeader() {
	fmt.Println()
	fmt.Println()
	HeaderColor.Printf(`
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║                    🎮 CRYPTOOL - MODE INTERACTIF                   ║
║                                                                    ║
║  Suivez les invites pour chiffrer ou déchiffrer vos fichiers       ║
║  Toutes les entrées seront validées avant exécution                ║
║                                                                    ║
║  Ctrl+C = Retour au menu | Ctrl+D = Quitter                        ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝
`)
	fmt.Println()
	fmt.Println()
}

// PrintEncryptHeader displays the encryption operation header.
//
// This header is shown before prompting for encryption inputs.
// It clearly indicates to the user that they are in the encryption workflow.
func PrintEncryptHeader() {
	fmt.Println()
	InfoColor.Println("🔐 CHIFFREMENT DE FICHIER")
	fmt.Println("────────────────────────────────────────")
	fmt.Println()
}

// PrintDecryptHeader displays the decryption operation header.
//
// This header is shown before prompting for decryption inputs.
// It clearly indicates to the user that they are in the decryption workflow.
func PrintDecryptHeader() {
	fmt.Println()
	InfoColor.Println("🔓 DÉCHIFFREMENT DE FICHIER")
	fmt.Println("────────────────────────────────────────")
	fmt.Println()
}

// PrintInteractiveGoodbye displays the farewell message when exiting interactive mode.
//
// This message is shown when the user chooses to exit or sends Ctrl+D.
// It provides a friendly closing experience with decorative ASCII art borders.
func PrintInteractiveGoodbye() {
	fmt.Println()
	fmt.Println()
	SuccessColor.Printf(`
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║              👋 Merci d'avoir utilisé CRYPTOOL !                   ║
║                                                                    ║
║              À bientôt pour vos prochains chiffrements !           ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝
`)
	fmt.Println()
	fmt.Println()
}
