// Package lang provides internationalization support with thread-safe message retrieval.
//
// This package manages the active language bundle and provides the main entry point
// T() for retrieving localized messages throughout the application.
package lang

import (
	"fmt"
	"strings"
	"sync"
)

// Language represents a supported UI language.
type Language string

const (
	// English represents the English language.
	English Language = "en"
	// French represents the French language.
	French Language = "fr"
)

var (
	currentBundle MessageBundle
	currentLang   Language
	mu            sync.RWMutex
)

func init() {
	SetLanguage(English)
}

// SetLanguage changes the active language for all subsequent T() calls.
//
// It is safe for concurrent use.
//
// Parameters:
//   - lang: the language to activate (English or French)
func SetLanguage(lang Language) {
	mu.Lock()
	defer mu.Unlock()

	currentLang = lang

	switch lang {
	case French:
		currentBundle = &FrenchBundle{}
	default:
		currentBundle = &EnglishBundle{}
	}
}

// GetLanguage returns the currently active language.
//
// It is safe for concurrent use.
//
// Returns:
//   - Language: the active language (English or French)
func GetLanguage() Language {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

// T returns a formatted message for the given key in the currently active language.
//
// This is the main entry point for message localization. It retrieves the message
// from the active language bundle and formats it with the provided arguments.
//
// If the message cannot be retrieved (e.g., bundle is nil or key not found),
// it falls back to the default English message from GetDefaultMessage().
//
// Parameters:
//   - key: the message key to look up
//   - args: optional formatting arguments (uses fmt.Sprintf semantics)
//
// Returns:
//   - string: the formatted localized message
//
// Example:
//
//	msg := lang.T(lang.ErrFileNotFound, "myfile.txt")
func T(key Key, args ...interface{}) string {
	mu.RLock()
	defer mu.RUnlock()

	var msg string
	if currentBundle == nil {
		msg = GetDefaultMessage(key)
	} else {
		msg = currentBundle.GetMessage(key, args...)
	}

	// Fallback to default English if the message equals the key itself.
	if msg == string(key) {
		msg = GetDefaultMessage(key)
	}

	if len(args) > 0 && strings.Contains(msg, "%") {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// SupportedLanguages returns a slice of all available languages.
//
// Returns:
//   - []Language: list of supported languages (English, French)
func SupportedLanguages() []Language {
	return []Language{English, French}
}
