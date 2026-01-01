# CHANGELOG

## Version 1.1.0 - Generic Rendering Architecture (Dec 31, 2024)

### 🎯 MAJOR REFACTOR - Data-Driven Architecture

#### Eliminated ALL hardcoded component types - Now completely generic

**Philosophy Change:**
The atomic design principle means components should assemble **automatically** from their definitions, not through hardcoded switch statements.

**Before (WRONG):**
```go
// MoleculeRenderer with hardcoded types
func (mr *MoleculeRenderer) Render() (string, error) {
    switch mr.molecule.Type {
    case "linked_image":      // ← Hardcoded behavior!
        return mr.renderLinkedImage()
    case "search_form":       // ← Hardcoded behavior!
        return mr.renderSearchForm()
    default:
        return mr.renderGenericMolecule()
    }
}
```

**After (CORRECT):**
```go
// MoleculeRenderer - completely generic
func (mr *MoleculeRenderer) Render() (string, error) {
    // ALWAYS generic - assembles atoms according to JSON definition
    return mr.renderGenericMolecule()
}
```

### 🔧 Changes Made:

#### MoleculeRenderer (pkg/renderers/molecule_renderer.go)
- ❌ Removed: `renderLinkedImage()` - hardcoded
- ❌ Removed: `renderSearchForm()` - hardcoded
- ✅ Now: Only `renderGenericMolecule()` - data-driven
- Molecules are assembled from their atom definitions
- Type field is metadata, not behavior
- Semantic HTML tags chosen intelligently (form, nav, article, div)

#### OrganismRenderer (pkg/renderers/organism_renderer.go)
- ❌ Removed: `renderSiteHeader()` - hardcoded (96 lines)
- ❌ Removed: `renderHeroSection()` - hardcoded (79 lines)
- ❌ Removed: `renderCarousel()` - hardcoded (96 lines)
- ❌ Removed: `renderSiteFooter()` - hardcoded (45 lines)
- ✅ Now: Only `renderGenericOrganism()` - data-driven
- **Reduced from 610 lines to 360 lines** (41% code reduction)

### 🏗️ New Generic Architecture:

```
JSON Definition (data)
       ↓
Generic Renderer (logic)
       ↓
React Component (output)
```

All composition is driven by the JSON structure:
- Atoms → assembled from subatoms
- Molecules → assembled from atoms
- Organisms → assembled from atoms + molecules
- Pages → assembled from organisms

### ✨ Benefits:

1. **Zero maintenance per component type** - Add new types without code changes
2. **True atomic design** - Components compose naturally
3. **Database-ready** - When moving to PostgreSQL, same generic logic applies
4. **Scalable** - BCH, BCH2, BCH3... all use same renderer
5. **No "add more types"** - The question itself was wrong

### 📊 Code Metrics:

| File | Before | After | Reduction |
|------|--------|-------|-----------|
| molecule_renderer.go | 156 lines | 82 lines | 47% |
| organism_renderer.go | 610 lines | 360 lines | 41% |
| **Total** | **766 lines** | **442 lines** | **42%** |

### 🎯 Validation:

The generator still produces the **exact same working output**, but now:
- ✅ No hardcoded types
- ✅ Completely data-driven
- ✅ Ready for PostgreSQL migration
- ✅ True atomic design principles

---

## Version 1.0.5 - Inline Styles Fix (Dec 30, 2024)

### 🐛 Bug Fix

#### Fixed React inline style syntax
**Error:**
```javascript
<footer style={ backgroundColor: 'var(--color-primary)', ... }>
// Syntax error: expected "}"
```

**Problem:**
The `ToInlineStyle()` method was generating single curly braces for inline styles:
```javascript
style={ ... }  // ❌ Invalid React syntax
```

React requires **double curly braces** for inline styles:
- First `{` = JavaScript expression delimiter
- Second `{` = JavaScript object literal

**Fix:**
Changed `ToInlineStyle()` to generate double braces:

```go
// Before
return "{ " + strings.Join(styleStrings, ", ") + " }"

// After
return "{{ " + strings.Join(styleStrings, ", ") + " }}"
```

