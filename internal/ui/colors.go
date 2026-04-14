// Package ui provides terminal user interface utilities for aescryptool.
//
// This package centralizes all terminal UI concerns including:
//   - Color management for consistent visual feedback
//   - Progress bars for long-running operations
//   - Interactive prompts for user input
//   - Banner and header display
//   - Success and error message formatting
//
// All UI elements are designed to be user-friendly and provide clear
// visual feedback for encryption and decryption operations.
package ui

import "github.com/fatih/color"

var (
	// InfoColor is used for informational messages.
	//
	// Displays text in cyan with bold formatting.
	// Used for: progress updates, file paths, operation headers.
	InfoColor = color.New(color.FgCyan, color.Bold)

	// SuccessColor is used for success messages.
	//
	// Displays text in green with bold formatting.
	// Used for: successful operations, completion messages.
	SuccessColor = color.New(color.FgGreen, color.Bold)

	// ErrorColor is used for error messages.
	//
	// Displays text in red with bold formatting.
	// Used for: operation failures, validation errors, critical issues.
	ErrorColor = color.New(color.FgRed, color.Bold)

	// WarningColor is used for warning messages.
	//
	// Displays text in yellow (normal weight).
	// Used for: overwrite confirmations, non-critical issues.
	WarningColor = color.New(color.FgYellow)

	// HeaderColor is used for header text and banners.
	//
	// Displays text in magenta with bold formatting.
	// Used for: application banners, section separators.
	HeaderColor = color.New(color.FgMagenta, color.Bold)
)
