# Makefile for managing Go project tasks such as running, building, testing, and maintaining dependencies.
.PHONY: dev run build test tidy deps-upgrade deps-clean-cache clean

LOG_DIR  := logs
BIN_NAME := replify
BIN_DIR  := bin

# Detect OS for binary extension
ifeq ($(OS),Windows_NT)
	BIN_EXT := .exe
else
	BIN_EXT :=
endif

BIN_OUT := $(BIN_DIR)/$(BIN_NAME)$(BIN_EXT)

# ==============================================================================
# Development
# Prints CLI help output for quick reference during development.
# ==============================================================================
dev:
	@mkdir -p ./main
	go run ./main/main.go

# ==============================================================================
# Running the main application
# Executes the main.go file, useful for development and quick testing
run:
	go run ./cmd/replify

# Building the application
# Compiles the main.go file into an executable, for production deployment
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_OUT) ./cmd/replify

# ==============================================================================
# Module support and testing
# Runs tests across all packages in the project, showing code coverage
test:
	go test -cover ./...

# Cleaning and maintaining dependencies
# Cleans up the module by removing unused dependencies
# Copies all dependencies into the vendor directory, ensuring reproducibility
tidy:
	go mod tidy
	go mod vendor

# Upgrading dependencies
# Updates all dependencies to their latest minor or patch versions
# Cleans up the module after upgrade
# Re-vendors dependencies after upgrade
deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

# Cleaning up the module cache
# Removes all items from the Go module cache
deps-clean-cache:
	go clean -modcache

# Running code coverage
# Generates code coverage report and logs the results
coverage:
	sh ./sh/go_deps.sh

# Generating project file tree
# Creates a text file representing the project's directory structure, excluding certain directories
tree:
	@mkdir -p $(LOG_DIR)
	tree -I ".gradle|.idea|build|logs|.vscode|.git|.github|vendor" > ./$(LOG_DIR)/tree_source_oss.txt
	cat ./$(LOG_DIR)/tree_source_oss.txt

# ==============================================================================
# Clean
# Removes the build directory and log directory.
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(LOG_DIR)
