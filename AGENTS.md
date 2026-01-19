# Agent Guidelines

## Build, Lint, and Test Commands

### Build
```bash
# Standard build
go build -o banner-gen

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o banner-gen

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o banner-gen-linux
GOOS=windows GOARCH=amd64 go build -o banner-gen.exe
GOOS=darwin GOARCH=arm64 go build -o banner-gen-mac
GOOS=darwin GOARCH=amd64 go build -o banner-gen-mac-intel
```

### Lint
```bash
# Format code
go fmt ./...

# Static analysis
go vet ./...

# Check for common mistakes
go vet -all ./...
```

### Test
```bash
# Manual test with sample project
mkdir -p test-project
echo '<!-- banner-title: Test Project -->' > test-project/README.md
echo '<!-- banner-tagline: Testing banner generation -->' >> test-project/README.md
./banner-gen test-project

# Test all themes
./banner-gen test-project light center
./banner-gen test-project muted center
./banner-gen test-project dark center

# Test all alignments
./banner-gen test-project light left
./banner-gen test-project light center
./banner-gen test-project light right

# Verify PNG output
file test-project/banner.png  # Should show: PNG image data, 1600 x 600
```

### Run Single Test
```bash
# Run from source with specific arguments
go run . <project-dir> [theme] [align]

# Example
go run . ./test-project dark left
```

---

## Code Style Guidelines

