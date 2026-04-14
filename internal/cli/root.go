// Package cli provides the command-line interface for aescryptool.
//
// It implements the Cobra-based CLI with support for encryption, decryption,
// interactive mode, and version display. The package handles flag parsing,
// command routing, and delegates business logic to the service layer.
package cli

import (
	"strings"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// CLIConfig holds configuration for CLI commands.
type CLIConfig struct {
	Pass    string
	Workers int
	Force   bool
	Quiet   bool
	Lang    string
}

// GlobalConfig is the shared configuration for all commands.
var GlobalConfig = &CLIConfig{}

// applyLanguage sets the active language based on the flag value.
//
// It supports case-insensitive values "en", "english", "fr", "french".
// If an invalid language is provided, it falls back to English and prints a warning.
//
// Parameters:
//   - langFlag: language flag value from CLI (empty string uses default English)
func applyLanguage(langFlag string) {
	if langFlag == "" {
		lang.SetLanguage(lang.English)
		return
	}

	switch strings.ToLower(langFlag) {
	case "en", "english":
		lang.SetLanguage(lang.English)
	case "fr", "french":
		lang.SetLanguage(lang.French)
	default:
		ui.ErrorColor.Printf("⚠️ Invalid language '%s', using English (en). Supported: en, fr\n", langFlag)
		lang.SetLanguage(lang.English)
	}
}

// Execute runs the root command and returns any error encountered.
//
// This is the main entry point called from cmd/aescryptool/main.go.
// It ensures the language is initialized before parsing commands.
//
// Returns:
//   - error: any error from command execution, or nil on success
func Execute() error {
	if lang.GetLanguage() == "" {
		lang.SetLanguage(lang.English)
	}
	return rootCmd.Execute()
}

// rootCmd is the base command for aescryptool.
//
// It displays the application header and help information when called
// without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "aescryptool",
	Short: lang.T(lang.RootShortDesc),
	Long: ui.HeaderColor.Sprint(`
╔══════════════════════════════════════════════════════════════╗
║                 🔐 AESCRYPTOOL - AES-GCM                     ║
╚══════════════════════════════════════════════════════════════╝
`) + "\n\n" + ui.InfoColor.Sprint(lang.T(lang.RootUsage)) + ` aescryptool [command] [flags]

` + ui.InfoColor.Sprint(lang.T(lang.RootCommandsTitle)) + `
  encrypt   ` + lang.T(lang.CmdEncryptShort) + `
  decrypt   ` + lang.T(lang.CmdDecryptShort) + `
  interact  ` + lang.T(lang.InteractiveTitle) + `
  version   ` + lang.T(lang.VersionShortDesc) + `
  help      Display help about any command

` + ui.InfoColor.Sprint(lang.T(lang.RootPasswordManagement)) + `
  For both encrypt and decrypt commands, you can either:
    - Provide --pass flag (visible in process list)
    - Omit the flag and enter password interactively (recommended)

  For encryption, interactive mode includes password confirmation and strength validation.

` + ui.InfoColor.Sprint(lang.T(lang.RootExamplesTitle)) + `
  ` + lang.T(lang.RootExampleEncrypt) + `
  aescryptool encrypt secret.txt secret.enc

  ` + lang.T(lang.RootExampleDecrypt) + `
  aescryptool decrypt secret.enc secret.txt

  ` + lang.T(lang.RootExamplePassFlag) + `
  aescryptool encrypt secret.txt secret.enc --pass "myPassword"

  ` + lang.T(lang.RootExampleWorkers) + `
  aescryptool encrypt largefile.mp4 encrypted.enc --workers 8

  ` + lang.T(lang.RootExampleForce) + `
  aescryptool encrypt data.txt data.enc --force
`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		return cmd.Help()
	},
}

// init registers all subcommands with the root command.
//
// This runs automatically when the package is imported, setting up the
// complete command tree before execution.
func init() {
	rootCmd.PersistentFlags().StringVar(&GlobalConfig.Lang, "lang", "", "Language for UI (en, fr) - default: en")
	rootCmd.AddCommand(NewEncryptCmd())
	rootCmd.AddCommand(NewDecryptCmd())
	rootCmd.AddCommand(NewInteractCmd())
	rootCmd.AddCommand(NewVersionCmd())
}
