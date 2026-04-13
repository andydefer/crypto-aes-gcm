# ===================================================
# Generic Makefile with binary exclusion
# ===================================================

# ---------------------------------------------------
# Source Configuration
# ---------------------------------------------------
SOURCE_DIRS = src config database tests pkg cmd internal
IGNORED_FILES = CHANGED_FILES.md FILES_CHECKLIST.md Makefile .gitkeep
IGNORED_EXTENSIONS = .bin .exe .enc .jpg .png .pdf .zip .tar .gz .so .dll .o .a
IGNORED_FILES_PATTERNS = cryptool cryptool-* test.enc test.bin all.txt

# ---------------------------------------------------
# Version Control Operations
# ---------------------------------------------------

.PHONY: git-commit-push
git-commit-push:
	@read -p "Enter commit message: " commit_message; \
	if [ -z "$$commit_message" ]; then \
		echo "❌ Error: Commit message cannot be empty"; \
		exit 1; \
	fi; \
	git add .; \
	git commit -m "$$commit_message"; \
	git push

.PHONY: git-tag
git-tag:
	@bash -c '\
	read -p "Tag type (major/minor/patch): " tag_type; \
	last_tag=$$(git tag --sort=-v:refname | head -n 1); \
	if [ -z "$$last_tag" ]; then last_tag="0.0.0"; fi; \
	major=$$(echo $$last_tag | cut -d. -f1); \
	minor=$$(echo $$last_tag | cut -d. -f2); \
	patch=$$(echo $$last_tag | cut -d. -f3); \
	case "$$tag_type" in \
		major) major=$$((major + 1)); minor=0; patch=0;; \
		minor) minor=$$((minor + 1)); patch=0;; \
		patch) patch=$$((patch + 1));; \
		*) echo "❌ Invalid tag type: $$tag_type"; exit 1;; \
	esac; \
	new_tag="$$major.$$minor.$$patch"; \
	git tag -a "$$new_tag" -m "Release $$new_tag"; \
	git push origin "$$new_tag"; \
	echo "✅ Released new tag: $$new_tag"; \
	'

