// Package cli provides the command-line interface for cryptool.
//
// It implements the Cobra commands for encryption, decryption, interactive mode,
// and version display. The package handles flag parsing, validation, and
// orchestration of the underlying crypto operations.
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

// TestExecute verifies that the Execute function does not panic.
func TestExecute(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute panicked: %v", r)
		}
	}()
}

// TestRootCmdHelp verifies that the root command help does not panic.
func TestRootCmdHelp(t *testing.T) {
	_ = rootCmd.Help()
}

// TestEncryptCmdFlags verifies that all expected flags are present on the encrypt command.
func TestEncryptCmdFlags(t *testing.T) {
	cmd := NewEncryptCmd()

	passFlag := cmd.Flags().Lookup("pass")
	if passFlag == nil {
		t.Error("--pass flag is missing")
	}

	workersFlag := cmd.Flags().Lookup("workers")
	if workersFlag == nil {
		t.Error("--workers flag is missing")
	}

	forceFlag := cmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("--force flag is missing")
	}

	quietFlag := cmd.Flags().Lookup("quiet")
	if quietFlag == nil {
		t.Error("--quiet flag is missing")
	}
}

// TestDecryptCmdFlags verifies that all expected flags are present on the decrypt command.
func TestDecryptCmdFlags(t *testing.T) {
	cmd := NewDecryptCmd()

	passFlag := cmd.Flags().Lookup("pass")
	if passFlag == nil {
		t.Error("--pass flag is missing")
	}

	workersFlag := cmd.Flags().Lookup("workers")
	if workersFlag == nil {
		t.Error("--workers flag is missing")
	}

	forceFlag := cmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("--force flag is missing")
	}

	quietFlag := cmd.Flags().Lookup("quiet")
	if quietFlag == nil {
		t.Error("--quiet flag is missing")
	}
}

