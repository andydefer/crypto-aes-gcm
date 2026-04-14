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
# Directories
# ---------------------------------------------------
PRIVATE_DIR = private
DIFF_DIR = $(PRIVATE_DIR)/diff
CONCAT_DIR = $(PRIVATE_DIR)/concat-all

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
	@mkdir -p $(DIFF_DIR)
	@timestamp=$$(date +"%Y%m%d_%H%M%S"); \
	read -p "📁 Enter directory/path(s) to include in the diff (space-separated, leave empty for all changes): " DIR_PATHS; \
	if [ -z "$$DIR_PATHS" ]; then \
		echo "📝 Generating git diff for ALL changes into $(DIFF_DIR)/diff_$${timestamp}.txt..."; \
		echo "Tu es un expert en revue de code et en conventions de commits (Conventional Commits)." > $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "À partir du diff Git ci-dessous, fais les choses suivantes :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "1. Propose un nom de commit clair et concis en anglais" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   avec le format <type>(<scope>): <description>," >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   en respectant les Conventional Commits" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   (ex: feat:, fix:, refactor:, test:, chore:, docs:)." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "2. Rédige un résumé du travail effectué en quelques phrases," >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   orienté métier et technique." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "3. Donne une liste d'exemples concrets de changements, en t'appuyant sur le diff :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - méthodes ajoutées, modifiées ou supprimées" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - responsabilités déplacées ou clarifiées" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - améliorations de validation, de logique ou de structure" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - impacts fonctionnels éventuels" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "Contraintes :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Ne décris que ce qui est réellement visible dans le diff" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Sois précis, factuel et structuré" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Évite les suppositions" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Utilise un ton professionnel" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "4. SI et SEULEMENT SI les changements sont cassants (breaking changes) :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Génère une entrée de CHANGELOG conforme à Keep a Changelog et SemVer." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Le changelog doit apparaître APRES les recommandations ci-dessus." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Utilise STRICTEMENT la structure suivante :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ## [X.0.0] - YYYY-MM-DD" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ### Changed" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     - Description claire du changement cassant" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ### Removed (si applicable)" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     - API, méthode ou comportement supprimé" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ### Security (si applicable)" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     - Impact sécurité lié au changement" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Ne génère PAS de changelog si aucun breaking change n'est détecté." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - N'invente PAS de version." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "Voici le diff :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		git diff HEAD -- . ':!*.phpunit.result.cache' ':!$(PRIVATE_DIR)/*' ':!*.enc' ':!*.bin' ':!cryptool*' ':!build/*' ':!tests/test_data/*' >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "✅ Clean diff generated successfully: $(DIFF_DIR)/diff_$${timestamp}.txt"; \
	else \
		echo "📝 Generating clean git diff for paths: $${DIR_PATHS} into $(DIFF_DIR)/diff_$${timestamp}.txt..."; \
		echo "Tu es un expert en revue de code et en conventions de commits (Conventional Commits)." > $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "À partir du diff Git ci-dessous, fais les choses suivantes :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "1. Propose un nom de commit clair et concis en anglais" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   avec le format <type>(<scope>): <description>," >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   en respectant les Conventional Commits" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   (ex: feat:, fix:, refactor:, test:, chore:, docs:)." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "2. Rédige un résumé du travail effectué en quelques phrases," >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   orienté métier et technique." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "3. Donne une liste d'exemples concrets de changements, en t'appuyant sur le diff :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - méthodes ajoutées, modifiées ou supprimées" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - responsabilités déplacées ou clarifiées" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - améliorations de validation, de logique ou de structure" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - impacts fonctionnels éventuels" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "Contraintes :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Ne décris que ce qui est réellement visible dans le diff" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Sois précis, factuel et structuré" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Évite les suppositions" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Utilise un ton professionnel" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "4. SI et SEULEMENT SI les changements sont cassants (breaking changes) :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Génère une entrée de CHANGELOG conforme à Keep a Changelog et SemVer." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Le changelog doit apparaître APRES les recommandations ci-dessus." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Utilise STRICTEMENT la structure suivante :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ## [X.0.0] - YYYY-MM-DD" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ### Changed" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     - Description claire du changement cassant" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ### Removed (si applicable)" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     - API, méthode ou comportement supprimé" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     ### Security (si applicable)" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "     - Impact sécurité lié au changement" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - Ne génère PAS de changelog si aucun breaking change n'est détecté." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "   - N'invente PAS de version." >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "Voici le diff :" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "" >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		git diff HEAD -- $$DIR_PATHS ':!*.phpunit.result.cache' ':!$(PRIVATE_DIR)/*' ':!*.enc' ':!*.bin' ':!cryptool*' ':!build/*' ':!tests/test_data/*' >> $(DIFF_DIR)/diff_$${timestamp}.txt; \
		echo "✅ Clean diff generated successfully: $(DIFF_DIR)/diff_$${timestamp}.txt"; \
	fi