.PHONY: generate-ai-diff
generate-ai-diff:
	@mkdir -p diff
	@timestamp=$$(date +"%Y%m%d_%H%M%S"); \
	read -p "📁 Enter directory/path(s) to include in the diff (space-separated, leave empty for all changes): " DIR_PATHS; \
	if [ -z "$$DIR_PATHS" ]; then \
		echo "📝 Generating git diff for ALL changes into diff/diff_$${timestamp}.txt..."; \
		echo "Tu es un expert en revue de code et en conventions de commits (Conventional Commits)." > diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "À partir du diff Git ci-dessous, fais les choses suivantes :" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "1. Propose un nom de commit clair et concis en anglais" >> diff/diff_$${timestamp}.txt; \
		echo "   avec le format <type>(<scope>): <description>," >> diff/diff_$${timestamp}.txt; \
		echo "   en respectant les Conventional Commits" >> diff/diff_$${timestamp}.txt; \
		echo "   (ex: feat:, fix:, refactor:, test:, chore:, docs:)." >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "2. Rédige un résumé du travail effectué en quelques phrases," >> diff/diff_$${timestamp}.txt; \
		echo "   orienté métier et technique." >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "3. Donne une liste d'exemples concrets de changements, en t'appuyant sur le diff :" >> diff/diff_$${timestamp}.txt; \
		echo "   - méthodes ajoutées, modifiées ou supprimées" >> diff/diff_$${timestamp}.txt; \
		echo "   - responsabilités déplacées ou clarifiées" >> diff/diff_$${timestamp}.txt; \
		echo "   - améliorations de validation, de logique ou de structure" >> diff/diff_$${timestamp}.txt; \
		echo "   - impacts fonctionnels éventuels" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "Contraintes :" >> diff/diff_$${timestamp}.txt; \
		echo "   - Ne décris que ce qui est réellement visible dans le diff" >> diff/diff_$${timestamp}.txt; \
		echo "   - Sois précis, factuel et structuré" >> diff/diff_$${timestamp}.txt; \
		echo "   - Évite les suppositions" >> diff/diff_$${timestamp}.txt; \
		echo "   - Utilise un ton professionnel" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "4. SI et SEULEMENT SI les changements sont cassants (breaking changes) :" >> diff/diff_$${timestamp}.txt; \
		echo "   - Génère une entrée de CHANGELOG conforme à Keep a Changelog et SemVer." >> diff/diff_$${timestamp}.txt; \
		echo "   - Le changelog doit apparaître APRES les recommandations ci-dessus." >> diff/diff_$${timestamp}.txt; \
		echo "   - Utilise STRICTEMENT la structure suivante :" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "     ## [X.0.0] - YYYY-MM-DD" >> diff/diff_$${timestamp}.txt; \
		echo "     ### Changed" >> diff/diff_$${timestamp}.txt; \
		echo "     - Description claire du changement cassant" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "     ### Removed (si applicable)" >> diff/diff_$${timestamp}.txt; \
		echo "     - API, méthode ou comportement supprimé" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "     ### Security (si applicable)" >> diff/diff_$${timestamp}.txt; \
		echo "     - Impact sécurité lié au changement" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "   - Ne génère PAS de changelog si aucun breaking change n'est détecté." >> diff/diff_$${timestamp}.txt; \
		echo "   - N'invente PAS de version." >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "Voici le diff :" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		git diff HEAD -- . ':!*.phpunit.result.cache' ':!diff/*' ':!*.enc' ':!*.bin' ':!cryptool*' >> diff/diff_$${timestamp}.txt; \
		echo "✅ Clean diff generated successfully: diff/diff_$${timestamp}.txt"; \
	else \
		echo "📝 Generating clean git diff for paths: $${DIR_PATHS} into diff/diff_$${timestamp}.txt..."; \
		echo "Tu es un expert en revue de code et en conventions de commits (Conventional Commits)." > diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "À partir du diff Git ci-dessous, fais les choses suivantes :" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "1. Propose un nom de commit clair et concis en anglais" >> diff/diff_$${timestamp}.txt; \
		echo "   avec le format <type>(<scope>): <description>," >> diff/diff_$${timestamp}.txt; \
		echo "   en respectant les Conventional Commits" >> diff/diff_$${timestamp}.txt; \
		echo "   (ex: feat:, fix:, refactor:, test:, chore:, docs:)." >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "2. Rédige un résumé du travail effectué en quelques phrases," >> diff/diff_$${timestamp}.txt; \
		echo "   orienté métier et technique." >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "3. Donne une liste d'exemples concrets de changements, en t'appuyant sur le diff :" >> diff/diff_$${timestamp}.txt; \
		echo "   - méthodes ajoutées, modifiées ou supprimées" >> diff/diff_$${timestamp}.txt; \
		echo "   - responsabilités déplacées ou clarifiées" >> diff/diff_$${timestamp}.txt; \
		echo "   - améliorations de validation, de logique ou de structure" >> diff/diff_$${timestamp}.txt; \
		echo "   - impacts fonctionnels éventuels" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "Contraintes :" >> diff/diff_$${timestamp}.txt; \
		echo "   - Ne décris que ce qui est réellement visible dans le diff" >> diff/diff_$${timestamp}.txt; \
		echo "   - Sois précis, factuel et structuré" >> diff/diff_$${timestamp}.txt; \
		echo "   - Évite les suppositions" >> diff/diff_$${timestamp}.txt; \
		echo "   - Utilise un ton professionnel" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "4. SI et SEULEMENT SI les changements sont cassants (breaking changes) :" >> diff/diff_$${timestamp}.txt; \
		echo "   - Génère une entrée de CHANGELOG conforme à Keep a Changelog et SemVer." >> diff/diff_$${timestamp}.txt; \
		echo "   - Le changelog doit apparaître APRES les recommandations ci-dessus." >> diff/diff_$${timestamp}.txt; \
		echo "   - Utilise STRICTEMENT la structure suivante :" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "     ## [X.0.0] - YYYY-MM-DD" >> diff/diff_$${timestamp}.txt; \
		echo "     ### Changed" >> diff/diff_$${timestamp}.txt; \
		echo "     - Description claire du changement cassant" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "     ### Removed (si applicable)" >> diff/diff_$${timestamp}.txt; \
		echo "     - API, méthode ou comportement supprimé" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "     ### Security (si applicable)" >> diff/diff_$${timestamp}.txt; \
		echo "     - Impact sécurité lié au changement" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "   - Ne génère PAS de changelog si aucun breaking change n'est détecté." >> diff/diff_$${timestamp}.txt; \
		echo "   - N'invente PAS de version." >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		echo "Voici le diff :" >> diff/diff_$${timestamp}.txt; \
		echo "" >> diff/diff_$${timestamp}.txt; \
		git diff HEAD -- $$DIR_PATHS ':!*.phpunit.result.cache' ':!diff/*' ':!*.enc' ':!*.bin' ':!cryptool*' >> diff/diff_$${timestamp}.txt; \
		echo "✅ Clean diff generated successfully: diff/diff_$${timestamp}.txt"; \
	fi

.PHONY: list-diffs
list-diffs:
	@echo "📁 Available diff files:"
	@ls -la diff/diff_*.txt 2>/dev/null || echo "No diff files found"

.PHONY: git-tag-republish
git-tag-republish:
	@bash -c '\
	last_tag=$$(git tag --sort=-v:refname | head -n 1); \
	if [ -z "$$last_tag" ]; then echo "❌ No tags found!"; exit 1; fi; \
	echo "Republishing last tag: $$last_tag"; \
	git push origin "$$last_tag" --force; \
	echo "✅ Tag $$last_tag republished"; \
	'

# ---------------------------------------------------
# File Management Operations
# ---------------------------------------------------

