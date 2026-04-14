// Package cli provides the command-line interface for cryptool.
package cli

import (
	"runtime"

	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command.
//
// This command displays version information including:
//   - Application name and version (v2.0.0)
//   - Cryptographic algorithms used (AES-256-GCM, Argon2id, Parallel)
//   - Go build version
//   - Operating system and architecture
//   - Number of available CPU cores
//
// Returns:
//   - *cobra.Command: Configured Cobra command that prints version info
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display cryptool version, build information, and system details",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
}

// printVersion displays the version banner and system information.
//
// The output includes:
//   - ASCII art banner with application name and version
//   - List of cryptographic algorithms used
//   - Go runtime version
//   - Target OS and architecture
//   - Number of available CPU cores
func printVersion() {
	ui.HeaderColor.Printf(`
╔═══════════════════════════════════════╗
║  🔐 CRYPTOOL - AES-GCM v2.0.0         ║
║  AES-256-GCM | Argon2id | Parallel    ║
╚═══════════════════════════════════════╝
`)
	ui.InfoColor.Printf("\n  📦 Build: %s\n", runtime.Version())
	ui.InfoColor.Printf("  🖥️  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	ui.InfoColor.Printf("  💻 CPUs: %d\n\n", runtime.NumCPU())
}
