package renderers

import (
	"fmt"
	"strings"

	"atomic-generator/pkg/models"
)

// SubatomRenderer generates base React components (Image, Heading, Link, Button, etc.)
type SubatomRenderer struct {
	atom      *models.Atom
	converter *StyleConverter
}

func NewSubatomRenderer(atom *models.Atom) *SubatomRenderer {
	return &SubatomRenderer{
		atom:      atom,
		converter: NewStyleConverter(),
	}
}

// Render generates the JSX for a subatomic component
func (sr *SubatomRenderer) Render() (string, error) {
	switch sr.atom.Subatom {
	case "Image":
		return sr.renderImage()
	case "Heading":
		return sr.renderHeading()
	case "Link":
		return sr.renderLink()
	case "Button":
		return sr.renderButton()
	case "Input":
		return sr.renderInput()
	case "Text":
		return sr.renderText()
	default:
		return "", fmt.Errorf("unknown subatom type: %s", sr.atom.Subatom)
	}
}

func (sr *SubatomRenderer) renderImage() (string, error) {
	var attrs []string

	// Extract config
	if src, ok := sr.atom.Config["src"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`src="%s"`, src))
	}
	if alt, ok := sr.atom.Config["alt"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`alt="%s"`, alt))
	}
	if loading, ok := sr.atom.Config["loading"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`loading="%s"`, loading))
	}
	if width, ok := sr.atom.Config["width"]; ok {
		attrs = append(attrs, fmt.Sprintf(`width="%v"`, width))
	}
	if height, ok := sr.atom.Config["height"]; ok {
		attrs = append(attrs, fmt.Sprintf(`height="%v"`, height))
	}

	// Add styles
	if len(sr.atom.Styles) > 0 {
		styleStr := sr.converter.ToInlineStyle(sr.atom.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	return fmt.Sprintf("<img %s />", strings.Join(attrs, " ")), nil
}

func (sr *SubatomRenderer) renderHeading() (string, error) {
	level := 1
	if l, ok := sr.atom.Config["level"].(float64); ok {
		level = int(l)
	}
	
	content := ""
	if c, ok := sr.atom.Config["content"].(string); ok {
		content = c
	}

	var attrs []string
	
	// Add styles
	if len(sr.atom.Styles) > 0 {
		styleStr := sr.converter.ToInlineStyle(sr.atom.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	tag := fmt.Sprintf("h%d", level)
	if len(attrs) > 0 {
		return fmt.Sprintf("<%s %s>%s</%s>", tag, strings.Join(attrs, " "), content, tag), nil
	}
	return fmt.Sprintf("<%s>%s</%s>", tag, content, tag), nil
}

func (sr *SubatomRenderer) renderLink() (string, error) {
	var attrs []string

	if href, ok := sr.atom.Config["href"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`href="%s"`, href))
	}
	if target, ok := sr.atom.Config["target"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`target="%s"`, target))
	}
	if ariaLabel, ok := sr.atom.Config["ariaLabel"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`aria-label="%s"`, ariaLabel))
	}

	content := ""
	if c, ok := sr.atom.Config["content"].(string); ok {
		content = c
	}

	// Add styles
	if len(sr.atom.Styles) > 0 {
		styleStr := sr.converter.ToInlineStyle(sr.atom.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	return fmt.Sprintf("<a %s>%s</a>", strings.Join(attrs, " "), content), nil
}

func (sr *SubatomRenderer) renderButton() (string, error) {
	var attrs []string

	if btnType, ok := sr.atom.Config["type"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`type="%s"`, btnType))
	}
	if dataAction, ok := sr.atom.Config["dataAction"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`data-action="%s"`, dataAction))
	}
	if dataTarget, ok := sr.atom.Config["dataTarget"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`data-target="%s"`, dataTarget))
	}

	content := ""
	if c, ok := sr.atom.Config["content"].(string); ok {
		content = c
	}

	// Add styles
	if len(sr.atom.Styles) > 0 {
		styleStr := sr.converter.ToInlineStyle(sr.atom.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	return fmt.Sprintf("<button %s>%s</button>", strings.Join(attrs, " "), content), nil
}

func (sr *SubatomRenderer) renderInput() (string, error) {
	var attrs []string

	if inputType, ok := sr.atom.Config["type"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`type="%s"`, inputType))
	}
	if name, ok := sr.atom.Config["name"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`name="%s"`, name))
	}
	if placeholder, ok := sr.atom.Config["placeholder"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`placeholder="%s"`, placeholder))
	}
	if ariaLabel, ok := sr.atom.Config["ariaLabel"].(string); ok {
		attrs = append(attrs, fmt.Sprintf(`aria-label="%s"`, ariaLabel))
	}

	// Add styles
	if len(sr.atom.Styles) > 0 {
		styleStr := sr.converter.ToInlineStyle(sr.atom.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	return fmt.Sprintf("<input %s />", strings.Join(attrs, " ")), nil
}

func (sr *SubatomRenderer) renderText() (string, error) {
	content := ""
	if c, ok := sr.atom.Config["content"].(string); ok {
		content = c
	}

	tag := "span"
	if t, ok := sr.atom.Config["tag"].(string); ok {
		tag = t
	}

	var attrs []string
	
	// Add styles
	if len(sr.atom.Styles) > 0 {
		styleStr := sr.converter.ToInlineStyle(sr.atom.Styles)
		attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
	}

	if len(attrs) > 0 {
		return fmt.Sprintf("<%s %s>%s</%s>", tag, strings.Join(attrs, " "), content, tag), nil
	}
	return fmt.Sprintf("<%s>%s</%s>", tag, content, tag), nil
}

// RenderAsComponent generates a full React component for the atom
func (sr *SubatomRenderer) RenderAsComponent() (string, error) {
	componentName := ToPascalCase(sr.atom.ID)
	jsx, err := sr.Render()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`import React from 'react';

const %s = () => {
  return (
    %s
  );
};

export default %s;
`, componentName, jsx, componentName), nil
}
