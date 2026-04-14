// Package lang provides English message bundle implementation.
package lang

import "fmt"

// EnglishBundle implements MessageBundle for English language.
type EnglishBundle struct{}

var _ MessageBundle = (*EnglishBundle)(nil)

// GetMessage returns a formatted message for the given key.
func (e EnglishBundle) GetMessage(key Key, args ...interface{}) string {
	var format string

	switch key {
	// Argon2 errors
	case ErrMemoryTooLow:
		format = e.GetErrMemoryTooLow()
	case ErrMemoryTooHigh:
		format = e.GetErrMemoryTooHigh()
	case ErrThreadsMin:
		format = e.GetErrThreadsMin()
	case ErrThreadsMax:
		format = e.GetErrThreadsMax()
	case ErrThreadsExceed:
		format = e.GetErrThreadsExceed()
	case ErrTimeMin:
		format = e.GetErrTimeMin()
	case ErrTimeMax:
		format = e.GetErrTimeMax()
	case ErrKeyLenShort:
		format = e.GetErrKeyLenShort()
	case ErrKeyLenLong:
		format = e.GetErrKeyLenLong()

	// CLI messages
	case CliFileExists:
		format = e.GetCliFileExists()
	case CliOperationCancelled:
		format = e.GetCliOperationCancelled()
	case CliError:
		format = e.GetCliError()

	// Flag descriptions
	case FlagPassDesc:
		format = e.GetFlagPassDesc()
	case FlagWorkersDesc:
		format = e.GetFlagWorkersDesc()
	case FlagForceDesc:
		format = e.GetFlagForceDesc()
	case FlagQuietDesc:
		format = e.GetFlagQuietDesc()

	// Command descriptions
	case CmdEncryptShort:
		format = e.GetCmdEncryptShort()
	case CmdEncryptLong:
		format = e.GetCmdEncryptLong()
	case CmdDecryptShort:
		format = e.GetCmdDecryptShort()
	case CmdDecryptLong:
		format = e.GetCmdDecryptLong()

	// Interactive mode messages
	case InteractiveTitle:
		format = e.GetInteractiveTitle()
	case InteractiveEncryptFlow:
		format = e.GetInteractiveEncryptFlow()
	case InteractiveDecryptFlow:
		format = e.GetInteractiveDecryptFlow()
	case InteractiveInputFile:
		format = e.GetInteractiveInputFile()
	case InteractiveOutputFile:
		format = e.GetInteractiveOutputFile()
	case InteractivePassword:
		format = e.GetInteractivePassword()
	case InteractiveConfirm:
		format = e.GetInteractiveConfirm()
	case InteractiveWorkerCount:
		format = e.GetInteractiveWorkerCount()
	case InteractiveOverwrite:
		format = e.GetInteractiveOverwrite()
	case InteractiveCancel:
		format = e.GetInteractiveCancel()
	case InteractivePressEnter:
		format = e.GetInteractivePressEnter()
	case InteractiveFileToEncrypt:
		format = e.GetInteractiveFileToEncrypt()
	case InteractiveEncryptedFile:
		format = e.GetInteractiveEncryptedFile()
	case InteractivePasswordsNotMatch:
		format = e.GetInteractivePasswordsNotMatch()

	// Password handling messages
	case PasswordPrompt:
		format = e.GetPasswordPrompt()
	case PasswordConfirmPrompt:
		format = e.GetPasswordConfirmPrompt()
	case PasswordEmpty:
		format = e.GetPasswordEmpty()
	case PasswordReadError:
		format = e.GetPasswordReadError()
	case PasswordConfirmError:
		format = e.GetPasswordConfirmError()
	case PasswordNotMatch:
		format = e.GetPasswordNotMatch()
	case PasswordMinLength:
		format = e.GetPasswordMinLength()
	case PasswordUppercase:
		format = e.GetPasswordUppercase()
	case PasswordLowercase:
		format = e.GetPasswordLowercase()
	case PasswordDigit:
		format = e.GetPasswordDigit()

	// Root command messages
	case RootShortDesc:
		format = e.GetRootShortDesc()
	case RootLongDesc:
		format = e.GetRootLongDesc()
	case RootUsage:
		format = e.GetRootUsage()
	case RootCommandsTitle:
		format = e.GetRootCommandsTitle()
	case RootPasswordManagement:
		format = e.GetRootPasswordManagement()
	case RootExamplesTitle:
		format = e.GetRootExamplesTitle()
	case RootExampleEncrypt:
		format = e.GetRootExampleEncrypt()
	case RootExampleDecrypt:
		format = e.GetRootExampleDecrypt()
	case RootExamplePassFlag:
		format = e.GetRootExamplePassFlag()
	case RootExampleWorkers:
		format = e.GetRootExampleWorkers()
	case RootExampleForce:
		format = e.GetRootExampleForce()

	// Version command messages
	case VersionShortDesc:
		format = e.GetVersionShortDesc()
	case VersionLongDesc:
		format = e.GetVersionLongDesc()
	case VersionBuildInfo:
		format = e.GetVersionBuildInfo()
	case VersionOSArch:
		format = e.GetVersionOSArch()
	case VersionCPUs:
		format = e.GetVersionCPUs()

	// Crypto errors
	case ErrDestSliceTooShort:
		format = e.GetErrDestSliceTooShort()

	// File operation errors
	case ErrOpenFile:
		format = e.GetErrOpenFile()

	// Service validation errors
	case ErrFileAlreadyExists:
		format = e.GetErrFileAlreadyExists()
	case ErrFileNotFound:
		format = e.GetErrFileNotFound()
	case WarnWorkersReduced:
		format = e.GetWarnWorkersReduced()

	// UI banner messages
	case UIInteractiveHeader:
		format = e.GetUIInteractiveHeader()
	case UIEncryptHeader:
		format = e.GetUIEncryptHeader()
	case UIDecryptHeader:
		format = e.GetUIDecryptHeader()
	case UIHeaderSeparator:
		format = e.GetUIHeaderSeparator()
	case UIGoodbyeMessage:
		format = e.GetUIGoodbyeMessage()

	// UI prompts messages
	case UIPromptOperationLabel:
		format = e.GetUIPromptOperationLabel()
	case UIPromptEncryptOption:
		format = e.GetUIPromptEncryptOption()
	case UIPromptDecryptOption:
		format = e.GetUIPromptDecryptOption()
	case UIPromptExitOption:
		format = e.GetUIPromptExitOption()
	case UIPromptGoodbye:
		format = e.GetUIPromptGoodbye()
	case UIPromptPathEmpty:
		format = e.GetUIPromptPathEmpty()
	case UIPromptPathNotExist:
		format = e.GetUIPromptPathNotExist()
	case UIPromptPathSuccess:
		format = e.GetUIPromptPathSuccess()
	case UIPromptPasswordMinLength:
		format = e.GetUIPromptPasswordMinLength()
	case UIPromptPasswordUppercase:
		format = e.GetUIPromptPasswordUppercase()
	case UIPromptPasswordLowercase:
		format = e.GetUIPromptPasswordLowercase()
	case UIPromptPasswordDigit:
		format = e.GetUIPromptPasswordDigit()
	case UIPromptPasswordSuccess:
		format = e.GetUIPromptPasswordSuccess()
	case UIPromptWorkersLabel:
		format = e.GetUIPromptWorkersLabel()
	case UIPromptWorkersSuccess:
		format = e.GetUIPromptWorkersSuccess()
	case UIPromptWorkersInvalid:
		format = e.GetUIPromptWorkersInvalid()
	case UIPromptWorkersMax:
		format = e.GetUIPromptWorkersMax()
	case UIPromptConfirmLabel:
		format = e.GetUIPromptConfirmLabel()
	case UIPromptConfirmInvalid:
		format = e.GetUIPromptConfirmInvalid()

	// UI success messages
	case UISuccessOperation:
		format = e.GetUISuccessOperation()
	case UISuccessOutput:
		format = e.GetUISuccessOutput()
	case UISuccessSize:
		format = e.GetUISuccessSize()

	// Cryptolib decryption errors
	case CryptolibErrOpenInput:
		format = e.GetCryptolibErrOpenInput()
	case CryptolibErrCreateOutput:
		format = e.GetCryptolibErrCreateOutput()
	case CryptolibErrCreateCipher:
		format = e.GetCryptolibErrCreateCipher()
	case CryptolibErrCreateGCM:
		format = e.GetCryptolibErrCreateGCM()
	case CryptolibErrReadHeader:
		format = e.GetCryptolibErrReadHeader()
	case CryptolibErrReadHeaderHMAC:
		format = e.GetCryptolibErrReadHeaderHMAC()
	case CryptolibErrReadNonce:
		format = e.GetCryptolibErrReadNonce()
	case CryptolibErrUnexpectedEOF:
		format = e.GetCryptolibErrUnexpectedEOF()
	case CryptolibErrReadChunkLen:
		format = e.GetCryptolibErrReadChunkLen()
	case CryptolibErrReadCiphertext:
		format = e.GetCryptolibErrReadCiphertext()
	case CryptolibErrDeriveNonce:
		format = e.GetCryptolibErrDeriveNonce()
	case CryptolibErrWritePlaintext:
		format = e.GetCryptolibErrWritePlaintext()

	// Cryptolib encryption errors
	case CryptolibErrOpenInputEnc:
		format = e.GetCryptolibErrOpenInputEnc()
	case CryptolibErrCreateOutputEnc:
		format = e.GetCryptolibErrCreateOutputEnc()
	case CryptolibErrGenerateSalt:
		format = e.GetCryptolibErrGenerateSalt()
	case CryptolibErrWriteHeader:
		format = e.GetCryptolibErrWriteHeader()
	case CryptolibErrWriteHeaderHMAC:
		format = e.GetCryptolibErrWriteHeaderHMAC()
	case CryptolibErrCreateCipherEnc:
		format = e.GetCryptolibErrCreateCipherEnc()
	case CryptolibErrCreateGCMEnc:
		format = e.GetCryptolibErrCreateGCMEnc()
	case CryptolibErrGenerateNonce:
		format = e.GetCryptolibErrGenerateNonce()
	case CryptolibErrWriteNonce:
		format = e.GetCryptolibErrWriteNonce()
	case CryptolibErrReadChunk:
		format = e.GetCryptolibErrReadChunk()
	case CryptolibErrWriteEndMarker:
		format = e.GetCryptolibErrWriteEndMarker()
	case CryptolibErrNonceDerivation:
		format = e.GetCryptolibErrNonceDerivation()
	case CryptolibErrMissingChunks:
		format = e.GetCryptolibErrMissingChunks()
	case CryptolibErrTooManyPending:
		format = e.GetCryptolibErrTooManyPending()
	case CryptolibErrWriteChunkLen:
		format = e.GetCryptolibErrWriteChunkLen()
	case CryptolibErrWriteCiphertext:
		format = e.GetCryptolibErrWriteCiphertext()
	case CryptolibErrCloseInput:
		format = e.GetCryptolibErrCloseInput()
	case CryptolibErrCloseOutput:
		format = e.GetCryptolibErrCloseOutput()

	// Cryptolib sentinel errors
	case CryptolibErrInvalidMagic:
		format = e.GetCryptolibErrInvalidMagic()
	case CryptolibErrUnsupportedVersion:
		format = e.GetCryptolibErrUnsupportedVersion()
	case CryptolibErrHeaderAuthFailed:
		format = e.GetCryptolibErrHeaderAuthFailed()
	case CryptolibErrDecryptionFailed:
		format = e.GetCryptolibErrDecryptionFailed()

	// Cryptolib stream errors
	case CryptolibErrReadHeaderStream:
		format = e.GetCryptolibErrReadHeaderStream()
	case CryptolibErrReadHeaderHMACStream:
		format = e.GetCryptolibErrReadHeaderHMACStream()
	case CryptolibErrReadNonceStream:
		format = e.GetCryptolibErrReadNonceStream()
	case CryptolibErrUnexpectedEOFStream:
		format = e.GetCryptolibErrUnexpectedEOFStream()
	case CryptolibErrReadChunkLenStream:
		format = e.GetCryptolibErrReadChunkLenStream()
	case CryptolibErrReadCiphertextStream:
		format = e.GetCryptolibErrReadCiphertextStream()
	case CryptolibErrDeriveNonceStream:
		format = e.GetCryptolibErrDeriveNonceStream()
	case CryptolibErrWritePlaintextStream:
		format = e.GetCryptolibErrWritePlaintextStream()
	case CryptolibErrCreateCipherStream:
		format = e.GetCryptolibErrCreateCipherStream()
	case CryptolibErrCreateGCMStream:
		format = e.GetCryptolibErrCreateGCMStream()

	default:
		return GetDefaultMessage(key)
	}

	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

// Argon2 errors

func (e EnglishBundle) GetErrMemoryTooLow() string {
	return "memory too low: %d KiB (minimum 8192 KiB)"
}

func (e EnglishBundle) GetErrMemoryTooHigh() string {
	return "memory too high: %d KiB (maximum 1,048,576 KiB)"
}

func (e EnglishBundle) GetErrThreadsMin() string {
	return "threads must be at least 1"
}

func (e EnglishBundle) GetErrThreadsMax() string {
	return "threads too high: %d (maximum %d)"
}

func (e EnglishBundle) GetErrThreadsExceed() string {
	return "threads exceed system capacity: %d (max %d)"
}

func (e EnglishBundle) GetErrTimeMin() string {
	return "time must be at least 1"
}

func (e EnglishBundle) GetErrTimeMax() string {
	return "time too high: %d (maximum 100)"
}

func (e EnglishBundle) GetErrKeyLenShort() string {
	return "key length too short: %d bytes (minimum 16)"
}

func (e EnglishBundle) GetErrKeyLenLong() string {
	return "key length too long: %d bytes (maximum 64)"
}

// CLI messages

func (e EnglishBundle) GetCliFileExists() string {
	return "File '%s' already exists. Overwrite?"
}

func (e EnglishBundle) GetCliOperationCancelled() string {
	return "❌ Operation cancelled"
}

func (e EnglishBundle) GetCliError() string {
	return "❌ Error: %v"
}

// Flag descriptions

func (e EnglishBundle) GetFlagPassDesc() string {
	return "Passphrase used for encryption (optional - will prompt if omitted)"
}

func (e EnglishBundle) GetFlagWorkersDesc() string {
	return "Number of parallel workers"
}

func (e EnglishBundle) GetFlagForceDesc() string {
	return "Overwrite existing output file without confirmation"
}

func (e EnglishBundle) GetFlagQuietDesc() string {
	return "Suppress progress bar output"
}

// Command descriptions

func (e EnglishBundle) GetCmdEncryptShort() string {
	return "🔒 Encrypt a file"
}

func (e EnglishBundle) GetCmdEncryptLong() string {
	return `Encrypt a file using AES-256-GCM with Argon2id key derivation.

The encryption process:
  1. Generates a random salt and nonce
  2. Derives a 256-bit key using Argon2id
  3. Splits the input into chunks (default 1MB)
  4. Encrypts chunks in parallel using the specified number of workers
  5. Writes header, HMAC, nonce, and encrypted chunks to the output file

Password can be provided via:
  - --pass flag (visible in process list, not recommended for shared environments)
  - Interactive prompt (recommended for manual use)

Examples:
  aescryptool encrypt secret.txt secret.enc              # Prompts for password
  aescryptool encrypt secret.txt secret.enc --pass myPassword
  aescryptool encrypt data.txt output.enc --pass secure123 --force
  aescryptool encrypt large.bin result.enc --workers 8 --quiet`
}

func (e EnglishBundle) GetCmdDecryptShort() string {
	return "🔓 Decrypt a file"
}

func (e EnglishBundle) GetCmdDecryptLong() string {
	return `Decrypt a file that was encrypted with the encrypt command.

The decryption process:
  1. Validates the input file exists
  2. Reads and verifies the file header
  3. Derives the encryption key using Argon2id with the salt from header
  4. Streams and decrypts the data to the output file
  5. Verifies integrity of each chunk via GCM authentication

Password can be provided via:
  - --pass flag (visible in process list, not recommended for shared environments)
  - Interactive prompt (recommended for manual use)

Examples:
  aescryptool decrypt secret.enc secret.txt              # Prompts for password
  aescryptool decrypt secret.enc secret.txt --pass myPassword
  aescryptool decrypt data.enc output.txt --pass secure123 --force
  aescryptool decrypt large.enc result.bin --workers 8 --quiet`
}

// Interactive mode messages

func (e EnglishBundle) GetInteractiveTitle() string {
	return "Interactive Mode"
}

func (e EnglishBundle) GetInteractiveEncryptFlow() string {
	return "Encryption"
}

func (e EnglishBundle) GetInteractiveDecryptFlow() string {
	return "Decryption"
}

func (e EnglishBundle) GetInteractiveInputFile() string {
	return "📁 File to encrypt"
}

func (e EnglishBundle) GetInteractiveOutputFile() string {
	return "📂 Output file"
}

func (e EnglishBundle) GetInteractivePassword() string {
	return "🔑 Password"
}

func (e EnglishBundle) GetInteractiveConfirm() string {
	return "✅ Confirmation"
}

func (e EnglishBundle) GetInteractiveWorkerCount() string {
	return "⚙️ Workers"
}

func (e EnglishBundle) GetInteractiveOverwrite() string {
	return "⚠️  File already exists. Overwrite?"
}

func (e EnglishBundle) GetInteractiveCancel() string {
	return "❌ Operation cancelled"
}

func (e EnglishBundle) GetInteractivePressEnter() string {
	return "🔁 Press Enter to continue..."
}

func (e EnglishBundle) GetInteractiveFileToEncrypt() string {
	return "📁 File to encrypt"
}

func (e EnglishBundle) GetInteractiveEncryptedFile() string {
	return "📁 Encrypted file"
}

func (e EnglishBundle) GetInteractivePasswordsNotMatch() string {
	return "❌ Passwords do not match"
}

// Password handling messages

func (e EnglishBundle) GetPasswordPrompt() string {
	return "🔑 Password: "
}

func (e EnglishBundle) GetPasswordConfirmPrompt() string {
	return "✅ Confirm password: "
}

func (e EnglishBundle) GetPasswordEmpty() string {
	return "password cannot be empty"
}

func (e EnglishBundle) GetPasswordReadError() string {
	return "read password: %w"
}

func (e EnglishBundle) GetPasswordConfirmError() string {
	return "read confirmation: %w"
}

func (e EnglishBundle) GetPasswordNotMatch() string {
	return "passwords do not match"
}

func (e EnglishBundle) GetPasswordMinLength() string {
	return "minimum 8 characters required"
}

func (e EnglishBundle) GetPasswordUppercase() string {
	return "at least one uppercase letter required"
}

func (e EnglishBundle) GetPasswordLowercase() string {
	return "at least one lowercase letter required"
}

func (e EnglishBundle) GetPasswordDigit() string {
	return "at least one digit required"
}

// Root command messages

func (e EnglishBundle) GetRootShortDesc() string {
	return "🔐 Secure file encryption using AES-256-GCM"
}

func (e EnglishBundle) GetRootLongDesc() string {
	return ""
}

func (e EnglishBundle) GetRootUsage() string {
	return "Usage:"
}

func (e EnglishBundle) GetRootCommandsTitle() string {
	return "Commands:"
}

func (e EnglishBundle) GetRootPasswordManagement() string {
	return "Password Management:"
}

func (e EnglishBundle) GetRootExamplesTitle() string {
	return "Examples:"
}

func (e EnglishBundle) GetRootExampleEncrypt() string {
	return "# Encrypt with interactive password prompt (recommended)"
}

func (e EnglishBundle) GetRootExampleDecrypt() string {
	return "# Decrypt with interactive password prompt (recommended)"
}

func (e EnglishBundle) GetRootExamplePassFlag() string {
	return "# Encrypt with --pass flag (for scripts)"
}

func (e EnglishBundle) GetRootExampleWorkers() string {
	return "# With parallel processing (8 workers)"
}

func (e EnglishBundle) GetRootExampleForce() string {
	return "# Force overwrite without confirmation"
}

// Version command messages

func (e EnglishBundle) GetVersionShortDesc() string {
	return "Show version information"
}

func (e EnglishBundle) GetVersionLongDesc() string {
	return "Display aescryptool version, build information, and system details"
}

func (e EnglishBundle) GetVersionBuildInfo() string {
	return "📦 Build: %s"
}

func (e EnglishBundle) GetVersionOSArch() string {
	return "🖥️  OS/Arch: %s/%s"
}

func (e EnglishBundle) GetVersionCPUs() string {
	return "💻 CPUs: %d"
}

// Crypto errors

func (e EnglishBundle) GetErrDestSliceTooShort() string {
	return "dest slice too short: need %d, got %d"
}

// File operation errors

func (e EnglishBundle) GetErrOpenFile() string {
	return "open file '%s': %w"
}

// Service validation errors

func (e EnglishBundle) GetErrFileAlreadyExists() string {
	return "file already exists"
}

func (e EnglishBundle) GetErrFileNotFound() string {
	return "file '%s' not found"
}

func (e EnglishBundle) GetWarnWorkersReduced() string {
	return "⚠️ Workers reduced to %d\n"
}

// UI banner messages

func (e EnglishBundle) GetUIInteractiveHeader() string {
	return `
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║                 🎮 AESCRYPTOOL - INTERACTIVE MODE                  ║
║                                                                    ║
║         Follow the prompts to encrypt or decrypt your files        ║
║         All inputs will be validated before execution              ║
║                                                                    ║
║           Ctrl+C = Return to menu | Ctrl+D = Quit                  ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝`
}

func (e EnglishBundle) GetUIEncryptHeader() string {
	return "🔐 FILE ENCRYPTION"
}

func (e EnglishBundle) GetUIDecryptHeader() string {
	return "🔓 FILE DECRYPTION"
}

func (e EnglishBundle) GetUIHeaderSeparator() string {
	return "────────────────────────────────────────"
}

func (e EnglishBundle) GetUIGoodbyeMessage() string {
	return `
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║              👋 Thank you for using AESCRYPTOOL!                   ║
║                                                                    ║
║              See you next time for your encryption needs!          ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝`
}

// UI prompts messages

func (e EnglishBundle) GetUIPromptOperationLabel() string {
	return "What do you want to do"
}

func (e EnglishBundle) GetUIPromptEncryptOption() string {
	return "🔒  Encrypt a file"
}

func (e EnglishBundle) GetUIPromptDecryptOption() string {
	return "🔓  Decrypt a file"
}

func (e EnglishBundle) GetUIPromptExitOption() string {
	return "🚪  Exit"
}

func (e EnglishBundle) GetUIPromptGoodbye() string {
	return "👋 Thank you for using CRYPTOOL !"
}

func (e EnglishBundle) GetUIPromptPathEmpty() string {
	return "❌ Path cannot be empty"
}

func (e EnglishBundle) GetUIPromptPathNotExist() string {
	return "❌ File '%s' does not exist"
}

func (e EnglishBundle) GetUIPromptPathSuccess() string {
	return "   ✓ %s"
}

func (e EnglishBundle) GetUIPromptPasswordMinLength() string {
	return "❌ Minimum 8 characters required"
}

func (e EnglishBundle) GetUIPromptPasswordUppercase() string {
	return "❌ At least one uppercase letter required"
}

func (e EnglishBundle) GetUIPromptPasswordLowercase() string {
	return "❌ At least one lowercase letter required"
}

func (e EnglishBundle) GetUIPromptPasswordDigit() string {
	return "❌ At least one digit required"
}

func (e EnglishBundle) GetUIPromptPasswordSuccess() string {
	return "   ✓ %s"
}

func (e EnglishBundle) GetUIPromptWorkersLabel() string {
	return "⚙️  Workers (default: %d, max: %d)"
}

func (e EnglishBundle) GetUIPromptWorkersSuccess() string {
	return "   ✓ %d workers"
}

func (e EnglishBundle) GetUIPromptWorkersInvalid() string {
	return "❌ Valid number required (>=1)"
}

func (e EnglishBundle) GetUIPromptWorkersMax() string {
	return "❌ Maximum %d workers"
}

func (e EnglishBundle) GetUIPromptConfirmLabel() string {
	return "❓ %s [%s]: "
}

func (e EnglishBundle) GetUIPromptConfirmInvalid() string {
	return "❌ Please answer y/n"
}

// UI success messages

func (e EnglishBundle) GetUISuccessOperation() string {
	return "✅ Operation successful!"
}

func (e EnglishBundle) GetUISuccessOutput() string {
	return "📄 Output: %s"
}

func (e EnglishBundle) GetUISuccessSize() string {
	return "📏 Size:   %s"
}

// Cryptolib decryption errors

func (e EnglishBundle) GetCryptolibErrOpenInput() string {
	return "open input: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateOutput() string {
	return "create output: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateCipher() string {
	return "create cipher: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateGCM() string {
	return "create GCM: %w"
}

func (e EnglishBundle) GetCryptolibErrReadHeader() string {
	return "read header: %w"
}

func (e EnglishBundle) GetCryptolibErrReadHeaderHMAC() string {
	return "read header HMAC: %w"
}

func (e EnglishBundle) GetCryptolibErrReadNonce() string {
	return "read nonce: %w"
}

func (e EnglishBundle) GetCryptolibErrUnexpectedEOF() string {
	return "unexpected EOF: missing end marker"
}

func (e EnglishBundle) GetCryptolibErrReadChunkLen() string {
	return "read chunk length: %w"
}

func (e EnglishBundle) GetCryptolibErrReadCiphertext() string {
	return "read ciphertext chunk %d: %w"
}

func (e EnglishBundle) GetCryptolibErrDeriveNonce() string {
	return "derive nonce for chunk %d: %w"
}

func (e EnglishBundle) GetCryptolibErrWritePlaintext() string {
	return "write plaintext chunk %d: %w"
}

// Cryptolib encryption errors

func (e EnglishBundle) GetCryptolibErrOpenInputEnc() string {
	return "open input: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateOutputEnc() string {
	return "create output: %w"
}

func (e EnglishBundle) GetCryptolibErrGenerateSalt() string {
	return "generate salt: %w"
}

func (e EnglishBundle) GetCryptolibErrWriteHeader() string {
	return "write header: %w"
}

func (e EnglishBundle) GetCryptolibErrWriteHeaderHMAC() string {
	return "write header HMAC: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateCipherEnc() string {
	return "create cipher: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateGCMEnc() string {
	return "create GCM: %w"
}

func (e EnglishBundle) GetCryptolibErrGenerateNonce() string {
	return "generate nonce: %w"
}

func (e EnglishBundle) GetCryptolibErrWriteNonce() string {
	return "write nonce: %w"
}

func (e EnglishBundle) GetCryptolibErrReadChunk() string {
	return "read chunk: %w"
}

func (e EnglishBundle) GetCryptolibErrWriteEndMarker() string {
	return "write end marker: %w"
}

func (e EnglishBundle) GetCryptolibErrNonceDerivation() string {
	return "nonce derivation failed: %w"
}

func (e EnglishBundle) GetCryptolibErrMissingChunks() string {
	return "missing chunks: expected index %d, have %d pending"
}

func (e EnglishBundle) GetCryptolibErrTooManyPending() string {
	return "too many pending chunks (limit %d) - possible reordering attack"
}

func (e EnglishBundle) GetCryptolibErrWriteChunkLen() string {
	return "write chunk length: %w"
}

func (e EnglishBundle) GetCryptolibErrWriteCiphertext() string {
	return "write ciphertext: %w"
}

func (e EnglishBundle) GetCryptolibErrCloseInput() string {
	return "close input: %w"
}

func (e EnglishBundle) GetCryptolibErrCloseOutput() string {
	return "close output: %w"
}

// Cryptolib sentinel errors

func (e EnglishBundle) GetCryptolibErrInvalidMagic() string {
	return "invalid magic bytes: file not encrypted with this tool"
}

func (e EnglishBundle) GetCryptolibErrUnsupportedVersion() string {
	return "unsupported file version"
}

func (e EnglishBundle) GetCryptolibErrHeaderAuthFailed() string {
	return "header authentication failed: wrong passphrase or corrupted file"
}

func (e EnglishBundle) GetCryptolibErrDecryptionFailed() string {
	return "decryption failed: corrupted data or wrong key"
}

// Cryptolib stream errors

func (e EnglishBundle) GetCryptolibErrReadHeaderStream() string {
	return "read header: %w"
}

func (e EnglishBundle) GetCryptolibErrReadHeaderHMACStream() string {
	return "read header HMAC: %w"
}

func (e EnglishBundle) GetCryptolibErrReadNonceStream() string {
	return "read nonce: %w"
}

func (e EnglishBundle) GetCryptolibErrUnexpectedEOFStream() string {
	return "unexpected EOF: missing end marker"
}

func (e EnglishBundle) GetCryptolibErrReadChunkLenStream() string {
	return "read chunk length: %w"
}

func (e EnglishBundle) GetCryptolibErrReadCiphertextStream() string {
	return "read ciphertext chunk %d: %w"
}

func (e EnglishBundle) GetCryptolibErrDeriveNonceStream() string {
	return "derive nonce for chunk %d: %w"
}

func (e EnglishBundle) GetCryptolibErrWritePlaintextStream() string {
	return "write plaintext chunk %d: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateCipherStream() string {
	return "create cipher: %w"
}

func (e EnglishBundle) GetCryptolibErrCreateGCMStream() string {
	return "create GCM: %w"
}
