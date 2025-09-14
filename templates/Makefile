.PHONY: safety-check lint test coverage static-analysis

generate:
	go generate ./...

build:
	@go build -o bin/{{.Name}}

run: build
	go run main.go

lint:
	@ls -la .golangci.yml || echo "File not found"
	golangci-lint run --config .golangci.yml

static-analysis:
	go vet ./...
	staticcheck ./...
	gosec ./...

test:
	@echo "Running unit tests..."
	go test -v ./... -cover -short
test-all: test

coverage:
	go test -coverprofile=coverage.out ./... -short
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

check-rules:
	@echo "Checking for recursion..."
	@! grep -r "func.*(" . --include="*.go" | grep -v "_test.go" | xargs -I {} grep -l "func.*(" {} | xargs grep -n "return.*(" | grep -v "return.*\." || (echo "Potential recursion found" && exit 1)
	
	@echo "Checking function length..."
	@find . -name "*.go" -not -name "*_test.go" -exec awk '/^func / {start=NR} /^}$$/ {if(NR-start > 60) print FILENAME":"start":"NR-start" lines"}' {} \;

docker-run:
	docker compose up -d

docker-build:
	docker build -t {{.Name}} .

docker-push:
	docker push {{.Name}}

clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

install-deps:
	go mod download
	go mod tidy

install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

dev:
	@echo "Starting development server..."
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
	fi
	go run main.go

dev-watch:
	@echo "Starting development server with file watching..."
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
	fi
	air

setup:
	@echo "Setting up development environment..."
	make install-deps
	make install-tools
	@if [ ! -f .env ]; then \
		echo "Creating .env from .env.example..."; \
		cp .env.example .env; \
		echo "Please update .env with your configuration"; \
	fi
	@echo "Setup complete!"

help:
	@echo "Available commands:"
	@echo "  build              - Build the application"
	@echo "  run                - Run the application"
	@echo "  dev                - Start development server"
	@echo "  dev-watch          - Start development server with file watching"
	@echo "  test               - Run unit tests"
	@echo "  test-all           - Run all tests"
	@echo "  coverage           - Generate coverage report"
	@echo "  lint               - Run linter"
	@echo "  static-analysis    - Run static analysis tools"
	@echo "  setup              - Setup development environment"
	@echo "  clean              - Clean build artifacts"
	@echo "  install-deps       - Install Go dependencies"
	@echo "  install-tools      - Install development tools"
	@echo "  docker-build       - Build Docker image"
	@echo "  docker-run         - Run with Docker Compose"
	@echo "  help               - Show this help message"
