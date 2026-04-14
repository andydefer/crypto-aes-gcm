#!/bin/bash

# Script de scénarios de test avancés pour cryptool
# Usage: ./test_scenarios.sh [--short] [--verbose]

# Couleurs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CRYPTOOL_BIN="$PROJECT_ROOT/build/cryptool"
TEST_DIR="$SCRIPT_DIR/test_data"
INPUT_DIR="$TEST_DIR/input"
ENCRYPTED_DIR="$TEST_DIR/encrypted"
DECRYPTED_DIR="$TEST_DIR/decrypted"
LOG_FILE="$TEST_DIR/test_scenarios.log"

# Options
SHORT_MODE="false"
VERBOSE="false"

# Compteurs
TOTAL_SCENARIOS=0
PASSED_SCENARIOS=0
FAILED_SCENARIOS=0

# Fonctions d'affichage
print_header() {
    echo ""
    echo -e "${MAGENTA}╔════════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${MAGENTA}║                         🎬 CRYPTOOL - SCÉNARIOS DE TEST                         ║${NC}"
    echo -e "${MAGENTA}╚════════════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_scenario() {
    echo ""
    echo -e "${CYAN}┌────────────────────────────────────────────────────────────────────────────────┐${NC}"
    echo -e "${CYAN}│ 📋 SCÉNARIO: $1${NC}"
    echo -e "${CYAN}└────────────────────────────────────────────────────────────────────────────────┘${NC}"
}

print_step() {
    echo -e "${BLUE}  ▶ $1${NC}"
}

print_success() {
    echo -e "${GREEN}  ✅ $1${NC}"
}

print_error() {
    echo -e "${RED}  ❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}  ⚠️  $1${NC}"
}

print_info() {
    echo -e "${WHITE}  ℹ️  $1${NC}"
}

print_debug() {
    if [ "$VERBOSE" = "true" ]; then
        echo -e "${CYAN}  🔍 $1${NC}"
    fi
}

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" >> "$LOG_FILE"
}

# Vérification du binaire
check_binary() {
    if [ ! -f "$CRYPTOOL_BIN" ]; then
        print_error "Binaire non trouvé: $CRYPTOOL_BIN"
        print_info "Lancez d'abord ./run_tests.sh pour compiler"
        return 1
    fi
    return 0
}

# Vérification des fichiers de test
check_test_files() {
    if [ ! -d "$INPUT_DIR" ] || [ -z "$(ls -A "$INPUT_DIR" 2>/dev/null)" ]; then
        print_error "Fichiers de test non trouvés"
        print_info "Lancez d'abord ./generate_test_files.sh"
        return 1
    fi
    return 0
}

# Calcul du hash MD5
get_hash() {
    local file="$1"
    if [ -f "$file" ]; then
        md5sum "$file" | cut -d' ' -f1
    else
        echo ""
    fi
}

# Taille du fichier
get_size() {
    local file="$1"
    if [ -f "$file" ]; then
        stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null
    else
        echo "0"
    fi
}

# Formatage taille
format_size() {
    numfmt --to=iec "$1" 2>/dev/null || echo "$1 bytes"
}

# Exécution d'un test avec vérification
run_test() {
    local test_name="$1"
    local expected_exit_code="$2"
    shift 2
    local cmd="$@"

    print_debug "Commande: $cmd"
    log "Test: $test_name - Commande: $cmd"

    eval "$cmd" > /dev/null 2>&1
    local exit_code=$?

    if [ $exit_code -eq $expected_exit_code ]; then
        print_debug "Exit code OK: $exit_code"
        return 0
    else
        print_error "Exit code attendu: $expected_exit_code, obtenu: $exit_code"
        log "ÉCHEC: $test_name - Exit code $exit_code (attendu $expected_exit_code)"
        return 1
    fi
}

# ============================================================================
# SCÉNARIOS DE TEST
# ============================================================================

