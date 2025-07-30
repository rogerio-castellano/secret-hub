APP_NAME = secret-hub
BUILD_DIR = ./bin
BUILD_PATH = $(BUILD_DIR)/$(APP_NAME)
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
OS := $(shell go env GOOS)

EXT :=
ifeq ($(OS),windows)
    EXT := .exe
endif

# Default target
.PHONY: all
all: build

# 🔨 Build the CLI
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_PATH)$(EXT) .

# 🧪 Run unit + integration tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# 🧪 Only run integration test
.PHONY: integration-test
integration-test:
	@echo "Running integration test..."
	go test -v integration_test.go

# 🧼 Clean binaries and test artifacts
.PHONY: clean
clean:
	@echo "Cleaning build and test files..."
	rm -rf $(BUILD_DIR)
	rm -f test-key.bin secret.enc recovered.txt secrets.json test_secrets.json

# 📦 Create a release binary
.PHONY: release
release:
	@echo "Creating release binary..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_PATH) .

# 🧪 Lint (requires go install golang.org/x/lint/golint)
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# 🆘 Help message
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build            - Compile the app into ./bin/$(APP_NAME)"
	@echo "  make test             - Run all unit and integration tests"
	@echo "  make integration-test - Run integration tests only"
	@echo "  make clean            - Remove build and temp files"
	@echo "  make release          - Cross-compile a stripped binary"
	@echo "  make lint             - Run linter"
	@echo "  make help             - Show this message"


.PHONY: sync
sync:
	git pull
	git push
