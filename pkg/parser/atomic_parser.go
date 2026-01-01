package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"atomic-generator/pkg/models"
)

type AtomicParser struct {
	filePath string
}

func NewAtomicParser(filePath string) *AtomicParser {
	return &AtomicParser{
		filePath: filePath,
	}
}

// Parse reads the JSON file and returns the atomic structure
func (p *AtomicParser) Parse() (*models.AtomicStructure, error) {
	// Read the file
	data, err := os.ReadFile(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", p.filePath, err)
	}

	// Parse JSON
	var structure models.AtomicStructure
	if err := json.Unmarshal(data, &structure); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	// Validate
	if err := p.validate(&structure); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	return &structure, nil
}

// validate performs basic validation on the structure
func (p *AtomicParser) validate(s *models.AtomicStructure) error {
	if s.Project.ID == "" {
		return fmt.Errorf("project.id is required")
	}
	if s.Project.Name == "" {
		return fmt.Errorf("project.name is required")
	}
	if s.Page.ID == "" {
		return fmt.Errorf("page.id is required")
	}
	if s.Layout.ID == "" {
		return fmt.Errorf("layout.id is required")
	}
	
	return nil
}

// GetAtomByID finds an atom by its ID across all atom categories
func (p *AtomicParser) GetAtomByID(structure *models.AtomicStructure, atomID string) *models.Atom {
	// Search in all atom categories
	allAtoms := [][]models.Atom{
		structure.Atoms.Images,
		structure.Atoms.Headings,
		structure.Atoms.Links,
		structure.Atoms.Buttons,
		structure.Atoms.Inputs,
		structure.Atoms.Text,
	}

	for _, atomList := range allAtoms {
		for i := range atomList {
			if atomList[i].ID == atomID {
				return &atomList[i]
			}
		}
	}
	return nil
}

// GetMoleculeByID finds a molecule by its ID
func (p *AtomicParser) GetMoleculeByID(structure *models.AtomicStructure, moleculeID string) *models.Molecule {
	for i := range structure.Molecules {
		if structure.Molecules[i].ID == moleculeID {
			return &structure.Molecules[i]
		}
	}
	return nil
}

// GetOrganismByID finds an organism by its ID
func (p *AtomicParser) GetOrganismByID(structure *models.AtomicStructure, organismID string) *models.Organism {
	for i := range structure.Organisms {
		if structure.Organisms[i].ID == organismID {
			return &structure.Organisms[i]
		}
	}
	return nil
}
