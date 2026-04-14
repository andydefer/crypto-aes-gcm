// Package ui provides terminal user interface utilities for aescryptool.
//
// This package handles all user interaction including:
//   - Colored output for different message types (info, success, error, warning)
//   - Progress bars for long-running operations
//   - Interactive prompts for file paths, passwords, and confirmations
//   - Banner displays for interactive mode
//
// All UI functions are designed to work consistently across different terminals
// and operating systems.
package ui

import (
	"testing"
)

// TestColors verifies that all color variables are properly initialized.
//
// This test ensures that the color package initialized correctly and that
// all exported color variables (InfoColor, SuccessColor, ErrorColor,
// WarningColor, HeaderColor) are non-nil.
func TestColors(t *testing.T) {
	if InfoColor == nil {
		t.Error("InfoColor is nil")
	}
	if SuccessColor == nil {
		t.Error("SuccessColor is nil")
	}
	if ErrorColor == nil {
		t.Error("ErrorColor is nil")
	}
	if WarningColor == nil {
		t.Error("WarningColor is nil")
	}
	if HeaderColor == nil {
		t.Error("HeaderColor is nil")
	}
}

// TestCreateProgressBar verifies that progress bar creation returns a non-nil value.
//
// This test ensures that CreateProgressBar properly initializes a progress bar
// with the given total size and description, and that the returned object
// implements the ProgressBar interface.
func TestCreateProgressBar(t *testing.T) {
	bar := CreateProgressBar(1024, "Testing")
	if bar == nil {
		t.Error("CreateProgressBar returned nil")
	}
}

// TestPrintFunctions verifies that all print functions execute without panicking.
//
// This test calls each UI print function to ensure they handle their output
// correctly and don't cause runtime panics due to nil pointers or other issues.
// The functions tested include:
//   - PrintInteractiveHeader - Interactive mode welcome banner
//   - PrintEncryptHeader - Encryption operation header
//   - PrintDecryptHeader - Decryption operation header
//   - PrintInteractiveGoodbye - Exit message
//   - PrintSuccess - Operation success message with file info
func TestPrintFunctions(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Print function panicked: %v", r)
		}
	}()

	PrintInteractiveHeader()
	PrintEncryptHeader()
	PrintDecryptHeader()
	PrintInteractiveGoodbye()
	PrintSuccess("test.txt", 1024)
}
