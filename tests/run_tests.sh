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
CRYPTOOL_BIN="$BUILD_DIR/cryptool"
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
build_cryptool() {
    print_info "Compilation de cryptool..."

    mkdir -p "$BUILD_DIR"
    cd "$PROJECT_ROOT"

    go build -o "$CRYPTOOL_BIN" ./cmd/cryptool

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
        build_cryptool
        if [ $? -ne 0 ]; then
            return 1
        fi
    fi

    # Vérifier la version
    print_info "Version de cryptool:"
    "$CRYPTOOL_BIN" version

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

# Génération des fichiers de test - Utilise le script externe
generate_test_files() {
    print_info "Génération des fichiers de test..."

    # Vérifier si le script de génération existe
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

# Test simple: encrypt + decrypt
test_simple_encrypt_decrypt() {
    print_test_header "Encrypt/Decrypt simple"

    local input="$TEST_DIR/input/small.txt"
    local encrypted="$TEST_DIR/encrypted/small.enc"
    local decrypted="$TEST_DIR/decrypted/small.txt"
    local password="test-password-123"

    # Vérifier que le fichier existe
    if [ ! -f "$input" ]; then
        print_error "Fichier test non trouvé: $input"
        return 1
    fi

    # Hash original
    local original_hash=$(md5sum "$input" | cut -d' ' -f1)

    # Encryption
    print_info "Encryption: $input -> $encrypted"
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -ne 0 ] || [ ! -f "$encrypted" ]; then
        print_error "Encryption échouée"
        return 1
    fi
    print_success "Encryption réussie"

    # Décryption
    print_info "Décryption: $encrypted -> $decrypted"
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -ne 0 ] || [ ! -f "$decrypted" ]; then
        print_error "Décryption échouée"
        return 1
    fi
    print_success "Décryption réussie"

    # Vérification
    local decrypted_hash=$(md5sum "$decrypted" | cut -d' ' -f1)

    if [ "$original_hash" = "$decrypted_hash" ]; then
        print_success "Hashs identiques: $original_hash"
        return 0
    else
        print_error "Hashs différents: original=$original_hash decrypted=$decrypted_hash"
        return 1
    fi
}

# Test avec différents workers
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

        print_info "Test avec $workers workers"

        # Hash original
        local original_hash=$(md5sum "$input" | cut -d' ' -f1)

        # Encryption
        time_start=$(date +%s%N)
        "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null
        time_end=$(date +%s%N)
        time_ms=$(( (time_end - time_start) / 1000000 ))

        if [ $? -ne 0 ]; then
            print_error "Encryption échouée avec $workers workers"
            return 1
        fi

        print_info "Encryption time: ${time_ms}ms"

        # Décryption
        "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null

        if [ $? -ne 0 ]; then
            print_error "Décryption échouée avec $workers workers"
            return 1
        fi

        # Vérification
        local decrypted_hash=$(md5sum "$decrypted" | cut -d' ' -f1)

        if [ "$original_hash" = "$decrypted_hash" ]; then
            print_success "Workers $workers: OK (${time_ms}ms)"
        else
            print_error "Workers $workers: Hash mismatch"
            return 1
        fi
    done

    return 0
}

# Test mauvais mot de passe
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

    # Tentative de décryption avec mauvais mot de passe
    local decrypted="$TEST_DIR/decrypted/wrong_pass.txt"

    # NE PAS rediriger stderr pour voir l'erreur
    # NE PAS utiliser --force (inutile pour la décryption)
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$wrong_password" 2>&1 | grep -q "authentication failed"
    local exit_code=$?

    if [ $exit_code -eq 0 ]; then
        print_success "Décryption a échoué (comme attendu)"
        return 0
    else
        print_error "ERREUR: La décryption aurait dû échouer avec mauvais mot de passe!"
        return 1
    fi
}

# Test fichier vide
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
    print_success "Encryption du fichier vide réussie"

    # Décryption
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -ne 0 ]; then
        print_error "Décryption du fichier vide échouée"
        return 1
    fi

    # Vérification taille
    local size=$(stat -c%s "$decrypted" 2>/dev/null || echo "0")
    if [ "$size" -eq 0 ]; then
        print_success "Fichier décrypté est vide (taille 0)"
        return 0
    else
        print_error "Fichier décrypté n'est pas vide (taille: $size)"
        return 1
    fi
}

# Test force overwrite
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

    # Deuxième encryption avec --force
    "$CRYPTOOL_BIN" encrypt "$input" "$output" --pass "$password" --force --quiet 2>/dev/null

    if [ $? -eq 0 ]; then
        print_success "Force overwrite fonctionne"
        return 0
    else
        print_error "Force overwrite a échoué"
        return 1
    fi
}