# Scénario 1: Chiffrement/Déchiffrement basique
scenario_basic_encrypt_decrypt() {
    print_scenario "Chiffrement/Déchiffrement basique"

    local input="$INPUT_DIR/small.txt"
    local encrypted="$ENCRYPTED_DIR/basic.enc"
    local decrypted="$DECRYPTED_DIR/basic.txt"
    local password="basic-test-password"

    # Hash original
    local original_hash=$(get_hash "$input")
    print_step "Hash original: $original_hash"

    # Chiffrement
    print_step "Chiffrement en cours..."
    if ! run_test "basic_encrypt" 0 "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet; then
        print_error "Échec du chiffrement"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Vérification fichier chiffré
    if [ ! -f "$encrypted" ] || [ $(get_size "$encrypted") -eq 0 ]; then
        print_error "Fichier chiffré invalide"
        return 1
    fi
    print_info "Taille chiffrée: $(format_size $(get_size "$encrypted"))"

    # Déchiffrement
    print_step "Déchiffrement en cours..."
    if ! run_test "basic_decrypt" 0 "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet; then
        print_error "Échec du déchiffrement"
        return 1
    fi
    print_success "Déchiffrement réussi"

    # Vérification
    local decrypted_hash=$(get_hash "$decrypted")
    print_step "Hash déchiffré: $decrypted_hash"

    if [ "$original_hash" = "$decrypted_hash" ]; then
        print_success "✅ Hashs identiques - Test réussi"
        return 0
    else
        print_error "❌ Hashs différents"
        return 1
    fi
}

# Scénario 2: Mauvais mot de passe
scenario_wrong_password() {
    print_scenario "Tentative avec mauvais mot de passe"

    local input="$INPUT_DIR/small.txt"
    local encrypted="$ENCRYPTED_DIR/wrong.enc"
    local correct_password="correct-password-123"
    local wrong_password="wrong-password-456"

    # Chiffrement avec bon mot de passe
    print_step "Chiffrement avec bon mot de passe..."
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$correct_password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec du chiffrement"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Tentative de déchiffrement avec mauvais mot de passe
    print_step "Tentative de déchiffrement avec mauvais mot de passe..."
    local decrypted="$DECRYPTED_DIR/wrong.txt"
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$wrong_password" --force --quiet 2>/dev/null

    if [ $? -eq 0 ]; then
        print_error "❌ Le déchiffrement a réussi avec un mauvais mot de passe!"
        return 1
    else
        print_success "✅ Le déchiffrement a échoué (comme attendu)"
        return 0
    fi
}

# Scénario 3: Fichiers de différentes tailles
scenario_different_sizes() {
    print_scenario "Test avec différentes tailles de fichiers"

    local password="size-test-password"
    local failed=0

    local test_files=(
        "one_byte.bin:1 byte"
        "single_line.txt:single line"
        "small.txt:2KB"
        "config.json:1KB"
        "exact_chunk.bin:1MB exact"
        "random1.bin:1MB"
        "pattern.txt:5MB"
    )

    for file_info in "${test_files[@]}"; do
        local filename="${file_info%%:*}"
        local description="${file_info##*:}"
        local input="$INPUT_DIR/$filename"

        if [ ! -f "$input" ]; then
            print_warning "Fichier non trouvé: $filename (ignoré)"
            continue
        fi

        print_step "Test: $description ($filename)"

        local encrypted="$ENCRYPTED_DIR/size_${filename}.enc"
        local decrypted="$DECRYPTED_DIR/size_${filename}"
        local original_hash=$(get_hash "$input")

        # Chiffrement
        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec chiffrement"
            failed=$((failed + 1))
            continue
        fi

        # Déchiffrement
        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec déchiffrement"
            failed=$((failed + 1))
            continue
        fi

        # Vérification
        local decrypted_hash=$(get_hash "$decrypted")
        if [ "$original_hash" = "$decrypted_hash" ]; then
            print_success "  ✅ $description: OK"
        else
            print_error "  ❌ $description: Hash mismatch"
            failed=$((failed + 1))
        fi
    done

    if [ $failed -eq 0 ]; then
        print_success "✅ Tous les fichiers de test ont réussi"
        return 0
    else
        print_error "❌ $failed test(s) ont échoué"
        return 1
    fi
}

