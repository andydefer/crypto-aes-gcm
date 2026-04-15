// Package lang provides French message bundle implementation.
package lang

import "fmt"

// FrenchBundle implements MessageBundle for French language.
type FrenchBundle struct{}

var _ MessageBundle = (*FrenchBundle)(nil)

// GetMessage returns a formatted message for the given key.
func (f FrenchBundle) GetMessage(key Key, args ...interface{}) string {
	var format string

	switch key {
	// Argon2 errors
	case ErrMemoryTooLow:
		format = f.GetErrMemoryTooLow()
	case ErrMemoryTooHigh:
		format = f.GetErrMemoryTooHigh()
	case ErrThreadsMin:
		format = f.GetErrThreadsMin()
	case ErrThreadsMax:
		format = f.GetErrThreadsMax()
	case ErrThreadsExceed:
		format = f.GetErrThreadsExceed()
	case ErrTimeMin:
		format = f.GetErrTimeMin()
	case ErrTimeMax:
		format = f.GetErrTimeMax()
	case ErrKeyLenShort:
		format = f.GetErrKeyLenShort()
	case ErrKeyLenLong:
		format = f.GetErrKeyLenLong()

	// CLI messages
	case CliFileExists:
		format = f.GetCliFileExists()
	case CliOperationCancelled:
		format = f.GetCliOperationCancelled()
	case CliError:
		format = f.GetCliError()

	// Flag descriptions
	case FlagPassDesc:
		format = f.GetFlagPassDesc()
	case FlagWorkersDesc:
		format = f.GetFlagWorkersDesc()
	case FlagForceDesc:
		format = f.GetFlagForceDesc()
	case FlagQuietDesc:
		format = f.GetFlagQuietDesc()

	// Command descriptions
	case CmdEncryptShort:
		format = f.GetCmdEncryptShort()
	case CmdEncryptLong:
		format = f.GetCmdEncryptLong()
	case CmdDecryptShort:
		format = f.GetCmdDecryptShort()
	case CmdDecryptLong:
		format = f.GetCmdDecryptLong()

	// Interactive mode messages
	case InteractiveTitle:
		format = f.GetInteractiveTitle()
	case InteractiveEncryptFlow:
		format = f.GetInteractiveEncryptFlow()
	case InteractiveDecryptFlow:
		format = f.GetInteractiveDecryptFlow()
	case InteractiveInputFile:
		format = f.GetInteractiveInputFile()
	case InteractiveOutputFile:
		format = f.GetInteractiveOutputFile()
	case InteractivePassword:
		format = f.GetInteractivePassword()
	case InteractiveConfirm:
		format = f.GetInteractiveConfirm()
	case InteractiveWorkerCount:
		format = f.GetInteractiveWorkerCount()
	case InteractiveOverwrite:
		format = f.GetInteractiveOverwrite()
	case InteractiveOverwriteConfirm:
		format = f.GetInteractiveOverwriteConfirm()
	case InteractiveOverwriteCancelled:
		format = f.GetInteractiveOverwriteCancelled()
	case InteractiveCancel:
		format = f.GetInteractiveCancel()
	case InteractiveCancelOperation:
		format = f.GetInteractiveCancelOperation()
	case InteractivePressEnter:
		format = f.GetInteractivePressEnter()
	case InteractiveFileToEncrypt:
		format = f.GetInteractiveFileToEncrypt()
	case InteractiveEncryptedFile:
		format = f.GetInteractiveEncryptedFile()
	case InteractivePasswordsNotMatch:
		format = f.GetInteractivePasswordsNotMatch()
	case InteractiveCheckExists:
		format = f.GetInteractiveCheckExists()

	// Password handling messages
	case PasswordPrompt:
		format = f.GetPasswordPrompt()
	case PasswordConfirmPrompt:
		format = f.GetPasswordConfirmPrompt()
	case PasswordEmpty:
		format = f.GetPasswordEmpty()
	case PasswordReadError:
		format = f.GetPasswordReadError()
	case PasswordConfirmError:
		format = f.GetPasswordConfirmError()
	case PasswordNotMatch:
		format = f.GetPasswordNotMatch()
	case PasswordMinLength:
		format = f.GetPasswordMinLength()
	case PasswordUppercase:
		format = f.GetPasswordUppercase()
	case PasswordLowercase:
		format = f.GetPasswordLowercase()
	case PasswordDigit:
		format = f.GetPasswordDigit()

	// Root command messages
	case RootShortDesc:
		format = f.GetRootShortDesc()
	case RootLongDesc:
		format = f.GetRootLongDesc()
	case RootUsage:
		format = f.GetRootUsage()
	case RootCommandsTitle:
		format = f.GetRootCommandsTitle()
	case RootPasswordManagement:
		format = f.GetRootPasswordManagement()
	case RootExamplesTitle:
		format = f.GetRootExamplesTitle()
	case RootExampleEncrypt:
		format = f.GetRootExampleEncrypt()
	case RootExampleDecrypt:
		format = f.GetRootExampleDecrypt()
	case RootExamplePassFlag:
		format = f.GetRootExamplePassFlag()
	case RootExampleWorkers:
		format = f.GetRootExampleWorkers()
	case RootExampleForce:
		format = f.GetRootExampleForce()

	// Version command messages
	case VersionShortDesc:
		format = f.GetVersionShortDesc()
	case VersionLongDesc:
		format = f.GetVersionLongDesc()
	case VersionBuildInfo:
		format = f.GetVersionBuildInfo()
	case VersionOSArch:
		format = f.GetVersionOSArch()
	case VersionCPUs:
		format = f.GetVersionCPUs()

	// Crypto errors
	case ErrDestSliceTooShort:
		format = f.GetErrDestSliceTooShort()

	// File operation errors
	case ErrOpenFile:
		format = f.GetErrOpenFile()

	// Service validation errors
	case ErrFileAlreadyExists:
		format = f.GetErrFileAlreadyExists()
	case ErrFileNotFound:
		format = f.GetErrFileNotFound()
	case WarnWorkersReduced:
		format = f.GetWarnWorkersReduced()

	// UI banner messages
	case UIInteractiveHeader:
		format = f.GetUIInteractiveHeader()
	case UIEncryptHeader:
		format = f.GetUIEncryptHeader()
	case UIDecryptHeader:
		format = f.GetUIDecryptHeader()
	case UIHeaderSeparator:
		format = f.GetUIHeaderSeparator()
	case UIGoodbyeMessage:
		format = f.GetUIGoodbyeMessage()

	// UI prompts messages
	case UIPromptOperationLabel:
		format = f.GetUIPromptOperationLabel()
	case UIPromptEncryptOption:
		format = f.GetUIPromptEncryptOption()
	case UIPromptDecryptOption:
		format = f.GetUIPromptDecryptOption()
	case UIPromptExitOption:
		format = f.GetUIPromptExitOption()
	case UIPromptGoodbye:
		format = f.GetUIPromptGoodbye()
	case UIPromptPathEmpty:
		format = f.GetUIPromptPathEmpty()
	case UIPromptPathNotExist:
		format = f.GetUIPromptPathNotExist()
	case UIPromptPathSuccess:
		format = f.GetUIPromptPathSuccess()
	case UIPromptPasswordMinLength:
		format = f.GetUIPromptPasswordMinLength()
	case UIPromptPasswordUppercase:
		format = f.GetUIPromptPasswordUppercase()
	case UIPromptPasswordLowercase:
		format = f.GetUIPromptPasswordLowercase()
	case UIPromptPasswordDigit:
		format = f.GetUIPromptPasswordDigit()
	case UIPromptPasswordSuccess:
		format = f.GetUIPromptPasswordSuccess()
	case UIPromptWorkersLabel:
		format = f.GetUIPromptWorkersLabel()
	case UIPromptWorkersSuccess:
		format = f.GetUIPromptWorkersSuccess()
	case UIPromptWorkersInvalid:
		format = f.GetUIPromptWorkersInvalid()
	case UIPromptWorkersMax:
		format = f.GetUIPromptWorkersMax()
	case UIPromptConfirmLabel:
		format = f.GetUIPromptConfirmLabel()
	case UIPromptConfirmInvalid:
		format = f.GetUIPromptConfirmInvalid()

	// UI success messages
	case UISuccessOperation:
		format = f.GetUISuccessOperation()
	case UISuccessOutput:
		format = f.GetUISuccessOutput()
	case UISuccessSize:
		format = f.GetUISuccessSize()

	// Cryptolib decryption errors
	case CryptolibErrOpenInput:
		format = f.GetCryptolibErrOpenInput()
	case CryptolibErrCreateOutput:
		format = f.GetCryptolibErrCreateOutput()
	case CryptolibErrCreateCipher:
		format = f.GetCryptolibErrCreateCipher()
	case CryptolibErrCreateGCM:
		format = f.GetCryptolibErrCreateGCM()
	case CryptolibErrReadHeader:
		format = f.GetCryptolibErrReadHeader()
	case CryptolibErrReadHeaderHMAC:
		format = f.GetCryptolibErrReadHeaderHMAC()
	case CryptolibErrReadNonce:
		format = f.GetCryptolibErrReadNonce()
	case CryptolibErrUnexpectedEOF:
		format = f.GetCryptolibErrUnexpectedEOF()
	case CryptolibErrReadChunkLen:
		format = f.GetCryptolibErrReadChunkLen()
	case CryptolibErrReadCiphertext:
		format = f.GetCryptolibErrReadCiphertext()
	case CryptolibErrDeriveNonce:
		format = f.GetCryptolibErrDeriveNonce()
	case CryptolibErrWritePlaintext:
		format = f.GetCryptolibErrWritePlaintext()

	// Cryptolib encryption errors
	case CryptolibErrOpenInputEnc:
		format = f.GetCryptolibErrOpenInputEnc()
	case CryptolibErrCreateOutputEnc:
		format = f.GetCryptolibErrCreateOutputEnc()
	case CryptolibErrGenerateSalt:
		format = f.GetCryptolibErrGenerateSalt()
	case CryptolibErrWriteHeader:
		format = f.GetCryptolibErrWriteHeader()
	case CryptolibErrWriteHeaderHMAC:
		format = f.GetCryptolibErrWriteHeaderHMAC()
	case CryptolibErrCreateCipherEnc:
		format = f.GetCryptolibErrCreateCipherEnc()
	case CryptolibErrCreateGCMEnc:
		format = f.GetCryptolibErrCreateGCMEnc()
	case CryptolibErrGenerateNonce:
		format = f.GetCryptolibErrGenerateNonce()
	case CryptolibErrWriteNonce:
		format = f.GetCryptolibErrWriteNonce()
	case CryptolibErrReadChunk:
		format = f.GetCryptolibErrReadChunk()
	case CryptolibErrWriteEndMarker:
		format = f.GetCryptolibErrWriteEndMarker()
	case CryptolibErrNonceDerivation:
		format = f.GetCryptolibErrNonceDerivation()
	case CryptolibErrMissingChunks:
		format = f.GetCryptolibErrMissingChunks()
	case CryptolibErrTooManyPending:
		format = f.GetCryptolibErrTooManyPending()
	case CryptolibErrWriteChunkLen:
		format = f.GetCryptolibErrWriteChunkLen()
	case CryptolibErrWriteCiphertext:
		format = f.GetCryptolibErrWriteCiphertext()
	case CryptolibErrCloseInput:
		format = f.GetCryptolibErrCloseInput()
	case CryptolibErrCloseOutput:
		format = f.GetCryptolibErrCloseOutput()

	// Cryptolib sentinel errors
	case CryptolibErrInvalidMagic:
		format = f.GetCryptolibErrInvalidMagic()
	case CryptolibErrUnsupportedVersion:
		format = f.GetCryptolibErrUnsupportedVersion()
	case CryptolibErrHeaderAuthFailed:
		format = f.GetCryptolibErrHeaderAuthFailed()
	case CryptolibErrDecryptionFailed:
		format = f.GetCryptolibErrDecryptionFailed()
	case CryptolibErrChunkTooLarge:
		format = f.GetCryptolibErrChunkTooLarge()

	// Cryptolib stream errors
	case CryptolibErrReadHeaderStream:
		format = f.GetCryptolibErrReadHeaderStream()
	case CryptolibErrReadHeaderHMACStream:
		format = f.GetCryptolibErrReadHeaderHMACStream()
	case CryptolibErrReadNonceStream:
		format = f.GetCryptolibErrReadNonceStream()
	case CryptolibErrUnexpectedEOFStream:
		format = f.GetCryptolibErrUnexpectedEOFStream()
	case CryptolibErrReadChunkLenStream:
		format = f.GetCryptolibErrReadChunkLenStream()
	case CryptolibErrReadCiphertextStream:
		format = f.GetCryptolibErrReadCiphertextStream()
	case CryptolibErrDeriveNonceStream:
		format = f.GetCryptolibErrDeriveNonceStream()
	case CryptolibErrWritePlaintextStream:
		format = f.GetCryptolibErrWritePlaintextStream()
	case CryptolibErrCreateCipherStream:
		format = f.GetCryptolibErrCreateCipherStream()
	case CryptolibErrCreateGCMStream:
		format = f.GetCryptolibErrCreateGCMStream()

	default:
		return GetDefaultMessage(key)
	}

	if len(args) > 0 {
		return fmt.Sprintf(format, args...)
	}
	return format
}

