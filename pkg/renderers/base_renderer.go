package renderers

import (
	"fmt"
	"strings"
)

// Renderer is the base interface for all renderers
type Renderer interface {
	Render() (string, error)
}

// StyleConverter converts style maps to CSS-in-JS or inline styles
type StyleConverter struct{}

func NewStyleConverter() *StyleConverter {
	return &StyleConverter{}
}

// ToInlineStyle converts a style map to React inline style format
func (sc *StyleConverter) ToInlineStyle(styles map[string]interface{}) string {
	if len(styles) == 0 {
		return ""
	}

	var styleStrings []string
	for key, value := range styles {
		jsKey := sc.toJSProperty(key)
		jsValue := sc.formatValue(value)
		styleStrings = append(styleStrings, fmt.Sprintf("%s: %s", jsKey, jsValue))
	}

	// React inline styles need double braces: style={{ ... }}
	return "{{ " + strings.Join(styleStrings, ", ") + " }}"
}

// ToCSSModule converts styles to CSS module format
func (sc *StyleConverter) ToCSSModule(className string, styles map[string]interface{}) string {
	if len(styles) == 0 {
		return ""
	}

	var cssLines []string
	cssLines = append(cssLines, fmt.Sprintf(".%s {", className))

	for key, value := range styles {
		cssKey := sc.toCSSProperty(key)
		cssValue := sc.formatValue(value)
		cssLines = append(cssLines, fmt.Sprintf("  %s: %s;", cssKey, cssValue))
	}

	cssLines = append(cssLines, "}")
	return strings.Join(cssLines, "\n")
}

// toJSProperty converts CSS property names to camelCase for React
func (sc *StyleConverter) toJSProperty(prop string) string {
	// Convert kebab-case to camelCase
	parts := strings.Split(prop, "-")
	if len(parts) == 1 {
		return prop
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return result
}

// toCSSProperty converts camelCase to kebab-case for CSS
func (sc *StyleConverter) toCSSProperty(prop string) string {
	var result []rune
	for i, r := range prop {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '-')
			result = append(result, r+32) // convert to lowercase
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// formatValue formats a value for CSS/JS
func (sc *StyleConverter) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// Check if it's a CSS variable or function
		if strings.HasPrefix(v, "var(") || strings.HasPrefix(v, "clamp(") || 
		   strings.HasPrefix(v, "calc(") || strings.HasPrefix(v, "rgba(") ||
		   strings.HasPrefix(v, "rgb(") {
			return fmt.Sprintf("'%s'", v)
		}
		return fmt.Sprintf("'%s'", v)
	case int, int64:
		return fmt.Sprintf("%d", v)
	case float64:
		// Check if it's a whole number
		if v == float64(int(v)) {
			return fmt.Sprintf("%d", int(v))
		}
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}

// IndentCode adds indentation to code blocks
func IndentCode(code string, levels int) string {
	indent := strings.Repeat("  ", levels)
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "" {
			lines[i] = indent + line
		}
	}
	return strings.Join(lines, "\n")
}

// Capitalize capitalizes the first letter of a string
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// ToPascalCase converts a string to PascalCase
func ToPascalCase(s string) string {
	// Handle snake_case and kebab-case
	s = strings.ReplaceAll(s, "_", "-")
	parts := strings.Split(s, "-")
	
	var result string
	for _, part := range parts {
		if len(part) > 0 {
			result += Capitalize(part)
		}
	}
	return result
}
