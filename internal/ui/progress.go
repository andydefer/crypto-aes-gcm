// Package ui provides terminal user interface utilities.
package ui

import (
	"os"

	"github.com/schollz/progressbar/v3"
)

// ProgressBar interface defines the methods needed for progress tracking.
type ProgressBar interface {
	Set64(int64) error
	Finish() error
	Clear() error
}

// CreateProgressBar initializes a progress bar for file operations.
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
