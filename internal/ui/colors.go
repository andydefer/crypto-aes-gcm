// Package ui provides terminal user interface utilities.
package ui

import "github.com/fatih/color"

var (
	// InfoColor is used for informational messages.
	InfoColor = color.New(color.FgCyan, color.Bold)

	// SuccessColor is used for success messages.
	SuccessColor = color.New(color.FgGreen, color.Bold)

	// ErrorColor is used for error messages.
	ErrorColor = color.New(color.FgRed, color.Bold)

	// WarningColor is used for warning messages.
	WarningColor = color.New(color.FgYellow)

	// HeaderColor is used for header text.
	HeaderColor = color.New(color.FgMagenta, color.Bold)
)