# Test avec tous les types de fichiers
test_all_file_types() {
    print_test_header "Tous les types de fichiers"

    local password="all-types-test"
    local failed=0

    for input in "$TEST_DIR/input/"*; do
        if [ -f "$input" ]; then
            local filename=$(basename "$input")
            local encrypted="$TEST_DIR/encrypted/${filename}.enc"
            local decrypted="$TEST_DIR/decrypted/${filename}"

            print_info "Test: $filename"

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

# Test performance
test_performance() {
    if [ "$SHORT_MODE" = "true" ]; then
        print_warning "Test performance ignoré (mode court)"
        return 0
    fi

    print_test_header "Performance"

    local input="$TEST_DIR/input/random10.bin"
    local password="perf-test"
    local size_mb=10

    if [ ! -f "$input" ]; then
        print_warning "Fichier random10.bin non trouvé, test ignoré"
        return 0
    fi

    echo ""
    printf "%-15s %-15s %-15s %-15s\n" "Workers" "Encrypt(s)" "Decrypt(s)" "Speed(MB/s)"
    echo "-------------------------------------------------------------"

    for workers in 1 2 4 8; do
        local encrypted="$TEST_DIR/temp/perf_${workers}.enc"
        local decrypted="$TEST_DIR/temp/perf_${workers}.dec"

        # Encryption avec mesure de temps
        local enc_time=$( { time "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null; } 2>&1 | grep real | awk '{print $2}' | sed 's/[^0-9.]//g')

        if [ -z "$enc_time" ]; then
            enc_time=$( { time "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers "$workers" --force --quiet 2>/dev/null; } 2>&1 | grep real | awk '{print $2}' | sed 's/m/*60+/g' | sed 's/s//' | bc 2>/dev/null || echo "0")
        fi

        if [ -z "$enc_time" ]; then
            enc_time="0"
        fi

        # Décryption avec mesure de temps
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

        printf "%-15s %-15s %-15s %-15s\n" "$workers" "$enc_time" "$dec_time" "$enc_speed"
    done

    echo ""
    print_success "Test performance terminé"
    return 0
}

# Test intégrité après corruption contrôlée
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

    # Copier et corrompre
    local corrupted="$TEST_DIR/encrypted/corrupted.enc"
    cp "$encrypted" "$corrupted"

    # Corrompre un byte dans le ciphertext (après le header)
    # Header: Magic(4) + Version(1) + Salt(16) + ChunkSize(4) + HMAC(32) + Nonce(12) = 69 bytes
    # On corrompt à l'offset 100 (dans le premier chunk)
    dd if=/dev/zero of="$corrupted" bs=1 count=1 seek=100 conv=notrunc 2>/dev/null

    # Tentative de décryption - doit échouer car GCM détecte la corruption
    local decrypted="$TEST_DIR/decrypted/corrupted.txt"

    # Laisser stderr pour voir l'erreur GCM
    "$CRYPTOOL_BIN" decrypt "$corrupted" "$decrypted" --pass "$password" 2>&1 | grep -q "decryption failed"
    local exit_code=$?

    if [ $exit_code -eq 0 ]; then
        print_success "Corruption détectée (décryption échouée)"
        return 0
    else
        print_error "ERREUR: Décryption a réussi malgré la corruption!"
        return 1
    fi
}

# Test gros fichier
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

    print_info "Taille du fichier: $(du -h "$input" | cut -f1)"

    # Encryption
    print_info "Encryption en cours..."
    local enc_start=$(date +%s)
    "$CRYPTOOL_BIN" encrypt "$input" "$encrypted" --pass "$password" --workers 8 --force --quiet 2>/dev/null
    local enc_end=$(date +%s)

    if [ $? -ne 0 ]; then
        print_error "Encryption du gros fichier échouée"
        return 1
    fi
    print_success "Encryption terminée en $((enc_end - enc_start)) secondes"

    # Décryption
    print_info "Décryption en cours..."
    local dec_start=$(date +%s)
    "$CRYPTOOL_BIN" decrypt "$encrypted" "$decrypted" --pass "$password" --force --quiet 2>/dev/null
    local dec_end=$(date +%s)

    if [ $? -ne 0 ]; then
        print_error "Décryption du gros fichier échouée"
        return 1
    fi
    print_success "Décryption terminée en $((dec_end - dec_start)) secondes"

    # Vérification
    local decrypted_hash=$(md5sum "$decrypted" | cut -d' ' -f1)

    if [ "$original_hash" = "$decrypted_hash" ]; then
        print_success "Gros fichier: vérification OK"
        return 0
    else
        print_error "Gros fichier: hash mismatch"
        return 1
    fi
}

# Fonction principale
main() {
    print_header

    # Parse arguments
    if [ "$1" = "--short" ] || [ "$1" = "-s" ]; then
        SHORT_MODE="true"
        print_info "Mode court activé (tests lourds ignorés)"
    fi

    # Setup
    check_binary
    if [ $? -ne 0 ]; then
        print_error "Impossible de continuer sans le binaire"
        exit 1
    fi

    setup_directories
    cleanup_before_test
    generate_test_files

    # Exécution des tests
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

    # Liste des tests
    run_test "Test simple encrypt/decrypt" test_simple_encrypt_decrypt
    run_test "Test workers parallèles" test_workers_parallel
    run_test "Test mauvais mot de passe" test_wrong_password
    run_test "Test fichier vide" test_empty_file
    run_test "Test force overwrite" test_force_overwrite
    run_test "Test tous types de fichiers" test_all_file_types
    run_test "Test détection corruption" test_corruption_detection
    run_test "Test gros fichier" test_large_file
    run_test "Test performance" test_performance

    # Résumé
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

# Exécution
main "$@"
