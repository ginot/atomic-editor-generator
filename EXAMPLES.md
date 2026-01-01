# Atomic Generator - Examples & Extension Guide

## 📚 Complete Examples

### Example 1: Simple Landing Page

```json
{
  "project": {
    "id": "simple-landing",
    "name": "Simple Landing Page",
    "version": "1.0.0",
    "brand": {
      "colors": {
        "primary": "#2563eb",
        "secondary": "#10b981",
        "text": "#1f2937",
        "background": "#ffffff"
      },
      "typography": {
        "fontFamily": {
          "primary": "'Inter', sans-serif"
        },
        "fontSizes": {
          "h1": "3rem",
          "body": "1rem"
        }
      },
      "spacing": {
        "sm": "0.5rem",
        "md": "1rem",
        "lg": "2rem"
      }
    }
  },
  "atoms": {
    "headings": [
      {
        "id": "hero_title",
        "subatom": "Heading",
        "config": {
          "level": 1,
          "content": "Welcome to Our Product"
        },
        "styles": {
          "fontSize": "var(--font-size-h1)",
          "color": "var(--color-primary)",
          "textAlign": "center"
        }
      }
    ],
    "buttons": [
      {
        "id": "cta_button",
        "subatom": "Button",
        "config": {
          "content": "Get Started",
          "type": "button"
        },
        "styles": {
          "backgroundColor": "var(--color-primary)",
          "color": "#ffffff",
          "padding": "1rem 2rem",
          "borderRadius": "0.5rem"
        }
      }
    ]
  },
  "organisms": [
    {
      "id": "hero_section",
      "type": "hero_section",
      "atoms": {
        "heading": "hero_title",
        "button": "cta_button"
      },
      "styles": {
        "display": "flex",
        "flexDirection": "column",
        "alignItems": "center",
        "padding": "4rem 2rem",
        "minHeight": "80vh"
      }
    }
  ],
  "page": {
    "id": "home",
    "route": "/",
    "title": "Home - Simple Landing"
  },
  "layout": {
    "id": "main_layout",
    "structure": [
      {
        "section": "main",
        "organism": "hero_section"
      }
    ]
  }
}
```

Generate it:
```bash
./bin/atomic-generator -input simple-landing.json -output ./simple-landing-app
```

### Example 2: Blog with Multiple Organisms

```json
{
  "project": {
    "id": "blog",
    "name": "My Blog",
    "version": "1.0.0"
  },
  "atoms": {
    "images": [
      {
        "id": "post_thumbnail",
        "subatom": "Image",
        "config": {
          "src": "/images/post.jpg",
          "alt": "Post thumbnail",
          "loading": "lazy"
        }
      }
    ],
    "headings": [
      {
        "id": "post_title",
        "subatom": "Heading",
        "config": {
          "level": 2,
          "content": "Blog Post Title"
        }
      }
    ],
    "text": [
      {
        "id": "post_excerpt",
        "subatom": "Text",
        "config": {
          "tag": "p",
          "content": "This is a blog post excerpt..."
        }
      }
    ]
  },
  "molecules": [
    {
      "id": "post_card",
      "type": "card",
      "atoms": {
        "image": "post_thumbnail",
        "title": "post_title",
        "excerpt": "post_excerpt"
      },
      "styles": {
        "border": "1px solid #e5e7eb",
        "borderRadius": "0.5rem",
        "padding": "1rem",
        "maxWidth": "400px"
      }
    }
  ],
  "organisms": [
    {
      "id": "blog_grid",
      "type": "grid",
      "molecules": ["post_card", "post_card", "post_card"],
      "styles": {
        "display": "grid",
        "gridTemplateColumns": "repeat(auto-fit, minmax(300px, 1fr))",
        "gap": "2rem",
        "padding": "2rem"
      }
    }
  ]
}
```

## 🔧 Extending the Generator

### Adding a New Subatom Type

1. **Update models.go** (if needed for new config fields)

2. **Add rendering logic to subatom_renderer.go**:

```go
func (sr *SubatomRenderer) renderVideo() (string, error) {
    var attrs []string

    if src, ok := sr.atom.Config["src"].(string); ok {
        attrs = append(attrs, fmt.Sprintf(`src="%s"`, src))
    }
    if controls, ok := sr.atom.Config["controls"].(bool); ok && controls {
        attrs = append(attrs, "controls")
    }
    if autoplay, ok := sr.atom.Config["autoplay"].(bool); ok && autoplay {
        attrs = append(attrs, "autoplay")
    }

    if len(sr.atom.Styles) > 0 {
        styleStr := sr.converter.ToInlineStyle(sr.atom.Styles)
        attrs = append(attrs, fmt.Sprintf("style=%s", styleStr))
    }

    return fmt.Sprintf("<video %s />", strings.Join(attrs, " ")), nil
}
```

