package renderers

import (
	"fmt"

	"atomic-generator/pkg/models"
)

// AtomRenderer generates React components from Atoms
type AtomRenderer struct {
	atom      *models.Atom
	structure *models.AtomicStructure
}

func NewAtomRenderer(atom *models.Atom, structure *models.AtomicStructure) *AtomRenderer {
	return &AtomRenderer{
		atom:      atom,
		structure: structure,
	}
}

// Render generates the JSX for an atom
func (ar *AtomRenderer) Render() (string, error) {
	// Use SubatomRenderer to render the base component
	subatomRenderer := NewSubatomRenderer(ar.atom)
	return subatomRenderer.Render()
}

// RenderAsComponent generates a full React component for the atom
func (ar *AtomRenderer) RenderAsComponent() (string, error) {
	componentName := ToPascalCase(ar.atom.ID)
	
	// Generate JSX
	jsx, err := ar.Render()
	if err != nil {
		return "", err
	}

	// Generate state management if needed (for hover, focus, etc.)
	var stateCode string
	var eventHandlers string
	
	if len(ar.atom.States) > 0 {
		stateCode = ar.generateStateManagement()
		eventHandlers = ar.generateEventHandlers()
	}

	component := fmt.Sprintf(`import React, { useState } from 'react';

const %s = () => {
%s
  return (
    %s%s
  );
};

export default %s;
`, componentName, stateCode, IndentCode(jsx, 2), eventHandlers, componentName)

	return component, nil
}

// generateStateManagement creates state hooks for interactive atoms
func (ar *AtomRenderer) generateStateManagement() string {
	if len(ar.atom.States) == 0 {
		return ""
	}

	var stateVars []string
	for stateName := range ar.atom.States {
		stateVars = append(stateVars, fmt.Sprintf("  const [is%s, setIs%s] = useState(false);", 
			Capitalize(stateName), Capitalize(stateName)))
	}

	return "\n" + fmt.Sprintf("%s\n", stateVars[0]) // For now, just return first state
}

// generateEventHandlers creates event handlers for state changes
func (ar *AtomRenderer) generateEventHandlers() string {
	if len(ar.atom.States) == 0 {
		return ""
	}

	// For hover state
	if _, hasHover := ar.atom.States["hover"]; hasHover {
		return `
      onMouseEnter={() => setIsHover(true)}
      onMouseLeave={() => setIsHover(false)}`
	}

	// For focus state
	if _, hasFocus := ar.atom.States["focus"]; hasFocus {
		return `
      onFocus={() => setIsFocus(true)}
      onBlur={() => setIsFocus(false)}`
	}

	return ""
}

// RenderInline generates inline JSX without wrapping component
func (ar *AtomRenderer) RenderInline() (string, error) {
	return ar.Render()
}