# Scénario 4: Force overwrite
scenario_force_overwrite() {
    print_scenario "Overwrite forcé (--force)"

    local input="$INPUT_DIR/small.txt"
    local output="$ENCRYPTED_DIR/force.enc"
    local password="force-test"

    # Premier chiffrement
    print_step "Premier chiffrement..."
    "$CRYPTOOL_BIN" encrypt "$input" "$output" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec premier chiffrement"
        return 1
    fi
    print_success "Premier fichier créé"

    # Deuxième chiffrement avec force
    print_step "Overwrite avec --force..."
    "$CRYPTOOL_BIN" encrypt "$input" "$output" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -eq 0 ]; then
        print_success "✅ Overwrite avec --force réussi"
        return 0
    else
        print_error "❌ Overwrite avec --force a échoué"
        return 1
    fi
}

# Scénario 5: Différents workers
scenario_different_workers() {
    print_scenario "Test avec différents nombres de workers"

    local input="$INPUT_DIR/random10.bin"
    local password="workers-test"

    if [ ! -f "$input" ]; then
        print_warning "Fichier 10MB non trouvé, utilisation d'un fichier plus petit"
        input="$INPUT_DIR/random1.bin"
    fi

    local original_hash=$(get_hash "$input")
    local workers_list=(1 2 4 8)
    local failed=0

    for workers in "${workers_list[@]}"; do
        print_step "Test avec $workers workers"

        local encrypted="$ENCRYPTED_DIR/workers_${workers}.enc"
        local decrypted="$DECRYPTED_DIR/workers_${workers}.txt"

        # Chiffrement
        local enc_start=$(date +%s%N)
        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null
        local enc_status=$?
        local enc_end=$(date +%s%N)
        local enc_time=$(( (enc_end - enc_start) / 1000000 ))

        if [ $enc_status -ne 0 ]; then
            print_error "  Échec chiffrement avec $workers workers"
            failed=$((failed + 1))
            continue
        fi

        # Déchiffrement
        local dec_start=$(date +%s%N)
        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
        local dec_status=$?
        local dec_end=$(date +%s%N)
        local dec_time=$(( (dec_end - dec_start) / 1000000 ))

        if [ $dec_status -ne 0 ]; then
            print_error "  Échec déchiffrement avec $workers workers"
            failed=$((failed + 1))
            continue
        fi

        # Vérification
        local decrypted_hash=$(get_hash "$decrypted")
        if [ "$original_hash" = "$decrypted_hash" ]; then
            print_success "  ✅ $workers workers: chiffrement ${enc_time}ms, déchiffrement ${dec_time}ms"
        else
            print_error "  ❌ $workers workers: Hash mismatch"
            failed=$((failed + 1))
        fi
    done

    if [ $failed -eq 0 ]; then
        return 0
    else
        return 1
    fi
}

# Scénario 6: Fichier corrompu
scenario_corrupted_file() {
    print_scenario "Détection de corruption"

    local input="$INPUT_DIR/small.txt"
    local encrypted="$ENCRYPTED_DIR/corrupt.enc"
    local password="corrupt-test"

    # Chiffrement
    print_step "Chiffrement du fichier original..."
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec du chiffrement"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Types de corruption
    local corruptions=(
        "header_magic:1:0xFF"
        "header_version:5:0x00"
        "header_hmac:25:0x00"
        "nonce:57:0xFF"
        "ciphertext:100:0x00"
        "chunk_length:69:0xFFFF"
    )

    local failed=0

    for corruption in "${corruptions[@]}"; do
        local name=$(echo "$corruption" | cut -d':' -f1)
        local offset=$(echo "$corruption" | cut -d':' -f2)
        local value=$(echo "$corruption" | cut -d':' -f3)

        print_step "Corruption: $name (offset $offset)"

        # Copie et corruption
        local corrupted="$ENCRYPTED_DIR/corrupted_${name}.enc"
        cp "$encrypted" "$corrupted"

        # Appliquer la corruption
        printf "$value" | dd of="$corrupted" bs=1 seek=$offset count=1 conv=notrunc 2>/dev/null

        # Tentative de déchiffrement
        local decrypted="$DECRYPTED_DIR/corrupted_${name}.txt"
        "$CRYPTOOL_BIN" decrypt "$corrupted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null

        if [ $? -eq 0 ]; then
            print_error "  ❌ Corruption $name non détectée (déchiffrement réussi)"
            failed=$((failed + 1))
        else
            print_success "  ✅ Corruption $name détectée"
        fi
    done

    if [ $failed -eq 0 ]; then
        print_success "✅ Toutes les corruptions ont été détectées"
        return 0
    else
        print_error "❌ $failed corruption(s) non détectée(s)"
        return 1
    fi
}