// ============================================================================
// Argon2 errors
// ============================================================================

func (f FrenchBundle) GetErrMemoryTooLow() string {
	return "mémoire trop basse: %d KiB (minimum 8192 KiB)"
}

func (f FrenchBundle) GetErrMemoryTooHigh() string {
	return "mémoire trop haute: %d KiB (maximum 1,048,576 KiB)"
}

func (f FrenchBundle) GetErrThreadsMin() string {
	return "le nombre de threads doit être au moins 1"
}

func (f FrenchBundle) GetErrThreadsMax() string {
	return "trop de threads: %d (maximum %d)"
}

func (f FrenchBundle) GetErrThreadsExceed() string {
	return "les threads dépassent la capacité système: %d (max %d)"
}

func (f FrenchBundle) GetErrTimeMin() string {
	return "le temps doit être au moins 1"
}

func (f FrenchBundle) GetErrTimeMax() string {
	return "temps trop élevé: %d (maximum 100)"
}

func (f FrenchBundle) GetErrKeyLenShort() string {
	return "clé trop courte: %d bytes (minimum 16)"
}

func (f FrenchBundle) GetErrKeyLenLong() string {
	return "clé trop longue: %d bytes (maximum 64)"
}

// ============================================================================
// CLI messages
// ============================================================================

