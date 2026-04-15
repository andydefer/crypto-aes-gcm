// Package ui provides terminal user interface utilities.
//
// This package contains UI components for the aescryptool CLI including:
//   - Colorized output for different message types (info, success, error, warning)
//   - Progress bars for long-running operations
//   - Interactive prompts for file paths, passwords, and user confirmations
//   - Banner displays for application headers and interactive mode
package ui

import (
	"os"

	"github.com/schollz/progressbar/v3"
)

// ProgressBar defines the interface for progress tracking operations.
//
// This interface abstracts the underlying progress bar implementation,
// allowing for different progress bar implementations or a no-op version
// for quiet mode operations.
type ProgressBar interface {
	// Set64 updates the progress bar to the specified value.
	// Returns an error if the operation fails.
	Set64(int64) error

	// Finish completes the progress bar, rendering it as 100% complete.
	// Returns an error if the operation fails.
	Finish() error

	// Clear removes the progress bar from the terminal.
	// Returns an error if the operation fails.
	Clear() error
}

// CreateProgressBar initializes a progress bar for file operations.
//
// The progress bar displays:
//   - A description of the current operation
//   - A visual bar with percentage completion
//   - Current and total counts (bytes processed)
//   - Operations per second (its/s)
//
// Parameters:
//   - total: Total number of bytes to process (used as 100%)
//   - description: Text displayed before the progress bar (e.g., "🔒 Encrypting")
//
// Returns:
//   - ProgressBar: A configured progress bar ready for use
//
// Example:
//
//	bar := CreateProgressBar(constants.MB, "🔒 Encrypting")
//	defer bar.Finish()
//	bar.Set64(512*1024) // 50% complete
func CreateProgressBar(total int64, description string) ProgressBar {
	return progressbar.NewOptions64(
		total,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(40),
		progressbar.OptionThrottle(65),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)
}
