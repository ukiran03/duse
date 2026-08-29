# Default target shows available recipes
default:
    @just --list

# Build the duse binary
build:
    go build -o duse

# Run tests
test:
    go test -v

# Run tests with coverage
test-coverage:
    go test -v -cover

# Clean build artifacts
clean:
    rm -f duse

# Install to /usr/local/bin (requires sudo)
install: build
    sudo mv duse /usr/local/bin/

# Run benchmarks
bench:
    go test -bench=. -benchmem

# Format code
fmt:
    go fmt ./...

# Run linter
lint:
    golangci-lint run
