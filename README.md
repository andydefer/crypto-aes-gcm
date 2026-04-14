# 🔐 Crypto-AES-GCM

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-200%20passed-brightgreen.svg)](.)
[![Cobra](https://img.shields.io/badge/cli-cobra-blue)](https://github.com/spf13/cobra)
[![GoDoc](https://godoc.org/github.com/andydefer/crypto-aes-gcm?status.svg)](https://godoc.org/github.com/andydefer/crypto-aes-gcm)

Un outil de chiffrement sécurisé et performant utilisant **AES-256-GCM** en mode streaming avec dérivation de clé **Argon2id** et interface CLI moderne.

## ✨ Fonctionnalités

- 🔐 **Chiffrement ultra-sécurisé** avec AES-256-GCM (authentification par chunk)
- 🎮 **Mode interactif** - Interface guidée sans ligne de commande complexe
- 🔑 **Prompt interactif pour mot de passe** - Plus besoin d'exposer le mot de passe dans la ligne de commande
- 🚀 **Traitement parallèle** configurable (jusqu'à 2×CPU cores)
- 📦 **Mode streaming pur** - mémoire constante (O(1))
- 💻 **Interface colorée** avec barre de progression
- 🔄 **Deux modes d'utilisation** : CLI classique ou interactif
- 🛡️ **Validation des mots de passe** (8+ caractères, maj/min, chiffre)
- ✅ **Intégrité vérifiée** par chunk (GCM AEAD)
- 🧂 **Dérivation Argon2id** résistante aux attaques GPU/ASIC

---

# 👤 Pour les utilisateurs finaux

## Installation rapide

### Linux / macOS

```bash
# Compilation depuis les sources
git clone https://github.com/andydefer/crypto-aes-gcm.git
cd crypto-aes-gcm
make build
sudo cp build/aescryptool /usr/local/bin/

# Ou installation automatique via Makefile
make install              # Installation globale (/usr/local/bin)
make install-local        # Installation locale (~/.local/bin)

# Ou télécharger le binaire (releases)
wget https://github.com/andydefer/crypto-aes-gcm/releases/latest/download/aescryptool-linux-amd64
chmod +x aescryptool-linux-amd64
sudo mv aescryptool-linux-amd64 /usr/local/bin/aescryptool
```

### Windows

```powershell
# Télécharger aescryptool-windows-amd64.exe
# Le placer dans C:\Windows\System32\
```

### Vérifier l'installation

```bash
aescryptool version
make install-check         # Vérifie si l'installation est correcte
```

## 🎮 Mode interactif

Le mode interactif guide l'utilisateur étape par étape sans mémoriser les options :

```bash
aescryptool interact
```

### Menu interactif :

```
╔════════════════════════════════════════════════════════════════════╗
║                 🎮 AESCRYPTOOL - MODE INTERACTIF                   ║
║                                                                    ║
║  Suivez les invites pour chiffrer ou déchiffrer vos fichiers       ║
║  Toutes les entrées seront validées avant exécution                ║
║                                                                    ║
║  Ctrl+C = Retour au menu | Ctrl+D = Quitter                        ║
║                                                                    ║
╚════════════════════════════════════════════════════════════════════╝

📋 Que souhaitez-vous faire:
  ▸ 🔒  Chiffrer un fichier
    🔓  Déchiffrer un fichier
    🚪  Quitter
```

### Exemple de session interactive (chiffrement) :

```
🔐 CHIFFREMENT DE FICHIER
────────────────────────────────────────

📁 Fichier à chiffrer: document.pdf
   ✓ document.pdf

📂 Fichier de sortie: document.pdf.enc
   ✓ document.pdf.enc

🔑 Mot de passe: **********
   ✓ **********

✅ Confirmation: **********
   ✓ **********

⚙️  Workers (défaut: 4, max: 16): 4
   ✓ 4 workers

❓ ⚠️  Le fichier existe déjà. Écraser ? [Y/n]: y

🔒 Chiffrement [████████████████████] 100%
✅ Fichier chiffré : document.pdf.enc
```

## Commandes de base

### Mode interactif (recommandé pour usage manuel)

```bash
# Lance le mode interactif avec prompts guidés
aescryptool interact
```

### Mode non-interactif (CLI classique)

Le flag `--pass` est désormais **optionnel**. Si omis, le mot de passe est demandé interactivement :

```bash
# Avec prompt interactif (recommandé)
aescryptool encrypt monfichier.txt monfichier.enc
aescryptool decrypt monfichier.enc monfichier.txt

# Avec flag --pass (pour scripts/automatisation)
aescryptool encrypt monfichier.txt monfichier.enc --pass "monMotDePasse"
aescryptool decrypt monfichier.enc monfichier.txt --pass "monMotDePasse"

# Aide
aescryptool --help
aescryptool encrypt --help
```

## Options importantes

| Option | Description |
|--------|-------------|
| `--pass, -p` | Mot de passe (optionnel - sera demandé interactivement si omis) |
| `--workers, -w` | Accélération pour gros fichiers (défaut: 4, max: 2×CPU) |
| `--force, -f` | Écraser sans demander confirmation |
| `--quiet, -q` | Mode silencieux (pas de barre de progression) |

## Exemples quotidiens

### 📄 Documents personnels

```bash
# Recommandé - avec prompt interactif
aescryptool encrypt declaration-2024.pdf declaration-2024.pdf.enc

# Pour scripts - avec flag --pass
aescryptool encrypt declaration-2024.pdf declaration-2024.pdf.enc --pass "MotDePasseFort123!"

# Mode interactif complet
aescryptool interact
# → Choisir "Chiffrer" → suivre les invites
```

### 🎥 Vidéos (gros fichiers)

```bash
# Avec optimisation parallèle (8 workers) et prompt interactif
aescryptool encrypt video.mp4 video.mp4.enc --workers 8

# Avec flag --pass pour scripts
aescryptool encrypt video.mp4 video.mp4.enc --pass "Vacances2024!" --workers 8
```

### 📦 Chiffrement de dossier

```bash
# Compresser + chiffrer (avec prompt)
tar czf - dossier/ | aescryptool encrypt /dev/stdin backup.enc

# Déchiffrer + décompresser (avec prompt)
aescryptool decrypt backup.enc /dev/stdout | tar xzf -
```

## Dépannage rapide

| Problème | Cause probable | Solution |
|----------|---------------|----------|
| `le fichier n'existe pas` | Fichier source introuvable | Vérifier le chemin |
| `le fichier existe déjà` | Fichier destination existe | Utiliser `--force` ou confirmer l'écrasement |
| `le mot de passe ne correspond pas` | Confirmation erronée (chiffrement) | Ressaisir correctement |
| `le mot de passe ne peut pas être vide` | Aucun mot de passe fourni | Saisir un mot de passe valide |
| `header authentication failed` | Mot de passe incorrect | Vérifier la casse et les caractères spéciaux |

### Exigences des mots de passe (chiffrement interactif)

Lorsque vous chiffrez un fichier sans utiliser le flag `--pass`, le mode interactif valide la force du mot de passe :

- Minimum **8 caractères**
- Au moins **une majuscule** (A-Z)
- Au moins **une minuscule** (a-z)
- Au moins **un chiffre** (0-9)

> 💡 **Note** : Le flag `--pass` contourne cette validation pour les scripts.

### Sécurité des mots de passe

| Méthode | Sécurité | Usage recommandé |
|---------|----------|------------------|
| Prompt interactif | ✅✅✅ Très bonne | Usage manuel |
| Flag `--pass` | ⚠️ Visible dans `ps aux` | Scripts (environnement contrôlé) |

> ⚠️ **Avertissement** : Évitez le flag `--pass` sur les systèmes multi-utilisateurs car la ligne de commande est visible par tous les processus (`ps aux`, `/proc/.../cmdline`). Préférez le prompt interactif pour une meilleure sécurité.

---

# 👨‍💻 Pour les développeurs

## Architecture

```
crypto-aes-gcm/
├── cmd/
│   └── aescryptool/           # Application CLI (Cobra)
│       ├── main.go            # Point d'entrée (bootstrap pur)
│       └── main_test.go       # Tests unitaires
├── internal/
│   ├── argon2/                # Dérivation de clé Argon2id
│   ├── cli/                   # Commandes Cobra (encrypt, decrypt, interact)
│   │   └── password.go        # Gestion interactive des mots de passe
│   ├── header/                # Sérialisation et validation des headers
│   ├── service/               # Orchestration métier
│   ├── ui/                    # Interface utilisateur (couleurs, prompts, progress)
│   └── utils/                 # Utilitaires (formatage taille)
├── pkg/
│   └── cryptolib/             # Bibliothèque exportable (cœur crypto)
│       ├── encrypt.go         # Chiffrement parallèle
│       ├── decrypt.go         # Déchiffrement streaming
│       ├── stream.go          # API streaming simplifiée (DecryptStream)
│       ├── types.go           # Types et constantes
│       ├── errors.go          # Erreurs sentinelles
│       └── *_test.go          # Tests complets (>150 tests)
├── tests/                     # Tests shell et scénarios
│   ├── run_tests.sh           # Scripts de test réalistes
│   ├── test_scenarios.sh      # Scénarios avancés
│   └── generate_test_files.sh
├── build/                     # Binaires compilés
├── private/                   # Artefacts générés (diffs, concat)
├── Makefile                   # Commandes de build, tests, release, install
└── README.md
```

## Installation pour les développeurs

### Prérequis

- Go 1.23 ou supérieur
- `make` (optionnel)
- `gotestsum` (optionnel)

### Depuis les sources

```bash
git clone https://github.com/andydefer/crypto-aes-gcm.git
cd crypto-aes-gcm

# Build
make build

# Tests
make test

# Installation (choisissez l'une des options)
make install              # Installation globale (/usr/local/bin) - nécessite sudo
make install-local        # Installation locale (~/.local/bin) - utilisateur uniquement
go install ./cmd/aescryptool

# Vérification
make install-check
aescryptool version
```

## API Reference

### Types exportés

```go
// Encryptor handles parallel streaming encryption
type Encryptor struct {
    // champs non exportés
}

// Decryptor handles streaming decryption
type Decryptor struct {
    // champs non exportés
}

// EncryptorConfig holds configuration options for the Encryptor
type EncryptorConfig struct {
    Workers          int  // Parallel workers (default: 4)
    ChunkSize        int  // Chunk size in bytes (default: 1MB)
    MaxPendingChunks int  // Max out-of-order chunks buffered (default: 100)
}

// FileHeader represents the encrypted file header
type FileHeader struct {
    Magic     [4]byte    // "CRYP"
    Version   byte       // Format version (2)
    Salt      [16]byte   // Argon2id salt
    ChunkSize uint32     // Chunk size in bytes (1MB default)
}
```

### Constantes

```go
const (
    Magic                  = "CRYP"
    Version                = 2
    SaltSize               = 16
    NonceSize              = 12
    KeySize                = 32
    DefaultChunkSize       = 1024 * 1024  // 1MB
    DefaultWorkers         = 4
    DefaultMaxPendingChunks = 100         // Anti-DoS
    MaxMaxPendingChunks     = 1000        // Absolute maximum
)
```

### Fonctions principales

```go
// Encryptor
func NewEncryptor(workers int) (*Encryptor, error)
func NewEncryptorWithConfig(config EncryptorConfig) (*Encryptor, error)
func (e *Encryptor) EncryptFile(inputPath, outputPath, passphrase string) error
func (e *Encryptor) Encrypt(r io.Reader, w io.Writer, passphrase string) error

// Decryptor
func NewDecryptor(passphrase string, salt []byte) (*Decryptor, error)
func (d *Decryptor) DecryptFile(inputPath, outputPath string) error
func (d *Decryptor) Decrypt(r io.Reader, w io.Writer) error

// Convenience function
func DecryptStream(r io.Reader, w io.Writer, passphrase string) error

// Default configuration
func DefaultEncryptorConfig() EncryptorConfig
```

### Erreurs

```go
var (
    ErrInvalidMagic       = errors.New("invalid magic bytes: file not encrypted with this tool")
    ErrUnsupportedVersion = errors.New("unsupported file version")
    ErrHeaderAuthFailed   = errors.New("header authentication failed: wrong passphrase or corrupted file")
    ErrDecryptionFailed   = errors.New("decryption failed: corrupted data or wrong key")
)
```

## Makefile Commands

```bash
# Aide complète
make help

# 🚀 Run (exécution directe sans build)
make run-interact          # Mode interactif
make run-version           # Affiche la version
make run ARGS="encrypt test.txt test.enc"  # Arguments personnalisés (prompt pour mot de passe)
make run-encrypt INPUT=file.txt OUTPUT=file.enc PASS=secret  # Avec flag --pass
make run-decrypt INPUT=file.enc OUTPUT=file.txt PASS=secret

# 🔨 Build
make build                 # Build pour plateforme courante
make build-all             # Build multi-plateformes (Linux, Windows, macOS)

# 📦 Installation
make install               # Installation globale (/usr/local/bin) - nécessite sudo
make install-local         # Installation locale (~/.local/bin) - utilisateur uniquement
make uninstall             # Désinstallation globale et locale
make reinstall             # Réinstallation complète
make install-check         # Vérifie si l'installation est correcte

# 🧪 Tests Go
make test                  # Tous les tests Go
make test-short            # Tests rapides
make test-coverage         # Tests avec couverture
make gotestsum             # Tests avec formateur

# 🧪 Fuzzing
make fuzz                  # Fuzz tests (1 minute each)
make fuzz-short            # Fuzz tests (10 seconds each)

# 📊 Benchmarks
make bench                 # Tous les benchmarks
make bench-cpu             # Benchmarks avec profiling CPU
make bench-mem             # Benchmarks avec profiling mémoire

# 📋 Tests shell réalistes
make test-scripts          # Scripts de test
make test-scenarios        # Scénarios avancés
make test-all              # Tous les tests (Go + scripts)

# 📁 Génération fichiers de test
make generate-test-files       # Génère fichiers (dont 50MB+)
make generate-test-files-short # Mode court

# 🧹 Nettoyage
make clean                # Fichiers temporaires
make clean-test-data      # Données de test
make clean-all            # Nettoyage complet

# 🔄 Version Control
make git-commit-push      # Commit et push
make git-tag              # Créer un tag
make generate-ai-diff     # Génère diff pour revue AI
make release              # Créer une release
```

## Tests

```bash
# Tous les tests Go (>150 tests)
make test

# Tests avec gotestsum (formatage amélioré)
make gotestsum

# Fuzzing (détection de bugs)
make fuzz

# Benchmarks de performance
make bench

# Tests shell réalistes
make test-scripts

# Scénarios avancés
make test-scenarios

# Tous les tests (Go + scripts + scénarios)
make test-all

# Mode court (ignore gros fichiers)
make test-all-short

# Avec couverture
make test-coverage
go tool cover -html=coverage.out
```

## Format du fichier chiffré (v2.0.0)

Le format de fichier est le suivant :

```
┌─────────────────────────────────────────────────────────────┐
│                        FILE FORMAT                          │
├─────────────────────────────────────────────────────────────┤
│ HEADER (25 bytes)                                           │
│  ├─ Magic: "CRYP" (4 bytes)                                 │
│  ├─ Version: 2 (1 byte)                                     │
│  ├─ Salt: 16 bytes (Argon2id salt)                          │
│  └─ ChunkSize: 4 bytes (uint32)                             │
├─────────────────────────────────────────────────────────────┤
│ HEADER HMAC (32 bytes) - HMAC-SHA256 du header              │
├─────────────────────────────────────────────────────────────┤
│ BASE NONCE (12 bytes) - Généré aléatoirement                │
├─────────────────────────────────────────────────────────────┤
│ CHUNK 1                                                     │
│  ├─ Length: 4 bytes (uint32)                                │
│  └─ Ciphertext: variable (GCM sealed + authentifié)         │
├─────────────────────────────────────────────────────────────┤
│ CHUNK 2...N (même structure)                                │
├─────────────────────────────────────────────────────────────┤
│ END MARKER (4 bytes) = 0                                    │
└─────────────────────────────────────────────────────────────┘
```

**Authentification :**
- L'en-tête est protégé par **HMAC-SHA256** (vérifié avant tout déchiffrement)
- Chaque chunk de données est authentifié par **AES-256-GCM** (AEAD)
- La clé est dérivée du mot de passe avec **Argon2id** (time=4, memory=64MB, threads=4)

## Performance

Tests sur Intel i7-1165G7 @ 2.80GHz, SSD NVMe

| Fichier | Workers | Encrypt | Decrypt | Vitesse |
|---------|---------|---------|---------|---------|
| 10 MB | 1 | 0.13s | 0.14s | ~75 MB/s |
| 10 MB | 4 | 0.11s | 0.12s | ~90 MB/s |
| 10 MB | 8 | 0.10s | 0.11s | ~100 MB/s |
| 100 MB | 4 | 0.65s | 0.60s | ~160 MB/s |
| 1 GB | 8 | 6.2s | 5.8s | ~165 MB/s |

## Sécurité - Détails techniques

| Algorithme | Paramètres | Justification |
|------------|------------|----------------|
| **Chiffrement** | AES-256-GCM | Authenticated encryption (AEAD) |
| **Dérivation** | Argon2id (time=4, memory=64MB, threads=4) | Résistant aux attaques GPU/ASIC |
| **Intégrité header** | HMAC-SHA256 | Vérification avant déchiffrement |
| **Authentification** | GCM par chunk | Chaque chunk authentifié individuellement |
| **Nonce** | 12 bytes (base + XOR avec index) | Safe pour 2^64 chunks |
| **Salt** | 16 bytes (crypto/rand) | Unique par fichier |
| **Chunk size** | 1 MB | Bon compromis mémoire/performance |

## Contribution

### Conventions de commit

Nous suivons [Conventional Commits](https://www.conventionalcommits.org/) :

- `feat:` nouvelle fonctionnalité
- `fix:` correction de bug
- `docs:` documentation
- `test:` tests
- `refactor:` refactorisation
- `chore:` maintenance
- `perf:` amélioration de performance

### Processus

1. Fork le projet
2. Créez votre branche (`git checkout -b feature/amazing`)
3. Committez (`git commit -m 'feat: add amazing feature'`)
4. Push (`git push origin feature/amazing`)
5. Ouvrez une Pull Request

## 📝 License

MIT License

## ⚠️ Avertissement

Ce logiciel est fourni "tel quel". Pour des données extrêmement sensibles, consultez un expert en sécurité.

## 🙏 Remerciements

- [Argon2](https://github.com/P-H-C/phc-winner-argon2) - Password Hashing Competition winner
- [AES-GCM](https://csrc.nist.gov/publications/detail/sp/800-38d/final) - NIST standard
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts
- [progressbar](https://github.com/schollz/progressbar) - Progress bars

## 📞 Support

- 🐛 Issues: [GitHub Issues](https://github.com/andydefer/crypto-aes-gcm/issues)
- 📖 Documentation: [GoDoc](https://godoc.org/github.com/andydefer/crypto-aes-gcm)

---

**Made with 🔐 by andydefer**

*Version 2.0.0 - Mode interactif + Prompt mot de passe + Installation automatique + Streaming pur + Authentification par chunk*