.PHONY: list-diffs
list-diffs:
	@echo "📁 Available diff files:"
	@ls -la $(DIFF_DIR)/diff_*.txt 2>/dev/null || echo "No diff files found in $(DIFF_DIR)"

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
	@mkdir -p $(CONCAT_DIR)
	@read -p "📁 Enter the source directory path to scan (leave empty for default './pkg ./internal ./cmd'): " SOURCE_PATH; \
	if [ -z "$$SOURCE_PATH" ]; then \
		SOURCE_DIRS="./pkg ./internal ./cmd"; \
		echo "🔗 Concatenating all TEXT files from default directories: $${SOURCE_DIRS} into $(CONCAT_DIR)/all.txt..."; \
	else \
		SOURCE_DIRS="$$SOURCE_PATH"; \
		echo "🔗 Concatenating all TEXT files from directory: $${SOURCE_DIRS} into $(CONCAT_DIR)/all.txt..."; \
	fi; \
	timestamp=$$(date +"%Y%m%d_%H%M%S"); \
	output_file="$(CONCAT_DIR)/all_$${timestamp}.txt"; \
	> "$$output_file"; \
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
				-exec sh -c 'echo ""; echo "// ==== {} ==="; echo ""; cat "{}" 2>/dev/null || echo "⚠️  Cannot read: {}"' \; >> "$$output_file" 2>/dev/null; \
		else \
			echo "⚠️  Directory not found: $$dir"; \
		fi; \
	done; \
	echo "✅ File $$output_file generated successfully from: $${SOURCE_DIRS} (binary files excluded)"

.PHONY: list-concats
list-concats:
	@echo "📁 Available concatenated files:"
	@ls -la $(CONCAT_DIR)/all_*.txt 2>/dev/null || echo "No concatenated files found in $(CONCAT_DIR)"

# ---------------------------------------------------
# Clean temporary files
# ---------------------------------------------------

.PHONY: clean
clean:
	@echo "🧹 Cleaning temporary files..."
	@rm -f all.txt
	@rm -f test.enc test.decrypted.txt test.bin
	@echo "✅ Clean completed"

.PHONY: clean-private
clean-private:
	@echo "🧹 Cleaning private directory..."
	@rm -rf $(PRIVATE_DIR)
	@echo "✅ Private directory cleaned"

.PHONY: clean-all
clean-all: clean
	@echo "🧹 Cleaning all generated files..."
	@rm -f cryptool cryptool-*
	@rm -rf $(PRIVATE_DIR)
	@rm -rf build
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
	@mkdir -p build
	@go build -o build/cryptool ./cmd/cryptool
	@chmod +x build/cryptool
	@echo "✅ Build completed: build/cryptool"

.PHONY: build-all
build-all:
	@echo "🔨 Building for all platforms..."
	@mkdir -p build

	@echo "  📦 Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build -o build/cryptool-linux-amd64 ./cmd/cryptool && chmod +x build/cryptool-linux-amd64

	@echo "  📦 Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build -o build/cryptool-windows-amd64.exe ./cmd/cryptool

	@echo "  📦 Building for macOS Intel (amd64)..."
	@GOOS=darwin GOARCH=amd64 go build -o build/cryptool-darwin-amd64 ./cmd/cryptool && chmod +x build/cryptool-darwin-amd64

	@echo "  📦 Building for macOS Apple Silicon (arm64)..."
	@GOOS=darwin GOARCH=arm64 go build -o build/cryptool-darwin-arm64 ./cmd/cryptool && chmod +x build/cryptool-darwin-arm64

	@echo "  📦 Building for Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 go build -o build/cryptool-linux-arm64 ./cmd/cryptool && chmod +x build/cryptool-linux-arm64

	@echo "✅ Build completed for all platforms"
	@ls -lh build/cryptool-* 2>/dev/null || true

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


.PHONY: gotestsum
gotestsum:
	@echo "🧪 Running tests with gotestsum..."
	@PATH="$$(go env GOPATH)/bin:$$PATH" gotestsum --format testname ./...

