// Package utils provides utility functions for cryptool.
//
// This package contains helper functions that are used across multiple
// packages in the application, including:
//   - Human-readable file size formatting
//   - Common validation logic
//   - String manipulation utilities
//
// All functions in this package are pure and have no side effects,
// making them safe for concurrent use.
package utils

import "fmt"

// FormatFileSize converts a byte count into a human-readable string with
// appropriate units (B, KB, MB, GB).
//
// The function uses binary units (1024 bytes = 1 KB) which is standard for
// file sizes in computing. The output is formatted with two decimal places
// for KB, MB, and GB to provide sufficient precision while remaining readable.
//
// Conversion rules:
//   - bytes < 1024: display as raw bytes (e.g., "500 B")
//   - bytes < 1024^2: display as KB (e.g., "1.46 KB")
//   - bytes < 1024^3: display as MB (e.g., "5.00 MB")
//   - bytes >= 1024^3: display as GB (e.g., "3.00 GB")
//
// Parameters:
//   - bytes: Number of bytes to format (can be zero or positive)
//
// Returns:
//   - A formatted string with the appropriate unit suffix
//
// Examples:
//
//	FormatFileSize(500)           // returns "500 B"
//	FormatFileSize(1500)          // returns "1.46 KB"
//	FormatFileSize(5*1024*1024)   // returns "5.00 MB"
//	FormatFileSize(3*1024*1024*1024) // returns "3.00 GB"
//	FormatFileSize(0)             // returns "0 B"
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
