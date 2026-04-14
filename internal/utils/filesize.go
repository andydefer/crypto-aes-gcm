// Package utils provides utility functions.
package utils

import "fmt"

// FormatFileSize converts bytes to human-readable format.
func FormatFileSize(bytes int64) string {
	switch {
	case bytes > 1024*1024*1024:
		return fmt.Sprintf("%.2f GB", float64(bytes)/(1024*1024*1024))
	case bytes > 1024*1024:
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1024*1024))
	case bytes > 1024:
		return fmt.Sprintf("%.2f KB", float64(bytes)/1024)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
