// Package cli provides password handling with interactive prompts for encryption/decryption operations.
package cli

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// ResolvePassword retrieves the password from flags or prompts interactively.
//
// Parameters:
//   - flagPass: Password provided via command-line flag (empty if not provided)
//   - needConfirmation: If true, prompts twice and verifies match (encryption mode);
//     if false, prompts once (decryption mode)
//
// Returns:
//   - string: The resolved password
//   - error: If password validation fails or user input cannot be read
//
// When flagPass is non-empty, it's used directly without prompting.
// Flag-provided passwords skip confirmation to maintain script compatibility.
func ResolvePassword(flagPass string, needConfirmation bool) (string, error) {
	if flagPass != "" {
		return flagPass, nil
	}

	if needConfirmation {
		return promptPasswordWithConfirm()
	}
	return promptPassword()
}

// promptPassword asks for a password once (for decryption mode).
//
// Returns:
//   - string: The entered password (trimmed)
//   - error: If password reading fails or password is empty
func promptPassword() (string, error) {
	fmt.Print("🔑 Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}
	fmt.Println()

	password := strings.TrimSpace(string(bytePassword))

	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return password, nil
}

// promptPasswordWithConfirm asks for password twice and validates strength (for encryption mode).
//
// Returns:
//   - string: The confirmed password (trimmed)
//   - error: If password reading fails, validation fails, or passwords don't match
func promptPasswordWithConfirm() (string, error) {
	fmt.Print("🔑 Password: ")
	bytePassword1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}
	fmt.Println()

	password1 := strings.TrimSpace(string(bytePassword1))

	if err := validatePasswordStrength(password1); err != nil {
		return "", err
	}

	fmt.Print("✅ Confirm password: ")
	bytePassword2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("read confirmation: %w", err)
	}
	fmt.Println()

	password2 := strings.TrimSpace(string(bytePassword2))

	if password1 != password2 {
		return "", fmt.Errorf("passwords do not match")
	}

	return password1, nil
}

// validatePasswordStrength checks if password meets security requirements.
//
// Requirements:
//   - Minimum length: 8 characters
//   - At least one uppercase letter (A-Z)
//   - At least one lowercase letter (a-z)
//   - At least one digit (0-9)
//
// Returns an error with a descriptive message if any requirement is not met.
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("minimum 8 characters required")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("at least one uppercase letter required")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("at least one lowercase letter required")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("at least one digit required")
	}
	return nil
}