# Scénario 7: Fichier vide
scenario_empty_file() {
    print_scenario "Fichier vide"

    local input="$INPUT_DIR/empty.txt"
    local encrypted="$ENCRYPTED_DIR/empty.enc"
    local decrypted="$DECRYPTED_DIR/empty.txt"
    local password="empty-test"

    # Vérifier que le fichier est bien vide
    if [ $(get_size "$input") -ne 0 ]; then
        print_warning "Le fichier n'est pas vide, recréation..."
        > "$input"
    fi

    print_step "Taille originale: $(get_size "$input") bytes"

    # Chiffrement
    print_step "Chiffrement du fichier vide..."
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -ne 0 ]; then
        print_error "Échec du chiffrement du fichier vide"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Vérification taille fichier chiffré (devrait contenir header + HMAC + nonce)
    local enc_size=$(get_size "$encrypted")
    print_info "Taille fichier chiffré: $(format_size $enc_size)"

    if [ $enc_size -eq 0 ]; then
        print_error "Fichier chiffré vide (invalide)"
        return 1
    fi

    # Déchiffrement
    print_step "Déchiffrement..."
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -ne 0 ]; then
        print_error "Échec du déchiffrement du fichier vide"
        return 1
    fi

    # Vérification
    local dec_size=$(get_size "$decrypted")
    if [ $dec_size -eq 0 ]; then
        print_success "✅ Fichier déchiffré vide (taille $dec_size bytes)"
        return 0
    else
        print_error "❌ Fichier déchiffré non vide (taille $dec_size bytes)"
        return 1
    fi
}

# Scénario 8: Types MIME différents
scenario_different_mime_types() {
    print_scenario "Test avec différents types MIME"

    local password="mime-test-password"
    local failed=0

    # Déclaration des types de fichiers
    declare -A test_files=(
        ["small.txt"]="text/plain"
        ["config.json"]="application/json"
        ["data.xml"]="application/xml"
        ["data.csv"]="text/csv"
        ["test_script.sh"]="text/x-shellscript"
        ["random1.bin"]="application/octet-stream"
        ["special.txt"]="text/plain; charset=utf-8"
    )

    for filename in "${!test_files[@]}"; do
        local mime="${test_files[$filename]}"
        local input="$INPUT_DIR/$filename"

        if [ ! -f "$input" ]; then
            print_warning "Fichier non trouvé: $filename (ignoré)"
            continue
        fi

        print_step "Test: $filename ($mime)"

        local encrypted="$ENCRYPTED_DIR/mime_${filename}.enc"
        local decrypted="$DECRYPTED_DIR/mime_${filename}"
        local original_hash=$(get_hash "$input")

        # Chiffrement
        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec chiffrement"
            failed=$((failed + 1))
            continue
        fi

        # Déchiffrement
        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec déchiffrement"
            failed=$((failed + 1))
            continue
        fi

        # Vérification
        local decrypted_hash=$(get_hash "$decrypted")
        if [ "$original_hash" = "$decrypted_hash" ]; then
            print_success "  ✅ $filename: OK"
        else
            print_error "  ❌ $filename: Hash mismatch"
            failed=$((failed + 1))
        fi
    done

    if [ $failed -eq 0 ]; then
        print_success "✅ Tous les types MIME ont réussi"
        return 0
    else
        print_error "❌ $failed type(s) MIME ont échoué"
        return 1
    fi
}

