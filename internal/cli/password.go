// Package cli provides password handling with interactive prompts for encryption/decryption operations.
package cli

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"syscall"

	"github.com/andydefer/crypto-aes-gcm/internal/lang"
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
	fmt.Print(lang.T(lang.PasswordPrompt))
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("%s: %w", lang.T(lang.PasswordReadError), err)
	}
	fmt.Println()

	password := strings.TrimSpace(string(bytePassword))

	if password == "" {
		return "", errors.New(lang.T(lang.PasswordEmpty))
	}

	return password, nil
}

// promptPasswordWithConfirm asks for password twice and validates strength (for encryption mode).
//
// Returns:
//   - string: The confirmed password (trimmed)
//   - error: If password reading fails, validation fails, or passwords don't match
func promptPasswordWithConfirm() (string, error) {
	fmt.Print(lang.T(lang.PasswordPrompt))
	bytePassword1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("%s: %w", lang.T(lang.PasswordReadError), err)
	}
	fmt.Println()

	password1 := strings.TrimSpace(string(bytePassword1))

	if err := validatePasswordStrength(password1); err != nil {
		return "", err
	}

	fmt.Print(lang.T(lang.PasswordConfirmPrompt))
	bytePassword2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("%s: %w", lang.T(lang.PasswordConfirmError), err)
	}
	fmt.Println()

	password2 := strings.TrimSpace(string(bytePassword2))

	if password1 != password2 {
		return "", errors.New(lang.T(lang.PasswordNotMatch))
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
		return errors.New(lang.T(lang.PasswordMinLength))
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New(lang.T(lang.PasswordUppercase))
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New(lang.T(lang.PasswordLowercase))
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New(lang.T(lang.PasswordDigit))
	}
	return nil
}
