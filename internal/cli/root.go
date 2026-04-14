// Package cli provides the command-line interface for cryptool.
//
// It implements the Cobra-based CLI with support for encryption, decryption,
// interactive mode, and version display. The package handles flag parsing,
// command routing, and delegates business logic to the service layer.
package cli

import (
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// Global flags shared across all commands.
var (
	pass    string // Password for encryption/decryption
	workers int    // Number of parallel workers (encryption only)
	force   bool   // Force overwrite existing output file
	quiet   bool   // Suppress progress output
)

// Execute runs the root command and returns any error encountered.
//
// This is the main entry point called from cmd/cryptool/main.go.
// It parses command-line arguments, executes the appropriate command,
// and returns an error that can be used to set the exit code.
func Execute() error {
	return rootCmd.Execute()
}

// rootCmd is the base command for cryptool.
//
// It displays the application header and help information when called
// without any subcommands.
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
  interact  Interactive mode with guided prompts
  version   Show version information
  help      Display help about any command

Examples:
  # Encrypt a file (non-interactive)
  cryptool encrypt secret.txt secret.enc --pass "myPassword"

  # Decrypt a file (non-interactive)
  cryptool decrypt secret.enc output.txt --pass "myPassword"

  # Interactive mode (guided prompts)
  cryptool interact

  # With parallel processing (8 workers)
  cryptool encrypt largefile.mp4 encrypted.enc --pass "secure" --workers 8

  # Force overwrite without confirmation
  cryptool encrypt data.txt data.enc --pass "pass" --force

  # Silent mode (no progress bar)
  cryptool encrypt log.txt log.enc --pass "secret" --quiet
`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// init registers all subcommands with the root command.
//
// This runs automatically when the package is imported, setting up the
// complete command tree before execution.
func init() {
	rootCmd.AddCommand(NewEncryptCmd())
	rootCmd.AddCommand(NewDecryptCmd())
	rootCmd.AddCommand(NewInteractCmd())
	rootCmd.AddCommand(NewVersionCmd())
}