# Scénario 9: Gros fichier (streaming)
scenario_large_file_streaming() {
    if [ "$SHORT_MODE" = "true" ]; then
        print_warning "Scénario gros fichier ignoré (mode court)"
        return 0
    fi

    print_scenario "Gros fichier avec streaming (50MB+)"

    local input="$INPUT_DIR/random50.bin"

    if [ ! -f "$input" ]; then
        print_warning "Fichier 50MB non trouvé, génération..."
        dd if=/dev/urandom of="$input" bs=1M count=50 2>/dev/null
    fi

    local encrypted="$ENCRYPTED_DIR/large_stream.enc"
    local decrypted="$DECRYPTED_DIR/large_stream.bin"
    local password="large-stream-password"

    local original_hash=$(get_hash "$input")
    local original_size=$(get_size "$input")

    print_info "Taille du fichier: $(format_size $original_size)"

    # Chiffrement avec monitoring mémoire
    print_step "Chiffrement streaming en cours..."

    local mem_before=$(free -m | awk 'NR==2{print $3}')
    local time_start=$(date +%s)

    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers 8 --force --quiet 2>/dev/null

    local time_end=$(date +%s)
    local mem_after=$(free -m | awk 'NR==2{print $3}')
    local time_total=$((time_end - time_start))

    if [ $? -ne 0 ]; then
        print_error "Échec du chiffrement streaming"
        return 1
    fi

    print_success "Chiffrement terminé en ${time_total}s"
    print_info "Mémoire utilisée: ~$((mem_after - mem_before)) MB"

    # Déchiffrement streaming
    print_step "Déchiffrement streaming en cours..."

    local time_start=$(date +%s)

    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null

    local time_end=$(date +%s)
    local time_total=$((time_end - time_start))

    if [ $? -ne 0 ]; then
        print_error "Échec du déchiffrement streaming"
        return 1
    fi

    print_success "Déchiffrement terminé en ${time_total}s"

    # Vérification
    local decrypted_hash=$(get_hash "$decrypted")
    local decrypted_size=$(get_size "$decrypted")

    if [ "$original_hash" = "$decrypted_hash" ] && [ "$original_size" -eq "$decrypted_size" ]; then
        print_success "✅ Gros fichier: vérification OK"
        if [ $time_total -gt 0 ]; then
            print_info "Taux: $((original_size / 1024 / 1024 / time_total)) MB/s"
        else
            print_info "Taux: N/A (temps trop court)"
        fi
        return 0
    else
        print_error "❌ Gros fichier: vérification échouée"
        return 1
    fi
}

# Scénario 10: Chemins avec espaces
scenario_paths_with_spaces() {
    print_scenario "Chemins contenant des espaces"

    local input_dir="$TEST_DIR/temp/input with spaces"
    local output_dir="$TEST_DIR/temp/output with spaces"

    mkdir -p "$input_dir"
    mkdir -p "$output_dir"

    local input="$input_dir/test file.txt"
    local encrypted="$output_dir/encrypted file.enc"
    local decrypted="$output_dir/decrypted file.txt"
    local password="spaces-test"

    # Création fichier test
    echo "Test file with spaces in path" > "$input"
    local original_hash=$(get_hash "$input")

    print_step "Fichier avec espaces dans le chemin"

    # Chiffrement
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec chiffrement avec espaces"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Déchiffrement
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec déchiffrement avec espaces"
        return 1
    fi

    # Vérification
    local decrypted_hash=$(get_hash "$decrypted")
    if [ "$original_hash" = "$decrypted_hash" ]; then
        print_success "✅ Chemins avec espaces: OK"
        return 0
    else
        print_error "❌ Chemins avec espaces: Hash mismatch"
        return 1
    fi
}