func (f FrenchBundle) GetCliFileExists() string {
	return "Le fichier '%s' existe déjà. Écraser ?"
}

func (f FrenchBundle) GetCliOperationCancelled() string {
	return "❌ Opération annulée"
}

func (f FrenchBundle) GetCliError() string {
	return "❌ Erreur: %v"
}

// ============================================================================
// Flag descriptions
// ============================================================================

func (f FrenchBundle) GetFlagPassDesc() string {
	return "Phrase de passe utilisée pour le chiffrement (optionnel - sera demandée si omis)"
}

func (f FrenchBundle) GetFlagWorkersDesc() string {
	return "Nombre de workers parallèles"
}

func (f FrenchBundle) GetFlagForceDesc() string {
	return "Écraser le fichier de sortie existant sans confirmation"
}

func (f FrenchBundle) GetFlagQuietDesc() string {
	return "Supprimer l'affichage de la barre de progression"
}

// ============================================================================
// Command descriptions
// ============================================================================

func (f FrenchBundle) GetCmdEncryptShort() string {
	return "🔒 Chiffrer un fichier"
}

func (f FrenchBundle) GetCmdEncryptLong() string {
	return `Chiffrer un fichier en utilisant AES-256-GCM avec dérivation de clé Argon2id.

Le processus de chiffrement:
  1. Génère un sel et un nonce aléatoires
  2. Dérive une clé 256-bit en utilisant Argon2id
  3. Divise l'entrée en chunks (1MB par défaut)
  4. Chiffre les chunks en parallèle avec le nombre de workers spécifié
  5. Écrit l'en-tête, HMAC, nonce, et les chunks chiffrés dans le fichier de sortie

Le mot de passe peut être fourni via:
  - Flag --pass (visible dans la liste des processus, déconseillé pour les environnements partagés)
  - Invite interactive (recommandé pour une utilisation manuelle)

Exemples:
  aescryptool encrypt secret.txt secret.enc              # Demande le mot de passe
  aescryptool encrypt secret.txt secret.enc --pass myPassword
  aescryptool encrypt data.txt output.enc --pass secure123 --force
  aescryptool encrypt large.bin result.enc --workers 8 --quiet`
}

