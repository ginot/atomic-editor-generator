#!/bin/bash

# Atomic Generator - Go Run Helper
# Ejecuta el generador directamente con go run (sin compilar)

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     ATOMIC GENERATOR (go run mode)       ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════╝${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}❌ Go no está instalado${NC}"
    echo ""
    echo "Instala Go 1.21+:"
    echo "  Ubuntu/Debian: sudo apt install golang-go"
    echo "  macOS: brew install go"
    echo "  O descarga de: https://go.dev/dl/"
    exit 1
fi

echo -e "${GREEN}✅ Go $(go version | awk '{print $3}') encontrado${NC}"
echo ""

# Parse command
case "${1:-example}" in
    version)
        echo -e "${BLUE}🔍 Mostrando versión...${NC}"
        go run ./cmd/generator -version
        ;;
    
    example)
        echo -e "${BLUE}🚀 Generando Barcelona Culinary Hub...${NC}"
        echo ""
        go run ./cmd/generator \
            -input examples/bch_complete_atomic_structure.json \
            -output output/bch-app
        
        echo ""
        echo -e "${GREEN}✅ ¡Listo! Tu app está en: output/bch-app${NC}"
        echo ""
        echo "Para correrla:"
        echo "  cd output/bch-app"
        echo "  npm install"
        echo "  npm run dev"
        ;;
    
    custom)
        if [ -z "$2" ] || [ -z "$3" ]; then
            echo -e "${YELLOW}Uso: ./run.sh custom <input.json> <output-dir>${NC}"
            exit 1
        fi
        
        echo -e "${BLUE}🚀 Generando desde $2 hacia $3...${NC}"
        echo ""
        go run ./cmd/generator -input "$2" -output "$3"
        
        echo ""
        echo -e "${GREEN}✅ ¡Listo! Tu app está en: $3${NC}"
        ;;
    
    help|*)
        echo "Uso: ./run.sh [comando]"
        echo ""
        echo "Comandos:"
        echo "  version    - Muestra la versión del generador"
        echo "  example    - Genera el ejemplo BCH (por defecto)"
        echo "  custom     - Genera desde un JSON personalizado"
        echo "  help       - Muestra esta ayuda"
        echo ""
        echo "Ejemplos:"
        echo "  ./run.sh                                    # Genera ejemplo"
        echo "  ./run.sh version                            # Ver versión"
        echo "  ./run.sh custom mi-app.json ./mi-output    # Custom"
        ;;
esac
