// Package utils provides utility functions for aescryptool.
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

import (
	"testing"

	"github.com/andydefer/crypto-aes-gcm/internal/constants"
)

// TestFormatFileSize verifies that byte counts are correctly formatted into
// human-readable strings with appropriate units (B, KB, MB, GB).
//
// This test covers:
//   - Byte-level precision (0-1023 bytes)
//   - Kilobyte conversion (1024 bytes and above)
//   - Megabyte conversion (1,048,576 bytes and above)
//   - Gigabyte conversion (1,073,741,824 bytes and above)
//   - Edge cases (zero bytes, single byte, boundary values)
//
// The test uses table-driven testing to ensure consistent formatting
// across all size ranges.
func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"bytes", 500, "500 B"},
		{"kilobytes", 1500, "1.46 KB"},
		{"megabytes", 5 * constants.MB, "5.00 MB"},
		{"gigabytes", 3 * constants.GB, "3.00 GB"},
		{"zero bytes", 0, "0 B"},
		{"1 byte", 1, "1 B"},
		{"1023 bytes", 1023, "1023 B"},
		{"1025 bytes", 1025, "1.00 KB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFileSize(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatFileSize(%d) = %q, want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}
