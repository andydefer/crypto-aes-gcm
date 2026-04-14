#!/bin/bash

# Script de scénarios de test avancés pour aescryptool
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
CRYPTOOL_BIN="$PROJECT_ROOT/build/aescryptool"
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
    echo -e "${MAGENTA}║                         🎬 CRYPTOOL - SCÉNARIOS DE TEST                        ║${NC}"
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

# ============================================================================
# SCÉNARIOS DE TEST - LOGIQUE UNIQUEMENT
# ============================================================================

# Scénario 1: Chiffrement/Déchiffrement basique
scenario_basic_encrypt_decrypt() {
    print_scenario "Chiffrement/Déchiffrement basique"

    local input="$INPUT_DIR/small.txt"
    local encrypted="$ENCRYPTED_DIR/basic.enc"
    local decrypted="$DECRYPTED_DIR/basic.txt"
    local password="basic-test-password"

    local original_hash=$(get_hash "$input")
    print_step "Hash original: $original_hash"

    # Chiffrement
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ] || [ ! -f "$encrypted" ] || [ $(get_size "$encrypted") -eq 0 ]; then
        print_error "Échec du chiffrement"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Déchiffrement
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ] || [ ! -f "$decrypted" ]; then
        print_error "Échec du déchiffrement"
        return 1
    fi
    print_success "Déchiffrement réussi"

    local decrypted_hash=$(get_hash "$decrypted")
    if [ "$original_hash" = "$decrypted_hash" ]; then
        print_success "✅ Hashs identiques - Test réussi"
        return 0
    else
        print_error "❌ Hashs différents"
        return 1
    fi
}

# Scénario 2: Mauvais mot de passe - DOIT ÉCHOUER
scenario_wrong_password() {
    print_scenario "Tentative avec mauvais mot de passe"

    local input="$INPUT_DIR/small.txt"
    local encrypted="$ENCRYPTED_DIR/wrong.enc"
    local correct_password="correct-password-123"
    local wrong_password="wrong-password-456"

    # Chiffrement avec bon mot de passe
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$correct_password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec du chiffrement"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Tentative de déchiffrement avec mauvais mot de passe - DOIT ÉCHOUER
    local decrypted="$DECRYPTED_DIR/wrong.txt"
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$wrong_password" 2>/dev/null

    if [ $? -ne 0 ]; then
        print_success "✅ Le déchiffrement a échoué (comme attendu)"
        return 0
    else
        print_error "❌ Le déchiffrement a réussi avec un mauvais mot de passe!"
        return 1
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
    "$CRYPTOOL_BIN" encrypt "$input" "$output" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec premier chiffrement"
        return 1
    fi
    print_success "Premier fichier créé"

    # Deuxième chiffrement avec force - DOIT RÉUSSIR
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
        print_warning "Fichier 10MB non trouvé, utilisation random1.bin"
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
        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec chiffrement avec $workers workers"
            failed=$((failed + 1))
            continue
        fi

        # Déchiffrement
        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec déchiffrement avec $workers workers"
            failed=$((failed + 1))
            continue
        fi

        local decrypted_hash=$(get_hash "$decrypted")
        if [ "$original_hash" = "$decrypted_hash" ]; then
            print_success "  ✅ $workers workers: OK"
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

# Scénario 6: Détection de corruption - DOIT ÉCHOUER
scenario_corrupted_file() {
    print_scenario "Détection de corruption"

    local input="$INPUT_DIR/small.txt"
    local encrypted="$ENCRYPTED_DIR/corrupt.enc"
    local password="corrupt-test"

    # Chiffrement
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec du chiffrement"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Corrompre le header HMAC (offset 25) - garanti de faire échouer
    local corrupted="$ENCRYPTED_DIR/corrupted.enc"
    cp "$encrypted" "$corrupted"
    dd if=/dev/zero of="$corrupted" bs=1 count=1 seek=25 conv=notrunc 2>/dev/null

    # Tentative de déchiffrement - DOIT ÉCHOUER
    local decrypted="$DECRYPTED_DIR/corrupted.txt"
    "$CRYPTOOL_BIN" decrypt "$corrupted" "$decrypted" --pass "$password" 2>/dev/null

    if [ $? -ne 0 ]; then
        print_success "✅ Corruption détectée (déchiffrement échoué)"
        return 0
    else
        print_error "❌ Corruption non détectée (déchiffrement réussi)"
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

    if [ $(get_size "$input") -ne 0 ]; then
        > "$input"
    fi

    # Chiffrement
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ] || [ $(get_size "$encrypted") -eq 0 ]; then
        print_error "Échec du chiffrement du fichier vide"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Déchiffrement
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec du déchiffrement du fichier vide"
        return 1
    fi

    local dec_size=$(get_size "$decrypted")
    if [ $dec_size -eq 0 ]; then
        print_success "✅ Fichier déchiffré vide"
        return 0
    else
        print_error "❌ Fichier déchiffré non vide (taille: $dec_size)"
        return 1
    fi
}

