// Package cli provides the command-line interface for aescryptool.
//
// It implements CLI commands including version display with formatted output,
// colored banners, and system information.
package cli

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	appName     = "AESCRYPTOOL"
	appVersion  = "v2.0.0"
	cryptoAlgos = "AES-256-GCM | Argon2id | Parallel"
	bannerWidth = 52
)

// NewVersionCmd creates and returns the version command.
//
// The command displays application version, build information, and system details
// including Go version, operating system, architecture, and CPU count.
//
// Returns:
//   - *cobra.Command: Configured cobra command for version display
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Long:  "Display aescryptool version, build information, and system details",
		RunE: func(cmd *cobra.Command, _ []string) error {
			printVersion(cmd.OutOrStdout())
			return nil
		},
	}
}

// printVersion writes the version banner and system information to w.
//
// Parameters:
//   - w: Destination writer for version output (typically stdout)
func printVersion(w io.Writer) {
	printBanner(w)
	printSystemInfo(w)
}

// printBanner writes the ASCII art banner to w.
//
// The banner includes the application name, version, and supported algorithms
// centered within a box of width bannerWidth (52 characters).
func printBanner(w io.Writer) {
	headerColor := color.New(color.FgMagenta, color.Bold)
	headerColor.SetWriter(w)

	horizontalLine := strings.Repeat("═", bannerWidth)

	line1 := fmt.Sprintf("🔐 %s - %s", appName, appVersion)
	line2 := cryptoAlgos

	banner := fmt.Sprintf(`
╔%s╗
%s
%s
╚%s╝
`,
		horizontalLine,
		centerText(line1, bannerWidth),
		centerText(line2, bannerWidth),
		horizontalLine,
	)

	headerColor.Fprint(w, banner)
}

// printSystemInfo writes the runtime and system information to w.
//
// Parameters:
//   - w: Destination writer for system information (typically stdout)
//
// Displays Go version, operating system, architecture, and CPU count.
func printSystemInfo(w io.Writer) {
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

// centerText centers a string within a specified width.
//
// Parameters:
//   - text: String to center
//   - width: Target width for centering
//
// Returns:
//   - string: Padded string centered within the width
//
// If the text length exceeds width, the original text is returned unchanged.
func centerText(text string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-textLen)
}
