.PHONY: build run test clean docker-build docker-run install fmt vet lint

# Variables
BINARY_NAME=nexttrace_exporter
VERSION?=0.1.0
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}"

# Default target
all: build

# Build the binary
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINARY_NAME}-linux-amd64 .

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .

# Run the application
run: build
	./${BINARY_NAME} --config.file=examples/config.yml

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -f ${BINARY_NAME}
	rm -f ${BINARY_NAME}-*
	rm -rf dist/
	rm -f coverage.out

# Build Docker image
docker-build:
	docker build -t ${BINARY_NAME}:${VERSION} .
	docker tag ${BINARY_NAME}:${VERSION} ${BINARY_NAME}:latest

# Run Docker container
docker-run:
	docker run -d \
		-p 9101:9101 \
		-v $(PWD)/examples/config.yml:/etc/nexttrace_exporter/config.yml \
		--cap-add=NET_RAW \
		--name ${BINARY_NAME} \
		${BINARY_NAME}:latest

# Stop and remove Docker container
docker-stop:
	docker stop ${BINARY_NAME} || true
	docker rm ${BINARY_NAME} || true

# Install the binary
install: build
	sudo cp ${BINARY_NAME} /usr/local/bin/

# Uninstall the binary
uninstall:
	sudo rm -f /usr/local/bin/${BINARY_NAME}

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  build-linux   - Build for Linux"
	@echo "  build-all     - Build for all platforms"
	@echo "  run           - Build and run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  deps          - Install dependencies"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-stop   - Stop and remove Docker container"
	@echo "  install       - Install binary to /usr/local/bin"
	@echo "  uninstall     - Uninstall binary"
	@echo "  help          - Show this help message"