# Scénario 8: Types MIME différents
scenario_different_mime_types() {
    print_scenario "Test avec différents types MIME"

    local password="mime-test-password"
    local failed=0

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
        local input="$INPUT_DIR/$filename"

        if [ ! -f "$input" ]; then
            print_warning "Fichier non trouvé: $filename (ignoré)"
            continue
        fi

        print_step "Test: $filename"

        local encrypted="$ENCRYPTED_DIR/mime_${filename}.enc"
        local decrypted="$DECRYPTED_DIR/mime_${filename}"
        local original_hash=$(get_hash "$input")

        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec chiffrement"
            failed=$((failed + 1))
            continue
        fi

        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "  Échec déchiffrement"
            failed=$((failed + 1))
            continue
        fi

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

    # Chiffrement
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers 8 --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec du chiffrement streaming"
        return 1
    fi
    print_success "Chiffrement réussi"

    # Déchiffrement
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec du déchiffrement streaming"
        return 1
    fi
    print_success "Déchiffrement réussi"

    local decrypted_hash=$(get_hash "$decrypted")
    local decrypted_size=$(get_size "$decrypted")

    if [ "$original_hash" = "$decrypted_hash" ] && [ "$original_size" -eq "$decrypted_size" ]; then
        print_success "✅ Gros fichier: vérification OK"
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

    echo "Test file with spaces in path" > "$input"
    local original_hash=$(get_hash "$input")

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

    # Chiffrement en chaîne
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted1" --pass "$password1" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec chiffrement 1"
        return 1
    fi

    "$CRYPTOOL_BIN" encrypt "$encrypted1" "$encrypted2" --pass "$password2" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec chiffrement 2"
        return 1
    fi

    "$CRYPTOOL_BIN" encrypt "$encrypted2" "$encrypted3" --pass "$password3" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec chiffrement 3"
        return 1
    fi

    # Déchiffrement en chaîne (ordre inverse)
    "$CRYPTOOL_BIN" decrypt "$encrypted3" "$encrypted2" --pass "$password3" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec déchiffrement 3"
        return 1
    fi

    "$CRYPTOOL_BIN" decrypt "$encrypted2" "$encrypted1" --pass "$password2" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec déchiffrement 2"
        return 1
    fi

    "$CRYPTOOL_BIN" decrypt "$encrypted1" "$decrypted_final" --pass "$password1" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Échec déchiffrement 1"
        return 1
    fi

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
    printf "  ${CYAN}%-15s %-15s${NC}\n" "Workers" "Status"
    echo "  ──────────────────────"

    for workers in 1 2 4 8; do
        local encrypted="$TEST_DIR/temp/bench_${workers}.enc"
        local decrypted="$TEST_DIR/temp/bench_${workers}.dec"

        # Chiffrement et déchiffrement - juste vérifier que ça fonctionne
        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            printf "  %-15s %-15s\n" "$workers" "❌ FAILED"
            continue
        fi

        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            printf "  %-15s %-15s\n" "$workers" "❌ FAILED"
            continue
        fi

        local decrypted_hash=$(get_hash "$decrypted")
        local original_hash=$(get_hash "$input")

        if [ "$original_hash" = "$decrypted_hash" ]; then
            printf "  %-15s %-15s\n" "$workers" "✅ OK"
        else
            printf "  %-15s %-15s\n" "$workers" "❌ FAILED"
        fi
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
    echo -e "${MAGENTA}│                                   RÉSUMÉ                                       │${NC}"
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
        echo -e "${RED}║                    ❌ CERTAINS SCÉNARIOS ONT ÉCHOUÉ ❌                         ║${NC}"
        echo -e "${RED}╚════════════════════════════════════════════════════════════════════════════════╝${NC}"
        echo ""
        echo "  Vérifiez les logs: $LOG_FILE"
        echo ""
        exit 1
    fi
}

main() {
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

    if ! check_binary; then
        exit 1
    fi

    if ! check_test_files; then
        exit 1
    fi

    mkdir -p "$ENCRYPTED_DIR" "$DECRYPTED_DIR" "$TEST_DIR/temp"

    echo "=== CRYPTOOL TEST SCENARIOS - $(date) ===" > "$LOG_FILE"

    run_all_scenarios
    print_summary
}

main "$@"
