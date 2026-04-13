// ==== ./cmd/cryptool/main.go ====
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	pass    string
	workers int
	force   bool
	quiet   bool

	// Color styles
	infoColor    = color.New(color.FgCyan, color.Bold)
	successColor = color.New(color.FgGreen, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow)
	headerColor  = color.New(color.FgMagenta, color.Bold)
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "cryptool",
		Short: "🔐 Secure file encryption using AES-256-GCM",
		Long: headerColor.Sprint(`
╔══════════════════════════════════════════════════════════════╗
║                    🔐 CRYPTOOL - AES-GCM                     ║
║                                                              ║
║  Secure file encryption with Argon2id key derivation         ║
║  and parallel streaming encryption.                          ║
╚══════════════════════════════════════════════════════════════╝
`) + "\n\n" + infoColor.Sprint("Usage:") + ` cryptool [command] [flags]

Commands:
  encrypt   Encrypt a file
  decrypt   Decrypt a file
  version   Show version information
  help      Help about any command

Examples:
  # Encrypt a file
  cryptool encrypt secret.txt secret.enc --pass "myPassword"

  # Decrypt a file
  cryptool decrypt secret.enc output.txt --pass "myPassword"

  # With custom workers (parallel processing)
  cryptool encrypt largefile.mp4 encrypted.enc --pass "secure" --workers 8

  # Force overwrite without confirmation
  cryptool encrypt data.txt data.enc --pass "pass" --force
`,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	// Encrypt command
	encryptCmd := &cobra.Command{
		Use:   "encrypt [input] [output]",
		Short: "🔒 Encrypt a file",
		Long:  "Encrypt a file using AES-256-GCM with Argon2id key derivation.",
		Args:  cobra.ExactArgs(2),
		Run:   runEncrypt,
	}

	// Decrypt command
	decryptCmd := &cobra.Command{
		Use:   "decrypt [input] [output]",
		Short: "🔓 Decrypt a file",
		Long:  "Decrypt a file encrypted with the encrypt command.",
		Args:  cobra.ExactArgs(2),
		Run:   runDecrypt,
	}

	// Version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}

	// Add global flags to both commands
	for _, cmd := range []*cobra.Command{encryptCmd, decryptCmd} {
		cmd.Flags().StringVarP(&pass, "pass", "p", "", "Passphrase for encryption/decryption (required)")
		cmd.Flags().IntVarP(&workers, "workers", "w", cryptolib.DefaultWorkers, "Number of parallel workers (encryption only)")
		cmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite existing output file")
		cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress progress output")
		cmd.MarkFlagRequired("pass")
	}

	rootCmd.AddCommand(encryptCmd, decryptCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runEncrypt(cmd *cobra.Command, args []string) {
	input := args[0]
	output := args[1]

	// Validate worker count
	workerCount := workers
	if workerCount <= 0 {
		workerCount = cryptolib.DefaultWorkers
	}
	if workerCount > runtime.NumCPU()*2 {
		workerCount = runtime.NumCPU() * 2
		if !quiet {
			warningColor.Printf("⚠️  Workers reduced to %d (max 2×CPU cores)\n", workerCount)
		}
	}

	// Check input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		errorColor.Fprintf(os.Stderr, "❌ Error: input file '%s' does not exist\n", input)
		os.Exit(1)
	}

	// Check output file overwrite
	if err := checkOverwrite(output, force); err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Get file size for progress bar
	fileInfo, _ := os.Stat(input)
	fileSize := fileInfo.Size()

	// Show info header
	if !quiet {
		printHeader("ENCRYPT", input, output, workerCount)
	}

	// Create progress bar
	bar := createProgressBar(fileSize, "🔒 Encrypting")

	// Create encryptor
	encryptor, err := cryptolib.NewEncryptor(workerCount)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Failed to create encryptor: %v\n", err)
		os.Exit(1)
	}

	// Wrap reader with progress tracking
	reader := &progressReader{
		r:     mustOpenFile(input),
		bar:   bar,
		total: fileSize,
	}

	// Create output file
	outFile, err := os.Create(output)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Perform encryption
	if err := encryptor.Encrypt(reader, outFile, pass); err != nil {
		bar.Clear()
		errorColor.Fprintf(os.Stderr, "❌ Encryption failed: %v\n", err)
		os.Exit(1)
	}

	bar.Finish()
	printSuccess(output, fileSize)
}

