// Package cli provides the command-line interface for aesaesaescryptool.
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
// This is the main entry point called from cmd/aesaesaescryptool/main.go.
// It parses command-line arguments, executes the appropriate command,
// and returns an error that can be used to set the exit code.
func Execute() error {
	return rootCmd.Execute()
}

// rootCmd is the base command for aesaesaescryptool.
//
// It displays the application header and help information when called
// without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "aesaesaescryptool",
	Short: "🔐 Secure file encryption using AES-256-GCM",
	Long: ui.HeaderColor.Sprint(`
╔══════════════════════════════════════════════════════════════╗
║                 🔐 AESCRYPTOOL - AES-GCM                     ║
╚══════════════════════════════════════════════════════════════╝
`) + "\n\n" + ui.InfoColor.Sprint("Usage:") + ` aesaesaescryptool [command] [flags]

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
  aesaesaescryptool encrypt secret.txt secret.enc

  # Decrypt with interactive password prompt (recommended)
  aesaesaescryptool decrypt secret.enc secret.txt

  # Encrypt with --pass flag (for scripts)
  aesaesaescryptool encrypt secret.txt secret.enc --pass "myPassword"

  # With parallel processing (8 workers)
  aesaesaescryptool encrypt largefile.mp4 encrypted.enc --workers 8

  # Force overwrite without confirmation
  aesaesaescryptool encrypt data.txt data.enc --force
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
