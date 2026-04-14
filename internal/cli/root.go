// Package cli provides the command-line interface for aescryptool.
//
// It implements the Cobra-based CLI with support for encryption, decryption,
// interactive mode, and version display. The package handles flag parsing,
// command routing, and delegates business logic to the service layer.
package cli

import (
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// CLIConfig holds configuration for CLI commands.
// This replaces the global variables with a structured approach.
type CLIConfig struct {
	Pass    string
	Workers int
	Force   bool
	Quiet   bool
}

// GlobalConfig is the shared configuration for all commands.
// While still a global, it's encapsulated in a struct for better maintainability.
var GlobalConfig = &CLIConfig{}

// Execute runs the root command and returns any error encountered.
//
// This is the main entry point called from cmd/aescryptool/main.go.
// It parses command-line arguments, executes the appropriate command,
// and returns an error that can be used to set the exit code.
func Execute() error {
	return rootCmd.Execute()
}

// rootCmd is the base command for aescryptool.
//
// It displays the application header and help information when called
// without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "aescryptool",
	Short: "🔐 Secure file encryption using AES-256-GCM",
	Long: ui.HeaderColor.Sprint(`
╔══════════════════════════════════════════════════════════════╗
║                 🔐 AESCRYPTOOL - AES-GCM                     ║
╚══════════════════════════════════════════════════════════════╝
`) + "\n\n" + ui.InfoColor.Sprint("Usage:") + ` aescryptool [command] [flags]

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
  aescryptool encrypt secret.txt secret.enc

  # Decrypt with interactive password prompt (recommended)
  aescryptool decrypt secret.enc secret.txt

  # Encrypt with --pass flag (for scripts)
  aescryptool encrypt secret.txt secret.enc --pass "myPassword"

  # With parallel processing (8 workers)
  aescryptool encrypt largefile.mp4 encrypted.enc --workers 8

  # Force overwrite without confirmation
  aescryptool encrypt data.txt data.enc --force
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
