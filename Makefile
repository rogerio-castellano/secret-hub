APP_NAME = secret-hub
BUILD_DIR = ./bin
BUILD_PATH = $(BUILD_DIR)/$(APP_NAME)

SECRET_NAME = DEMO_SECRET
SECRET_VALUE = myS3cr3t!
KEY_FILE = test-key.bin
STORE_FILE = secrets.json
ENC_FILE = secret.enc
PLAIN_FILE = recovered.txt

# Default target
.PHONY: all
all: build

.PHONY: build
build:
	@echo "🚀 Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_PATH) .

.PHONY: test
test:
	@echo "🧪 Running all tests..."
	go test ./...

.PHONY: integration-test
integration-test:
	@echo "🧩 Running integration test..."
	go test -v integration_test.go

.PHONY: clean
clean:
	@echo "🧹 Cleaning build and temp files..."
	rm -rf $(BUILD_DIR) $(KEY_FILE) $(STORE_FILE) $(ENC_FILE) $(PLAIN_FILE) secret.txt

.PHONY: release
release:
	@echo "📦 Creating release binary..."
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_PATH) .

.PHONY: lint
lint:
	@echo "🔍 Linting code..."
	golangci-lint run

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make build             - Compile the app into ./bin/$(APP_NAME)"
	@echo "  make test              - Run all unit + integration tests"
	@echo "  make integration-test  - Run integration test only"
	@echo "  make demo              - Run full demo: generate key, store, get, delete"
	@echo "  make encrypt-demo      - Run encrypt/decrypt file demo"
	@echo "  make clean             - Remove build and temp files"
	@echo "  make release           - Cross-compile a small binary"
	@echo "  make lint              - Run linter"
	@echo "  make help              - Show this message"

# 🧪 Run demo lifecycle for CLI
.PHONY: demo
demo: build
	@echo "🔑 Generating key..."
	$(BUILD_PATH) generate-key --out $(KEY_FILE)
	@echo "🔐 Storing secret..."
	$(BUILD_PATH) store --name $(SECRET_NAME) --value "$(SECRET_VALUE)" --key $(KEY_FILE) --storage $(STORE_FILE)
	@echo "📋 Listing secrets..."
	$(BUILD_PATH) list --storage $(STORE_FILE)
	@echo "📤 Retrieving secret..."
	$(BUILD_PATH) get --name $(SECRET_NAME) --key $(KEY_FILE) --storage $(STORE_FILE)
	@echo "🗑️ Deleting secret..."
	$(BUILD_PATH) delete --name $(SECRET_NAME) --storage $(STORE_FILE)
	@echo "✅ Demo completed."

# 🔐 Encrypt/decrypt file demo
.PHONY: encrypt-demo
encrypt-demo: build
	@echo "🔑 Generating key..."
	$(BUILD_PATH) generate-key --out $(KEY_FILE)
	@echo "Writing test file..."
	echo "super-secret-token" > secret.txt
	@echo "🔒 Encrypting file..."
	$(BUILD_PATH) encrypt --in secret.txt --out $(ENC_FILE) --key $(KEY_FILE)
	@echo "🔓 Decrypting file..."
	$(BUILD_PATH) decrypt --in $(ENC_FILE) --out $(PLAIN_FILE) --key $(KEY_FILE)
	@echo "📄 Decrypted content:"
	cat $(PLAIN_FILE)

.PHONY: sync
sync:
	git pull
	git push
