.PHONY: help install build run-api run-worker test lint fmt clean migrate docker-up docker-down

# Default target
help:
@echo "Typecraft - Makefile Commands"
@echo ""
@echo "  make install      - Install Go dependencies"
@echo "  make build        - Build binaries"
@echo "  make run-api      - Run API server"
@echo "  make run-worker   - Run async worker"
@echo "  make test         - Run tests"
@echo "  make lint         - Run linter"
@echo "  make fmt          - Format code"
@echo "  make migrate      - Run database migrations"
@echo "  make docker-up    - Start Docker services"
@echo "  make docker-down  - Stop Docker services"
@echo "  make clean        - Clean build artifacts"

install:
@echo "📦 Installing dependencies..."
go mod download
go mod tidy

build:
@echo "🔨 Building binaries..."
go build -o bin/api ./cmd/api
go build -o bin/worker ./cmd/worker

run-api:
@echo "🚀 Starting API server..."
go run ./cmd/api/main.go

run-worker:
@echo "⚙️  Starting worker..."
go run ./cmd/worker/main.go

test:
@echo "🧪 Running tests..."
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

lint:
@echo "🔍 Running linter..."
golangci-lint run ./...

fmt:
@echo "✨ Formatting code..."
gofmt -s -w .
goimports -w .

migrate:
@echo "🗄️  Running migrations..."
# TODO: Implement migrations
@echo "Migrations not yet implemented"

docker-up:
@echo "🐳 Starting Docker services..."
docker compose up -d
@echo "✅ Services running: PostgreSQL, Redis, MinIO"

docker-down:
@echo "🛑 Stopping Docker services..."
docker compose down

clean:
@echo "🧹 Cleaning build artifacts..."
rm -rf bin/
rm -rf dist/
rm -f coverage.out coverage.html
go clean
