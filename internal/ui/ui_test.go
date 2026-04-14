package ui

import (
	"testing"
)

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

func TestCreateProgressBar(t *testing.T) {
	bar := CreateProgressBar(1024, "Testing")
	if bar == nil {
		t.Error("CreateProgressBar returned nil")
	}
}

func TestPrintFunctions(t *testing.T) {
	// These functions should not panic
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
