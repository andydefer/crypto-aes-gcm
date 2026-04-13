# 🔐 Crypto-AES-GCM

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-passing-brightgreen.svg)](.)
[![Cobra](https://img.shields.io/badge/cli-cobra-blue)](https://github.com/spf13/cobra)
[![GoDoc](https://godoc.org/github.com/andydefer/crypto-aes-gcm?status.svg)](https://godoc.org/github.com/andydefer/crypto-aes-gcm)

Un outil de chiffrement sécurisé et performant utilisant **AES-256-GCM** en mode streaming avec dérivation de clé **Argon2id** et interface CLI moderne.

## ✨ Fonctionnalités

- 🔐 **Chiffrement AES-256-GCM** authentifié
- 🧂 **Dérivation de clé Argon2id** résistante aux attaques GPU
- 🚀 **Traitement parallèle** pour les gros fichiers (configurable)
- 📦 **Mode streaming** - aucun chargement complet en mémoire
- ✅ **Intégrité vérifiée** avec HMAC-SHA256 (header + contenu)
- 💻 **Interface CLI moderne** avec Cobra, couleurs et barre de progression
- 📚 **Bibliothèque Go** réutilisable
- 🔄 **Concurrent** - chiffrement multiple en parallèle
- 🎨 **Interface utilisateur interactive** avec prompts de confirmation

## 🏗️ Architecture

```
crypto-aes-gcm/
├── cmd/
│   └── cryptool/           # Application CLI (Cobra)
│       ├── main.go         # Point d'entrée
│       └── version.go      # Version info
├── internal/
│   ├── argon2/             # Dérivation de clé Argon2id
│   └── header/             # Utilitaires de header et HMAC
├── pkg/
│   └── cryptolib/          # Bibliothèque exportable
│       ├── encrypt.go      # Chiffrement parallèle
│       ├── decrypt.go      # Déchiffrement
│       ├── types.go        # Types et constantes
│       └── errors.go       # Erreurs sentinelles
├── go.mod
├── go.sum
├── Makefile                # Commandes de build et tests
└── README.md
```

## 📦 Installation

### Prérequis

- Go 1.23 ou supérieur
- Dépendances automatiquement gérées par Go modules

### Depuis les sources

```bash
git clone https://github.com/andydefer/crypto-aes-gcm.git
cd crypto-aes-gcm
make build
# ou
go build -o cryptool ./cmd/cryptool
```

### Installation système

```bash
sudo cp cryptool /usr/local/bin/
```

### Via go install

```bash
go install github.com/andydefer/crypto-aes-gcm/cmd/cryptool@latest
```

## 🚀 Utilisation CLI

### Syntaxe moderne (recommandée)

```bash
# Aide générale
cryptool --help
cryptool help

# Chiffrer un fichier
cryptool encrypt fichier.txt fichier.enc --pass "monMotDePasse"

# Déchiffrer un fichier
cryptool decrypt fichier.enc fichier.dec.txt --pass "monMotDePasse"

# Version
cryptool version
```

### Options disponibles

| Option | Description | Défaut |
|--------|-------------|--------|
| `--pass, -p` | Mot de passe (requis) | - |
| `--workers, -w` | Nombre de workers parallèles | 4 |
| `--force, -f` | Écraser sans confirmation | false |
| `--quiet, -q` | Mode silencieux (pas de barre de progression) | false |

### Exemples avancés

```bash
# Gros fichier avec workers optimisés (2× CPU cores max)
cryptool encrypt video.mp4 video.enc --pass "secure" --workers 8

# Mode silencieux pour les scripts
cryptool encrypt data.bin data.enc --pass "secret" --quiet

# Forcer l'écrasement sans confirmation
cryptool encrypt output.txt output.enc --pass "pass" --force

# Dans un pipeline (avec /dev/stdin)
cat secret.txt | cryptool encrypt /dev/stdin data.enc --pass "pass"

# Compression + chiffrement
tar czf - dossier/ | cryptool encrypt /dev/stdin backup.tar.gz.enc --pass "pass"
```

### Exemple d'exécution

```bash
$ cryptool encrypt test.txt test.enc --pass "Hello@0405" --workers 8

🔐 Crypto-AES-GCM - ENCRYPT MODE
──────────────────────────────────────────────────
📁 Input:   test.txt
📂 Output:  test.enc
⚙️  Workers: 8
──────────────────────────────────────────────────

🔒 Encrypting [████████████████████████████████████████] 100% (245 B/245 B)

✅ Encryption successful!
📄 Output: test.enc
📏 Size:   245 B
```

## 🔄 Migration depuis l'ancienne syntaxe

L'ancienne syntaxe avec flags (`-mode`, `-in`, `-out`) est toujours supportée pour la rétrocompatibilité, mais nous recommandons d'utiliser la nouvelle syntaxe avec sous-commandes.

| Ancienne syntaxe | Nouvelle syntaxe |
|-----------------|------------------|
| `cryptool -mode encrypt -in file.txt -out file.enc -pass pwd` | `cryptool encrypt file.txt file.enc --pass pwd` |
| `cryptool -mode decrypt -in file.enc -out dec.txt -pass pwd` | `cryptool decrypt file.enc dec.txt --pass pwd` |
| `cryptool -mode encrypt -in file.txt -out file.enc -pass pwd -workers 8` | `cryptool encrypt file.txt file.enc --pass pwd --workers 8` |
| `cryptool -mode encrypt -in file.txt -out file.enc -pass pwd -force` | `cryptool encrypt file.txt file.enc --pass pwd --force` |

## 📚 Utilisation comme bibliothèque

### Installation

```bash
go get github.com/andydefer/crypto-aes-gcm
```

### Exemple basique

```go
package main

import (
    "log"
    "github.com/andydefer/crypto-aes-gcm/pkg/cryptolib"
)

func main() {
    // Chiffrement
    encryptor, err := cryptolib.NewEncryptor(4)
    if err != nil {
        log.Fatal(err)
    }

    if err := encryptor.EncryptFile("input.txt", "output.enc", "password"); err != nil {
        log.Fatal(err)
    }

    // Pour déchiffrer, vous avez besoin du salt présent dans le fichier
    // Lisez d'abord le header pour extraire le salt
    // Voir l'exemple complet dans la documentation
}
```

### Streaming

```go
func encryptStream(r io.Reader, w io.Writer, pass string) error {
    encryptor, err := cryptolib.NewEncryptor(4)
    if err != nil {
        return err
    }
    return encryptor.Encrypt(r, w, pass)
}

func decryptStream(r io.Reader, w io.Writer, pass string) error {
    return cryptolib.DecryptStream(r, w, pass)
}
```

### Traitement concurrent de plusieurs fichiers

```go
func encryptFiles(files []string, pass string) error {
    encryptor, err := cryptolib.NewEncryptor(4)
    if err != nil {
        return err
    }

    var wg sync.WaitGroup
    errChan := make(chan error, len(files))

    for _, f := range files {
        wg.Add(1)
        go func(file string) {
            defer wg.Done()
            if err := encryptor.EncryptFile(file, file+".enc", pass); err != nil {
                errChan <- err
            }
        }(f)
    }

    wg.Wait()
    close(errChan)

    for err := range errChan {
        if err != nil {
            return err
        }
    }
    return nil
}
```

### Lecture du header pour extraction du salt

```go
func getSaltFromEncryptedFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    var header cryptolib.FileHeader
    if err := binary.Read(f, binary.BigEndian, &header); err != nil {
        return nil, err
    }

    return header.Salt[:], nil
}
```

## 🔒 Sécurité

### Paramètres cryptographiques

| Algorithme | Paramètres | Justification |
|------------|------------|----------------|
| **Chiffrement** | AES-256-GCM | Authenticated encryption with associated data (AEAD) |
| **Dérivation** | Argon2id (time=4, memory=64MB, threads=4) | Résistant aux attaques GPU et ASIC |
| **Intégrité header** | HMAC-SHA256 | Vérification avant déchiffrement |
| **Intégrité contenu** | HMAC-SHA256 global | Détection de corruption/modification |
| **Nonce** | 12 bytes (aléatoire + counter) | Counter sur 8 bytes, safe pour 2^64 chunks |
| **Salt** | 16 bytes (aléatoire) | Unique par fichier, prévient les rainbow tables |
| **Chunk size** | 1 MB par défaut | Bon compromis mémoire/performance |

### Format du fichier chiffré

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
│ HEADER HMAC (32 bytes)                                      │
├─────────────────────────────────────────────────────────────┤
│ BASE NONCE (12 bytes)                                       │
├─────────────────────────────────────────────────────────────┤
│ CHUNK 1                                                     │
│  ├─ Length: 4 bytes (uint32)                                │
│  └─ Ciphertext: variable (GCM sealed)                       │
├─────────────────────────────────────────────────────────────┤
│ CHUNK 2...N                                                 │
│  ├─ Length: 4 bytes                                         │
│  └─ Ciphertext: variable                                    │
├─────────────────────────────────────────────────────────────┤
│ END MARKER (4 bytes) = 0                                    │
├─────────────────────────────────────────────────────────────┤
│ GLOBAL HMAC (32 bytes) - HMAC-SHA256 of all ciphertexts     │
└─────────────────────────────────────────────────────────────┘
```

## 🧪 Tests

```bash
# Tous les tests
make test
# ou
go test -v ./...

# Avec couverture
make test-coverage
# ou
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Mode short (ignore les gros fichiers)
make test-short
# ou
go test -short -v ./...

# Avec race detector
go test -race ./...

# Benchmarks
go test -bench=. ./...
```

## 📊 Performance

Tests réalisés sur **Intel i7-1165G7 @ 2.80GHz**, SSD NVMe, 16GB RAM.

| Fichier | Workers | Temps encrypt | Temps decrypt | Vitesse |
|---------|---------|---------------|---------------|---------|
| 10 MB | 1 | 0.15s | 0.14s | ~70 MB/s |
| 10 MB | 4 | 0.08s | 0.07s | ~140 MB/s |
| 100 MB | 4 | 0.65s | 0.60s | ~160 MB/s |
| 1 GB | 8 | 6.2s | 5.8s | ~165 MB/s |
| 10 GB | 8 | 62s | 58s | ~165 MB/s |

*Les performances sont limitées par l'I/O disque au-delà de 1 GB.*

## 🛠️ Makefile Commands

```bash
# Aide
make help

# Build pour la plateforme courante
make build

# Build pour toutes les plateformes (Linux, Windows, macOS)
make build-all

# Tests
make test              # Tous les tests
make test-short        # Tests rapides
make test-coverage     # Tests avec couverture

# Nettoyage
make clean             # Fichiers temporaires
make clean-all         # Nettoyage complet

# Gestion de version
make git-commit-push   # Commit et push
make git-tag           # Créer un tag
make release           # Créer une release

# Concaténation des sources
make concat-all        # Génère all.txt
```

## 📖 API Reference

### Types exportés

```go
type Encryptor struct {
    // contient des champs non exportés
}

type Decryptor struct {
    // contient des champs non exportés
}

type FileHeader struct {
    Magic     [4]byte
    Version   byte
    Salt      [SaltSize]byte
    ChunkSize uint32
}
```

### Constantes

```go
const (
    Magic            = "CRYP"
    Version          = 2
    SaltSize         = 16
    NonceSize        = 12
    KeySize          = 32
    DefaultChunkSize = 1024 * 1024  // 1MB
    DefaultWorkers   = 4
)
```

### Fonctions - Encryptor

```go
// NewEncryptor crée un nouvel encrypteur avec le nombre de workers spécifié.
// Le nombre est automatiquement limité entre 1 et 2×CPU cores.
func NewEncryptor(workers int) (*Encryptor, error)

// EncryptFile chiffre un fichier sur le disque.
func (e *Encryptor) EncryptFile(inputPath, outputPath, passphrase string) error

// Encrypt lit depuis un io.Reader et écrit les données chiffrées dans un io.Writer.
func (e *Encryptor) Encrypt(r io.Reader, w io.Writer, passphrase string) error
```

### Fonctions - Decryptor

```go
// NewDecryptor crée un nouveau décrypteur avec la passphrase et le salt fournis.
func NewDecryptor(passphrase string, salt []byte) (*Decryptor, error)

// DecryptFile déchiffre un fichier sur le disque.
func (d *Decryptor) DecryptFile(inputPath, outputPath string) error

// Decrypt lit depuis un io.Reader et écrit les données déchiffrées dans un io.Writer.
func (d *Decryptor) Decrypt(r io.Reader, w io.Writer) error
```

### Fonctions utilitaires

```go
// DecryptStream est une fonction pratique pour déchiffrer directement depuis un Reader.
// Combine la lecture du header, la vérification HMAC et le déchiffrement.
func DecryptStream(r io.Reader, w io.Writer, passphrase string) error
```

### Erreurs

```go
var (
    ErrInvalidMagic       = errors.New("invalid magic bytes")
    ErrUnsupportedVersion = errors.New("unsupported file version")
    ErrHeaderAuthFailed   = errors.New("header authentication failed")
    ErrGlobalHMACFailed   = errors.New("global HMAC verification failed")
    ErrDecryptionFailed   = errors.New("decryption failed")
)
```

## 🔧 Dépannage

### Erreur : "invalid magic bytes"

**Cause** : Le fichier n'a pas été chiffré avec cet outil ou est corrompu.

**Solution** : Vérifiez que vous utilisez le bon fichier et qu'il a été chiffré avec cryptool.

### Erreur : "header authentication failed"

**Cause** : Mot de passe incorrect ou fichier corrompu.

**Solution** : Vérifiez votre mot de passe. Si le mot de passe est correct, le fichier est probablement corrompu.

### Erreur : "decryption failed (chunk X)"

**Cause** : Fichier corrompu ou clé incorrecte.

**Solution** : Le fichier a été modifié après chiffrement. Utilisez une sauvegarde si disponible.

## 🤝 Contribution

1. Fork le projet
2. Créez votre branche (`git checkout -b feature/amazing`)
3. Committez vos changements (`git commit -m 'feat: add amazing feature'`)
4. Push vers la branche (`git push origin feature/amazing`)
5. Ouvrez une Pull Request

### Conventions de commit

Nous suivons [Conventional Commits](https://www.conventionalcommits.org/) :

- `feat:` nouvelle fonctionnalité
- `fix:` correction de bug
- `docs:` documentation
- `test:` tests
- `refactor:` refactorisation
- `chore:` maintenance

## 📝 License

MIT License - voir le fichier [LICENSE](LICENSE) pour plus de détails.

## ⚠️ Avertissement

Ce logiciel est fourni "tel quel". Bien que des efforts aient été faits pour assurer sa sécurité, utilisez-le à vos propres risques. Pour des données extrêmement sensibles, consultez un expert en sécurité.

**Recommandations de sécurité :**
- Utilisez des mots de passe longs et complexes (minimum 12 caractères)
- Ne partagez jamais vos mots de passe
- Sauvegardez vos fichiers chiffrés
- Vérifiez l'intégrité des fichiers après déchiffrement

## 🙏 Remerciements

- [Argon2](https://github.com/P-H-C/phc-winner-argon2) - Password Hashing Competition winner
- [AES-GCM](https://csrc.nist.gov/publications/detail/sp/800-38d/final) - NIST standard
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Color](https://github.com/fatih/color) - Terminal colors
- [Progressbar](https://github.com/schollz/progressbar) - Progress bars

## 📞 Support

- 📧 Email: votre-email@example.com
- 🐛 Issues: [GitHub Issues](https://github.com/andydefer/crypto-aes-gcm/issues)
- 📖 Documentation: [GoDoc](https://godoc.org/github.com/andydefer/crypto-aes-gcm)

## 📈 Roadmap

- [ ] Support du chiffrement asymétrique (age)
- [ ] Compression automatique avant chiffrement
- [ ] Mode archive (multiple fichiers)
- [ ] Interface TUI avec BubbleTea
- [ ] Support YubiKey/PIV
- [ ] Déchiffrement parallèle

---

**Made with 🔐 by andydefer**

*Version 2.0.0 - Interface CLI moderne avec Cobra*