# Scénario 11: Chiffrement/déchiffrement en chaîne
scenario_chain_encryption() {
    print_scenario "Chiffrement/déchiffrement en chaîne (multiples opérations)"

    local input="$INPUT_DIR/small.txt"
    local password1="first-password"
    local password2="second-password"
    local password3="third-password"

    local encrypted1="$ENCRYPTED_DIR/chain1.enc"
    local encrypted2="$ENCRYPTED_DIR/chain2.enc"
    local encrypted3="$ENCRYPTED_DIR/chain3.enc"
    local decrypted_final="$DECRYPTED_DIR/chain_final.txt"

    local original_hash=$(get_hash "$input")

    print_step "Opération 1: Chiffrement avec password1"
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted1" --pass "$password1" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec chiffrement 1"
        return 1
    fi

    print_step "Opération 2: Chiffrement avec password2"
    "$CRYPTOOL_BIN" encrypt "$encrypted1" "$encrypted2" --pass "$password2" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec chiffrement 2"
        return 1
    fi

    print_step "Opération 3: Chiffrement avec password3"
    "$CRYPTOOL_BIN" encrypt "$encrypted2" "$encrypted3" --pass "$password3" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec chiffrement 3"
        return 1
    fi

    print_step "Opération 4: Déchiffrement avec password3"
    "$CRYPTOOL_BIN" decrypt "$encrypted3" "$encrypted2" --pass "$password3" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec déchiffrement 3"
        return 1
    fi

    print_step "Opération 5: Déchiffrement avec password2"
    "$CRYPTOOL_BIN" decrypt "$encrypted2" "$encrypted1" --pass "$password2" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec déchiffrement 2"
        return 1
    fi

    print_step "Opération 6: Déchiffrement avec password1"
    "$CRYPTOOL_BIN" decrypt "$encrypted1" "$decrypted_final" --pass "$password1" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec déchiffrement 1"
        return 1
    fi

    # Vérification
    local final_hash=$(get_hash "$decrypted_final")
    if [ "$original_hash" = "$final_hash" ]; then
        print_success "✅ Chiffrement en chaîne: OK"
        return 0
    else
        print_error "❌ Chiffrement en chaîne: Hash mismatch"
        return 1
    fi
}

