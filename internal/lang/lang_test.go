// Package lang provides internationalization support tests.
package lang

import (
	"testing"
)

// TestAllBundlesImplementAllKeys verifies that both language bundles
// implement all required methods. Compile-time interface checks provide
// the primary protection; runtime checks ensure non-empty return values.
func TestAllBundlesImplementAllKeys(t *testing.T) {
	// Compile-time interface implementation checks
	var _ MessageBundle = (*EnglishBundle)(nil)
	var _ MessageBundle = (*FrenchBundle)(nil)

	// Runtime validation of English bundle
	eng := &EnglishBundle{}
	testBundle(t, eng, "English")

	// Runtime validation of French bundle
	fr := &FrenchBundle{}
	testBundle(t, fr, "French")
}

// testBundle verifies that all message getters return non-empty strings.
func testBundle(t *testing.T, bundle MessageBundle, name string) {
	t.Run(name, func(t *testing.T) {
		// Argon2 key derivation errors
		if bundle.GetErrMemoryTooLow() == "" {
			t.Error("GetErrMemoryTooLow returned empty")
		}
		if bundle.GetErrMemoryTooHigh() == "" {
			t.Error("GetErrMemoryTooHigh returned empty")
		}
		if bundle.GetErrThreadsMin() == "" {
			t.Error("GetErrThreadsMin returned empty")
		}
		if bundle.GetErrThreadsMax() == "" {
			t.Error("GetErrThreadsMax returned empty")
		}
		if bundle.GetErrThreadsExceed() == "" {
			t.Error("GetErrThreadsExceed returned empty")
		}
		if bundle.GetErrTimeMin() == "" {
			t.Error("GetErrTimeMin returned empty")
		}
		if bundle.GetErrTimeMax() == "" {
			t.Error("GetErrTimeMax returned empty")
		}
		if bundle.GetErrKeyLenShort() == "" {
			t.Error("GetErrKeyLenShort returned empty")
		}
		if bundle.GetErrKeyLenLong() == "" {
			t.Error("GetErrKeyLenLong returned empty")
		}
	})
}

// TestTFunction verifies that the T() helper returns correctly formatted
// messages for both English and French languages.
func TestTFunction(t *testing.T) {
	// English language test
	SetLanguage(English)
	result := T(ErrMemoryTooLow, 4096)
	expected := "memory too low: 4096 KiB (minimum 8192 KiB)"
	if result != expected {
		t.Errorf("English: got %q, want %q", result, expected)
	}

	// French language test
	SetLanguage(French)
	result = T(ErrMemoryTooLow, 4096)
	expected = "mémoire trop basse: 4096 KiB (minimum 8192 KiB)"
	if result != expected {
		t.Errorf("French: got %q, want %q", result, expected)
	}
}
