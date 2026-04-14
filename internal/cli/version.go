// Package cli provides the command-line interface for aesaesaescryptool.
package cli

import (
	"fmt"
	"io"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	appName     = "AESCRYPTOOL"
	appVersion  = "v2.0.0"
	cryptoAlgos = "AES-256-GCM | Argon2id | Parallel"
)

// NewVersionCmd creates and returns the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display aesaesaescryptool version, build information, and system details",
		RunE: func(cmd *cobra.Command, _ []string) error {
			printVersion(cmd.OutOrStdout())
			return nil
		},
	}
}

// printVersion writes the version banner and system information to w.
func printVersion(w io.Writer) {
	printBanner(w)
	printBuildInfo(w)
}

// printBanner writes the ASCII art banner to w.
func printBanner(w io.Writer) {
	headerColor := color.New(color.FgMagenta, color.Bold)
	headerColor.SetWriter(w)

	line := "════════════════════════════════════════"

	banner := fmt.Sprintf(`
╔%s╗
║  🔐 %s - %s ║
║  %s ║
╚%s╝
`,
		line,
		appName,
		appVersion,
		cryptoAlgos,
		line,
	)

	headerColor.Fprint(w, banner)
}

// printBuildInfo writes the runtime and system information to w.
func printBuildInfo(w io.Writer) {
	infoColor := color.New(color.FgCyan, color.Bold)
	infoColor.SetWriter(w)

	info := fmt.Sprintf(
		"\n  📦 Build: %s\n  🖥️  OS/Arch: %s/%s\n  💻 CPUs: %d\n\n",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
	)

	infoColor.Fprint(w, info)
}