# Scénario 12: Performance benchmark
scenario_performance_benchmark() {
    if [ "$SHORT_MODE" = "true" ]; then
        print_warning "Benchmark ignoré (mode court)"
        return 0
    fi

    print_scenario "Benchmark de performance"

    local input="$INPUT_DIR/random10.bin"

    if [ ! -f "$input" ]; then
        print_warning "Fichier 10MB non trouvé, utilisation random1.bin"
        input="$INPUT_DIR/random1.bin"
    fi

    local password="benchmark-password"
    local size_mb=$(($(get_size "$input") / 1024 / 1024))

    echo ""
    printf "  ${CYAN}%-15s %-15s %-15s %-15s${NC}\n" "Workers" "Encrypt(s)" "Decrypt(s)" "Speed(MB/s)"
    echo "  ───────────────────────────────────────────────────────────────"

    for workers in 1 2 4 8; do
        local encrypted="$TEST_DIR/temp/bench_${workers}.enc"
        local decrypted="$TEST_DIR/temp/bench_${workers}.dec"

        # Chiffrement
        local enc_time=$( { time "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null; } 2>&1 | grep real | awk '{print $2}' | sed 's/[^0-9.]//g')

        if [ -z "$enc_time" ]; then
            enc_time=$( { time "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null; } 2>&1 | grep real | awk '{print $2}' | sed 's/m/*60+/g' | sed 's/s//' | bc 2>/dev/null || echo "0")
        fi

        if [ -z "$enc_time" ]; then
            enc_time="0"
        fi

        # Déchiffrement
        local dec_time=$( { time "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null; } 2>&1 | grep real | awk '{print $2}' | sed 's/[^0-9.]//g')

        if [ -z "$dec_time" ]; then
            dec_time=$( { time "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null; } 2>&1 | grep real | awk '{print $2}' | sed 's/m/*60+/g' | sed 's/s//' | bc 2>/dev/null || echo "0")
        fi

        if [ -z "$dec_time" ]; then
            dec_time="0"
        fi

        # Calcul vitesse
        local enc_speed="0"
        if [ "$enc_time" != "0" ] && [ "$enc_time" != "" ]; then
            enc_speed=$(echo "scale=2; $size_mb / $enc_time" | bc 2>/dev/null || echo "0")
        fi

        printf "  %-15s %-15s %-15s %-15s\n" "$workers" "$enc_time" "$dec_time" "$enc_speed"
    done

    echo ""
    print_success "Benchmark terminé"
    return 0
}

# ============================================================================
# EXÉCUTION PRINCIPALE
# ============================================================================

run_all_scenarios() {
    local scenarios=(
        "Basique:scenario_basic_encrypt_decrypt"
        "Mauvais mot de passe:scenario_wrong_password"
        "Différentes tailles:scenario_different_sizes"
        "Force overwrite:scenario_force_overwrite"
        "Différents workers:scenario_different_workers"
        "Détection corruption:scenario_corrupted_file"
        "Fichier vide:scenario_empty_file"
        "Types MIME:scenario_different_mime_types"
        "Chemins avec espaces:scenario_paths_with_spaces"
        "Chiffrement en chaîne:scenario_chain_encryption"
        "Gros fichier streaming:scenario_large_file_streaming"
        "Benchmark performance:scenario_performance_benchmark"
    )

    TOTAL_SCENARIOS=${#scenarios[@]}

    for scenario in "${scenarios[@]}"; do
        local name="${scenario%%:*}"
        local func="${scenario##*:}"

        echo ""
        $func
        if [ $? -eq 0 ]; then
            PASSED_SCENARIOS=$((PASSED_SCENARIOS + 1))
            log "SUCCÈS: $name"
        else
            FAILED_SCENARIOS=$((FAILED_SCENARIOS + 1))
            log "ÉCHEC: $name"
        fi
    done
}

print_summary() {
    echo ""
    echo -e "${MAGENTA}╔════════════════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${MAGENTA}│                                   RÉSUMÉ                                      │${NC}"
    echo -e "${MAGENTA}╚════════════════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "  📊 Scénarios total:   ${TOTAL_SCENARIOS}"
    echo -e "  ${GREEN}✅ Réussis:           ${PASSED_SCENARIOS}${NC}"
    echo -e "  ${RED}❌ Échoués:           ${FAILED_SCENARIOS}${NC}"

    if [ $FAILED_SCENARIOS -eq 0 ]; then
        echo ""
        echo -e "${GREEN}╔════════════════════════════════════════════════════════════════════════════════╗${NC}"
        echo -e "${GREEN}║                    🎉 TOUS LES SCÉNARIOS ONT RÉUSSI ! 🎉                       ║${NC}"
        echo -e "${GREEN}╚════════════════════════════════════════════════════════════════════════════════╝${NC}"
        echo ""
        exit 0
    else
        echo ""
        echo -e "${RED}╔════════════════════════════════════════════════════════════════════════════════╗${NC}"
        echo -e "${RED}║                    ❌ CERTAINS SCÉNARIOS ONT ÉCHOUÉ ❌                          ║${NC}"
        echo -e "${RED}╚════════════════════════════════════════════════════════════════════════════════╝${NC}"
        echo ""
        echo "  Vérifiez les logs: $LOG_FILE"
        echo ""
        exit 1
    fi
}

# Fonction principale
main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --short|-s)
                SHORT_MODE="true"
                shift
                ;;
            --verbose|-v)
                VERBOSE="true"
                shift
                ;;
            *)
                shift
                ;;
        esac
    done

    print_header

    # Vérifications préalables
    if ! check_binary; then
        exit 1
    fi

    if ! check_test_files; then
        exit 1
    fi

    # Création des répertoires
    mkdir -p "$ENCRYPTED_DIR" "$DECRYPTED_DIR" "$TEST_DIR/temp"

    # Initialisation log
    echo "=== CRYPTOOL TEST SCENARIOS - $(date) ===" > "$LOG_FILE"

    # Exécution
    run_all_scenarios

    # Résumé
    print_summary
}

# Exécution
main "$@"
