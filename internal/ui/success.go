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
	"fmt"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
	"github.com/andydefer/crypto-aes-gcm/internal/utils"
)

// PrintSuccess displays completion information with file size.
//
// This function prints a success message with the output file path and its
// human-readable size. The output is formatted with colors for better visibility:
//   - Success message in green
//   - File information in cyan
//
// The file size is automatically formatted using utils.FormatFileSize to
// display in B, KB, MB, or GB as appropriate.
//
// Parameters:
//   - output: Path to the output file that was created
//   - size: Size of the output file in bytes
func PrintSuccess(output string, size int64) {
	fmt.Println()
	SuccessColor.Println(lang.T(lang.UISuccessOperation))
	InfoColor.Printf(lang.T(lang.UISuccessOutput), output)
	InfoColor.Printf(lang.T(lang.UISuccessSize), utils.FormatFileSize(size))
	fmt.Println()
}
