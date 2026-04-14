// Package cli provides the command-line interface for cryptool.
package cli

import (
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	pass    string
	workers int
	force   bool
	quiet   bool
)

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "cryptool",
	Short: "🔐 Secure file encryption using AES-256-GCM",
	Long: ui.HeaderColor.Sprint(`
╔══════════════════════════════════════════════════════════════╗
║                    🔐 CRYPTOOL - AES-GCM                     ║
╚══════════════════════════════════════════════════════════════╝
`) + "\n\n" + ui.InfoColor.Sprint("Usage:") + ` cryptool [command] [flags]

Commands:
  encrypt   Encrypt a file
  decrypt   Decrypt a file
  interact  Interactive mode
  version   Show version
  help      Help

Examples:
  cryptool encrypt secret.txt secret.enc --pass "myPassword"
  cryptool decrypt secret.enc output.txt --pass "myPassword"
  cryptool interact
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(NewEncryptCmd())
	rootCmd.AddCommand(NewDecryptCmd())
	rootCmd.AddCommand(NewInteractCmd())
	rootCmd.AddCommand(NewVersionCmd())
}