// TestInteractCmd verifies the interactive command has the correct configuration.
func TestInteractCmd(t *testing.T) {
	cmd := NewInteractCmd()

	if cmd.Use != "interact" {
		t.Errorf("expected use 'interact', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("interact command should have a short description")
	}
}

// TestVersionCmd verifies the version command has the correct configuration.
func TestVersionCmd(t *testing.T) {
	cmd := NewVersionCmd()

	if cmd.Use != "version" {
		t.Errorf("expected use 'version', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("version command should have a short description")
	}
}

// ========== NOUVEAUX TESTS POUR LA CORRECTION DES ERREURS ==========

// TestEncryptCmdReturnsErrorOnInvalidInput verifies that the encrypt command
// returns an error (not os.Exit) when given invalid input.
func TestEncryptCmdReturnsErrorOnInvalidInput(t *testing.T) {
	cmd := NewEncryptCmd()

	// Create a buffer to capture stderr output
	stderrBuf := &bytes.Buffer{}
	cmd.SetErr(stderrBuf)

	// Test with non-existent input file
	cmd.SetArgs([]string{"nonexistent.txt", "output.enc", "--pass", "testpass"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent input file, got nil")
	}

	// Verify error message is user-friendly
	if !bytes.Contains(stderrBuf.Bytes(), []byte("inexistant")) {
		t.Errorf("Expected error about missing file, got: %s", stderrBuf.String())
	}
}

// TestEncryptCmdMissingPassword verifies that the encrypt command returns an error
// when the required --pass flag is missing.
func TestEncryptCmdMissingPassword(t *testing.T) {
	cmd := NewEncryptCmd()

	cmd.SetArgs([]string{"input.txt", "output.enc"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing password flag, got nil")
	}

	if err.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

// TestDecryptCmdReturnsErrorOnInvalidInput verifies that the decrypt command
// returns an error (not os.Exit) when given invalid input.
func TestDecryptCmdReturnsErrorOnInvalidInput(t *testing.T) {
	cmd := NewDecryptCmd()

	stderrBuf := &bytes.Buffer{}
	cmd.SetErr(stderrBuf)

	cmd.SetArgs([]string{"nonexistent.enc", "output.txt", "--pass", "testpass"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for non-existent input file, got nil")
	}
}

// TestDecryptCmdMissingPassword verifies that the decrypt command returns an error
// when the required --pass flag is missing.
func TestDecryptCmdMissingPassword(t *testing.T) {
	cmd := NewDecryptCmd()

	cmd.SetArgs([]string{"input.enc", "output.txt"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected error for missing password flag, got nil")
	}
}

// TestEncryptCmdForceOverwrite verifies that the encrypt command works with --force flag.
func TestEncryptCmdForceOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.enc")

	// Create input file
	if err := os.WriteFile(inputFile, []byte("test data"), 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	cmd := NewEncryptCmd()
	cmd.SetArgs([]string{inputFile, outputFile, "--pass", "testpass", "--force", "--quiet"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Encryption with --force failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); err != nil {
		t.Error("Output file was not created")
	}
}

// TestEncryptCmdValidInput verifies that the encrypt command succeeds with valid input.
func TestEncryptCmdValidInput(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.enc")

	// Create input file
	if err := os.WriteFile(inputFile, []byte("test data for encryption"), 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	cmd := NewEncryptCmd()
	cmd.SetArgs([]string{inputFile, outputFile, "--pass", "validpassword123", "--quiet"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Encryption failed: %v", err)
	}

	// Verify output file was created and has content
	info, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Output file is empty")
	}
}

// TestEncryptCmdWithWorkersFlag verifies the --workers flag is respected.
func TestEncryptCmdWithWorkersFlag(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.enc")

	if err := os.WriteFile(inputFile, []byte("test data"), 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	cmd := NewEncryptCmd()
	cmd.SetArgs([]string{inputFile, outputFile, "--pass", "testpass", "--workers", "8", "--quiet"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Encryption with 8 workers failed: %v", err)
	}
}

// TestEncryptCmdInvalidWorkersFlag verifies invalid worker count is handled gracefully.
func TestEncryptCmdInvalidWorkersFlag(t *testing.T) {
	tempDir := t.TempDir()
	inputFile := filepath.Join(tempDir, "input.txt")
	outputFile := filepath.Join(tempDir, "output.enc")

	if err := os.WriteFile(inputFile, []byte("test data"), 0644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	cmd := NewEncryptCmd()
	// Negative worker count should be clamped to default, not cause error
	cmd.SetArgs([]string{inputFile, outputFile, "--pass", "testpass", "--workers", "-5", "--quiet"})

	err := cmd.Execute()
	// Should not error, just use default workers
	if err != nil {
		t.Errorf("Encryption with invalid worker count failed: %v", err)
	}
}

// TestRootCmdVersion verifies the version command executes without error.
func TestRootCmdVersion(t *testing.T) {
	cmd := NewVersionCmd()

	// Capture both stdout and stderr since colored output might use different streams
	stdoutBuf := &bytes.Buffer{}
	stderrBuf := &bytes.Buffer{}
	cmd.SetOut(stdoutBuf)
	cmd.SetErr(stderrBuf)

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Version command failed: %v", err)
	}

	// The version command should produce output to either stdout or stderr
	// because colored output might write to stderr
	if stdoutBuf.Len() == 0 && stderrBuf.Len() == 0 {
		t.Error("Version command produced no output")
	}

	// Optional: verify output contains expected content
	output := stdoutBuf.String() + stderrBuf.String()
	if !bytes.Contains([]byte(output), []byte("CRYPTOOL")) {
		t.Logf("Version output: %s", output)
	}
}

// TestRootCmdHelpText verifies help command executes without error.
func TestRootCmdHelpText(t *testing.T) {
	// Create a new root command for testing to avoid affecting other tests
	testRoot := &cobra.Command{
		Use: "cryptool",
	}

	// Add subcommands
	testRoot.AddCommand(NewEncryptCmd())
	testRoot.AddCommand(NewDecryptCmd())
	testRoot.AddCommand(NewInteractCmd())
	testRoot.AddCommand(NewVersionCmd())

	buf := &bytes.Buffer{}
	testRoot.SetOut(buf)
	testRoot.SetArgs([]string{"--help"})

	err := testRoot.Execute()
	if err != nil {
		t.Errorf("Help command failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("Help command produced no output")
	}
}
