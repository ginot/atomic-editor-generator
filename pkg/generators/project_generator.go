package generators

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"atomic-generator/pkg/models"
	"atomic-generator/pkg/parser"
	"atomic-generator/pkg/renderers"
)

// ProjectGenerator orchestrates the generation of the complete React project
type ProjectGenerator struct {
	structure  *models.AtomicStructure
	parser     *parser.AtomicParser
	outputDir  string
}

func NewProjectGenerator(structure *models.AtomicStructure, outputDir string) *ProjectGenerator {
	return &ProjectGenerator{
		structure: structure,
		parser:    &parser.AtomicParser{},
		outputDir: outputDir,
	}
}

// Generate creates the complete React project structure
func (pg *ProjectGenerator) Generate() error {
	// Create base directories
	if err := pg.createDirectories(); err != nil {
		return fmt.Errorf("error creating directories: %w", err)
	}

	// Generate configuration files
	if err := pg.generateConfigFiles(); err != nil {
		return fmt.Errorf("error generating config files: %w", err)
	}

	// Generate global styles
	if err := pg.generateGlobalStyles(); err != nil {
		return fmt.Errorf("error generating global styles: %w", err)
	}

	// Generate atoms
	if err := pg.generateAtoms(); err != nil {
		return fmt.Errorf("error generating atoms: %w", err)
	}

	// Generate molecules
	if err := pg.generateMolecules(); err != nil {
		return fmt.Errorf("error generating molecules: %w", err)
	}

	// Generate organisms
	if err := pg.generateOrganisms(); err != nil {
		return fmt.Errorf("error generating organisms: %w", err)
	}

	// Generate pages
	if err := pg.generatePages(); err != nil {
		return fmt.Errorf("error generating pages: %w", err)
	}

	// Generate App.jsx and routing
	if err := pg.generateAppComponent(); err != nil {
		return fmt.Errorf("error generating app component: %w", err)
	}

	// Generate index files
	if err := pg.generateIndexFiles(); err != nil {
		return fmt.Errorf("error generating index files: %w", err)
	}

	fmt.Printf("✅ Project generated successfully in %s\n", pg.outputDir)
	return nil
}

