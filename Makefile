.PHONY: all build test e2e clean

BINARY_NAME=create-harness-app
CMD_PATH=./cmd/create-harness-app

all: test build

build:
	@echo "🔨 Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(CMD_PATH)
	@echo "✅ Build completed successfully."

test:
	@echo "🧪 Running unit tests..."
	go test ./... -v

e2e: build
	@echo "🚀 Running E2E Harness Verification..."
	./$(BINARY_NAME) -template antigravity test-e2e-app
	cd test-e2e-app && ../$(BINARY_NAME) status
	cd test-e2e-app && ../$(BINARY_NAME) next
	./$(BINARY_NAME) extract test-e2e-app --template antigravity
	rm -rf test-e2e-app
	@echo "✅ E2E Verification Passed!"

clean:
	@echo "🧹 Cleaning binary..."
	rm -f $(BINARY_NAME)
	rm -rf test-e2e-app test-app test-dev-app
