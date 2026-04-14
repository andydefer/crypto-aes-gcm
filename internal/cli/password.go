// Package cli provides password handling with interactive prompts.
package cli

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// resolvePassword retrieves the password from flags or prompts interactively.
// For encryption (needConfirmation=true), it prompts twice and verifies match.
// For decryption (needConfirmation=false), it prompts once.
func resolvePassword(flagPass string, needConfirmation bool) (string, error) {
	// If password provided via flag, use it directly
	if flagPass != "" {
		if needConfirmation {
			// For encryption, even flag-provided passwords should be confirmed?
			// Skip confirmation for flag to maintain script compatibility
			return flagPass, nil
		}
		return flagPass, nil
	}

	// No flag provided - prompt interactively
	if needConfirmation {
		return promptPasswordWithConfirm()
	}
	return promptPassword()
}

// promptPassword asks for a password once (for decryption)
func promptPassword() (string, error) {
	fmt.Print("🔑 Mot de passe: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()

	password := strings.TrimSpace(string(bytePassword))

	// Basic validation for decryption (non-empty)
	if password == "" {
		return "", fmt.Errorf("le mot de passe ne peut pas être vide")
	}

	return password, nil
}

// promptPasswordWithConfirm asks for password twice and validates strength (for encryption)
func promptPasswordWithConfirm() (string, error) {
	// First password entry
	fmt.Print("🔑 Mot de passe: ")
	bytePassword1, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()

	password1 := strings.TrimSpace(string(bytePassword1))

	// Validate password strength
	if err := validatePasswordStrength(password1); err != nil {
		return "", err
	}

	// Second password entry for confirmation
	fmt.Print("✅ Confirmation du mot de passe: ")
	bytePassword2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	fmt.Println()

	password2 := strings.TrimSpace(string(bytePassword2))

	// Check if passwords match
	if password1 != password2 {
		return "", fmt.Errorf("les mots de passe ne correspondent pas")
	}

	return password1, nil
}

// validatePasswordStrength checks if password meets security requirements.
// Requirements:
//   - Minimum length: 8 characters
//   - At least one uppercase letter (A-Z)
//   - At least one lowercase letter (a-z)
//   - At least one digit (0-9)
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("8 caractères minimum requis")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return fmt.Errorf("au moins une majuscule requise")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return fmt.Errorf("au moins une minuscule requise")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return fmt.Errorf("au moins un chiffre requis")
	}
	return nil
}

// promptPasswordFromReader is used for testing with mocked input
func promptPasswordFromReader(reader *bufio.Reader, needConfirmation bool) (string, error) {
	if needConfirmation {
		fmt.Print("🔑 Mot de passe: ")
		password1, _ := reader.ReadString('\n')
		password1 = strings.TrimSpace(password1)

		if err := validatePasswordStrength(password1); err != nil {
			return "", err
		}

		fmt.Print("✅ Confirmation du mot de passe: ")
		password2, _ := reader.ReadString('\n')
		password2 = strings.TrimSpace(password2)

		if password1 != password2 {
			return "", fmt.Errorf("les mots de passe ne correspondent pas")
		}
		return password1, nil
	}

	fmt.Print("🔑 Mot de passe: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if password == "" {
		return "", fmt.Errorf("le mot de passe ne peut pas être vide")
	}
	return password, nil
}