func (pg *ProjectGenerator) createDirectories() error {
	dirs := []string{
		"src",
		"src/components",
		"src/components/atoms",
		"src/components/molecules",
		"src/components/organisms",
		"src/pages",
		"src/styles",
		"src/assets",
		"public",
	}

	for _, dir := range dirs {
		path := filepath.Join(pg.outputDir, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (pg *ProjectGenerator) generateConfigFiles() error {
	// Generate package.json
	packageJSON := pg.generatePackageJSON()
	if err := pg.writeFile("package.json", packageJSON); err != nil {
		return err
	}

	// Generate vite.config.js
	viteConfig := pg.generateViteConfig()
	if err := pg.writeFile("vite.config.js", viteConfig); err != nil {
		return err
	}

	// Generate .gitignore
	gitignore := `node_modules
dist
.DS_Store
*.log
.env
.env.local
`
	if err := pg.writeFile(".gitignore", gitignore); err != nil {
		return err
	}

	// Generate README.md
	readme := pg.generateReadme()
	if err := pg.writeFile("README.md", readme); err != nil {
		return err
	}

	return nil
}

func (pg *ProjectGenerator) generatePackageJSON() string {
	return fmt.Sprintf(`{
  "name": "%s",
  "version": "%s",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "lint": "eslint src"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-router-dom": "^6.20.0",
    "react-helmet-async": "^2.0.4"
  },
  "devDependencies": {
    "@types/react": "^18.2.43",
    "@types/react-dom": "^18.2.17",
    "@vitejs/plugin-react": "^4.2.1",
    "vite": "^5.0.8",
    "eslint": "^8.55.0",
    "eslint-plugin-react": "^7.33.2"
  }
}
`, pg.structure.Project.ID, pg.structure.Project.Version)
}

func (pg *ProjectGenerator) generateViteConfig() string {
	return `import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000
  }
})
`
}

func (pg *ProjectGenerator) generateReadme() string {
	return fmt.Sprintf("# %s\n\n%s\n\n## Generated with Atomic Generator\n\nThis project was automatically generated from an atomic design structure.\n\n## Getting Started\n\n```bash\n# Install dependencies\nnpm install\n\n# Run development server\nnpm run dev\n\n# Build for production\nnpm run build\n```\n\n## Project Structure\n\n- `src/components/atoms` - Basic UI elements\n- `src/components/molecules` - Combinations of atoms\n- `src/components/organisms` - Complex UI sections\n- `src/pages` - Page components\n- `src/styles` - Global styles and CSS\n\n## Version\n\n%s\n", pg.structure.Project.Name, pg.structure.Project.Name, pg.structure.Project.Version)
}

func (pg *ProjectGenerator) generateGlobalStyles() error {
	// Generate CSS variables from brand
	cssVars := pg.generateCSSVariables()
	
	// Generate global CSS
	globalCSS := fmt.Sprintf(`/* Global Styles */
:root {
%s
}

/* Reset */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: %s;
  font-size: %s;
  color: %s;
  background-color: %s;
  line-height: 1.6;
}

/* Normalize */
img {
  max-width: 100%%;
  height: auto;
}

a {
  color: inherit;
  text-decoration: none;
}

button {
  font-family: inherit;
  cursor: pointer;
}
`, cssVars, 
	pg.getCSSVar("font-family-primary"),
	pg.getCSSVar("font-size-body"),
	pg.getCSSVar("color-text"),
	pg.getCSSVar("color-background"))

	return pg.writeFile("src/styles/global.css", globalCSS)
}

func (pg *ProjectGenerator) generateCSSVariables() string {
	var vars []string

	// Colors
	for key, value := range pg.structure.Project.Brand.Colors {
		varName := strings.ReplaceAll(key, "_", "-")
		vars = append(vars, fmt.Sprintf("  --color-%s: %s;", varName, value))
	}

	// Typography - Font Families
	for key, value := range pg.structure.Project.Brand.Typography.FontFamily {
		varName := strings.ReplaceAll(key, "_", "-")
		vars = append(vars, fmt.Sprintf("  --font-family-%s: %s;", varName, value))
	}

	// Typography - Font Sizes
	for key, value := range pg.structure.Project.Brand.Typography.FontSizes {
		varName := strings.ReplaceAll(key, "_", "-")
		vars = append(vars, fmt.Sprintf("  --font-size-%s: %s;", varName, value))
	}

	// Typography - Font Weights
	for key, value := range pg.structure.Project.Brand.Typography.FontWeights {
		varName := strings.ReplaceAll(key, "_", "-")
		vars = append(vars, fmt.Sprintf("  --font-weight-%s: %v;", varName, value))
	}

	// Spacing
	for key, value := range pg.structure.Project.Brand.Spacing {
		varName := strings.ReplaceAll(key, "_", "-")
		vars = append(vars, fmt.Sprintf("  --spacing-%s: %s;", varName, value))
	}

	// Breakpoints
	for key, value := range pg.structure.Project.Brand.Breakpoints {
		varName := strings.ReplaceAll(key, "_", "-")
		vars = append(vars, fmt.Sprintf("  --breakpoint-%s: %s;", varName, value))
	}

	return strings.Join(vars, "\n")
}

func (pg *ProjectGenerator) getCSSVar(name string) string {
	return fmt.Sprintf("var(--%s)", name)
}

func (pg *ProjectGenerator) generateAtoms() error {
	allAtoms := [][]models.Atom{
		pg.structure.Atoms.Images,
		pg.structure.Atoms.Headings,
		pg.structure.Atoms.Links,
		pg.structure.Atoms.Buttons,
		pg.structure.Atoms.Inputs,
		pg.structure.Atoms.Text,
	}

	for _, atomList := range allAtoms {
		for _, atom := range atomList {
			renderer := renderers.NewAtomRenderer(&atom, pg.structure)
			component, err := renderer.RenderAsComponent()
			if err != nil {
				return err
			}

			filename := fmt.Sprintf("src/components/atoms/%s.jsx", renderers.ToPascalCase(atom.ID))
			if err := pg.writeFile(filename, component); err != nil {
				return err
			}
		}
	}

	fmt.Printf("✅ Generated %d atoms\n", pg.countAtoms())
	return nil
}

func (pg *ProjectGenerator) generateMolecules() error {
	for _, molecule := range pg.structure.Molecules {
		renderer := renderers.NewMoleculeRenderer(&molecule, pg.structure)
		component, err := renderer.RenderAsComponent()
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("src/components/molecules/%s.jsx", renderers.ToPascalCase(molecule.ID))
		if err := pg.writeFile(filename, component); err != nil {
			return err
		}
	}

	fmt.Printf("✅ Generated %d molecules\n", len(pg.structure.Molecules))
	return nil
}

func (pg *ProjectGenerator) generateOrganisms() error {
	for _, organism := range pg.structure.Organisms {
		renderer := renderers.NewOrganismRenderer(&organism, pg.structure)
		component, err := renderer.RenderAsComponent()
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("src/components/organisms/%s.jsx", renderers.ToPascalCase(organism.ID))
		if err := pg.writeFile(filename, component); err != nil {
			return err
		}
	}

	fmt.Printf("✅ Generated %d organisms\n", len(pg.structure.Organisms))
	return nil
}

func (pg *ProjectGenerator) generatePages() error {
	renderer := renderers.NewPageRenderer(&pg.structure.Page, &pg.structure.Layout, pg.structure)
	component, err := renderer.Render()
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("src/pages/%s.jsx", renderers.ToPascalCase(pg.structure.Page.ID))
	if err := pg.writeFile(filename, component); err != nil {
		return err
	}

	fmt.Printf("✅ Generated page: %s\n", pg.structure.Page.ID)
	return nil
}

func (pg *ProjectGenerator) generateAppComponent() error {
	pageName := renderers.ToPascalCase(pg.structure.Page.ID)
	route := pg.structure.Page.Route

	appComponent := fmt.Sprintf(`import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';
import %s from './pages/%s';
import './styles/global.css';

function App() {
  return (
    <HelmetProvider>
      <BrowserRouter>
        <Routes>
          <Route path="%s" element={<%s />} />
        </Routes>
      </BrowserRouter>
    </HelmetProvider>
  );
}

export default App;
`, pageName, pageName, route, pageName)

	return pg.writeFile("src/App.jsx", appComponent)
}

func (pg *ProjectGenerator) generateIndexFiles() error {
	// Generate main index.html
	indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>%s</title>
%s
  </head>
  <body>
    <div id="root"></div>
    <script type="module" src="/src/main.jsx"></script>
  </body>
</html>
`, pg.structure.Project.Name, pg.generateFontLinks())

	if err := pg.writeFile("index.html", indexHTML); err != nil {
		return err
	}

	// Generate main.jsx
	mainJSX := `import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'

ReactDOM.createRoot(document.getElementById('root')).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
`

	return pg.writeFile("src/main.jsx", mainJSX)
}

func (pg *ProjectGenerator) generateFontLinks() string {
	if pg.structure.Project.ThirdParty.Fonts == nil {
		return ""
	}

	var links []string
	for _, fontURL := range pg.structure.Project.ThirdParty.Fonts.Google {
		links = append(links, fmt.Sprintf(`    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link href="%s" rel="stylesheet" />`, fontURL))
	}

	return strings.Join(links, "\n")
}

func (pg *ProjectGenerator) countAtoms() int {
	return len(pg.structure.Atoms.Images) +
		len(pg.structure.Atoms.Headings) +
		len(pg.structure.Atoms.Links) +
		len(pg.structure.Atoms.Buttons) +
		len(pg.structure.Atoms.Inputs) +
		len(pg.structure.Atoms.Text)
}

func (pg *ProjectGenerator) writeFile(path, content string) error {
	fullPath := filepath.Join(pg.outputDir, path)
	
	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, []byte(content), 0644)
}