func (f FrenchBundle) GetCmdDecryptShort() string {
	return "🔓 Déchiffrer un fichier"
}

func (f FrenchBundle) GetCmdDecryptLong() string {
	return `Déchiffrer un fichier qui a été chiffré avec la commande encrypt.

Le processus de déchiffrement:
  1. Vérifie que le fichier d'entrée existe
  2. Lit et vérifie l'en-tête du fichier
  3. Dérive la clé de chiffrement avec Argon2id en utilisant le sel de l'en-tête
  4. Déchiffre et écrit les données en streaming
  5. Vérifie l'intégrité de chaque chunk via l'authentification GCM

Le mot de passe peut être fourni via:
  - Flag --pass (visible dans la liste des processus, déconseillé pour les environnements partagés)
  - Invite interactive (recommandé pour une utilisation manuelle)

Exemples:
  aescryptool decrypt secret.enc secret.txt              # Demande le mot de passe
  aescryptool decrypt secret.enc secret.txt --pass myPassword
  aescryptool decrypt data.enc output.txt --pass secure123 --force
  aescryptool decrypt large.enc result.bin --workers 8 --quiet`
}

// ============================================================================
// Interactive mode messages
// ============================================================================

func (f FrenchBundle) GetInteractiveTitle() string {
	return "Mode Interactif"
}

