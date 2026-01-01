.PHONY: build run clean install test example

# Build the generator
build:
	@echo "Building atomic-generator..."
	@go build -o bin/atomic-generator ./cmd/generator

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build -o bin/atomic-generator-linux-amd64 ./cmd/generator
	@GOOS=darwin GOARCH=amd64 go build -o bin/atomic-generator-darwin-amd64 ./cmd/generator
	@GOOS=darwin GOARCH=arm64 go build -o bin/atomic-generator-darwin-arm64 ./cmd/generator
	@GOOS=windows GOARCH=amd64 go build -o bin/atomic-generator-windows-amd64.exe ./cmd/generator
	@echo "✅ Built for all platforms in bin/"

# Run the generator with example file
run: build
	@echo "Running generator..."
	@./bin/atomic-generator -input examples/bch_complete_atomic_structure.json -output output/bch-app

# Example: Generate from BCH structure
example: build
	@echo "Generating Barcelona Culinary Hub app..."
	@./bin/atomic-generator -input examples/bch_complete_atomic_structure.json -output output/bch-app
	@echo ""
	@echo "Project generated! To run it:"
	@echo "  cd output/bch-app"
	@echo "  npm install"
	@echo "  npm run dev"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf output/
	@echo "✅ Cleaned"

# Install the generator to $GOPATH/bin
install:
	@echo "Installing atomic-generator..."
	@go install ./cmd/generator
	@echo "✅ Installed to $(go env GOPATH)/bin/generator"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Show version
version: build
	@./bin/atomic-generator -version

# Help
help:
	@echo "Atomic Generator - Makefile commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       - Build the generator binary"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  run         - Build and run with example file"
	@echo "  example     - Generate the BCH example app"
	@echo "  clean       - Remove build artifacts"
	@echo "  install     - Install to GOPATH/bin"
	@echo "  test        - Run tests"
	@echo "  version     - Show version"
	@echo "  help        - Show this help message"
