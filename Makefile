# Development commands (like package.json scripts)
.PHONY: dev build run clean test

# Run in development mode with auto-reload
dev:
	@echo "🚀 Starting development server..."
	go run cmd/server/main.go

# Build the application
build:
	@echo "🏗️  Building application..."
	go build -o bin/server cmd/server/main.go

# Run the built application
run: build
	@echo "▶️  Running application..."
	./bin/server

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/

# Run tests
test:
	@echo "🧪 Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	golint ./...