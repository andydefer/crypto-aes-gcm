package ui

import "fmt"

// PrintInteractiveHeader displays the interactive mode header.
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

// PrintEncryptHeader displays the encryption header.
func PrintEncryptHeader() {
	fmt.Println()
	InfoColor.Println("🔐 CHIFFREMENT DE FICHIER")
	fmt.Println("────────────────────────────────────────")
	fmt.Println()
}

// PrintDecryptHeader displays the decryption header.
func PrintDecryptHeader() {
	fmt.Println()
	InfoColor.Println("🔓 DÉCHIFFREMENT DE FICHIER")
	fmt.Println("────────────────────────────────────────")
	fmt.Println()
}

// PrintInteractiveGoodbye displays the goodbye message.
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
