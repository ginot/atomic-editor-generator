package renderers

import (
	"fmt"
	"strings"

	"atomic-generator/pkg/models"
	"atomic-generator/pkg/parser"
)

// PageRenderer generates a complete React page from layout definition
type PageRenderer struct {
	page      *models.Page
	layout    *models.Layout
	structure *models.AtomicStructure
	parser    *parser.AtomicParser
}

func NewPageRenderer(page *models.Page, layout *models.Layout, structure *models.AtomicStructure) *PageRenderer {
	return &PageRenderer{
		page:      page,
		layout:    layout,
		structure: structure,
		parser:    &parser.AtomicParser{},
	}
}

// Render generates the complete page component
func (pr *PageRenderer) Render() (string, error) {
	var imports []string
	var sections []string

	// Track which components we need to import
	componentImports := make(map[string]bool)

	// Process each layout section
	for _, layoutSection := range pr.layout.Structure {
		sectionJSX, componentNames, err := pr.renderLayoutSection(layoutSection)
		if err != nil {
			return "", err
		}

		// Add all component names to imports (componentNames can be multiple)
		for _, name := range componentNames {
			if name != "" {
				componentImports[name] = true
			}
		}

		sections = append(sections, sectionJSX)
	}

	// Build imports - one per component with correct paths
	for componentName := range componentImports {
		// Determine component type (organism, molecule, or atom)
		componentType := pr.getComponentType(componentName)
		
		// Pages are in src/pages/, components in src/components/
		importPath := fmt.Sprintf("../components/%s/%s", componentType, componentName)
		imports = append(imports, fmt.Sprintf("import %s from '%s';", componentName, importPath))
	}

	importsStr := ""
	if len(imports) > 0 {
		importsStr = strings.Join(imports, "\n") + "\n\n"
	}

	// Build page component
	componentName := ToPascalCase(pr.page.ID)
	sectionsJSX := strings.Join(sections, "\n    ")

	// Generate metadata component
	metaComponent := pr.generateMetadata()

	component := fmt.Sprintf(`import React from 'react';
%s%s

const %s = () => {
  return (
    <>
      %s
    </>
  );
};

export default %s;
`, importsStr, metaComponent, componentName, sectionsJSX, componentName)

	return component, nil
}

func (pr *PageRenderer) renderLayoutSection(section models.LayoutSection) (string, []string, error) {
	var jsx string
	var componentNames []string

	// Single organism
	if section.Organism != "" {
		organism := pr.parser.GetOrganismByID(pr.structure, section.Organism)
		if organism == nil {
			return "", nil, fmt.Errorf("organism not found: %s", section.Organism)
		}

		componentName := ToPascalCase(organism.ID)
		componentNames = append(componentNames, componentName)
		jsx = fmt.Sprintf("<%s />", componentName)
	}

	// Multiple organisms
	if len(section.Organisms) > 0 {
		var organismComponents []string

		for _, orgID := range section.Organisms {
			organism := pr.parser.GetOrganismByID(pr.structure, orgID)
			if organism != nil {
				compName := ToPascalCase(organism.ID)
				organismComponents = append(organismComponents, fmt.Sprintf("<%s />", compName))
				componentNames = append(componentNames, compName)
			}
		}

		jsx = fmt.Sprintf(`<main>
      %s
    </main>`, strings.Join(organismComponents, "\n      "))
	}

	return jsx, componentNames, nil
}

// getComponentType determines if a component is an organism, molecule, or atom
func (pr *PageRenderer) getComponentType(componentName string) string {
	// Convert PascalCase component name back to ID format (snake_case)
	// MainHeader -> main_header
	componentID := ""
	for i, r := range componentName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			componentID += "_"
		}
		componentID += string(r)
	}
	componentID = strings.ToLower(componentID)
	
	// Check if it's an organism
	for _, organism := range pr.structure.Organisms {
		if organism.ID == componentID {
			return "organisms"
		}
	}
	
	// Check if it's a molecule
	for _, molecule := range pr.structure.Molecules {
		if molecule.ID == componentID {
			return "molecules"
		}
	}
	
	// Default to atoms (though pages typically use organisms)
	return "atoms"
}

func (pr *PageRenderer) generateMetadata() string {
	// Generate Helmet or react-helmet-async metadata
	var metaTags []string

	metaTags = append(metaTags, fmt.Sprintf(`<title>%s</title>`, pr.page.Title))

	if desc, ok := pr.page.Meta["description"]; ok {
		metaTags = append(metaTags, fmt.Sprintf(`<meta name="description" content="%s" />`, desc))
	}

	if keywords, ok := pr.page.Meta["keywords"]; ok {
		metaTags = append(metaTags, fmt.Sprintf(`<meta name="keywords" content="%s" />`, keywords))
	}

	if ogImage, ok := pr.page.Meta["ogImage"]; ok {
		metaTags = append(metaTags, fmt.Sprintf(`<meta property="og:image" content="%s" />`, ogImage))
	}

	if lang, ok := pr.page.Meta["language"]; ok {
		metaTags = append(metaTags, fmt.Sprintf(`<meta httpEquiv="content-language" content="%s" />`, lang))
	}

	metaTagsStr := strings.Join(metaTags, "\n        ")

	return fmt.Sprintf(`import { Helmet } from 'react-helmet-async';

const PageMetadata = () => (
  <Helmet>
    %s
  </Helmet>
);
`, metaTagsStr)
}

// RenderWithRouter generates the page component with React Router integration
func (pr *PageRenderer) RenderWithRouter() (string, error) {
	component, err := pr.Render()
	if err != nil {
		return "", err
	}

	// Wrap with route information
	route := pr.page.Route
	componentName := ToPascalCase(pr.page.ID)

	routeConfig := fmt.Sprintf(`// Route configuration for %s
export const %sRoute = {
  path: '%s',
  component: %s,
  title: '%s'
};
`, componentName, componentName, route, componentName, pr.page.Title)

	return component + "\n\n" + routeConfig, nil
}