3. **Update the Render switch statement**:

```go
func (sr *SubatomRenderer) Render() (string, error) {
    switch sr.atom.Subatom {
    case "Image":
        return sr.renderImage()
    case "Video":  // NEW
        return sr.renderVideo()
    // ... other cases
    default:
        return "", fmt.Errorf("unknown subatom type: %s", sr.atom.Subatom)
    }
}
```

4. **Use it in JSON**:

```json
{
  "atoms": {
    "videos": [
      {
        "id": "intro_video",
        "subatom": "Video",
        "config": {
          "src": "/videos/intro.mp4",
          "controls": true,
          "autoplay": false
        },
        "styles": {
          "width": "100%",
          "maxWidth": "800px"
        }
      }
    ]
  }
}
```

### Adding a New Organism Type

1. **Add rendering logic to organism_renderer.go**:

```go
func (or *OrganismRenderer) renderPricingTable() (string, error) {
    // Custom rendering logic for pricing tables
    var pricingCards []string

    if molecules, ok := or.organism.Molecules.([]interface{}); ok {
        for _, molID := range molecules {
            if molIDStr, ok := molID.(string); ok {
                molecule := or.parser.GetMoleculeByID(or.structure, molIDStr)
                if molecule != nil {
                    renderer := NewMoleculeRenderer(molecule, or.structure)
                    jsx, err := renderer.Render()
                    if err != nil {
                        return "", err
                    }
                    pricingCards = append(pricingCards, jsx)
                }
            }
        }
    }

    // Build pricing table layout
    cardsJSX := strings.Join(pricingCards, "\n      ")
    
    return fmt.Sprintf(`<div className="pricing-table">
      %s
    </div>`, cardsJSX), nil
}
```

2. **Update the Render switch**:

```go
func (or *OrganismRenderer) Render() (string, error) {
    switch or.organism.Type {
    case "pricing_table":  // NEW
        return or.renderPricingTable()
    // ... other cases
    default:
        return or.renderGenericOrganism()
    }
}
```

### Adding Custom Behaviors

Add interactive behaviors to organisms:

```json
{
  "organisms": [
    {
      "id": "tabs_section",
      "type": "tabs",
      "behavior": {
        "type": "tabs",
        "defaultTab": 0,
        "animated": true,
        "transition": "fade"
      },
      "molecules": ["tab_1", "tab_2", "tab_3"]
    }
  ]
}
```

Implement in organism_renderer.go:

```go
func (or *OrganismRenderer) renderTabs() (string, error) {
    // Get behavior config
    behavior := or.organism.Behavior
    defaultTab := 0
    if behavior != nil && behavior.DefaultTab != nil {
        defaultTab = *behavior.DefaultTab
    }

    // Generate state management
    stateCode := fmt.Sprintf(`
  const [activeTab, setActiveTab] = useState(%d);
`, defaultTab)

    // Generate tab content
    var tabContents []string
    if molecules, ok := or.organism.Molecules.([]interface{}); ok {
        for i, molID := range molecules {
            // ... render each tab
        }
    }

    // Return component with tabs
    return fmt.Sprintf(`<div className="tabs-container">
      <div className="tabs-header">
        {/* Tab buttons */}
      </div>
      <div className="tabs-content">
        %s
      </div>
    </div>`, strings.Join(tabContents, "\n"))
}
```

## 🎨 Advanced Styling Patterns

### Using CSS Modules

Generate CSS module files alongside components:

```go
// In component_generator.go
func (cg *ComponentGenerator) GenerateWithCSSModule(component string, styles map[string]interface{}) error {
    // Generate .module.css file
    cssContent := cg.converter.ToCSSModule("container", styles)
    
    // Write CSS file
    cssPath := fmt.Sprintf("%s.module.css", componentName)
    if err := cg.writeFile(cssPath, cssContent); err != nil {
        return err
    }
    
    // Update component to import CSS module
    component = fmt.Sprintf(`import styles from './%s.module.css';\n\n%s`, 
        componentName, component)
    
    return nil
}
```

### Theme Support

Add theme switching:

