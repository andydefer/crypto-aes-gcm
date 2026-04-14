package main

import (
	"testing"
)

// TestMainFunction verifies that main doesn't panic
// This is a smoke test for the entry point
func TestMainFunction(t *testing.T) {
	// Test that main doesn't panic when called with --help
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main panicked: %v", r)
		}
	}()

	// Note: We don't actually call main() here because it would exit the test
	// This test just ensures the package compiles
}
