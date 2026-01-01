package renderers

import (
	"fmt"
	"strings"

	"atomic-generator/pkg/models"
	"atomic-generator/pkg/parser"
)

// MoleculeRenderer generates React components from Molecules
type MoleculeRenderer struct {
	molecule  *models.Molecule
	structure *models.AtomicStructure
	parser    *parser.AtomicParser
	converter *StyleConverter
}

func NewMoleculeRenderer(molecule *models.Molecule, structure *models.AtomicStructure) *MoleculeRenderer {
	return &MoleculeRenderer{
		molecule:  molecule,
		structure: structure,
		parser:    &parser.AtomicParser{},
		converter: NewStyleConverter(),
	}
}

// Render generates the JSX for a molecule - completely generic
func (mr *MoleculeRenderer) Render() (string, error) {
	// All molecules are rendered generically based on their atom composition
	return mr.renderGenericMolecule()
}

func (mr *MoleculeRenderer) renderGenericMolecule() (string, error) {
	var children []string

	// Render all atoms in the molecule
	for atomKey, atomID := range mr.molecule.Atoms {
		atom := mr.parser.GetAtomByID(mr.structure, atomID)
		if atom != nil {
			renderer := NewAtomRenderer(atom, mr.structure)
			jsx, err := renderer.Render()
			if err != nil {
				return "", fmt.Errorf("error rendering atom %s in molecule %s: %w", atomID, mr.molecule.ID, err)
			}
			
			// Wrap with a key for React if needed (for lists)
			children = append(children, jsx)
		} else {
			// Log warning but continue - atom might be optional
			fmt.Printf("Warning: atom %s (key: %s) not found in molecule %s\n", atomID, atomKey, mr.molecule.ID)
		}
	}

	if len(children) == 0 {
		// Empty molecule - return empty div
		return "<div></div>", nil
	}

	// Build wrapper with molecule's styles
	var attrs []string
	if len(mr.molecule.Styles) > 0 {
		styleStr := mr.converter.ToInlineStyle(mr.molecule.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	// Add className if specified
	if mr.molecule.Type != "" {
		attrs = append(attrs, fmt.Sprintf(`className="molecule-%s"`, mr.molecule.Type))
	}

	wrapperAttrs := ""
	if len(attrs) > 0 {
		wrapperAttrs = " " + strings.Join(attrs, " ")
	}

	childrenJSX := strings.Join(children, "\n      ")

	// Use semantic HTML tag if it makes sense, otherwise div
	tag := mr.getSemanticTag()

	return fmt.Sprintf(`<%s%s>
      %s
    </%s>`, tag, wrapperAttrs, childrenJSX, tag), nil
}

// getSemanticTag returns appropriate HTML tag based on molecule type
func (mr *MoleculeRenderer) getSemanticTag() string {
	// Map common types to semantic HTML
	switch mr.molecule.Type {
	case "search_form", "contact_form", "login_form":
		return "form"
	case "navigation", "nav_menu":
		return "nav"
	case "article_card", "content_card":
		return "article"
	default:
		return "div"
	}
}

// RenderAsComponent generates a full React component for the molecule
func (mr *MoleculeRenderer) RenderAsComponent() (string, error) {
	componentName := ToPascalCase(mr.molecule.ID)
	
	// Generate JSX
	jsx, err := mr.Render()
	if err != nil {
		return "", err
	}

	// Generate state management for responsive behavior
	var stateCode string
	var hookCode string
	
	if len(mr.molecule.Responsive) > 0 {
		stateCode = `  const [breakpoint, setBreakpoint] = useState('desktop');

  useEffect(() => {
    const updateBreakpoint = () => {
      const width = window.innerWidth;
      if (width < 768) setBreakpoint('mobile');
      else if (width < 1024) setBreakpoint('tablet');
      else setBreakpoint('desktop');
    };

    updateBreakpoint();
    window.addEventListener('resize', updateBreakpoint);
    return () => window.removeEventListener('resize', updateBreakpoint);
  }, []);

`
		hookCode = ", useEffect"
	}

	component := fmt.Sprintf(`import React, { useState%s } from 'react';

const %s = () => {
%s  return (
    %s
  );
};

export default %s;
`, hookCode, componentName, stateCode, IndentCode(jsx, 2), componentName)

	return component, nil
}