.PHONY: concat-all
concat-all:
	@read -p "📁 Enter the source directory path to scan (leave empty for default './pkg ./internal ./cmd'): " SOURCE_PATH; \
	if [ -z "$$SOURCE_PATH" ]; then \
		SOURCE_DIRS="./pkg ./internal ./cmd"; \
		echo "🔗 Concatenating all TEXT files from default directories: $${SOURCE_DIRS} into all.txt..."; \
	else \
		SOURCE_DIRS="$$SOURCE_PATH"; \
		echo "🔗 Concatenating all TEXT files from directory: $${SOURCE_DIRS} into all.txt..."; \
	fi; \
	> all.txt; \
	for dir in $${SOURCE_DIRS}; do \
		if [ -d "$$dir" ]; then \
			find "$$dir" -type f \
				-not -name "*.bin" \
				-not -name "*.enc" \
				-not -name "*.exe" \
				-not -name "*.so" \
				-not -name "*.dll" \
				-not -name "*.o" \
				-not -name "*.a" \
				-not -name "cryptool" \
				-not -name "cryptool-*" \
				-not -name "test.enc" \
				-not -name "test.bin" \
				-not -name "all.txt" \
				-exec sh -c 'echo ""; echo "// ==== {} ==="; echo ""; cat "{}" 2>/dev/null || echo "⚠️  Cannot read: {}"' \; >> all.txt 2>/dev/null; \
		else \
			echo "⚠️  Directory not found: $$dir"; \
		fi; \
	done; \
	echo "✅ File all.txt generated successfully from: $${SOURCE_DIRS} (binary files excluded)"

# ---------------------------------------------------
# Clean temporary files
# ---------------------------------------------------

.PHONY: clean
clean:
	@echo "🧹 Cleaning temporary files..."
	@rm -f all.txt
	@rm -f test.enc test.decrypted.txt test.bin
	@echo "✅ Clean completed"

.PHONY: clean-all
clean-all: clean
	@echo "🧹 Cleaning all generated files..."
	@rm -f cryptool cryptool-*
	@rm -rf diff/*.txt
	@echo "✅ Deep clean completed"

# ---------------------------------------------------
# Release Management Workflow
# ---------------------------------------------------

.PHONY: release
release:
	@echo "🚀 Creating release..."
	@make git-tag
	@echo "✅ Release created successfully"

# ---------------------------------------------------
# Build Commands
# ---------------------------------------------------

.PHONY: build
build:
	@echo "🔨 Building cryptool for current platform..."
	@go build -o cryptool ./cmd/cryptool
	@chmod +x cryptool
	@echo "✅ Build completed: ./cryptool (executable with chmod +x)"

.PHONY: build-all
build-all:
	@echo "🔨 Building for all platforms..."
	@echo "  📦 Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -o cryptool-linux-amd64 ./cmd/cryptool && chmod +x cryptool-linux-amd64
	@echo "  📦 Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build -o cryptool-windows-amd64.exe ./cmd/cryptool
	@echo "  📦 Building for macOS Intel (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build -o cryptool-darwin-amd64 ./cmd/cryptool && chmod +x cryptool-darwin-amd64
	@echo "  📦 Building for macOS Apple Silicon (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build -o cryptool-darwin-arm64 ./cmd/cryptool && chmod +x cryptool-darwin-arm64
	@echo "  📦 Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build -o cryptool-linux-arm64 ./cmd/cryptool && chmod +x cryptool-linux-arm64
	@echo "✅ Build completed for all platforms"
	@ls -lh cryptool-* 2>/dev/null || true

# ---------------------------------------------------
# Test Commands
# ---------------------------------------------------

.PHONY: test
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

.PHONY: test-short
test-short:
	@echo "🧪 Running short tests..."
	@go test -short -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@go test -cover ./...
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# ---------------------------------------------------
# Help & Documentation
# ---------------------------------------------------

.PHONY: help
help:
	@echo "📚 Available commands:"
	@echo ""
	@echo "🚀 Version Control:"
	@echo "  git-commit-push       Commit and push all changes"
	@echo "  git-tag               Create and push a new version tag"
	@echo "  generate-ai-diff      Generate clean diff for AI review"
	@echo "  git-tag-republish     Force push the last tag"
	@echo ""
	@echo "📁 File Management:"
	@echo "  concat-all            Concatenate all TEXT files (excludes binaries)"
	@echo "  clean                 Remove temporary files"
	@echo "  clean-all             Remove all generated files including binaries"
	@echo ""
	@echo "🔨 Build:"
	@echo "  build                 Build cryptool for current platform (with chmod +x)"
	@echo "  build-all             Build cryptool for all platforms (Linux, Windows, macOS)"
	@echo ""
	@echo "🧪 Test:"
	@echo "  test                  Run all tests"
	@echo "  test-short            Run short tests"
	@echo "  test-coverage         Run tests with coverage report"
	@echo ""
	@echo "🔄 Release Management:"
	@echo "  release               Create new release (includes pre-release)"
	@echo ""
	@echo "❓ Help:"
	@echo "  help                  Display this help message"

# ---------------------------------------------------
# Default Target
# ---------------------------------------------------
.DEFAULT_GOAL := help
