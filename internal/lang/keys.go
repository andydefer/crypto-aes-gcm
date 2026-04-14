// Package lang provides internationalization support with type-safe message keys.
//
// This package defines all message keys used throughout the application and
// the MessageBundle interface that each language implementation must satisfy.
package lang

// Key represents a localized message key.
//
// Each key corresponds to a specific user-facing message in the application.
// Keys are organized by domain (argon2, cli, ui, cryptolib, etc.) for clarity.
type Key string

// All message keys in the application.
const (
	// Argon2 key derivation errors
	ErrMemoryTooLow  Key = "argon2.err.memory_too_low"
	ErrMemoryTooHigh Key = "argon2.err.memory_too_high"
	ErrThreadsMin    Key = "argon2.err.threads_min"
	ErrThreadsMax    Key = "argon2.err.threads_max"
	ErrThreadsExceed Key = "argon2.err.threads_exceed"
	ErrTimeMin       Key = "argon2.err.time_min"
	ErrTimeMax       Key = "argon2.err.time_max"
	ErrKeyLenShort   Key = "argon2.err.key_len_short"
	ErrKeyLenLong    Key = "argon2.err.key_len_long"

	// CLI user interactions (shared between encrypt/decrypt)
	CliFileExists         Key = "cli.file_exists"
	CliOperationCancelled Key = "cli.operation_cancelled"
	CliError              Key = "cli.error"

	// Command-line flags descriptions
	FlagPassDesc    Key = "cli.flag.pass_desc"
	FlagWorkersDesc Key = "cli.flag.workers_desc"
	FlagForceDesc   Key = "cli.flag.force_desc"
	FlagQuietDesc   Key = "cli.flag.quiet_desc"

	// Command descriptions
	CmdEncryptShort Key = "cli.cmd.encrypt_short"
	CmdEncryptLong  Key = "cli.cmd.encrypt_long"
	CmdDecryptShort Key = "cli.cmd.decrypt_short"
	CmdDecryptLong  Key = "cli.cmd.decrypt_long"

	// Interactive mode prompts and messages
	InteractiveTitle             Key = "cli.interactive.title"
	InteractiveEncryptFlow       Key = "cli.interactive.encrypt_flow"
	InteractiveDecryptFlow       Key = "cli.interactive.decrypt_flow"
	InteractiveInputFile         Key = "cli.interactive.input_file"
	InteractiveOutputFile        Key = "cli.interactive.output_file"
	InteractivePassword          Key = "cli.interactive.password"
	InteractiveConfirm           Key = "cli.interactive.confirm"
	InteractiveWorkerCount       Key = "cli.interactive.worker_count"
	InteractiveOverwrite         Key = "cli.interactive.overwrite"
	InteractiveCancel            Key = "cli.interactive.cancel"
	InteractivePressEnter        Key = "cli.interactive.press_enter"
	InteractiveFileToEncrypt     Key = "cli.interactive.file_to_encrypt"
	InteractiveEncryptedFile     Key = "cli.interactive.encrypted_file"
	InteractivePasswordsNotMatch Key = "cli.interactive.passwords_not_match"

	// Password validation and prompts
	PasswordPrompt        Key = "cli.password.prompt"
	PasswordConfirmPrompt Key = "cli.password.confirm_prompt"
	PasswordEmpty         Key = "cli.password.empty"
	PasswordReadError     Key = "cli.password.read_error"
	PasswordConfirmError  Key = "cli.password.confirm_error"
	PasswordNotMatch      Key = "cli.password.not_match"
	PasswordMinLength     Key = "cli.password.min_length"
	PasswordUppercase     Key = "cli.password.uppercase"
	PasswordLowercase     Key = "cli.password.lowercase"
	PasswordDigit         Key = "cli.password.digit"

	// Root command help sections
	RootShortDesc          Key = "cli.root.short_desc"
	RootLongDesc           Key = "cli.root.long_desc"
	RootUsage              Key = "cli.root.usage"
	RootCommandsTitle      Key = "cli.root.commands_title"
	RootPasswordManagement Key = "cli.root.password_management"
	RootExamplesTitle      Key = "cli.root.examples_title"
	RootExampleEncrypt     Key = "cli.root.example_encrypt"
	RootExampleDecrypt     Key = "cli.root.example_decrypt"
	RootExamplePassFlag    Key = "cli.root.example_pass_flag"
	RootExampleWorkers     Key = "cli.root.example_workers"
	RootExampleForce       Key = "cli.root.example_force"

	// Version command output
	VersionShortDesc Key = "cli.version.short_desc"
	VersionLongDesc  Key = "cli.version.long_desc"
	VersionBuildInfo Key = "cli.version.build_info"
	VersionOSArch    Key = "cli.version.os_arch"
	VersionCPUs      Key = "cli.version.cpus"

	// Cryptographic nonce errors
	ErrDestSliceTooShort Key = "crypto.err.dest_slice_too_short"

	// File system operation errors
	ErrOpenFile Key = "service.err.open_file"

	// Service layer validation
	ErrFileAlreadyExists Key = "service.err.file_already_exists"
	ErrFileNotFound      Key = "service.err.file_not_found"
	WarnWorkersReduced   Key = "service.warn.workers_reduced"

	// UI banners and headers
	UIInteractiveHeader Key = "ui.banner.interactive_header"
	UIEncryptHeader     Key = "ui.banner.encrypt_header"
	UIDecryptHeader     Key = "ui.banner.decrypt_header"
	UIHeaderSeparator   Key = "ui.banner.header_separator"
	UIGoodbyeMessage    Key = "ui.banner.goodbye_message"

	// UI prompt messages
	UIPromptOperationLabel    Key = "ui.prompt.operation_label"
	UIPromptEncryptOption     Key = "ui.prompt.encrypt_option"
	UIPromptDecryptOption     Key = "ui.prompt.decrypt_option"
	UIPromptExitOption        Key = "ui.prompt.exit_option"
	UIPromptGoodbye           Key = "ui.prompt.goodbye"
	UIPromptPathEmpty         Key = "ui.prompt.path_empty"
	UIPromptPathNotExist      Key = "ui.prompt.path_not_exist"
	UIPromptPathSuccess       Key = "ui.prompt.path_success"
	UIPromptPasswordMinLength Key = "ui.prompt.password_min_length"
	UIPromptPasswordUppercase Key = "ui.prompt.password_uppercase"
	UIPromptPasswordLowercase Key = "ui.prompt.password_lowercase"
	UIPromptPasswordDigit     Key = "ui.prompt.password_digit"
	UIPromptPasswordSuccess   Key = "ui.prompt.password_success"
	UIPromptWorkersLabel      Key = "ui.prompt.workers_label"
	UIPromptWorkersSuccess    Key = "ui.prompt.workers_success"
	UIPromptWorkersInvalid    Key = "ui.prompt.workers_invalid"
	UIPromptWorkersMax        Key = "ui.prompt.workers_max"
	UIPromptConfirmLabel      Key = "ui.prompt.confirm_label"
	UIPromptConfirmInvalid    Key = "ui.prompt.confirm_invalid"

	// UI success messages
	UISuccessOperation Key = "ui.success.operation"
	UISuccessOutput    Key = "ui.success.output"
	UISuccessSize      Key = "ui.success.size"

	// Cryptolib decryption errors
	CryptolibErrOpenInput      Key = "cryptolib.err.open_input"
	CryptolibErrCreateOutput   Key = "cryptolib.err.create_output"
	CryptolibErrCreateCipher   Key = "cryptolib.err.create_cipher"
	CryptolibErrCreateGCM      Key = "cryptolib.err.create_gcm"
	CryptolibErrReadHeader     Key = "cryptolib.err.read_header"
	CryptolibErrReadHeaderHMAC Key = "cryptolib.err.read_header_hmac"
	CryptolibErrReadNonce      Key = "cryptolib.err.read_nonce"
	CryptolibErrUnexpectedEOF  Key = "cryptolib.err.unexpected_eof"
	CryptolibErrReadChunkLen   Key = "cryptolib.err.read_chunk_len"
	CryptolibErrReadCiphertext Key = "cryptolib.err.read_ciphertext"
	CryptolibErrDeriveNonce    Key = "cryptolib.err.derive_nonce"
	CryptolibErrWritePlaintext Key = "cryptolib.err.write_plaintext"

	// Cryptolib encryption errors
	CryptolibErrOpenInputEnc    Key = "cryptolib.err.open_input_enc"
	CryptolibErrCreateOutputEnc Key = "cryptolib.err.create_output_enc"
	CryptolibErrGenerateSalt    Key = "cryptolib.err.generate_salt"
	CryptolibErrWriteHeader     Key = "cryptolib.err.write_header"
	CryptolibErrWriteHeaderHMAC Key = "cryptolib.err.write_header_hmac"
	CryptolibErrCreateCipherEnc Key = "cryptolib.err.create_cipher_enc"
	CryptolibErrCreateGCMEnc    Key = "cryptolib.err.create_gcm_enc"
	CryptolibErrGenerateNonce   Key = "cryptolib.err.generate_nonce"
	CryptolibErrWriteNonce      Key = "cryptolib.err.write_nonce"
	CryptolibErrReadChunk       Key = "cryptolib.err.read_chunk"
	CryptolibErrWriteEndMarker  Key = "cryptolib.err.write_end_marker"
	CryptolibErrNonceDerivation Key = "cryptolib.err.nonce_derivation"
	CryptolibErrMissingChunks   Key = "cryptolib.err.missing_chunks"
	CryptolibErrTooManyPending  Key = "cryptolib.err.too_many_pending"
	CryptolibErrWriteChunkLen   Key = "cryptolib.err.write_chunk_len"
	CryptolibErrWriteCiphertext Key = "cryptolib.err.write_ciphertext"
	CryptolibErrCloseInput      Key = "cryptolib.err.close_input"
	CryptolibErrCloseOutput     Key = "cryptolib.err.close_output"

	// Cryptolib sentinel errors
	CryptolibErrInvalidMagic       Key = "cryptolib.err.invalid_magic"
	CryptolibErrUnsupportedVersion Key = "cryptolib.err.unsupported_version"
	CryptolibErrHeaderAuthFailed   Key = "cryptolib.err.header_auth_failed"
	CryptolibErrDecryptionFailed   Key = "cryptolib.err.decryption_failed"

	// Cryptolib stream processing errors
	CryptolibErrReadHeaderStream     Key = "cryptolib.err.read_header_stream"
	CryptolibErrReadHeaderHMACStream Key = "cryptolib.err.read_header_hmac_stream"
	CryptolibErrReadNonceStream      Key = "cryptolib.err.read_nonce_stream"
	CryptolibErrUnexpectedEOFStream  Key = "cryptolib.err.unexpected_eof_stream"
	CryptolibErrReadChunkLenStream   Key = "cryptolib.err.read_chunk_len_stream"
	CryptolibErrReadCiphertextStream Key = "cryptolib.err.read_ciphertext_stream"
	CryptolibErrDeriveNonceStream    Key = "cryptolib.err.derive_nonce_stream"
	CryptolibErrWritePlaintextStream Key = "cryptolib.err.write_plaintext_stream"
	CryptolibErrCreateCipherStream   Key = "cryptolib.err.create_cipher_stream"
	CryptolibErrCreateGCMStream      Key = "cryptolib.err.create_gcm_stream"
)

