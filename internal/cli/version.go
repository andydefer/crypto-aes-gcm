// Package cli provides the command-line interface for cryptool.
package cli

import (
	"fmt"
	"io"
	"runtime"

	"github.com/fatih/color"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			printVersionToWriter(cmd.OutOrStdout())
			return nil
		},
	}
}

// printVersionToWriter displays the version banner and system information
// to the specified writer.
func printVersionToWriter(w io.Writer) {
	// Use color.New with the specific writer for colored output
	headerColor := color.New(color.FgMagenta, color.Bold)
	headerColor.SetWriter(w)

	header := `
╔═══════════════════════════════════════╗
║  🔐 CRYPTOOL - AES-GCM v2.0.0         ║
║  AES-256-GCM | Argon2id | Parallel    ║
╚═══════════════════════════════════════╝
`
	headerColor.Fprint(w, header)

	infoColor := color.New(color.FgCyan, color.Bold)
	infoColor.SetWriter(w)

	info := fmt.Sprintf("\n  📦 Build: %s\n  🖥️  OS/Arch: %s/%s\n  💻 CPUs: %d\n\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.NumCPU())

	infoColor.Fprint(w, info)
}

// printVersion maintains backward compatibility for existing code
func printVersion() {
	printVersionToWriter(color.Output)
}