```json
{
  "project": {
    "brand": {
      "themes": {
        "light": {
          "colors": {
            "background": "#ffffff",
            "text": "#000000"
          }
        },
        "dark": {
          "colors": {
            "background": "#1a1a1a",
            "text": "#ffffff"
          }
        }
      }
    }
  }
}
```

## 🚀 Production Optimization

### Code Splitting

Generate dynamic imports for routes:

```go
func (pg *PageRenderer) RenderWithCodeSplitting() (string, error) {
    return fmt.Sprintf(`
import React, { lazy, Suspense } from 'react';

const %s = lazy(() => import('./pages/%s'));

export default function App() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <%s />
    </Suspense>
  );
}
`, componentName, componentName, componentName), nil
}
```

### Image Optimization

Add responsive image support:

```json
{
  "atoms": {
    "images": [
      {
        "id": "hero_image",
        "subatom": "Image",
        "config": {
          "src": "/images/hero.jpg",
          "srcSet": "/images/hero-400.jpg 400w, /images/hero-800.jpg 800w",
          "sizes": "(max-width: 600px) 400px, 800px"
        }
      }
    ]
  }
}
```

## 📦 Integration Examples

### With TypeScript

Generate TypeScript files instead of JSX:

```go
// Add flag in main.go
typescript := flag.Bool("typescript", false, "Generate TypeScript files")

// In renderers
extension := ".jsx"
if *typescript {
    extension = ".tsx"
}

// Add type definitions
func (ar *AtomRenderer) RenderAsTypeScript() (string, error) {
    return fmt.Sprintf(`import React from 'react';

interface %sProps {
  className?: string;
  onClick?: () => void;
}

const %s: React.FC<%sProps> = ({ className, onClick }) => {
  return (%s);
};

export default %s;
`, componentName, componentName, componentName, jsx, componentName), nil
}
```

### With Styled Components

Generate styled-components instead of inline styles:

```go
func (sr *StyleConverter) ToStyledComponent(styles map[string]interface{}) string {
    var cssLines []string
    for key, value := range styles {
        cssLines = append(cssLines, fmt.Sprintf("  %s: %s;", 
            sr.toCSSProperty(key), sr.formatValue(value)))
    }
    return strings.Join(cssLines, "\n")
}

// In component generation
styledComponent := fmt.Sprintf(`
import styled from 'styled-components';

const StyledContainer = styled.div\`
%s
\`;
`, sc.ToStyledComponent(styles))
```

## 🧪 Testing Generated Code

Add test generation:

```go
func (pg *ProjectGenerator) generateTests() error {
    for _, organism := range pg.structure.Organisms {
        componentName := ToPascalCase(organism.ID)
        
        testContent := fmt.Sprintf(`
import { render, screen } from '@testing-library/react';
import %s from '../%s';

describe('%s', () => {
  it('renders without crashing', () => {
    render(<%s />);
  });
});
`, componentName, componentName, componentName, componentName)

        filename := fmt.Sprintf("src/components/organisms/%s.test.jsx", componentName)
        if err := pg.writeFile(filename, testContent); err != nil {
            return err
        }
    }
    return nil
}
```

## 📖 JSON Schema Definition

For IDE autocomplete and validation, create a JSON schema:

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Atomic Structure",
  "type": "object",
  "required": ["project", "page", "layout"],
  "properties": {
    "project": {
      "type": "object",
      "required": ["id", "name", "version"],
      "properties": {
        "id": { "type": "string" },
        "name": { "type": "string" },
        "version": { "type": "string" }
      }
    }
  }
}
```

Save as `atomic-schema.json` and reference in your IDE.

## 🎯 Best Practices

1. **Keep atoms simple** - One responsibility per atom
2. **Reuse molecules** - DRY principle
3. **Organisms should be self-contained** - Include all necessary state
4. **Use CSS variables** - For consistent theming
5. **Plan the hierarchy** - Before writing JSON

## 💡 Tips & Tricks

### Tip 1: Preview Mode
Add a preview flag that generates a demo page:

```bash
./bin/atomic-generator -input structure.json -output ./app -preview
```

### Tip 2: Watch Mode
Add file watching for rapid development:

```bash
./bin/atomic-generator -input structure.json -output ./app -watch
```

### Tip 3: Partial Generation
Generate only specific components:

```bash
./bin/atomic-generator -input structure.json -output ./app -only=organisms
```

---

**"No se puede saltar un abismo en dos saltos" - Generate everything at once, the Un Palo way! 🚀**
