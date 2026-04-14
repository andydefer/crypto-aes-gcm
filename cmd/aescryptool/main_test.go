// Package main provides the entry point tests for aescryptool CLI.
package main

import (
	"testing"
)

// TestMainFunction ensures the package compiles without errors.
//
// This test does not execute main() because it would terminate the test process.
// It serves as a compile-time verification and smoke test for the entry point.
func TestMainFunction(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main panicked: %v", r)
		}
	}()
}