// MessageBundle defines the contract that all language bundles must implement.
//
// Each language implementation (English, French, etc.) must provide concrete
// implementations for all message getters defined in this interface.
// This ensures compile-time safety and guarantees that no message is missing.
type MessageBundle interface {
	// GetMessage returns a formatted message for the given key.
	GetMessage(key Key, args ...interface{}) string

	// Argon2 key derivation errors
	GetErrMemoryTooLow() string
	GetErrMemoryTooHigh() string
	GetErrThreadsMin() string
	GetErrThreadsMax() string
	GetErrThreadsExceed() string
	GetErrTimeMin() string
	GetErrTimeMax() string
	GetErrKeyLenShort() string
	GetErrKeyLenLong() string

	// CLI user interactions
	GetCliFileExists() string
	GetCliOperationCancelled() string
	GetCliError() string

	// Command-line flags descriptions
	GetFlagPassDesc() string
	GetFlagWorkersDesc() string
	GetFlagForceDesc() string
	GetFlagQuietDesc() string

	// Command descriptions
	GetCmdEncryptShort() string
	GetCmdEncryptLong() string
	GetCmdDecryptShort() string
	GetCmdDecryptLong() string

	// Interactive mode messages
	GetInteractiveTitle() string
	GetInteractiveEncryptFlow() string
	GetInteractiveDecryptFlow() string
	GetInteractiveInputFile() string
	GetInteractiveOutputFile() string
	GetInteractivePassword() string
	GetInteractiveConfirm() string
	GetInteractiveWorkerCount() string
	GetInteractiveOverwrite() string
	GetInteractiveCancel() string
	GetInteractivePressEnter() string
	GetInteractiveFileToEncrypt() string
	GetInteractiveEncryptedFile() string
	GetInteractivePasswordsNotMatch() string

	// Password handling messages
	GetPasswordPrompt() string
	GetPasswordConfirmPrompt() string
	GetPasswordEmpty() string
	GetPasswordReadError() string
	GetPasswordConfirmError() string
	GetPasswordNotMatch() string
	GetPasswordMinLength() string
	GetPasswordUppercase() string
	GetPasswordLowercase() string
	GetPasswordDigit() string

	// Root command messages
	GetRootShortDesc() string
	GetRootLongDesc() string
	GetRootUsage() string
	GetRootCommandsTitle() string
	GetRootPasswordManagement() string
	GetRootExamplesTitle() string
	GetRootExampleEncrypt() string
	GetRootExampleDecrypt() string
	GetRootExamplePassFlag() string
	GetRootExampleWorkers() string
	GetRootExampleForce() string

	// Version command messages
	GetVersionShortDesc() string
	GetVersionLongDesc() string
	GetVersionBuildInfo() string
	GetVersionOSArch() string
	GetVersionCPUs() string

	// Crypto errors
	GetErrDestSliceTooShort() string

	// File operation errors
	GetErrOpenFile() string

	// Service validation errors
	GetErrFileAlreadyExists() string
	GetErrFileNotFound() string
	GetWarnWorkersReduced() string

	// UI banner messages
	GetUIInteractiveHeader() string
	GetUIEncryptHeader() string
	GetUIDecryptHeader() string
	GetUIHeaderSeparator() string
	GetUIGoodbyeMessage() string

	// UI prompts messages
	GetUIPromptOperationLabel() string
	GetUIPromptEncryptOption() string
	GetUIPromptDecryptOption() string
	GetUIPromptExitOption() string
	GetUIPromptGoodbye() string
	GetUIPromptPathEmpty() string
	GetUIPromptPathNotExist() string
	GetUIPromptPathSuccess() string
	GetUIPromptPasswordMinLength() string
	GetUIPromptPasswordUppercase() string
	GetUIPromptPasswordLowercase() string
	GetUIPromptPasswordDigit() string
	GetUIPromptPasswordSuccess() string
	GetUIPromptWorkersLabel() string
	GetUIPromptWorkersSuccess() string
	GetUIPromptWorkersInvalid() string
	GetUIPromptWorkersMax() string
	GetUIPromptConfirmLabel() string
	GetUIPromptConfirmInvalid() string

	// UI success messages
	GetUISuccessOperation() string
	GetUISuccessOutput() string
	GetUISuccessSize() string

	// Cryptolib decryption errors
	GetCryptolibErrOpenInput() string
	GetCryptolibErrCreateOutput() string
	GetCryptolibErrCreateCipher() string
	GetCryptolibErrCreateGCM() string
	GetCryptolibErrReadHeader() string
	GetCryptolibErrReadHeaderHMAC() string
	GetCryptolibErrReadNonce() string
	GetCryptolibErrUnexpectedEOF() string
	GetCryptolibErrReadChunkLen() string
	GetCryptolibErrReadCiphertext() string
	GetCryptolibErrDeriveNonce() string
	GetCryptolibErrWritePlaintext() string

	// Cryptolib encryption errors
	GetCryptolibErrOpenInputEnc() string
	GetCryptolibErrCreateOutputEnc() string
	GetCryptolibErrGenerateSalt() string
	GetCryptolibErrWriteHeader() string
	GetCryptolibErrWriteHeaderHMAC() string
	GetCryptolibErrCreateCipherEnc() string
	GetCryptolibErrCreateGCMEnc() string
	GetCryptolibErrGenerateNonce() string
	GetCryptolibErrWriteNonce() string
	GetCryptolibErrReadChunk() string
	GetCryptolibErrWriteEndMarker() string
	GetCryptolibErrNonceDerivation() string
	GetCryptolibErrMissingChunks() string
	GetCryptolibErrTooManyPending() string
	GetCryptolibErrWriteChunkLen() string
	GetCryptolibErrWriteCiphertext() string
	GetCryptolibErrCloseInput() string
	GetCryptolibErrCloseOutput() string

	// Cryptolib sentinel errors
	GetCryptolibErrInvalidMagic() string
	GetCryptolibErrUnsupportedVersion() string
	GetCryptolibErrHeaderAuthFailed() string
	GetCryptolibErrDecryptionFailed() string

	// Cryptolib stream errors
	GetCryptolibErrReadHeaderStream() string
	GetCryptolibErrReadHeaderHMACStream() string
	GetCryptolibErrReadNonceStream() string
	GetCryptolibErrUnexpectedEOFStream() string
	GetCryptolibErrReadChunkLenStream() string
	GetCryptolibErrReadCiphertextStream() string
	GetCryptolibErrDeriveNonceStream() string
	GetCryptolibErrWritePlaintextStream() string
	GetCryptolibErrCreateCipherStream() string
	GetCryptolibErrCreateGCMStream() string
}

// String returns the string representation of the key.
func (k Key) String() string {
	return string(k)
}
