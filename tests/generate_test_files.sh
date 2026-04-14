#!/bin/bash

# Script de génération des fichiers de test pour cryptool
# Usage: ./generate_test_files.sh [--short]

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="${SCRIPT_DIR}/test_data"
INPUT_DIR="${TEST_DIR}/input"

# Mode court (ignore les gros fichiers)
SHORT_MODE="false"
if [ "$1" = "--short" ] || [ "$1" = "-s" ]; then
    SHORT_MODE="true"
fi

print_header() {
    echo ""
    echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║           📁 GÉNÉRATION DES FICHIERS DE TEST                 ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_info() {
    echo -e "${BLUE}📁 $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Création des répertoires
setup_directories() {
    print_info "Création des répertoires..."
    mkdir -p "$INPUT_DIR"
    mkdir -p "${TEST_DIR}/encrypted"
    mkdir -p "${TEST_DIR}/decrypted"
    mkdir -p "${TEST_DIR}/temp"
    print_success "Répertoires créés"
}

# Génération des fichiers
generate_files() {
    print_info "Génération des fichiers de test..."

    # 1. Petit fichier texte (2KB)
    print_info "  → small.txt (2KB)"
    cat > "$INPUT_DIR/small.txt" << 'EOF'
╔═══════════════════════════════════════════════════════════════════════════════╗
║                          CRYPTOOL TEST FILE v1.0                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝

Ce fichier est utilisé pour tester les fonctionnalités de base de cryptool.

🔐 CRYPTOOL est un outil de chiffrement sécurisé utilisant :
  • AES-256-GCM pour le chiffrement
  • Argon2id pour la dérivation de clé
  • Chunking parallèle pour les gros fichiers
  • Streaming pour une utilisation mémoire optimisée

📋 Ce fichier contient :
  ✓ Texte en français
  ✓ Texte en anglais
  ✓ Caractères Unicode
  ✓ Emojis
  ✓ Lignes de test

═══════════════════════════════════════════════════════════════════════════════

ENGLISH SECTION:
This file is used to test the basic functionalities of cryptool.

CRYPTOOL is a secure encryption tool using:
  • AES-256-GCM for encryption
  • Argon2id for key derivation
  • Parallel chunking for large files
  • Streaming for optimized memory usage

═══════════════════════════════════════════════════════════════════════════════

TEST DATA:
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

LINE 25: Test line for chunk boundary verification
LINE 26: Another test line to ensure streaming works correctly
LINE 27: This file should be small enough for quick tests
LINE 28: But large enough to test multiple chunks if chunk size is small
LINE 29: End of test file - cryptool encryption verification
EOF

    # 2. Fichier JSON (config)
    print_info "  → config.json (1KB)"
    cat > "$INPUT_DIR/config.json" << 'EOF'
{
    "cryptool_test": {
        "version": "2.0.0",
        "description": "Test configuration file for cryptool",
        "encryption": {
            "algorithm": "AES-256-GCM",
            "key_derivation": "Argon2id",
            "recommended_workers": 4,
            "chunk_size": 1048576
        },
        "test_scenarios": [
            {
                "name": "basic_encrypt_decrypt",
                "enabled": true,
                "iterations": 10
            },
            {
                "name": "parallel_workers",
                "enabled": true,
                "workers": [1, 2, 4, 8, 16]
            },
            {
                "name": "large_file",
                "enabled": true,
                "sizes_mb": [1, 10, 50, 100]
            },
            {
                "name": "corruption_detection",
                "enabled": true
            }
        ],
        "performance": {
            "target_encrypt_mb_per_sec": 200,
            "target_decrypt_mb_per_sec": 200,
            "max_memory_mb": 512
        },
        "test_data": {
            "strings": [
                "Hello World",
                "Bonjour le monde",
                "Hola Mundo",
                "Ciao Mondo",
                "こんにちは世界"
            ],
            "numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
            "booleans": [true, false, true, false]
        }
    }
}
EOF

    # 3. Fichier XML
    print_info "  → data.xml (1KB)"
    cat > "$INPUT_DIR/data.xml" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<cryptool_tests>
    <metadata>
        <generated_at>2024-01-01T00:00:00Z</generated_at>
        <tool>cryptool-test-generator</tool>
        <version>2.0.0</version>
    </metadata>

    <test_suite name="encryption">
        <test_case id="TC001">
            <name>Basic Encryption</name>
            <description>Test basic encrypt/decrypt functionality</description>
            <expected_result>success</expected_result>
        </test_case>
        <test_case id="TC002">
            <name>Wrong Password</name>
            <description>Test decryption with incorrect password</description>
            <expected_result>failure</expected_result>
        </test_case>
        <test_case id="TC003">
            <name>Corrupted File</name>
            <description>Test detection of corrupted encrypted files</description>
            <expected_result>failure</expected_result>
        </test_case>
    </test_suite>

    <test_suite name="performance">
        <test_case id="PC001">
            <name>Parallel Workers</name>
            <parameters>
                <workers>1</workers>
                <workers>2</workers>
                <workers>4</workers>
                <workers>8</workers>
            </parameters>
        </test_case>
        <test_case id="PC002">
            <name>Large File Streaming</name>
            <parameters>
                <size_mb>10</size_mb>
                <size_mb>50</size_mb>
                <size_mb>100</size_mb>
            </parameters>
        </test_case>
    </test_suite>

    <test_suite name="integrity">
        <test_case id="IC001">
            <name>Header HMAC</name>
            <description>Verify header HMAC validation</description>
        </test_case>
        <test_case id="IC002">
            <name>Chunk Authentication</name>
            <description>Verify GCM authentication per chunk</description>
        </test_case>
    </test_suite>
</cryptool_tests>
EOF

    # 4. Fichier avec caractères spéciaux
    print_info "  → special.txt (2KB)"
    cat > "$INPUT_DIR/special.txt" << 'EOF'
╔══════════════════════════════════════════════════════════════════╗
║                    CARACTÈRES SPÉCIAUX                           ║
╚══════════════════════════════════════════════════════════════════╝

FRANÇAIS:
À Á Â Ã Ä Å Æ Ç È É Ê Ë Ì Í Î Ï Ð Ñ Ò Ó Ô Õ Ö Ø Ù Ú Û Ü Ý Þ ß
à á â ã ä å æ ç è é ê ë ì í î ï ð ñ ò ó ô õ ö ø ù ú û ü ý þ ÿ

ACCENTS:
é è ê ë
ç
à â æ
î ï
ô œ
ù û ü
ÿ

EMOJIS & SYMBOLS:
🚀 🔐 💻 📁 🔒 ✅ ❌ ⚠️ ℹ️ 🎉 📊 📈 📉 🔄 🔧 ⚙️ 🔨 🛠️ 📦 🗑️
😀 😃 😄 😁 😆 😅 😂 🤣 😊 😇 🙂 🙃 😉 😌 😍 🥰 😘 😗 😙 😚
❤️ 🧡 💛 💚 💙 💜 🖤 🤍 🤎 💔 ❣️ 💕 💞 💓 💗 💖 💘 💝

CURRENCIES:
$ € £ ¥ ₣ ₤ ₧ ₨ ₩ ₪ ₫ ₭ ₮ ₯ ₹ ₺ ₼ ₽ ₾ ₿

MATHEMATICAL:
∑ ∏ ∫ ∂ ∇ √ ∞ ∝ ∠ ∧ ∨ ∩ ∪ ⊂ ⊃ ⊆ ⊇ ∈ ∉ ∋ ∀ ∃ ∄ ∅
≈ ≠ ≡ ≤ ≥ ≪ ≫ ⊕ ⊗ ⊥ ⊤ ⊢ ⊨ ∴ ∵ ∼ ≃ ≅ ≈

ARROWS:
← ↑ → ↓ ↔ ↕ ↖ ↗ ↘ ↙ ↚ ↛ ↜ ↝ ↞ ↟ ↠ ↡ ↢ ↣ ↤ ↥ ↦ ↧
↨ ↩ ↪ ↫ ↬ ↭ ↮ ↯ ↰ ↱ ↲ ↳ ↴ ↵ ↶ ↷ ↸ ↹ ↺ ↻ ↼ ↽ ↾ ↿

BOX DRAWING:
┌ ┐ └ ┘ ├ ┤ ┬ ┴ ┼ ─ │ ━ ┃ ┏ ┓ ┗ ┛ ┣ ┫ ┳ ┻ ╋ ┠ ┨ ┯ ┷ ┿
╭ ╮ ╯ ╰ ╱ ╲ ╳ ╴ ╵ ╶ ╷ ╸ ╹ ╺ ╻ ╼ ╽ ╾ ╿

HTML ENTITIES:
&lt; &gt; &amp; &quot; &apos; &copy; &reg; &trade; &euro; &pound;
&yen; &sect; &para; &dagger; &Dagger; &bull; &hellip; &permil;

ESCAPED CHARACTERS:
\t \n \r \\ \' \" \0 \x00 \u0000 \U00000000

BINARY DATA REPRESENTATION:
\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0A\x0B\x0C\x0D\x0E\x0F
\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F
\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F

EOF

    # 5. Fichier CSV
    print_info "  → data.csv (2KB)"
    cat > "$INPUT_DIR/data.csv" << 'EOF'
id,name,email,phone,city,country,active,score
1,John Doe,john.doe@example.com,+1234567890,New York,USA,true,95.5
2,Jane Smith,jane.smith@example.com,+1234567891,Los Angeles,USA,true,87.3
3,Pierre Martin,pierre.martin@example.fr,+33123456789,Paris,France,true,92.8
4,Maria Garcia,maria.garcia@example.es,+34912345678,Madrid,Spain,false,76.2
5,Wei Zhang,wei.zhang@example.cn,+86123456789,Beijing,China,true,88.9
6,Anna Kowalski,anna.kowalski@example.pl,+48123456789,Warsaw,Poland,true,91.4
7,Carlos Lopez,carlos.lopez@example.mx,+52123456789,Mexico City,Mexico,false,69.7
8,Aisha Khan,aisha.khan@example.pk,+92123456789,Karachi,Pakistan,true,94.1
9,Hans Schmidt,hans.schmidt@example.de,+49123456789,Berlin,Germany,true,86.5
10,Sakura Tanaka,sakura.tanaka@example.jp,+81123456789,Tokyo,Japan,false,79.8
EOF

    # 6. Fichier binaire 1MB
    print_info "  → random1.bin (1MB)"
    dd if=/dev/urandom of="$INPUT_DIR/random1.bin" bs=1M count=1 2>/dev/null
    print_success "    random1.bin généré (1MB)"

    # 7. Fichier binaire 10MB
    print_info "  → random10.bin (10MB)"
    dd if=/dev/urandom of="$INPUT_DIR/random10.bin" bs=1M count=10 2>/dev/null
    print_success "    random10.bin généré (10MB)"

    # 8. Fichier binaire 50MB (sauf mode court)
    if [ "$SHORT_MODE" != "true" ]; then
        print_info "  → random50.bin (50MB)"
        dd if=/dev/urandom of="$INPUT_DIR/random50.bin" bs=1M count=50 2>/dev/null
        print_success "    random50.bin généré (50MB)"
    else
        print_warning "  → random50.bin ignoré (mode court)"
    fi

    # 9. Fichier binaire 100MB (sauf mode court)
    if [ "$SHORT_MODE" != "true" ]; then
        print_info "  → random100.bin (100MB)"
        dd if=/dev/urandom of="$INPUT_DIR/random100.bin" bs=1M count=100 2>/dev/null
        print_success "    random100.bin généré (100MB)"
    else
        print_warning "  → random100.bin ignoré (mode court)"
    fi

    # 10. Fichier avec répétition de pattern (pour test compression)
    print_info "  → pattern.txt (5MB)"
    for i in {1..1000}; do
        echo "Pattern line $i: ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
    done > "$INPUT_DIR/pattern.txt"
    print_success "    pattern.txt généré (5MB)"

    # 11. Fichier shell script (pour test exécutable)
    print_info "  → test_script.sh"
    cat > "$INPUT_DIR/test_script.sh" << 'EOF'
#!/bin/bash
# Test script for cryptool
echo "This is a test script"
echo "It should be encrypted and decrypted"
for i in {1..10}; do
    echo "Line $i"
done
EOF
    chmod +x "$INPUT_DIR/test_script.sh"
    print_success "    test_script.sh généré"

    # 12. Fichier vide
    print_info "  → empty.txt (0 bytes)"
    touch "$INPUT_DIR/empty.txt"
    print_success "    empty.txt généré"

    # 13. Fichier avec une seule ligne
    print_info "  → single_line.txt"
    echo "This is a single line file for testing edge cases" > "$INPUT_DIR/single_line.txt"
    print_success "    single_line.txt généré"

    # 14. Fichier avec très petite taille (1 byte)
    print_info "  → one_byte.bin"
    echo -n "A" > "$INPUT_DIR/one_byte.bin"
    print_success "    one_byte.bin généré"

    # 15. Fichier avec taille exacte d'un chunk (1MB)
    print_info "  → exact_chunk.bin (1MB exact)"
    dd if=/dev/urandom of="$INPUT_DIR/exact_chunk.bin" bs=1024 count=1024 2>/dev/null
    print_success "    exact_chunk.bin généré"
}

# Vérification des fichiers générés
verify_files() {
    print_info "Vérification des fichiers générés..."
    echo ""

    local total_size=0
    local file_count=0

    for file in "$INPUT_DIR"/*; do
        if [ -f "$file" ]; then
            local size=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null)
            local size_human=$(numfmt --to=iec $size 2>/dev/null || echo "$size bytes")
            printf "  %-20s %10s\n" "$(basename "$file")" "$size_human"
            total_size=$((total_size + size))
            file_count=$((file_count + 1))
        fi
    done

    echo ""
    local total_human=$(numfmt --to=iec $total_size 2>/dev/null || echo "$total_size bytes")
    print_success "$file_count fichiers générés, taille totale: $total_human"
}

# Génération des checksums
generate_checksums() {
    print_info "Génération des checksums MD5..."

    local checksum_file="$INPUT_DIR/../checksums.md5"
    > "$checksum_file"

    for file in "$INPUT_DIR"/*; do
        if [ -f "$file" ]; then
            md5sum "$file" >> "$checksum_file"
        fi
    done

    print_success "Checksums sauvegardés dans: $checksum_file"
}

# Fonction principale
main() {
    print_header

    setup_directories
    generate_files
    verify_files
    generate_checksums

    echo ""
    print_success "✨ Génération des fichiers de test terminée ! ✨"
    echo ""
    echo -e "${CYAN}📂 Répertoire: ${INPUT_DIR}${NC}"
    echo ""
}

# Exécution
main "$@"
