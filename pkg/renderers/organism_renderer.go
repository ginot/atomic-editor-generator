package renderers

import (
	"fmt"
	"strings"

	"atomic-generator/pkg/models"
	"atomic-generator/pkg/parser"
)

// OrganismRenderer generates React components from Organisms
type OrganismRenderer struct {
	organism  *models.Organism
	structure *models.AtomicStructure
	parser    *parser.AtomicParser
	converter *StyleConverter
}

func NewOrganismRenderer(organism *models.Organism, structure *models.AtomicStructure) *OrganismRenderer {
	return &OrganismRenderer{
		organism:  organism,
		structure: structure,
		parser:    &parser.AtomicParser{},
		converter: NewStyleConverter(),
	}
}

// getLayoutStyles safely extracts styles from layout field
// Layout can contain either direct styles or nested style maps
func (or *OrganismRenderer) getLayoutStyles(key string) map[string]interface{} {
	if or.organism.Layout == nil {
		return nil
	}
	
	value, ok := or.organism.Layout[key]
	if !ok {
		return nil
	}
	
	// If value is a map, return it
	if stylesMap, ok := value.(map[string]interface{}); ok {
		return stylesMap
	}
	
	// Otherwise return nil (value is a string or other type)
	return nil
}

// Render generates the JSX for an organism - completely generic and data-driven
func (or *OrganismRenderer) Render() (string, error) {
	// All organisms are rendered generically based on their composition
	return or.renderGenericOrganism()
}

