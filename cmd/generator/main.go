package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"atomic-generator/pkg/generators"
	"atomic-generator/pkg/models"
	"atomic-generator/pkg/parser"
)

const version = "1.0.0"

func main() {
	// Parse command line flags
	inputFile := flag.String("input", "", "Path to atomic structure JSON file")
	outputDir := flag.String("output", "./output", "Output directory for generated project")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("Atomic Generator v%s\n", version)
		fmt.Println("Generate React applications from atomic design structures")
		os.Exit(0)
	}

	// Validate input
	if *inputFile == "" {
		fmt.Println("Error: input file is required")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("  atomic-generator -input structure.json -output ./my-app")
		os.Exit(1)
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		log.Fatalf("Error: input file does not exist: %s", *inputFile)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	// Print banner
	printBanner()

	// Parse atomic structure
	fmt.Printf("📖 Reading atomic structure from: %s\n", *inputFile)
	atomicParser := parser.NewAtomicParser(*inputFile)
	structure, err := atomicParser.Parse()
	if err != nil {
		log.Fatalf("Error parsing atomic structure: %v", err)
	}

	fmt.Printf("✅ Parsed structure: %s v%s\n", structure.Project.Name, structure.Project.Version)
	fmt.Printf("   - Atoms: %d\n", countAtoms(structure))
	fmt.Printf("   - Molecules: %d\n", len(structure.Molecules))
	fmt.Printf("   - Organisms: %d\n", len(structure.Organisms))
	fmt.Printf("   - Pages: 1\n\n")

	// Generate project
	absOutputDir, err := filepath.Abs(*outputDir)
	if err != nil {
		log.Fatalf("Error resolving output directory: %v", err)
	}

	fmt.Printf("🚀 Generating React project in: %s\n\n", absOutputDir)
	
	projectGenerator := generators.NewProjectGenerator(structure, absOutputDir)
	if err := projectGenerator.Generate(); err != nil {
		log.Fatalf("Error generating project: %v", err)
	}

	// Print success message
	printSuccess(absOutputDir, structure.Project.Name)
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════╗
║                                           ║
║         ATOMIC GENERATOR                  ║
║         React from JSON                   ║
║         v` + version + `                            ║
║                                           ║
╚═══════════════════════════════════════════╝
`
	fmt.Println(banner)
}

func printSuccess(outputDir, projectName string) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("🎉 PROJECT GENERATED SUCCESSFULLY!")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("\nProject: %s\n", projectName)
	fmt.Printf("Location: %s\n", outputDir)
	fmt.Println("\nNext steps:")
	fmt.Printf("\n  1. cd %s\n", outputDir)
	fmt.Println("  2. npm install")
	fmt.Println("  3. npm run dev")
	fmt.Println("\nYour React application will be running at http://localhost:3000")
	fmt.Println(strings.Repeat("=", 50) + "\n")
}

func countAtoms(structure *models.AtomicStructure) int {
	return len(structure.Atoms.Images) +
		len(structure.Atoms.Headings) +
		len(structure.Atoms.Links) +
		len(structure.Atoms.Buttons) +
		len(structure.Atoms.Inputs) +
		len(structure.Atoms.Text)
}