func (f FrenchBundle) GetInteractiveEncryptFlow() string {
	return "Chiffrement"
}

func (f FrenchBundle) GetInteractiveDecryptFlow() string {
	return "Déchiffrement"
}

func (f FrenchBundle) GetInteractiveInputFile() string {
	return "📁 Fichier à chiffrer"
}

func (f FrenchBundle) GetInteractiveOutputFile() string {
	return "📂 Fichier de sortie"
}

func (f FrenchBundle) GetInteractivePassword() string {
	return "🔑 Mot de passe"
}

func (f FrenchBundle) GetInteractiveConfirm() string {
	return "✅ Confirmation"
}

func (f FrenchBundle) GetInteractiveWorkerCount() string {
	return "⚙️ Workers"
}

func (f FrenchBundle) GetInteractiveOverwrite() string {
	return "⚠️  Le fichier existe déjà. Écraser ?"
}

func (f FrenchBundle) GetInteractiveOverwriteConfirm() string {
	return "⚠️ Le fichier existe déjà. Écraser ?"
}

func (f FrenchBundle) GetInteractiveOverwriteCancelled() string {
	return "❌ Opération annulée"
}

func (f FrenchBundle) GetInteractiveCancel() string {
	return "❌ Opération annulée"
}

func (f FrenchBundle) GetInteractiveCancelOperation() string {
	return "opération annulée par l'utilisateur"
}

func (f FrenchBundle) GetInteractivePressEnter() string {
	return "🔁 Appuyez sur Entrée pour continuer..."
}

func (f FrenchBundle) GetInteractiveFileToEncrypt() string {
	return "📁 Fichier à chiffrer"
}

func (f FrenchBundle) GetInteractiveEncryptedFile() string {
	return "📁 Fichier chiffré"
}

func (f FrenchBundle) GetInteractivePasswordsNotMatch() string {
	return "❌ Les mots de passe ne correspondent pas"
}

func (f FrenchBundle) GetInteractiveCheckExists() string {
	return "vérification de l'existence du fichier"
}

// ============================================================================
// Password handling messages
// ============================================================================

func (f FrenchBundle) GetPasswordPrompt() string {
	return "🔑 Mot de passe: "
}

func (f FrenchBundle) GetPasswordConfirmPrompt() string {
	return "✅ Confirmation du mot de passe: "
}

func (f FrenchBundle) GetPasswordEmpty() string {
	return "le mot de passe ne peut pas être vide"
}

func (f FrenchBundle) GetPasswordReadError() string {
	return "lecture du mot de passe: %w"
}

func (f FrenchBundle) GetPasswordConfirmError() string {
	return "lecture de la confirmation: %w"
}

func (f FrenchBundle) GetPasswordNotMatch() string {
	return "les mots de passe ne correspondent pas"
}

func (f FrenchBundle) GetPasswordMinLength() string {
	return "8 caractères minimum requis"
}

func (f FrenchBundle) GetPasswordUppercase() string {
	return "au moins une lettre majuscule requise"
}

