package utils

import (
	"testing"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"bytes", 500, "500 B"},
		{"kilobytes", 1500, "1.46 KB"},
		{"megabytes", 5 * 1024 * 1024, "5.00 MB"},
		{"gigabytes", 3 * 1024 * 1024 * 1024, "3.00 GB"},
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
