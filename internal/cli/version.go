// Package cli provides the command-line interface for aescryptool.
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
)

// NewVersionCmd creates and returns the version command.
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
func printVersion(w io.Writer) {
	printBanner(w)
	printBuildInfo(w)
}

// printBanner writes the ASCII art banner to w (sans barres verticales).
func printBanner(w io.Writer) {
	headerColor := color.New(color.FgMagenta, color.Bold)
	headerColor.SetWriter(w)

	// Largeur totale de la boîte : 52 caractères
	width := 52

	// Construction de la ligne horizontale
	horizontalLine := strings.Repeat("═", width)

	// Centrage du texte dans une largeur donnée
	centerText := func(text string, width int) string {
		textLen := len(text)
		if textLen >= width {
			return text
		}
		padding := (width - textLen) / 2
		return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-textLen)
	}

	// Ligne 1: "🔐 AESCRYPTOOL - v2.0.0"
	line1 := fmt.Sprintf("🔐 %s - %s", appName, appVersion)
	line1Centered := centerText(line1, width)

	// Ligne 2: cryptoAlgos
	line2Centered := centerText(cryptoAlgos, width)

	banner := fmt.Sprintf(`
╔%s╗
%s
%s
╚%s╝
`,
		horizontalLine,
		line1Centered,
		line2Centered,
		horizontalLine,
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