func (f FrenchBundle) GetPasswordLowercase() string {
	return "au moins une lettre minuscule requise"
}

func (f FrenchBundle) GetPasswordDigit() string {
	return "au moins un chiffre requis"
}

// ============================================================================
// Root command messages
// ============================================================================

func (f FrenchBundle) GetRootShortDesc() string {
	return "🔐 Chiffrement sécurisé de fichiers avec AES-256-GCM"
}

func (f FrenchBundle) GetRootLongDesc() string {
	return ""
}

func (f FrenchBundle) GetRootUsage() string {
	return "Utilisation:"
}

func (f FrenchBundle) GetRootCommandsTitle() string {
	return "Commandes:"
}

func (f FrenchBundle) GetRootPasswordManagement() string {
	return "Gestion des mots de passe:"
}

func (f FrenchBundle) GetRootExamplesTitle() string {
	return "Exemples:"
}

func (f FrenchBundle) GetRootExampleEncrypt() string {
	return "# Chiffrement avec invite interactive (recommandé)"
}

func (f FrenchBundle) GetRootExampleDecrypt() string {
	return "# Déchiffrement avec invite interactive (recommandé)"
}

func (f FrenchBundle) GetRootExamplePassFlag() string {
	return "# Chiffrement avec flag --pass (pour les scripts)"
}

func (f FrenchBundle) GetRootExampleWorkers() string {
	return "# Avec traitement parallèle (8 workers)"
}

func (f FrenchBundle) GetRootExampleForce() string {
	return "# Forcer l'écrasement sans confirmation"
}

// ============================================================================
// Version command messages
// ============================================================================

func (f FrenchBundle) GetVersionShortDesc() string {
	return "Afficher les informations de version"
}

func (f FrenchBundle) GetVersionLongDesc() string {
	return "Afficher la version, les informations de build et les détails système d'aescryptool"
}

func (f FrenchBundle) GetVersionBuildInfo() string {
	return "📦 Build: %s"
}

func (f FrenchBundle) GetVersionOSArch() string {
	return "🖥️  OS/Arch: %s/%s"
}

func (f FrenchBundle) GetVersionCPUs() string {
	return "💻 CPUs: %d"
}

// ============================================================================
// Crypto errors
// ============================================================================

func (f FrenchBundle) GetErrDestSliceTooShort() string {
	return "tranche de destination trop courte: besoin de %d, reçu %d"
}

// ============================================================================
// File operation errors
// ============================================================================

func (f FrenchBundle) GetErrOpenFile() string {
	return "ouverture du fichier '%s': %w"
}

// ============================================================================
// Service validation errors
// ============================================================================

func (f FrenchBundle) GetErrFileAlreadyExists() string {
	return "fichier existe déjà"
}

func (f FrenchBundle) GetErrFileNotFound() string {
	return "fichier '%s' inexistant"
}

func (f FrenchBundle) GetWarnWorkersReduced() string {
	return "⚠️ Workers réduit à %d\n"
}

// ============================================================================
// UI banner messages
// ============================================================================

func (f FrenchBundle) GetUIInteractiveHeader() string {
	return `
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║                 🎮 AESCRYPTOOL - MODE INTERACTIF                   ║
║                                                                    ║
║   Suivez les invites pour chiffrer ou déchiffrer vos fichiers      ║
║       Toutes les entrées seront validées avant exécution           ║
║                                                                    ║
║          Ctrl+C = Retour au menu | Ctrl+D = Quitter                ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝`
}

func (f FrenchBundle) GetUIEncryptHeader() string {
	return "🔐 CHIFFREMENT DE FICHIER"
}

func (f FrenchBundle) GetUIDecryptHeader() string {
	return "🔓 DÉCHIFFREMENT DE FICHIER"
}

func (f FrenchBundle) GetUIHeaderSeparator() string {
	return "────────────────────────────────────────"
}

func (f FrenchBundle) GetUIGoodbyeMessage() string {
	return `
╔════════════════════════════════════════════════════════════════════╗
║                                                                    ║
║              👋 Merci d'avoir utilisé AESCRYPTOOL !                ║
║                                                                    ║
║              À bientôt pour vos prochains chiffrements !           ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝`
}

// ============================================================================
// UI prompts messages
// ============================================================================

