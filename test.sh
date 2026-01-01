#!/bin/bash

# Test script for Atomic Generator
# This script builds and tests the generator with the example JSON

set -e

echo "====================================="
echo "Atomic Generator - Test Script"
echo "====================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed"
    echo ""
    echo "Please install Go 1.21 or later:"
    echo "  - Download from: https://go.dev/dl/"
    echo "  - Or use your package manager:"
    echo "    Ubuntu/Debian: sudo apt install golang"
    echo "    Mac: brew install go"
    echo ""
    exit 1
fi

echo "✅ Go $(go version | awk '{print $3}') found"
echo ""

# Build the generator
echo "🔨 Building atomic-generator..."
go build -o bin/atomic-generator ./cmd/generator

if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi
echo ""

# Check if example file exists
if [ ! -f "examples/bch_complete_atomic_structure.json" ]; then
    echo "❌ Example file not found: examples/bch_complete_atomic_structure.json"
    exit 1
fi

echo "✅ Example file found"
echo ""

# Run the generator
echo "🚀 Generating Barcelona Culinary Hub app..."
echo ""

./bin/atomic-generator -input examples/bch_complete_atomic_structure.json -output output/bch-app

if [ $? -eq 0 ]; then
    echo ""
    echo "====================================="
    echo "✅ Generation completed successfully!"
    echo "====================================="
    echo ""
    echo "Your app is ready in: output/bch-app"
    echo ""
    echo "To run it:"
    echo "  cd output/bch-app"
    echo "  npm install"
    echo "  npm run dev"
    echo ""
else
    echo ""
    echo "❌ Generation failed"
    exit 1
fi
