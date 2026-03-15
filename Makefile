BINARY=agent-switcher
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")"

.PHONY: all build test lint clean run install fmt vet check tidy

all: lint test build

build:
	@echo "Building $(BINARY)..."
	go build $(LDFLAGS) -o $(BINARY) main.go

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	go clean

run: build
	./$(BINARY)

install: build
	go install

check: fmt vet lint test
	@echo "All checks passed!"

help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  test         - Run tests with race detection"
	@echo "  test-coverage- Run tests and generate HTML coverage report"
	@echo "  lint         - Run golangci-lint"
	@echo "  fmt          - Format code with gofmt"
	@echo "  vet          - Run go vet"
	@echo "  tidy         - Tidy go modules"
	@echo "  clean        - Remove build artifacts"
	@echo "  run          - Build and run"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  check        - Run fmt, vet, lint, and test"