func (f FrenchBundle) GetUIPromptOperationLabel() string {
	return "Que souhaitez-vous faire"
}

func (f FrenchBundle) GetUIPromptEncryptOption() string {
	return "🔒  Chiffrer un fichier"
}

func (f FrenchBundle) GetUIPromptDecryptOption() string {
	return "🔓  Déchiffrer un fichier"
}

func (f FrenchBundle) GetUIPromptExitOption() string {
	return "🚪  Quitter"
}

func (f FrenchBundle) GetUIPromptGoodbye() string {
	return "👋 Merci d'avoir utilisé CRYPTOOL !"
}

func (f FrenchBundle) GetUIPromptPathEmpty() string {
	return "❌ Le chemin ne peut pas être vide"
}

func (f FrenchBundle) GetUIPromptPathNotExist() string {
	return "❌ Le fichier '%s' n'existe pas"
}

func (f FrenchBundle) GetUIPromptPathSuccess() string {
	return "   ✓ %s"
}

func (f FrenchBundle) GetUIPromptPasswordMinLength() string {
	return "❌ 8 caractères minimum"
}

func (f FrenchBundle) GetUIPromptPasswordUppercase() string {
	return "❌ Une majuscule requise"
}

func (f FrenchBundle) GetUIPromptPasswordLowercase() string {
	return "❌ Une minuscule requise"
}

func (f FrenchBundle) GetUIPromptPasswordDigit() string {
	return "❌ Un chiffre requis"
}

func (f FrenchBundle) GetUIPromptPasswordSuccess() string {
	return "   ✓ %s"
}

func (f FrenchBundle) GetUIPromptWorkersLabel() string {
	return "⚙️  Workers (défaut: %d, max: %d)"
}

func (f FrenchBundle) GetUIPromptWorkersSuccess() string {
	return "   ✓ %d workers"
}

func (f FrenchBundle) GetUIPromptWorkersInvalid() string {
	return "❌ Nombre valide requis (>=1)"
}

func (f FrenchBundle) GetUIPromptWorkersMax() string {
	return "❌ Maximum %d workers"
}

func (f FrenchBundle) GetUIPromptConfirmLabel() string {
	return "❓ %s [%s]: "
}

func (f FrenchBundle) GetUIPromptConfirmInvalid() string {
	return "❌ Répondez par y/n"
}

// ============================================================================
// UI success messages
// ============================================================================

func (f FrenchBundle) GetUISuccessOperation() string {
	return "✅ Opération réussie !"
}

func (f FrenchBundle) GetUISuccessOutput() string {
	return "📄 Fichier: %s"
}

func (f FrenchBundle) GetUISuccessSize() string {
	return "📏 Taille:  %s"
}

// ============================================================================
// Cryptolib decryption errors
// ============================================================================

func (f FrenchBundle) GetCryptolibErrOpenInput() string {
	return "ouverture du fichier d'entrée: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateOutput() string {
	return "création du fichier de sortie: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateCipher() string {
	return "création du cipher: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateGCM() string {
	return "création du GCM: %w"
}

func (f FrenchBundle) GetCryptolibErrReadHeader() string {
	return "lecture de l'en-tête: %w"
}

func (f FrenchBundle) GetCryptolibErrReadHeaderHMAC() string {
	return "lecture du HMAC d'en-tête: %w"
}

func (f FrenchBundle) GetCryptolibErrReadNonce() string {
	return "lecture du nonce: %w"
}

func (f FrenchBundle) GetCryptolibErrUnexpectedEOF() string {
	return "EOF inattendu: marqueur de fin manquant"
}

func (f FrenchBundle) GetCryptolibErrReadChunkLen() string {
	return "lecture de la taille du chunk: %w"
}

func (f FrenchBundle) GetCryptolibErrReadCiphertext() string {
	return "lecture du chunk chiffré %d: %w"
}

func (f FrenchBundle) GetCryptolibErrDeriveNonce() string {
	return "dérivation du nonce pour le chunk %d: %w"
}

func (f FrenchBundle) GetCryptolibErrWritePlaintext() string {
	return "écriture du chunk clair %d: %w"
}

// ============================================================================
// Cryptolib encryption errors
// ============================================================================