func (or *OrganismRenderer) renderGenericOrganism() (string, error) {
	var elements []string

	// 1. Render atoms if present
	if len(or.organism.Atoms) > 0 {
		atomElements, err := or.renderAtoms()
		if err != nil {
			return "", err
		}
		elements = append(elements, atomElements...)
	}

	// 2. Render molecules if present
	if or.organism.Molecules != nil {
		moleculeElements, err := or.renderMolecules()
		if err != nil {
			return "", err
		}
		elements = append(elements, moleculeElements...)
	}

	// 3. Render sections if present (for complex organisms like footers)
	if len(or.organism.Sections) > 0 {
		sectionElements, err := or.renderSections()
		if err != nil {
			return "", err
		}
		elements = append(elements, sectionElements...)
	}

	// 4. If layout is specified, apply it
	if len(or.organism.Layout) > 0 {
		elements = or.applyLayout(elements)
	}

	// 5. Build wrapper with organism's styles
	var attrs []string
	if len(or.organism.Styles) > 0 {
		styleStr := or.converter.ToInlineStyle(or.organism.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	// Add className based on type
	if or.organism.Type != "" {
		attrs = append(attrs, fmt.Sprintf(`className="organism-%s"`, or.organism.Type))
	}

	wrapperAttrs := ""
	if len(attrs) > 0 {
		wrapperAttrs = " " + strings.Join(attrs, " ")
	}

	// Determine semantic tag
	tag := or.getSemanticTag()

	// Join elements
	elementsJSX := strings.Join(elements, "\n      ")

	return fmt.Sprintf(`<%s%s>
      %s
    </%s>`, tag, wrapperAttrs, elementsJSX, tag), nil
}

// renderAtoms renders all atoms in the organism
func (or *OrganismRenderer) renderAtoms() ([]string, error) {
	var elements []string

	for atomKey, atomID := range or.organism.Atoms {
		atom := or.parser.GetAtomByID(or.structure, atomID)
		if atom != nil {
			renderer := NewAtomRenderer(atom, or.structure)
			jsx, err := renderer.Render()
			if err != nil {
				return nil, fmt.Errorf("error rendering atom %s (key: %s): %w", atomID, atomKey, err)
			}
			elements = append(elements, jsx)
		} else {
			fmt.Printf("Warning: atom %s (key: %s) not found in organism %s\n", atomID, atomKey, or.organism.ID)
		}
	}

	return elements, nil
}

// renderMolecules renders all molecules in the organism
func (or *OrganismRenderer) renderMolecules() ([]string, error) {
	var elements []string

	// Molecules can be a map or an array
	switch molecules := or.organism.Molecules.(type) {
	case map[string]interface{}:
		// Map of molecules
		for molKey, molID := range molecules {
			if molIDStr, ok := molID.(string); ok {
				molecule := or.parser.GetMoleculeByID(or.structure, molIDStr)
				if molecule != nil {
					renderer := NewMoleculeRenderer(molecule, or.structure)
					jsx, err := renderer.Render()
					if err != nil {
						return nil, fmt.Errorf("error rendering molecule %s (key: %s): %w", molIDStr, molKey, err)
					}
					elements = append(elements, jsx)
				}
			}
		}
	case []interface{}:
		// Array of molecules (e.g., carousel items)
		for i, molID := range molecules {
			if molIDStr, ok := molID.(string); ok {
				molecule := or.parser.GetMoleculeByID(or.structure, molIDStr)
				if molecule != nil {
					renderer := NewMoleculeRenderer(molecule, or.structure)
					jsx, err := renderer.Render()
					if err != nil {
						return nil, fmt.Errorf("error rendering molecule %s at index %d: %w", molIDStr, i, err)
					}
					elements = append(elements, jsx)
				}
			}
		}
	}

	return elements, nil
}

// renderSections renders organism sections (used in complex organisms like footers)
func (or *OrganismRenderer) renderSections() ([]string, error) {
	var sections []string

	for _, section := range or.organism.Sections {
		var sectionMolecules []string

		for _, molID := range section.Molecules {
			molecule := or.parser.GetMoleculeByID(or.structure, molID)
			if molecule != nil {
				renderer := NewMoleculeRenderer(molecule, or.structure)
				jsx, err := renderer.Render()
				if err != nil {
					return nil, fmt.Errorf("error rendering molecule %s in section: %w", molID, err)
				}
				sectionMolecules = append(sectionMolecules, jsx)
			}
		}

		sectionJSX := strings.Join(sectionMolecules, "\n        ")
		sections = append(sections, fmt.Sprintf(`<div className="section-%s">
        %s
      </div>`, section.Type, sectionJSX))
	}

	return sections, nil
}

// applyLayout wraps elements according to layout specification
func (or *OrganismRenderer) applyLayout(elements []string) []string {
	// If layout specifies zones (like background, overlay, content), wrap accordingly
	// This is a simplified version - can be extended based on layout structure
	
	// For now, just wrap elements in layout containers if layout keys suggest it
	if len(or.organism.Layout) > 0 {
		// Check for common layout patterns
		hasBackground := or.getLayoutStyles("background") != nil
		hasOverlay := or.getLayoutStyles("overlay") != nil
		hasContent := or.getLayoutStyles("content") != nil

		if hasBackground || hasOverlay || hasContent {
			// This is a complex layout (like hero) - wrap elements
			var wrapped []string

			// Background layer
			if hasBackground {
				bgStyles := or.converter.ToInlineStyle(or.getLayoutStyles("background"))
				wrapped = append(wrapped, fmt.Sprintf(`<div style=%s className="layout-background"></div>`, bgStyles))
			}

			// Overlay layer
			if hasOverlay {
				overlayStyles := or.converter.ToInlineStyle(or.getLayoutStyles("overlay"))
				wrapped = append(wrapped, fmt.Sprintf(`<div style=%s className="layout-overlay"></div>`, overlayStyles))
			}

			// Content layer with elements
			if hasContent {
				contentStyles := or.converter.ToInlineStyle(or.getLayoutStyles("content"))
				contentJSX := strings.Join(elements, "\n          ")
				wrapped = append(wrapped, fmt.Sprintf(`<div style=%s className="layout-content">
          %s
        </div>`, contentStyles, contentJSX))
			} else {
				// No content wrapper, just add elements
				wrapped = append(wrapped, elements...)
			}

			return wrapped
		}
	}

	// No special layout, return elements as-is
	return elements
}

// getSemanticTag returns appropriate HTML tag based on organism type
func (or *OrganismRenderer) getSemanticTag() string {
	// Map common types to semantic HTML
	switch or.organism.Type {
	case "site_header", "page_header":
		return "header"
	case "site_footer", "page_footer":
		return "footer"
	case "hero_section", "content_section":
		return "section"
	case "navigation", "nav_menu":
		return "nav"
	case "article_content":
		return "article"
	case "sidebar":
		return "aside"
	default:
		return "div"
	}
}

// RenderAsComponent generates a full React component for the organism
func (or *OrganismRenderer) RenderAsComponent() (string, error) {
	componentName := ToPascalCase(or.organism.ID)
	
	// Generate JSX
	jsx, err := or.Render()
	if err != nil {
		return "", err
	}

	// Generate state and effects for interactive organisms
	var imports []string
	var stateCode string
	var effectCode string

	imports = append(imports, "useState")

	// Add state for carousel
	if or.organism.Behavior != nil && or.organism.Behavior.Type == "carousel" {
		imports = append(imports, "useEffect")
		moleculeCount := 0
		if moleculeList, ok := or.organism.Molecules.([]interface{}); ok {
			moleculeCount = len(moleculeList)
		}
		
		stateCode = `  const [currentSlide, setCurrentSlide] = useState(0);
`
		
		if or.organism.Behavior.Autoplay && moleculeCount > 0 {
			effectCode = fmt.Sprintf(`
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentSlide((prev) => (prev + 1) %% %d);
    }, %d);
    return () => clearInterval(interval);
  }, []);

  const nextSlide = () => setCurrentSlide((prev) => (prev + 1) %% %d);
  const prevSlide = () => setCurrentSlide((prev) => (prev - 1 + %d) %% %d);
`, moleculeCount, or.organism.Behavior.Interval, moleculeCount, moleculeCount, moleculeCount)
		}
	}

	// Add state for header scroll behavior
	if or.organism.Type == "site_header" && or.organism.Config["scrollBehavior"] != nil {
		imports = append(imports, "useEffect")
		stateCode += `  const [scrolled, setScrolled] = useState(false);
`
		effectCode += `
  useEffect(() => {
    const handleScroll = () => {
      setScrolled(window.scrollY > 50);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);
`
	}

	importStr := strings.Join(imports, ", ")

	component := fmt.Sprintf(`import React, { %s } from 'react';

const %s = () => {
%s%s  return (
    %s
  );
};

export default %s;
`, importStr, componentName, stateCode, effectCode, IndentCode(jsx, 2), componentName)

	return component, nil
}