func runDecrypt(cmd *cobra.Command, args []string) {
	input := args[0]
	output := args[1]

	// Check input file exists
	if _, err := os.Stat(input); os.IsNotExist(err) {
		errorColor.Fprintf(os.Stderr, "❌ Error: input file '%s' does not exist\n", input)
		os.Exit(1)
	}

	// Check output file overwrite
	if err := checkOverwrite(output, force); err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	// Get file size
	fileInfo, _ := os.Stat(input)
	fileSize := fileInfo.Size()

	// Show info header
	if !quiet {
		printHeader("DECRYPT", input, output, workers)
	}

	// Create progress bar
	bar := createProgressBar(fileSize, "🔓 Decrypting")

	// Read header first to get salt
	f, err := os.Open(input)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Failed to open input: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var header cryptolib.FileHeader
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Failed to read file header: %v\n", err)
		os.Exit(1)
	}

	// Create decryptor
	decryptor, err := cryptolib.NewDecryptor(pass, header.Salt[:])
	if err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Failed to create decryptor: %v\n", err)
		os.Exit(1)
	}

	// Re-open file for full decryption
	f.Seek(0, 0)
	reader := &progressReader{
		r:     f,
		bar:   bar,
		total: fileSize,
	}

	// Create output file
	outFile, err := os.Create(output)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Perform decryption
	if err := decryptor.Decrypt(reader, outFile); err != nil {
		bar.Clear()
		errorColor.Fprintf(os.Stderr, "❌ Decryption failed: %v\n", err)
		os.Exit(1)
	}

	bar.Finish()
	printSuccess(output, fileSize)
}

func checkOverwrite(output string, force bool) error {
	if force {
		return nil
	}

	if _, err := os.Stat(output); err == nil {
		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("File '%s' already exists. Overwrite", output),
			IsConfirm: true,
			Default:   "n",
		}

		result, err := prompt.Run()
		if err != nil || (result != "y" && result != "Y" && result != "yes" && result != "Yes") {
			return fmt.Errorf("operation cancelled")
		}
	}
	return nil
}

func printHeader(mode, input, output string, workers int) {
	infoColor.Printf("\n🔐 Crypto-AES-GCM - %s MODE\n", mode)
	fmt.Println(strings.Repeat("─", 50))
	infoColor.Printf("📁 Input:   %s\n", input)
	infoColor.Printf("📂 Output:  %s\n", output)
	infoColor.Printf("⚙️  Workers: %d\n", workers)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
}

func printSuccess(output string, size int64) {
	fmt.Println()
	successColor.Printf("✅ Decryption successful!\n")
	infoColor.Printf("📄 Output: %s\n", output)

	// Show file size
	var sizeStr string
	switch {
	case size > 1024*1024*1024:
		sizeStr = fmt.Sprintf("%.2f GB", float64(size)/(1024*1024*1024))
	case size > 1024*1024:
		sizeStr = fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	case size > 1024:
		sizeStr = fmt.Sprintf("%.2f KB", float64(size)/1024)
	default:
		sizeStr = fmt.Sprintf("%d B", size)
	}
	infoColor.Printf("📏 Size:   %s\n", sizeStr)
	fmt.Println()
}

func printVersion() {
	headerColor.Printf(`
╔═══════════════════════════════════════╗
║  🔐 CRYPTOOL - AES-GCM v2.0.0         ║
║                                       ║
║  AES-256-GCM | Argon2id | Parallel    ║
╚═══════════════════════════════════════╝
`)
	infoColor.Printf("\n  📦 Build:    %s\n", runtime.Version())
	infoColor.Printf("  🖥️  OS/Arch:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
	infoColor.Printf("  💻 CPUs:     %d\n\n", runtime.NumCPU())
}

func createProgressBar(total int64, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}

// progressReader wraps an io.Reader to track progress
type progressReader struct {
	r     io.ReadCloser
	bar   *progressbar.ProgressBar
	total int64
	read  int64
}

func (pr *progressReader) Read(p []byte) (n int, err error) {
	n, err = pr.r.Read(p)
	pr.read += int64(n)
	pr.bar.Set64(pr.read)
	return n, err
}

func (pr *progressReader) Close() error {
	return pr.r.Close()
}

func mustOpenFile(path string) io.ReadCloser {
	f, err := os.Open(path)
	if err != nil {
		errorColor.Fprintf(os.Stderr, "❌ Failed to open file: %v\n", err)
		os.Exit(1)
	}
	return f
}

// Helper for strings.Repeat (add at top if not imported)
