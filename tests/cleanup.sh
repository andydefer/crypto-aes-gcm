#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEST_DIR="$SCRIPT_DIR/test_data"
RESULT_DIR="$SCRIPT_DIR/results"
BUILD_DIR="$(cd "$SCRIPT_DIR/.." && pwd)/build"

echo "🧹 Nettoyage des fichiers de test..."

# Supprimer les répertoires de test
rm -rf "$TEST_DIR"
rm -rf "$RESULT_DIR"

# Optionnel: supprimer le binaire
if [ "$1" = "--all" ] || [ "$1" = "-a" ]; then
    echo "🗑️  Suppression du binaire..."
    rm -rf "$BUILD_DIR"
fi

echo "✅ Nettoyage terminé"