# ---------------------------------------------------
# Test Scripts Commands
# ---------------------------------------------------

.PHONY: test-scripts
test-scripts:
	@echo "🧪 Running test scripts..."
	@chmod +x tests/run_tests.sh
	@./tests/run_tests.sh

.PHONY: test-scripts-short
test-scripts-short:
	@echo "🧪 Running test scripts (short mode)..."
	@chmod +x tests/run_tests.sh
	@./tests/run_tests.sh --short

.PHONY: test-scenarios
test-scenarios:
	@echo "🎬 Running test scenarios..."
	@chmod +x tests/test_scenarios.sh
	@./tests/test_scenarios.sh

.PHONY: test-scenarios-short
test-scenarios-short:
	@echo "🎬 Running test scenarios (short mode)..."
	@chmod +x tests/test_scenarios.sh
	@./tests/test_scenarios.sh --short

.PHONY: test-scenarios-verbose
test-scenarios-verbose:
	@echo "🎬 Running test scenarios (verbose mode)..."
	@chmod +x tests/test_scenarios.sh
	@./tests/test_scenarios.sh --verbose

.PHONY: test-all
test-all: test-scripts test-scenarios
	@echo "✅ All tests completed"

.PHONY: test-all-short
test-all-short: test-scripts-short test-scenarios-short
	@echo "✅ All short tests completed"

# ---------------------------------------------------
# Generate Test Files
# ---------------------------------------------------

.PHONY: generate-test-files
generate-test-files:
	@echo "📁 Generating test files..."
	@chmod +x tests/generate_test_files.sh
	@./tests/generate_test_files.sh

.PHONY: generate-test-files-short
generate-test-files-short:
	@echo "📁 Generating test files (short mode)..."
	@chmod +x tests/generate_test_files.sh
	@./tests/generate_test_files.sh --short

# ---------------------------------------------------
# Clean Test Data
# ---------------------------------------------------

.PHONY: clean-test-data
clean-test-data:
	@echo "🧹 Cleaning test data..."
	@rm -rf tests/test_data
	@echo "✅ Test data cleaned"

.PHONY: clean-all-tests
clean-all-tests: clean-test-data
	@echo "🧹 Cleaning all test artifacts..."
	@rm -rf tests/results
	@rm -f tests/*.log
	@echo "✅ All test artifacts cleaned"

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
	@echo "  generate-ai-diff      Generate clean diff for AI review (saved in private/diff/)"
	@echo "  list-diffs            List all generated diff files"
	@echo "  git-tag-republish     Force push the last tag"
	@echo ""
	@echo "📁 File Management:"
	@echo "  concat-all            Concatenate all TEXT files (saved in private/concat-all/)"
	@echo "  list-concats          List all concatenated files"
	@echo "  clean                 Remove temporary files"
	@echo "  clean-private         Remove private directory"
	@echo "  clean-all             Remove all generated files including binaries"
	@echo ""
	@echo "🔨 Build:"
	@echo "  build                 Build cryptool for current platform"
	@echo "  build-all             Build cryptool for all platforms (Linux, Windows, macOS)"
	@echo ""
	@echo "🧪 Go Tests:"
	@echo "  test                  Run all Go tests"
	@echo "  test-short            Run short Go tests"
	@echo "  test-coverage         Run Go tests with coverage report"
	@echo "  gotestsum             Run tests with gotestsum formatter"
	@echo ""
	@echo "📋 Test Scripts (Realistic):"
	@echo "  test-scripts          Run realistic test scripts"
	@echo "  test-scripts-short    Run realistic test scripts (short mode)"
	@echo "  test-scenarios        Run advanced test scenarios"
	@echo "  test-scenarios-short  Run advanced test scenarios (short mode)"
	@echo "  test-scenarios-verbose Run advanced test scenarios (verbose mode)"
	@echo "  test-all              Run all test scripts and scenarios"
	@echo "  test-all-short        Run all tests in short mode"
	@echo ""
	@echo "📁 Test Files Generation:"
	@echo "  generate-test-files       Generate test files (including 50MB+ files)"
	@echo "  generate-test-files-short Generate test files (short mode, no large files)"
	@echo ""
	@echo "🧹 Test Cleanup:"
	@echo "  clean-test-data       Clean test data only"
	@echo "  clean-all-tests       Clean all test artifacts"
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
