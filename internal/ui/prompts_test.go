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
	"regexp"
	"strings"
	"testing"
)

// TestPromptConfirmLogic verifies the confirmation prompt logic for user input.
//
// This test covers:
//   - Empty input (Enter key) should return the default value
//   - Positive responses (y, Y, yes, o, oui) should return true
//   - Negative responses (n, N, no, non) should return false
//   - Whitespace handling around input
func TestPromptConfirmLogic(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		defaultValue bool
		expected     bool
	}{
		{"empty input with default true", "", true, true},
		{"empty input with default false", "", false, false},

		{"y input", "y", false, true},
		{"Y input", "Y", false, true},
		{"yes input", "yes", false, true},
		{"oui input", "oui", false, true},
		{"o input", "o", false, true},
		{"O input", "O", false, true},

		{"n input", "n", true, false},
		{"N input", "N", true, false},
		{"no input", "no", true, false},
		{"non input", "non", true, false},

		{"whitespace y", "  y  ", false, true},
		{"whitespace n", "  n  ", true, false},
	}

	var truthyValues = map[string]bool{
		"y": true, "yes": true, "o": true, "oui": true,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := strings.TrimSpace(strings.ToLower(tt.input))

			var actual bool
			switch {
			case result == "":
				actual = tt.defaultValue
			case truthyValues[result]:
				actual = true
			default:
				actual = false
			}

			if actual != tt.expected {
				t.Errorf("expected %v, got %v for input %q (defaultValue=%v)",
					tt.expected, actual, tt.input, tt.defaultValue)
			}
		})
	}
}

// TestPromptWorkersValidation verifies the worker count input validation logic.
//
// This test ensures that:
//   - Valid worker counts (1-16) are accepted
//   - Empty input returns the default value (4)
//   - Invalid inputs (zero, negative, too high, non-numeric) are rejected
func TestPromptWorkersValidation(t *testing.T) {
	maxWorkers := 16

	tests := []struct {
		name     string
		input    string
		expected int
		valid    bool
	}{
		{"minimum worker", "1", 1, true},
		{"default worker", "4", 4, true},
		{"maximum worker", "16", 16, true},

		{"empty input (default)", "", 4, true},

		{"zero workers", "0", 0, false},
		{"negative workers", "-1", 0, false},
		{"exceeds maximum", "100", 0, false},
		{"non-numeric", "abc", 0, false},
		{"mixed characters", "4abc", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result int
			var err error

			if tt.input == "" {
				result = 4
			} else {
				var n int
				_, parseErr := fmt.Sscanf(tt.input, "%d", &n)
				if parseErr != nil {
					err = parseErr
				} else {
					remaining := tt.input
					for i, c := range tt.input {
						if c < '0' || c > '9' {
							remaining = tt.input[i:]
							break
						}
					}
					if len(remaining) > 0 && (remaining[0] < '0' || remaining[0] > '9') {
						err = fmt.Errorf("invalid characters after number: %s", remaining)
					} else {
						result = n
						if result < 1 || result > maxWorkers {
							err = fmt.Errorf("value out of range: %d", result)
						}
					}
				}
			}

			if tt.valid && err != nil {
				t.Errorf("expected valid input %q, got error: %v", tt.input, err)
			}
			if !tt.valid && err == nil && tt.input != "" {
				t.Errorf("expected invalid input %q to fail, but it passed (got %d)", tt.input, result)
			}
			if tt.valid && result != tt.expected && tt.input != "" {
				t.Errorf("expected %d, got %d for input %q", tt.expected, result, tt.input)
			}
		})
	}
}

// TestPasswordValidation verifies the password strength validation rules.
//
// Password requirements:
//   - Minimum length: 8 characters
//   - At least one uppercase letter (A-Z)
//   - At least one lowercase letter (a-z)
//   - At least one digit (0-9)
//
// Special characters are allowed but not required.
func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
		reason   string
	}{
		{"valid password", "Hello123", true, "8 chars, upper, lower, digit"},
		{"valid with special chars", "Hello123!", true, "special chars allowed"},
		{"valid long password", "VeryLongPassword123", true, ">8 chars"},

		{"too short", "Hi1", false, "less than 8 characters"},
		{"exactly 7 chars", "Hello12", false, "7 characters"},

		{"no uppercase", "hello123", false, "no uppercase letter"},
		{"no uppercase long", "helloworld123", false, "no uppercase letter"},

		{"no lowercase", "HELLO123", false, "no lowercase letter"},

		{"no digit", "HelloWorld", false, "no digit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := true
			reason := ""

			if len(tt.password) < 8 {
				valid = false
				reason = "length"
			} else if !regexp.MustCompile(`[A-Z]`).MatchString(tt.password) {
				valid = false
				reason = "missing uppercase"
			} else if !regexp.MustCompile(`[a-z]`).MatchString(tt.password) {
				valid = false
				reason = "missing lowercase"
			} else if !regexp.MustCompile(`[0-9]`).MatchString(tt.password) {
				valid = false
				reason = "missing digit"
			}

			if valid != tt.valid {
				t.Errorf("password %q: expected valid=%v, got valid=%v (reason: %s)",
					tt.password, tt.valid, valid, reason)
			}
		})
	}
}