**Generated code now:**
```javascript
<footer style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-text-light)' }}>
//             ↑↑ double braces
```

**Files Changed:**
- `pkg/renderers/base_renderer.go` - Line 34: Fixed return statement

**Verification:**
```bash
go run ./cmd/generator -input examples/bch_complete_atomic_structure.json -output output/bch-app
cd output/bch-app
npm run dev
# ✅ Components now render with valid inline styles
```

---

## Version 1.0.4 - Import Path Fix (Dec 30, 2024)

### 🐛 Bug Fix

#### Fixed incorrect import paths in generated pages
**Error:**
```
Failed to resolve import "./components/MainHeader" from "src/pages/Homepage.jsx"
```

**Problem:**
Pages are generated in `src/pages/` but were importing components with relative path `./components/`, which would look for `src/pages/components/` (doesn't exist).

Components are actually in:
- `src/components/organisms/`
- `src/components/molecules/`
- `src/components/atoms/`

**Fix:**
1. Changed import paths from `./components/ComponentName` to `../components/{type}/ComponentName`
2. Added `getComponentType()` helper to determine if component is organism, molecule, or atom
3. Generate correct relative paths based on component type

**Generated imports now:**
```javascript
// Before (incorrect)
import MainHeader from './components/MainHeader';

// After (correct)
import MainHeader from '../components/organisms/MainHeader';
import HeroBanner from '../components/organisms/HeroBanner';
import QuotesCarousel from '../components/organisms/QuotesCarousel';
```

**Files Changed:**
- `pkg/renderers/page_renderer.go` - Lines 53-61 (import generation)
- `pkg/renderers/page_renderer.go` - New function `getComponentType()` (Lines 124-152)

**Verification:**
```bash
go run ./cmd/generator -input examples/bch_complete_atomic_structure.json -output output/bch-app
cd output/bch-app
npm run dev
# ✅ Imports now resolve correctly
```

---

## Version 1.0.3 - Import Fix (Dec 30, 2024)

### 🐛 Bug Fix

#### Fixed invalid JavaScript imports for multiple organisms
**Error:**
```javascript
import HeroBanner,QuotesCarousel from './components/HeroBanner,QuotesCarousel';
// Syntax error: expected "{"
```

**Problem:**
When a page layout had multiple organisms (e.g., main section with hero_banner, quotes_carousel, etc.), the generator was concatenating component names with commas and creating a single invalid import.

**Fix:**
1. Changed `renderLayoutSection()` to return `[]string` instead of `string` for component names
2. Updated `Render()` to use a map for tracking unique component names
3. Generate separate import statements for each component

**Generated code now:**
```javascript
import HeroBanner from './components/HeroBanner';
import QuotesCarousel from './components/QuotesCarousel';
import MainFooter from './components/MainFooter';
```

**Files Changed:**
- `pkg/renderers/page_renderer.go` - Lines 29-54 (Render function)
- `pkg/renderers/page_renderer.go` - Lines 84-121 (renderLayoutSection function)

**Verification:**
```bash
go run ./cmd/generator -input examples/bch_complete_atomic_structure.json -output output/bch-app
cd output/bch-app
npm install
npm run dev
# ✅ React app now compiles and runs successfully
```

---

## Version 1.0.2 - Critical Fix (Dec 30, 2024)

### 🐛 Critical Bug Fix

#### Fixed JSON parsing error with mixed Layout types
**Error:** 
```
json: cannot unmarshal string into Go struct field Organism.organisms.layout of type map[string]interface {}
```

**Problem:** 
The JSON structure uses `layout` field in two different ways:
- In `hero_banner`: nested maps for different sections
  ```json
  "layout": {
    "background": { "position": "absolute", ... },
    "overlay": { ... },
    "content": { ... }
  }
  ```
- In `main_footer`: direct CSS properties
  ```json
  "layout": {
    "display": "grid",
    "gridTemplateColumns": "repeat(...)",
    "gap": "var(--spacing-xl)"
  }
  ```

**Fix:**
1. Changed `Organism.Layout` type from `map[string]map[string]interface{}` to `map[string]interface{}`
2. Added `getLayoutStyles()` helper function to safely extract nested style maps
3. Updated all Layout accesses to use the helper function

**Files Changed:**
- `pkg/models/models.go` - Line 136: Changed Layout type
- `pkg/renderers/organism_renderer.go` - Added getLayoutStyles() helper and updated 3 usage sites

**Verification:**
```bash
go run ./cmd/generator -input examples/bch_complete_atomic_structure.json -output output/bch-app
# ✅ Now parses successfully
```

---

## Version 1.0.1 - Hotfixes (Dec 30, 2024)

### 🐛 Bug Fixes

#### 1. Fixed unused variable in `organism_renderer.go`
**Error:** 
```
pkg\renderers\organism_renderer.go:104:5: declared and not used: overlayStr
```

**Fix:**
```go
// ANTES
if overlayStr, ok := or.organism.Config["overlay"].(string); ok {
    // overlayStr nunca se usaba
}

// DESPUÉS
if _, ok := or.organism.Config["overlay"]; ok {
    // Solo verificamos existencia, no necesitamos el valor
}
```

**File:** `pkg/renderers/organism_renderer.go` line 104

---

#### 2. Fixed escaped backticks in `project_generator.go`
**Error:**
```
pkg\generators\project_generator.go:190:3: invalid character U+005C '\'
pkg\generators\project_generator.go:190:4: syntax error: unexpected literal `\`
```

**Problem:** Backticks were escaped with `\` inside raw string literals (backtick strings)

**Fix:** Changed from raw string literal (with backticks) to double-quoted string with proper escaping

```go
// ANTES - Raw string con backticks escapados (inválido)
return fmt.Sprintf(`# %s
\`\`\`bash
npm install
\`\`\`
`, name)

// DESPUÉS - String normal con escapes correctos
return fmt.Sprintf("# %s\n```bash\nnpm install\n```\n", name)
```

**File:** `pkg/generators/project_generator.go` line 179-213

---

### ✅ Verification

Both files now compile without errors:

```bash
go build -o bin/atomic-generator ./cmd/generator
# ✅ Success

go run ./cmd/generator -version
# ✅ Shows version

go run ./cmd/generator -input examples/bch_complete_atomic_structure.json -output output/bch-app
# ✅ Generates project successfully
```

---

### 📦 Files Updated

- `pkg/renderers/organism_renderer.go`
- `pkg/generators/project_generator.go`

### 🔄 How to Update

If you already downloaded the project:

```bash
# Re-download the fixed version
# The .tar.gz and .zip files have been updated

# Or manually apply the fixes:
# 1. Edit pkg/renderers/organism_renderer.go line 104
# 2. Edit pkg/generators/project_generator.go line 179-213
```

---

## Version 1.0.0 - Initial Release (Dec 30, 2024)

### ✨ Features

- Complete atomic design generator
- React + Vite project generation
- Support for atoms, molecules, organisms, pages
- CSS variable-based theming
- Responsive design support
- Interactive components (carousels, forms)
- SEO metadata generation
- Comprehensive documentation

### 📚 Documentation

- INDEX.md - Navigation guide
- SUMMARY.md - Executive summary
- QUICKSTART.md - 5-minute setup
- README.md - Full documentation
- ARCHITECTURE.md - System design
- EXAMPLES.md - Usage examples

### 🛠️ Tools

- `run.sh` - Quick runner script
- `test.sh` - Test script
- `Makefile` - Build commands
- `go.mod` - Dependencies

---

## Notes

### Go Compiler Strictness

Go is very strict about:
- ✅ Unused variables - must use or discard with `_`
- ✅ Unused imports - will fail to compile
- ✅ Syntax errors - no warnings, only errors

This ensures high code quality but requires fixing all issues before running.

### Testing Before Release

Always test compilation:
```bash
go build ./...
# or
go run ./cmd/generator -version
```

---

**All issues resolved! Ready to use.** 🚀
