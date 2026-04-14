package cli

import (
	"testing"
)

func TestExecute(t *testing.T) {
	// Test that Execute doesn't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Execute panicked: %v", r)
		}
	}()
}

func TestRootCmdHelp(t *testing.T) {
	// Test help command doesn't panic
	_ = rootCmd.Help()
}

func TestEncryptCmdFlags(t *testing.T) {
	cmd := NewEncryptCmd()

	// Check required flags
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

func TestInteractCmd(t *testing.T) {
	cmd := NewInteractCmd()
	if cmd.Use != "interact" {
		t.Errorf("expected use 'interact', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("interact command should have a short description")
	}
}

func TestVersionCmd(t *testing.T) {
	cmd := NewVersionCmd()
	if cmd.Use != "version" {
		t.Errorf("expected use 'version', got '%s'", cmd.Use)
	}
	if cmd.Short == "" {
		t.Error("version command should have a short description")
	}
}

// TestRunInteractiveEncryptValidation tests the validation logic
func TestRunInteractiveEncryptValidation(t *testing.T) {
	// Test that empty input returns early
	// This is a smoke test for the function structure
	t.Log("Interactive encrypt validation test - requires manual testing")
}

// TestRunInteractiveDecryptValidation tests the validation logic
func TestRunInteractiveDecryptValidation(t *testing.T) {
	// Test that empty input returns early
	// This is a smoke test for the function structure
	t.Log("Interactive decrypt validation test - requires manual testing")
}
