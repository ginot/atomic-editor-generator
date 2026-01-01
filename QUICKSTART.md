# Quick Start Guide

## Prerequisites

1. **Go 1.21 or later**
   ```bash
   # Check if Go is installed
   go version
   
   # If not installed:
   # Ubuntu/Debian
   sudo apt update && sudo apt install golang-go
   
   # macOS
   brew install go
   
   # Or download from: https://go.dev/dl/
   ```

2. **Node.js 18 or later** (for running the generated app)
   ```bash
   # Check Node version
   node --version
   
   # Install if needed
   # Ubuntu/Debian
   curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
   sudo apt-get install -y nodejs
   
   # macOS
   brew install node
   ```

## Installation

```bash
# Clone or navigate to the atomic-generator directory
cd atomic-generator

# Build the generator
make build

# Or manually:
go build -o bin/atomic-generator ./cmd/generator
```

## Usage

### Method 1: Using Make (Recommended)

```bash
# Generate the example Barcelona Culinary Hub app
make example

# The app will be generated in: output/bch-app
```

### Method 2: Using the Binary Directly

```bash
# Generate from any JSON structure
./bin/atomic-generator \
  -input examples/bch_complete_atomic_structure.json \
  -output ./my-custom-app
```

### Method 3: Using the Test Script

```bash
# Run the test script
./test.sh
```

## Running the Generated App

```bash
# Navigate to the generated app
cd output/bch-app

# Install dependencies
npm install

# Start development server
npm run dev

# Your app will be available at: http://localhost:3000
```

## Verify Installation

```bash
# Check generator version
./bin/atomic-generator -version

# Should output:
# Atomic Generator v1.0.0
# Generate React applications from atomic design structures
```

## Project Structure After Generation

```
output/bch-app/
├── src/
│   ├── components/
│   │   ├── atoms/       # Basic UI elements
│   │   ├── molecules/   # Combinations of atoms
│   │   └── organisms/   # Complex sections
│   ├── pages/           # Page components
│   ├── styles/          # Global styles
│   ├── App.jsx          # Main app component
│   └── main.jsx         # Entry point
├── public/              # Static assets
├── index.html           # HTML template
├── package.json         # Dependencies
└── vite.config.js       # Vite configuration
```

## Common Issues

### "go: not found"
Install Go from https://go.dev/dl/

### "npm: not found"
Install Node.js from https://nodejs.org/

### Port 3000 already in use
The generated app uses Vite which defaults to port 3000. If it's in use, Vite will automatically try the next available port.

### Build fails
Make sure you're in the atomic-generator directory and have Go installed correctly.

## Next Steps

1. **Modify the JSON**: Edit `examples/bch_complete_atomic_structure.json` to customize the generated app
2. **Add new components**: Define new atoms, molecules, or organisms in the JSON
3. **Customize styling**: Update brand colors, typography, and spacing in the project section
4. **Add pages**: Define multiple pages and routes in the JSON structure

## Support

For issues and questions:
1. Check the README.md for detailed documentation
2. Review the example JSON structure
3. Open an issue on GitHub

Happy generating! 🚀