### File Organization
- **main.go**: CLI entry point and argument parsing
- **generator.go**: SVG generation and PNG conversion logic
- **metadata.go**: README.md parsing
- **template.go**: Theme system and SVG template manipulation
- **templates/**: Embedded SVG templates (via `go:embed`)

### Imports
Standard library imports first, then external packages, separated by blank line:
```go
import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/kanrichan/resvg-go"
)
```

### Formatting
- **Formatter**: Always run `go fmt ./...` before commit
- **Indentation**: Tabs (Go standard)
- **Line length**: No strict limit, but keep it reasonable (~120 chars)
- **Brackets**: Opening brace on same line (Go standard)

### Types
- Never use type assertions without checking: `val, ok := x.(Type)`
- Prefer explicit types in function signatures
- Use type aliases for complex types: `type ThemePalette struct {...}`

### Naming Conventions
- **Variables/functions**: `camelCase` (Go standard)
- **Exported functions**: `PascalCase`
- **Constants**: `PascalCase` for exported, `camelCase` for internal
- **Files**: `kebab-case.go` or `snake_case.go`

### Error Handling
```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to parse metadata: %w", err)
}

// Never ignore errors
_ = someFunc()  // ❌ BAD

if err := someFunc(); err != nil {  // ✅ GOOD
    return err
}

// Never use empty catch blocks
if err != nil {
    // Handle or propagate
}
```

### Comments
- Use comments to explain "why", not "what"
- Document exported functions with GoDoc format:
```go
// generateSVG creates an SVG banner from metadata and theme.
// Returns the SVG as a string or an error if generation fails.
func generateSVG(metadata *Metadata, theme *ThemePalette, align string) (string, error) {
```

---

## Project-Specific Guidelines

### Architecture

#### Core Components
1. **CLI Layer** (`main.go`): Argument parsing and error handling
2. **Metadata Parser** (`metadata.go`): Extract banner-title and banner-tagline from README
3. **Theme Engine** (`template.go`): Color palette system (light/muted/dark)
4. **Template Engine** (`generator.go`): SVG variable substitution
5. **PNG Converter** (`generator.go`): Dual rendering (rsvg-convert + WASM fallback)

#### Key Design Patterns
- **Embedded Resources**: Templates are embedded via `go:embed` directive
- **Fail-fast**: Errors propagate up to main() for CLI handling
- **Zero Global State**: All config passed explicitly
- **Graceful Degradation**: PNG conversion falls back to WASM if rsvg not available

### Dependencies

**Current Dependencies**:
- **[github.com/kanrichan/resvg-go](https://github.com/kanrichan/resvg-go)**: Pure Go WASM SVG renderer (fallback)
- **[github.com/tetratelabs/wazero](https://github.com/tetratelabs/wazero)**: WebAssembly runtime

**External Tools** (optional):
- **rsvg-convert** (from librsvg): System tool for fast PNG rendering
- **Nerd Fonts**: For icon rendering in output

**When adding dependencies**:
- Prefer stdlib over external packages
- Avoid CGO dependencies (keep it pure Go)
- Check license compatibility (MIT/Apache-2.0/BSD)
- Verify cross-platform support

### Testing

**Manual Testing Pattern**:
```bash
# Create test project
mkdir -p test-dir
echo '<!-- banner-title: 🚀 Test Project -->' > test-dir/README.md
echo '<!-- banner-tagline: Testing all features -->' >> test-dir/README.md

# Test command
./banner-gen test-dir [theme] [align]

# Verify output
ls -lh test-dir/banner.{svg,png}
file test-dir/banner.png  # Should show PNG 1600x600
```

**What to Test**:
- All three themes: light, muted, dark
- All three alignments: center, left, right
- Unicode handling (Nerd Fonts, emoji)
- Missing README.md (should error)
- Missing banner-title (should error)
- Optional banner-tagline (should work)

### Git Workflow

- Write clear commit messages: `feat: add dark theme support`
- Use conventional commits: `fix:`, `feat:`, `docs:`, `refactor:`
- Keep commits focused and atomic
- Run `go fmt` and `go vet` before committing

---

## Frontend Guidelines (N/A)

This is a CLI tool with no frontend components.

---

## Backend Guidelines

### SVG Generation

**Template Variables**:
```
{{BG0}}           - Background gradient start
{{BG1}}           - Background gradient end
{{BG2}}           - Text background color
{{WAVE0}}         - Wave gradient start
{{WAVE1}}         - Wave gradient end
{{PROJECT_NAME}}  - Project title
{{TAGLINE}}       - Project tagline
{{BADGE_1}}       - First badge (unused)
{{BADGE_2}}       - Second badge (unused)
{{BADGE_3}}       - Third badge (unused)
```

**XML Escaping**:
- Characters `< > & " '` must be escaped
- Unicode > 127 converted to numeric entities: `&#XXXX;`
- See `escapeXML()` in `template.go`

### PNG Conversion

**Dual Rendering Strategy**:
1. **Primary**: Check for `rsvg-convert` via `exec.LookPath`
2. **Fallback**: Use `resvg-go` WASM renderer

**rsvg-convert Usage**:
```bash
rsvg-convert -f png < input.svg > output.png
```

**WASM Renderer Usage**:
```go
ctx, _ := resvg.NewContext(context.Background())
defer ctx.Close()
renderer, _ := ctx.NewRenderer()
defer renderer.Close()
renderer.LoadSystemFonts()
png, _ := renderer.Render(svgData)
```

---

## Additional Rules

### Performance Considerations

- **Startup**: ~10ms (vs. ~100ms Node.js)
- **Memory**: <10MB peak
- **Binary Size**: ~3MB with embedded templates
- **PNG Conversion**: rsvg-convert is 2-3x faster than WASM

### Cross-Compilation

Always test on target platform:
```bash
# Build for Linux from macOS
GOOS=linux GOARCH=amd64 go build -o banner-gen-linux

# Test in Docker
docker run --rm -v $(pwd):/app golang:1.21 /app/banner-gen-linux /app/test-project
```

### Font Handling

- **System Fonts**: Loaded automatically by WASM renderer
- **Nerd Fonts**: Must be installed system-wide
- **Locations**:
  - macOS: `~/Library/Fonts`, `/Library/Fonts`, `/System/Library/Fonts`
  - Linux: `/usr/share/fonts`, `~/.fonts`, `~/.local/share/fonts`
  - Windows: `C:\Windows\Fonts`, `%USERPROFILE%\AppData\Local\Microsoft\Windows\Fonts`

---

## Notes for AI Agents

### Key Invariants
- **NO CGO**: Keep project pure Go for easy cross-compilation
- **NO Runtime Filesystem Deps**: All templates embedded via `go:embed`
- **Fail-Fast Error Handling**: Errors propagate to main(), never silently fail
- **Graceful Degradation**: PNG conversion has automatic fallback

### Before Marking Task Complete
1. Run `go fmt ./...`
2. Run `go vet ./...`
3. Build successfully: `go build -o banner-gen`
4. Manual test: `./banner-gen test-project light center`
5. Verify PNG output: `file test-project/banner.png`

### Common Pitfalls
- **Don't add CGO dependencies** - breaks cross-compilation
- **Don't use external config files** - keep everything embedded
- **Don't ignore conversion errors** - always return error from `convertSVGToPNG`
- **Don't assume rsvg-convert exists** - always check with `exec.LookPath`

### When in Doubt
- Check existing code patterns in `generator.go` and `template.go`
- Follow Go standard library idioms
- Prefer explicitness over cleverness
- Test manually before claiming success
