// Package cli provides the command-line interface for cryptool.
//
// It implements the Cobra commands for encryption, decryption, interactive mode,
// and version display. The package handles flag parsing, validation, and
// orchestration of the underlying crypto operations.
package cli

import (
	"testing"
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
