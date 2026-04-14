#!/bin/bash

# Couleurs pour les outputs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_DIR="$PROJECT_ROOT/build"
CRYPTOOL_BIN="$BUILD_DIR/aescryptool"
TEST_DIR="$SCRIPT_DIR/test_data"
RESULT_DIR="$SCRIPT_DIR/results"

# Compteurs
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Options
SHORT_MODE="false"

# Fonction d'affichage
print_header() {
    echo ""
    echo -e "${MAGENTA}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${MAGENTA}║                    🔐 CRYPTOOL - TESTS RÉALISTES             ║${NC}"
    echo -e "${MAGENTA}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${CYAN}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_test_header() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}📋 TEST: $1${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
}

# Compilation du binaire
build_aescryptool() {
    print_info "Compilation de aescryptool..."

    mkdir -p "$BUILD_DIR"
    cd "$PROJECT_ROOT"

    go build -o "$CRYPTOOL_BIN" ./cmd/aescryptool

    if [ $? -eq 0 ] && [ -f "$CRYPTOOL_BIN" ]; then
        print_success "Compilation réussie: $CRYPTOOL_BIN"
        return 0
    else
        print_error "Échec de la compilation"
        return 1
    fi
}

# Vérification du binaire
check_binary() {
    if [ ! -f "$CRYPTOOL_BIN" ]; then
        print_error "Binaire non trouvé: $CRYPTOOL_BIN"
        print_info "Lancement de la compilation..."
        build_aescryptool
        if [ $? -ne 0 ]; then
            return 1
        fi
    fi

    print_info "Version de aescryptool:"
    "$CRYPTOOL_BIN" version 2>/dev/null

    return 0
}

# Création des répertoires
setup_directories() {
    print_info "Création des répertoires de test..."
    mkdir -p "$TEST_DIR"
    mkdir -p "$RESULT_DIR"
    mkdir -p "$TEST_DIR/input"
    mkdir -p "$TEST_DIR/encrypted"
    mkdir -p "$TEST_DIR/decrypted"
    mkdir -p "$TEST_DIR/temp"
    print_success "Répertoires créés"
}