func (f FrenchBundle) GetCryptolibErrOpenInputEnc() string {
	return "ouverture du fichier d'entrée: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateOutputEnc() string {
	return "création du fichier de sortie: %w"
}

func (f FrenchBundle) GetCryptolibErrGenerateSalt() string {
	return "génération du sel: %w"
}

func (f FrenchBundle) GetCryptolibErrWriteHeader() string {
	return "écriture de l'en-tête: %w"
}

func (f FrenchBundle) GetCryptolibErrWriteHeaderHMAC() string {
	return "écriture du HMAC d'en-tête: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateCipherEnc() string {
	return "création du cipher: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateGCMEnc() string {
	return "création du GCM: %w"
}

func (f FrenchBundle) GetCryptolibErrGenerateNonce() string {
	return "génération du nonce: %w"
}

func (f FrenchBundle) GetCryptolibErrWriteNonce() string {
	return "écriture du nonce: %w"
}

func (f FrenchBundle) GetCryptolibErrReadChunk() string {
	return "lecture du chunk: %w"
}

func (f FrenchBundle) GetCryptolibErrWriteEndMarker() string {
	return "écriture du marqueur de fin: %w"
}

func (f FrenchBundle) GetCryptolibErrNonceDerivation() string {
	return "échec de la dérivation du nonce: %w"
}

func (f FrenchBundle) GetCryptolibErrMissingChunks() string {
	return "chunks manquants: index attendu %d, %d en attente"
}

func (f FrenchBundle) GetCryptolibErrTooManyPending() string {
	return "trop de chunks en attente (limite %d) - possible attaque de réordonnancement"
}

func (f FrenchBundle) GetCryptolibErrWriteChunkLen() string {
	return "écriture de la taille du chunk: %w"
}

func (f FrenchBundle) GetCryptolibErrWriteCiphertext() string {
	return "écriture du texte chiffré: %w"
}

func (f FrenchBundle) GetCryptolibErrCloseInput() string {
	return "fermeture de l'entrée: %w"
}

func (f FrenchBundle) GetCryptolibErrCloseOutput() string {
	return "fermeture de la sortie: %w"
}

// ============================================================================
// Cryptolib sentinel errors
// ============================================================================

func (f FrenchBundle) GetCryptolibErrInvalidMagic() string {
	return "octets magiques invalides: fichier non chiffré avec cet outil"
}

func (f FrenchBundle) GetCryptolibErrUnsupportedVersion() string {
	return "version de fichier non supportée"
}

func (f FrenchBundle) GetCryptolibErrHeaderAuthFailed() string {
	return "échec d'authentification de l'en-tête: mauvais mot de passe ou fichier corrompu"
}

func (f FrenchBundle) GetCryptolibErrDecryptionFailed() string {
	return "échec du déchiffrement: données corrompues ou mauvaise clé"
}

func (f FrenchBundle) GetCryptolibErrChunkTooLarge() string {
	return "la taille du chunk dépasse la limite maximale autorisée"
}

// ============================================================================
// Cryptolib stream errors
// ============================================================================

func (f FrenchBundle) GetCryptolibErrReadHeaderStream() string {
	return "lecture de l'en-tête: %w"
}

func (f FrenchBundle) GetCryptolibErrReadHeaderHMACStream() string {
	return "lecture du HMAC d'en-tête: %w"
}

func (f FrenchBundle) GetCryptolibErrReadNonceStream() string {
	return "lecture du nonce: %w"
}

func (f FrenchBundle) GetCryptolibErrUnexpectedEOFStream() string {
	return "EOF inattendu: marqueur de fin manquant"
}

func (f FrenchBundle) GetCryptolibErrReadChunkLenStream() string {
	return "lecture de la taille du chunk: %w"
}

func (f FrenchBundle) GetCryptolibErrReadCiphertextStream() string {
	return "lecture du chunk chiffré %d: %w"
}

func (f FrenchBundle) GetCryptolibErrDeriveNonceStream() string {
	return "dérivation du nonce pour le chunk %d: %w"
}

func (f FrenchBundle) GetCryptolibErrWritePlaintextStream() string {
	return "écriture du chunk clair %d: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateCipherStream() string {
	return "création du cipher: %w"
}

func (f FrenchBundle) GetCryptolibErrCreateGCMStream() string {
	return "création du GCM: %w"
}
