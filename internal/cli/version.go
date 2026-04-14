package cli

import (
	"runtime"

	"github.com/andydefer/crypto-aes-gcm/internal/ui"
	"github.com/spf13/cobra"
)

// NewVersionCmd creates the version command.
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
}

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