# Nettoyage avant test
cleanup_before_test() {
    rm -rf "$TEST_DIR"/*
    mkdir -p "$TEST_DIR/input"
    mkdir -p "$TEST_DIR/encrypted"
    mkdir -p "$TEST_DIR/decrypted"
    mkdir -p "$TEST_DIR/temp"
}

# Génération des fichiers de test
generate_test_files() {
    print_info "Génération des fichiers de test..."

    if [ -f "$SCRIPT_DIR/generate_test_files.sh" ]; then
        chmod +x "$SCRIPT_DIR/generate_test_files.sh"
        if [ "$SHORT_MODE" = "true" ]; then
            "$SCRIPT_DIR/generate_test_files.sh" --short
        else
            "$SCRIPT_DIR/generate_test_files.sh"
        fi
    else
        print_error "Script generate_test_files.sh non trouvé"
        return 1
    fi
}

# Test 1: encrypt + decrypt - Vérifie que le cycle complet fonctionne
test_simple_encrypt_decrypt() {
    print_test_header "Encrypt/Decrypt simple"

    local input="$TEST_DIR/input/small.txt"
    local encrypted="$TEST_DIR/encrypted/small.enc"
    local decrypted="$TEST_DIR/decrypted/small.txt"
    local password="test-password-123"

    if [ ! -f "$input" ]; then
        print_error "Fichier test non trouvé: $input"
        return 1
    fi

    local original_hash=$(md5sum "$input" | cut -d' ' -f1)

    # Encryption
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ] || [ ! -f "$encrypted" ]; then
        print_error "Encryption échouée"
        return 1
    fi

    # Décryption
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ] || [ ! -f "$decrypted" ]; then
        print_error "Décryption échouée"
        return 1
    fi

    local decrypted_hash=$(md5sum "$decrypted" | cut -d' ' -f1)

    if [ "$original_hash" = "$decrypted_hash" ]; then
        print_success "Hashs identiques"
        return 0
    else
        print_error "Hashs différents"
        return 1
    fi
}

# Test 2: Workers parallèles - Vérifie que différents workers fonctionnent
test_workers_parallel() {
    print_test_header "Encryption avec workers parallèles"

    local input="$TEST_DIR/input/random10.bin"
    local password="workers-test"

    if [ ! -f "$input" ]; then
        print_warning "Fichier random10.bin non trouvé, utilisation random1.bin"
        input="$TEST_DIR/input/random1.bin"
    fi

    for workers in 1 2 4 8; do
        local encrypted="$TEST_DIR/encrypted/workers_${workers}.enc"
        local decrypted="$TEST_DIR/decrypted/workers_${workers}.txt"

        local original_hash=$(md5sum "$input" | cut -d' ' -f1)

        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "Encryption échouée avec $workers workers"
            return 1
        fi

        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            print_error "Décryption échouée avec $workers workers"
            return 1
        fi

        local decrypted_hash=$(md5sum "$decrypted" | cut -d' ' -f1)

        if [ "$original_hash" = "$decrypted_hash" ]; then
            print_success "Workers $workers: OK"
        else
            print_error "Workers $workers: Hash mismatch"
            return 1
        fi
    done

    return 0
}

# Test 3: Mauvais mot de passe - Vérifie que la décryption ÉCHOUE (code retour non nul)
test_wrong_password() {
    print_test_header "Mauvais mot de passe (doit échouer)"

    local input="$TEST_DIR/input/small.txt"
    local encrypted="$TEST_DIR/encrypted/wrong_pass.enc"
    local correct_password="correct-password"
    local wrong_password="wrong-password"

    # Encryption avec bon mot de passe
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$correct_password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Encryption échouée"
        return 1
    fi

    # Tentative de décryption avec mauvais mot de passe - DOIT ÉCHOUER
    local decrypted="$TEST_DIR/decrypted/wrong_pass.txt"
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$wrong_password" 2>/dev/null

    if [ $? -ne 0 ]; then
        print_success "Décryption a échoué (comme attendu)"
        return 0
    else
        print_error "ERREUR: La décryption a réussi avec un mauvais mot de passe!"
        return 1
    fi
}

# Test 4: Fichier vide - Vérifie que le cycle fonctionne avec un fichier vide
test_empty_file() {
    print_test_header "Fichier vide"

    local input="$TEST_DIR/input/empty.txt"
    local encrypted="$TEST_DIR/encrypted/empty.enc"
    local decrypted="$TEST_DIR/decrypted/empty.txt"
    local password="empty-test"

    # Encryption
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Encryption du fichier vide échouée"
        return 1
    fi

    # Décryption
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Décryption du fichier vide échouée"
        return 1
    fi

    local size=$(stat -c%s "$decrypted" 2>/dev/null || echo "0")
    if [ "$size" -eq 0 ]; then
        print_success "Fichier décrypté est vide"
        return 0
    else
        print_error "Fichier décrypté n'est pas vide (taille: $size)"
        return 1
    fi
}

# Test 5: Force overwrite - Vérifie que --force permet d'écraser
test_force_overwrite() {
    print_test_header "Force overwrite"

    local input="$TEST_DIR/input/small.txt"
    local output="$TEST_DIR/encrypted/overwrite.enc"
    local password="force-test"

    # Première encryption
    "$CRYPTOOL_BIN" encrypt "$input" "$output" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Première encryption échouée"
        return 1
    fi

    # Deuxième encryption avec --force - DOIT RÉUSSIR
    "$CRYPTOOL_BIN" encrypt "$input" "$output" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -eq 0 ]; then
        print_success "Force overwrite fonctionne"
        return 0
    else
        print_error "Force overwrite a échoué"
        return 1
    fi
}

# Test 6: Tous les types de fichiers - Vérifie l'intégrité pour chaque fichier
test_all_file_types() {
    print_test_header "Tous les types de fichiers"

    local password="all-types-test"
    local failed=0

    for input in "$TEST_DIR/input/"*; do
        if [ -f "$input" ]; then
            local filename=$(basename "$input")
            local encrypted="$TEST_DIR/encrypted/${filename}.enc"
            local decrypted="$TEST_DIR/decrypted/${filename}"

            local original_hash=$(md5sum "$input" | cut -d' ' -f1)

            "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
            if [ $? -ne 0 ]; then
                print_error "  Encryption échouée pour $filename"
                failed=$((failed + 1))
                continue
            fi

            "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
            if [ $? -ne 0 ]; then
                print_error "  Décryption échouée pour $filename"
                failed=$((failed + 1))
                continue
            fi

            local decrypted_hash=$(md5sum "$decrypted" | cut -d' ' -f1)

            if [ "$original_hash" = "$decrypted_hash" ]; then
                print_success "  $filename: OK"
            else
                print_error "  $filename: Hash mismatch"
                failed=$((failed + 1))
            fi
        fi
    done

    if [ $failed -eq 0 ]; then
        print_success "Tous les fichiers sont OK"
        return 0
    else
        print_error "$failed fichier(s) ont échoué"
        return 1
    fi
}

# Test 7: Détection de corruption - Vérifie que la corruption est DÉTECTÉE (échec)
test_corruption_detection() {
    print_test_header "Détection de corruption"

    local input="$TEST_DIR/input/small.txt"
    local encrypted="$TEST_DIR/encrypted/corrupt_test.enc"
    local password="corruption-test"

    # Encryption
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Encryption échouée"
        return 1
    fi

    # Copier et corrompre le HMAC (offset 25) - garanti de faire échouer
    local corrupted="$TEST_DIR/encrypted/corrupted.enc"
    cp "$encrypted" "$corrupted"
    dd if=/dev/zero of="$corrupted" bs=1 count=1 seek=25 conv=notrunc 2>/dev/null

    # Tentative de décryption - DOIT ÉCHOUER
    local decrypted="$TEST_DIR/decrypted/corrupted.txt"
    "$CRYPTOOL_BIN" decrypt "$corrupted" "$decrypted" --pass "$password" 2>/dev/null

    if [ $? -ne 0 ]; then
        print_success "Corruption détectée (décryption échouée)"
        return 0
    else
        print_error "ERREUR: Décryption a réussi malgré la corruption!"
        return 1
    fi
}

# Test 8: Gros fichier - Vérifie le streaming avec un gros fichier
test_large_file() {
    if [ "$SHORT_MODE" = "true" ]; then
        print_warning "Test gros fichier ignoré (mode court)"
        return 0
    fi

    print_test_header "Gros fichier (50MB)"

    local input="$TEST_DIR/input/random50.bin"

    if [ ! -f "$input" ]; then
        print_warning "Fichier 50MB non trouvé, test ignoré"
        return 0
    fi

    local encrypted="$TEST_DIR/encrypted/large.enc"
    local decrypted="$TEST_DIR/decrypted/large.bin"
    local password="large-file-password"

    local original_hash=$(md5sum "$input" | cut -d' ' -f1)

    # Encryption
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers 8 --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Encryption du gros fichier échouée"
        return 1
    fi

    # Décryption
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    if [ $? -ne 0 ]; then
        print_error "Décryption du gros fichier échouée"
        return 1
    fi

    local decrypted_hash=$(md5sum "$decrypted" | cut -d' ' -f1)

    if [ "$original_hash" = "$decrypted_hash" ]; then
        print_success "Gros fichier: vérification OK"
        return 0
    else
        print_error "Gros fichier: hash mismatch"
        return 1
    fi
}

# Test 9: Performance - Benchmark (optionnel, ne compte pas dans les échecs)
test_performance() {
    if [ "$SHORT_MODE" = "true" ]; then
        print_warning "Test performance ignoré (mode court)"
        return 0
    fi

    print_test_header "Performance"

    local input="$TEST_DIR/input/random10.bin"
    local password="perf-test"

    if [ ! -f "$input" ]; then
        print_warning "Fichier random10.bin non trouvé, test ignoré"
        return 0
    fi

    echo ""
    printf "%-15s %-15s\n" "Workers" "Status"
    echo "------------------------"

    for workers in 1 2 4 8; do
        local encrypted="$TEST_DIR/temp/perf_${workers}.enc"
        local decrypted="$TEST_DIR/temp/perf_${workers}.dec"

        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null
        if [ $? -ne 0 ]; then
            printf "%-15s %-15s\n" "$workers" "❌ FAILED"
        else
            "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
            if [ $? -ne 0 ]; then
                printf "%-15s %-15s\n" "$workers" "❌ FAILED"
            else
                printf "%-15s %-15s\n" "$workers" "✅ OK"
            fi
        fi
    done

    echo ""
    print_success "Test performance terminé"
    return 0
}

# Fonction principale
main() {
    print_header

    if [ "$1" = "--short" ] || [ "$1" = "-s" ]; then
        SHORT_MODE="true"
        print_info "Mode court activé (tests lourds ignorés)"
    fi

    check_binary
    if [ $? -ne 0 ]; then
        print_error "Impossible de continuer sans le binaire"
        exit 1
    fi

    setup_directories
    cleanup_before_test
    generate_test_files

    TOTAL_TESTS=0
    PASSED_TESTS=0
    FAILED_TESTS=0

    run_test() {
        local test_name=$1
        local test_func=$2
        TOTAL_TESTS=$((TOTAL_TESTS + 1))

        $test_func
        if [ $? -eq 0 ]; then
            PASSED_TESTS=$((PASSED_TESTS + 1))
            return 0
        else
            FAILED_TESTS=$((FAILED_TESTS + 1))
            return 1
        fi
    }

    run_test "Test simple encrypt/decrypt" test_simple_encrypt_decrypt
    run_test "Test workers parallèles" test_workers_parallel
    run_test "Test mauvais mot de passe" test_wrong_password
    run_test "Test fichier vide" test_empty_file
    run_test "Test force overwrite" test_force_overwrite
    run_test "Test tous types de fichiers" test_all_file_types
    run_test "Test détection corruption" test_corruption_detection
    run_test "Test gros fichier" test_large_file
    run_test "Test performance" test_performance

    echo ""
    echo -e "${MAGENTA}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${MAGENTA}                          RÉSULTATS                             ${NC}"
    echo -e "${MAGENTA}════════════════════════════════════════════════════════════════${NC}"
    echo -e "Total:  ${TOTAL_TESTS}"
    echo -e "Passed: ${GREEN}${PASSED_TESTS}${NC}"
    echo -e "Failed: ${RED}${FAILED_TESTS}${NC}"

    if [ $FAILED_TESTS -eq 0 ]; then
        echo ""
        echo -e "${GREEN}🎉 TOUS LES TESTS ONT RÉUSSI ! 🎉${NC}"
        exit 0
    else
        echo ""
        echo -e "${RED}❌ $FAILED_TESTS TEST(S) ONT ÉCHOUE(S) ❌${NC}"
        exit 1
    fi
}

main "$@"
