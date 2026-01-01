# Atomic Generator

Generate complete React applications from atomic design JSON structures.

## 🚀 Features

- **Atomic Design**: Follows Brad Frost's atomic design methodology
- **Complete Project Generation**: Creates full React project with Vite, routing, and styling
- **Type-Safe**: Written in Go for reliability and performance
- **Flexible**: Supports custom atoms, molecules, organisms, and pages
- **Production-Ready**: Generates optimized, modern React code

## 📦 Installation

### From Source

```bash
# Clone the repository
git clone <repository-url>
cd atomic-generator

# Build the generator
make build

# Or install to your GOPATH
make install
```

### Pre-built Binaries

Download the latest release for your platform from the releases page.

## 🎯 Quick Start

### 1. Create an Atomic Structure JSON

Create a JSON file defining your application structure:

```json
{
  "project": {
    "id": "myapp",
    "name": "My App",
    "version": "1.0.0",
    "brand": {
      "colors": {
        "primary": "#007bff",
        "secondary": "#6c757d"
      }
    }
  },
  "atoms": {
    "headings": [
      {
        "id": "main_heading",
        "subatom": "Heading",
        "config": {
          "level": 1,
          "content": "Welcome"
        }
      }
    ]
  }
}
```

### 2. Generate Your App

```bash
# Using the binary
./bin/atomic-generator -input structure.json -output ./my-app

# Or using make
make example
```

### 3. Run Your App

```bash
cd my-app
npm install
npm run dev
```

## 📖 Structure Overview

### Atomic Hierarchy

```
Project
  └── Clump (group of pages)
      └── Page
          └── Layout
              └── Organisms (complex sections)
                  └── Molecules (combinations)
                      └── Atoms (basic elements)
                          └── Subatoms (HTML elements)
```

### Supported Subatoms

- **Image**: `<img>` elements with responsive support
- **Heading**: `<h1>` through `<h6>` elements
- **Link**: `<a>` elements with routing
- **Button**: `<button>` elements with event handlers
- **Input**: `<input>` elements with validation
- **Text**: `<span>`, `<p>`, and text containers

### Supported Organisms

- **site_header**: Navigation headers with sticky/scroll behavior
- **hero_section**: Hero banners with background images and overlays
- **carousel**: Image/content carousels with autoplay
- **site_footer**: Multi-section footers with navigation

## 🎨 Brand System

Define your brand once, use everywhere:

```json
{
  "brand": {
    "colors": {
      "primary": "#1a1a1a",
      "secondary": "#d4af37"
    },
    "typography": {
      "fontFamily": {
        "primary": "'Montserrat', sans-serif"
      },
      "fontSizes": {
        "h1": "clamp(2rem, 5vw, 4rem)"
      }
    },
    "spacing": {
      "md": "1rem",
      "lg": "2rem"
    },
    "breakpoints": {
      "mobile": "320px",
      "tablet": "768px",
      "desktop": "1024px"
    }
  }
}
```

All values are converted to CSS variables for easy theming.

## 🔧 Advanced Features

### Responsive Molecules

```json
{
  "id": "logo_link",
  "type": "linked_image",
  "responsive": [
    {
      "breakpoint": "desktop",
      "atoms": { "image": "logo_large" }
    },
    {
      "breakpoint": "mobile",
      "atoms": { "image": "logo_small" }
    }
  ]
}
```

### Interactive Behaviors

```json
{
  "id": "carousel",
  "type": "carousel",
  "behavior": {
    "type": "carousel",
    "autoplay": true,
    "interval": 5000,
    "loop": true,
    "controls": true
  }
}
```

### Event Handlers

```json
{
  "events": {
    "onSubmit": {
      "action": "navigate",
      "target": "/search",
      "method": "GET"
    }
  }
}
```

## 📂 Generated Project Structure

```
my-app/
├── src/
│   ├── components/
│   │   ├── atoms/
│   │   │   ├── HeroHeading.jsx
│   │   │   ├── LogoIsologotipo.jsx
│   │   │   └── ...
│   │   ├── molecules/
│   │   │   ├── LogoLink.jsx
│   │   │   ├── SearchBox.jsx
│   │   │   └── ...
│   │   └── organisms/
│   │       ├── MainHeader.jsx
│   │       ├── HeroBanner.jsx
│   │       └── ...
│   ├── pages/
│   │   └── Homepage.jsx
│   ├── styles/
│   │   └── global.css
│   ├── App.jsx
│   └── main.jsx
├── public/
├── index.html
├── package.json
├── vite.config.js
└── README.md
```

## 🛠 CLI Usage

```bash
# Basic usage
atomic-generator -input structure.json -output ./my-app

# Show version
atomic-generator -version

# Show help
atomic-generator -help
```

### Flags

- `-input`: Path to atomic structure JSON file (required)
- `-output`: Output directory for generated project (default: `./output`)
- `-version`: Show version information

## 🧪 Example

A complete example is included in `examples/bch_complete_atomic_structure.json`.

Generate it with:

```bash
make example
```

This will create a full Barcelona Culinary Hub website with:
- Responsive header with navigation
- Hero section with background image
- Multiple content sections
- Interactive carousel
- Multi-section footer

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

MIT License - feel free to use this in your own projects.

## 🎓 Learn More

- [Atomic Design by Brad Frost](https://atomicdesign.bradfrost.com/)
- [React Documentation](https://react.dev/)
- [Vite Documentation](https://vitejs.dev/)

## ⚡ Performance

The generator is fast:
- Parses complex JSON structures in milliseconds
- Generates complete projects with 100+ components in under a second
- Produces optimized, tree-shakeable React code

## 🐛 Troubleshooting

### Generator fails to parse JSON

Make sure your JSON is valid and follows the atomic structure schema.

### Generated app won't start

1. Make sure you've run `npm install`
2. Check that all required dependencies are installed
3. Verify Node.js version (requires Node 18+)

### Components not rendering

Check the browser console for errors and verify all organism/molecule/atom IDs are correctly referenced.

## 📧 Contact

For questions and support, please open an issue on GitHub.

---

Built with ❤️ using Go and atomic design principles.
