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

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
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
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: lang.T(lang.VersionShortDesc),
		Long:  lang.T(lang.VersionLongDesc),
		RunE: func(cmd *cobra.Command, _ []string) error {
			printVersion(cmd.OutOrStdout())
			return nil
		},
	}
}

// printVersion writes the version banner and system information to w.
func printVersion(w io.Writer) {
	printBanner(w)
	printSystemInfo(w)
}

// printBanner writes the ASCII art banner to w.
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
func printSystemInfo(w io.Writer) {
	infoColor := color.New(color.FgCyan, color.Bold)
	infoColor.SetWriter(w)

	info := fmt.Sprintf(
		"\n  "+lang.T(lang.VersionBuildInfo)+"\n  "+lang.T(lang.VersionOSArch)+"\n  "+lang.T(lang.VersionCPUs)+"\n\n",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
	)

	infoColor.Fprint(w, info)
}

// centerText centers a string within a specified width.
func centerText(text string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-textLen)
}
