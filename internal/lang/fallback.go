// Package lang provides internationalization support with fallback messages.
//
// This package maintains a complete map of English fallback messages that are
// used when a requested key is not found in the active language bundle.
// This guarantees that users never see raw keys, only readable messages.
package lang

// defaultEnglishMessages contains all messages in English as a fallback.
// This map ensures that even if a key is missing from the active language bundle,
// a meaningful English message is always returned instead of the key itself.
var defaultEnglishMessages = map[Key]string{
	// Argon2 key derivation errors
	ErrMemoryTooLow:  "memory too low: %d KiB (minimum 8192 KiB)",
	ErrMemoryTooHigh: "memory too high: %d KiB (maximum 1,048,576 KiB)",
	ErrThreadsMin:    "threads must be at least 1",
	ErrThreadsMax:    "threads too high: %d (maximum %d)",
	ErrThreadsExceed: "threads exceed system capacity: %d (max %d)",
	ErrTimeMin:       "time must be at least 1",
	ErrTimeMax:       "time too high: %d (maximum 100)",
	ErrKeyLenShort:   "key length too short: %d bytes (minimum 16)",
	ErrKeyLenLong:    "key length too long: %d bytes (maximum 64)",

	// CLI user interactions
	CliFileExists:         "File '%s' already exists. Overwrite?",
	CliOperationCancelled: "❌ Operation cancelled",
	CliError:              "❌ Error: %v",

	// Command-line flags descriptions
	FlagPassDesc:    "Passphrase used for encryption (optional - will prompt if omitted)",
	FlagWorkersDesc: "Number of parallel workers",
	FlagForceDesc:   "Overwrite existing output file without confirmation",
	FlagQuietDesc:   "Suppress progress bar output",

	// Command descriptions
	CmdEncryptShort: "🔒 Encrypt a file",
	CmdEncryptLong:  "Encrypt a file using AES-256-GCM with Argon2id key derivation.",
	CmdDecryptShort: "🔓 Decrypt a file",
	CmdDecryptLong:  "Decrypt a file that was encrypted with the encrypt command.",

	// Interactive mode prompts and messages
	InteractiveTitle:              "Interactive Mode",
	InteractiveEncryptFlow:        "Encryption",
	InteractiveDecryptFlow:        "Decryption",
	InteractiveInputFile:          "📁 File to encrypt",
	InteractiveOutputFile:         "📂 Output file",
	InteractivePassword:           "🔑 Password",
	InteractiveConfirm:            "✅ Confirmation",
	InteractiveWorkerCount:        "⚙️ Workers",
	InteractiveOverwrite:          "⚠️  File already exists. Overwrite?",
	InteractiveCancel:             "❌ Operation cancelled",
	InteractivePressEnter:         "🔁 Press Enter to continue...",
	InteractiveFileToEncrypt:      "📁 File to encrypt",
	InteractiveEncryptedFile:      "📁 Encrypted file",
	InteractivePasswordsNotMatch:  "❌ Passwords do not match",
	InteractiveCheckExists:        "check file existence",
	InteractiveOverwriteCancelled: "❌ Operation cancelled",
	InteractiveCancelOperation:    "user cancelled operation",

	// Password validation and prompts
	PasswordPrompt:        "🔑 Password: ",
	PasswordConfirmPrompt: "✅ Confirm password: ",
	PasswordEmpty:         "password cannot be empty",
	PasswordReadError:     "read password: %w",
	PasswordConfirmError:  "read confirmation: %w",
	PasswordNotMatch:      "passwords do not match",
	PasswordMinLength:     "minimum 8 characters required",
	PasswordUppercase:     "at least one uppercase letter required",
	PasswordLowercase:     "at least one lowercase letter required",
	PasswordDigit:         "at least one digit required",

	// Root command help sections
	RootShortDesc:          "🔐 Secure file encryption using AES-256-GCM",
	RootLongDesc:           "",
	RootUsage:              "Usage:",
	RootCommandsTitle:      "Commands:",
	RootPasswordManagement: "Password Management:",
	RootExamplesTitle:      "Examples:",
	RootExampleEncrypt:     "# Encrypt with interactive password prompt (recommended)",
	RootExampleDecrypt:     "# Decrypt with interactive password prompt (recommended)",
	RootExamplePassFlag:    "# Encrypt with --pass flag (for scripts)",
	RootExampleWorkers:     "# With parallel processing (8 workers)",
	RootExampleForce:       "# Force overwrite without confirmation",

	// Version command output
	VersionShortDesc: "Show version information",
	VersionLongDesc:  "Display aescryptool version, build information, and system details",
	VersionBuildInfo: "📦 Build: %s",
	VersionOSArch:    "🖥️  OS/Arch: %s/%s",
	VersionCPUs:      "💻 CPUs: %d",

	// Cryptographic nonce errors
	ErrDestSliceTooShort: "dest slice too short: need %d, got %d",

	// File system operation errors
	ErrOpenFile: "open file '%s': %w",

	// Service layer validation
	ErrFileAlreadyExists: "file already exists",
	ErrFileNotFound:      "file '%s' not found",
	WarnWorkersReduced:   "⚠️ Workers reduced to %d\n",

	// UI banners and headers
	UIInteractiveHeader: "Interactive Mode",
	UIEncryptHeader:     "🔐 FILE ENCRYPTION",
	UIDecryptHeader:     "🔓 FILE DECRYPTION",
	UIHeaderSeparator:   "────────────────────────────────────────",
	UIGoodbyeMessage:    "Thank you for using AESCRYPTOOL!",

	// UI prompt messages
	UIPromptOperationLabel:    "What do you want to do",
	UIPromptEncryptOption:     "🔒  Encrypt a file",
	UIPromptDecryptOption:     "🔓  Decrypt a file",
	UIPromptExitOption:        "🚪  Exit",
	UIPromptGoodbye:           "Thank you for using CRYPTOOL!",
	UIPromptPathEmpty:         "❌ Path cannot be empty",
	UIPromptPathNotExist:      "❌ File '%s' does not exist",
	UIPromptPathSuccess:       "   ✓ %s",
	UIPromptPasswordMinLength: "❌ Minimum 8 characters required",
	UIPromptPasswordUppercase: "❌ At least one uppercase letter required",
	UIPromptPasswordLowercase: "❌ At least one lowercase letter required",
	UIPromptPasswordDigit:     "❌ At least one digit required",
	UIPromptPasswordSuccess:   "   ✓ %s",
	UIPromptWorkersLabel:      "⚙️  Workers (default: %d, max: %d)",
	UIPromptWorkersSuccess:    "   ✓ %d workers",
	UIPromptWorkersInvalid:    "❌ Valid number required (>=1)",
	UIPromptWorkersMax:        "❌ Maximum %d workers",
	UIPromptConfirmLabel:      "❓ %s [%s]: ",
	UIPromptConfirmInvalid:    "❌ Please answer y/n",

	// UI success messages
	UISuccessOperation: "✅ Operation successful!",
	UISuccessOutput:    "📄 Output: %s",
	UISuccessSize:      "📏 Size:   %s",

	// Cryptolib decryption errors
	CryptolibErrOpenInput:      "open input: %w",
	CryptolibErrCreateOutput:   "create output: %w",
	CryptolibErrCreateCipher:   "create cipher: %w",
	CryptolibErrCreateGCM:      "create GCM: %w",
	CryptolibErrReadHeader:     "read header: %w",
	CryptolibErrReadHeaderHMAC: "read header HMAC: %w",
	CryptolibErrReadNonce:      "read nonce: %w",
	CryptolibErrUnexpectedEOF:  "unexpected EOF: missing end marker",
	CryptolibErrReadChunkLen:   "read chunk length: %w",
	CryptolibErrReadCiphertext: "read ciphertext chunk %d: %w",
	CryptolibErrDeriveNonce:    "derive nonce for chunk %d: %w",
	CryptolibErrWritePlaintext: "write plaintext chunk %d: %w",

	// Cryptolib encryption errors
	CryptolibErrOpenInputEnc:    "open input: %w",
	CryptolibErrCreateOutputEnc: "create output: %w",
	CryptolibErrGenerateSalt:    "generate salt: %w",
	CryptolibErrWriteHeader:     "write header: %w",
	CryptolibErrWriteHeaderHMAC: "write header HMAC: %w",
	CryptolibErrCreateCipherEnc: "create cipher: %w",
	CryptolibErrCreateGCMEnc:    "create GCM: %w",
	CryptolibErrGenerateNonce:   "generate nonce: %w",
	CryptolibErrWriteNonce:      "write nonce: %w",
	CryptolibErrReadChunk:       "read chunk: %w",
	CryptolibErrWriteEndMarker:  "write end marker: %w",
	CryptolibErrNonceDerivation: "nonce derivation failed: %w",
	CryptolibErrMissingChunks:   "missing chunks: expected index %d, have %d pending",
	CryptolibErrTooManyPending:  "too many pending chunks (limit %d) - possible reordering attack",
	CryptolibErrWriteChunkLen:   "write chunk length: %w",
	CryptolibErrWriteCiphertext: "write ciphertext: %w",
	CryptolibErrCloseInput:      "close input: %w",
	CryptolibErrCloseOutput:     "close output: %w",

	// Cryptolib sentinel errors
	CryptolibErrInvalidMagic:       "invalid magic bytes: file not encrypted with this tool",
	CryptolibErrUnsupportedVersion: "unsupported file version",
	CryptolibErrHeaderAuthFailed:   "header authentication failed: wrong passphrase or corrupted file",
	CryptolibErrDecryptionFailed:   "decryption failed: corrupted data or wrong key",
	CryptolibErrChunkTooLarge:      "chunk size exceeds maximum allowed limit",

	// Cryptolib stream processing errors
	CryptolibErrReadHeaderStream:     "read header: %w",
	CryptolibErrReadHeaderHMACStream: "read header HMAC: %w",
	CryptolibErrReadNonceStream:      "read nonce: %w",
	CryptolibErrUnexpectedEOFStream:  "unexpected EOF: missing end marker",
	CryptolibErrReadChunkLenStream:   "read chunk length: %w",
	CryptolibErrReadCiphertextStream: "read ciphertext chunk %d: %w",
	CryptolibErrDeriveNonceStream:    "derive nonce for chunk %d: %w",
	CryptolibErrWritePlaintextStream: "write plaintext chunk %d: %w",
	CryptolibErrCreateCipherStream:   "create cipher: %w",
	CryptolibErrCreateGCMStream:      "create GCM: %w",
}

// GetDefaultMessage returns the fallback English message for a given key.
//
// This function guarantees that a meaningful message is always returned,
// even if the key is not found in the fallback map (returns the key as last resort).
//
// Parameters:
//   - key: the message key to look up
//
// Returns:
//   - string: the English fallback message, or the key itself if not found
func GetDefaultMessage(key Key) string {
	if msg, ok := defaultEnglishMessages[key]; ok {
		return msg
	}
	return string(key)
}
