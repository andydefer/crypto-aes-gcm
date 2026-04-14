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
	pass    string // Password for encryption/decryption (optional - will prompt if omitted)
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

Password Management:
  For both encrypt and decrypt commands, you can either:
    - Provide --pass flag (visible in process list)
    - Omit the flag and enter password interactively (recommended)

  For encryption, interactive mode includes password confirmation and strength validation.

Examples:
  # Encrypt with interactive password prompt (recommended)
  cryptool encrypt secret.txt secret.enc

  # Decrypt with interactive password prompt (recommended)
  cryptool decrypt secret.enc secret.txt

  # Encrypt with --pass flag (for scripts)
  cryptool encrypt secret.txt secret.enc --pass "myPassword"

  # With parallel processing (8 workers)
  cryptool encrypt largefile.mp4 encrypted.enc --workers 8

  # Force overwrite without confirmation
  cryptool encrypt data.txt data.enc --force
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
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
