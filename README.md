<!-- banner-title: 🚀 Banner Kit Go -->
<!-- banner-tagline: SVG/PNG banner generator for OSS projects -->

<img src="banner.svg" >

<p align="center">
  <strong>Generate beautiful SVG and PNG banners for your open-source projects</strong><br>
  Fast, portable, zero-dependency Go port of banner-kit
</p>

<p align="center">
  <a href="https://golang.org"><img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go" alt="Go 1.21+"></a>
  <a href="#features"><img src="https://img.shields.io/badge/Dual_Rendering-rsvg+WASM-blue" alt="Dual Rendering"></a>
  <a href="#features"><img src="https://img.shields.io/badge/Themes-3-purple" alt="3 Themes"></a>
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#development">Development</a>
</p>

---

## Why Banner Kit Go?

Transform your project README from plain text to eye-catching banners with zero hassle.

**Banner Kit Go** gives you:
- **🚀 Instant Banners**: Generate both SVG and PNG from simple HTML comments in your README
- **⚡ Smart Rendering**: Auto-detects `rsvg-convert` for speed, falls back to pure Go WASM
- **🎨 Beautiful Themes**: Three trans-pride inspired color palettes (light, muted, dark)
- **📦 Zero Runtime Deps**: Single binary with embedded templates, no npm/node_modules
- **🔧 Cross-platform**: Works on macOS, Linux, Windows (WSL)

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Features](#features)
- [Usage](#usage)
- [Development](#development)
- [Architecture](#architecture)

---

## Installation

### Option 1: Build from Source

**Requirements**: Go 1.21 or higher

```bash
# Clone the repository
git clone https://github.com/chhlga/banner-kit-go.git
cd banner-kit-go

# Build
go build -o banner-gen

# Optional: Install system-wide
sudo mv banner-gen /usr/local/bin/
```

### Option 2: Download Binary

Download pre-built binaries from [Releases](https://github.com/chhlga/banner-kit-go/releases) (coming soon).

### Optional Dependencies

#### rsvg-convert (Recommended for Best Performance)

For optimal PNG conversion speed and quality:

```bash
# macOS with Homebrew
brew install librsvg

# Debian/Ubuntu
apt-get install librsvg2-bin

# Fedora/RHEL
dnf install librsvg2-tools

# Arch Linux
pacman -S librsvg

# Alpine (Docker)
apk add librsvg
```

**Note**: If `rsvg-convert` is not installed, the tool automatically falls back to a pure Go WASM-based renderer (slightly slower but requires no external dependencies).

#### Nerd Fonts (For Icon Support)

For proper icon rendering in banners:

```bash
# macOS with Homebrew
brew install --cask font-hack-nerd-font

# Or download from https://www.nerdfonts.com/
```

---

## Quick Start

### 1. Add Banner Metadata to Your README

Add HTML comment markers anywhere in your project's `README.md`:

```markdown
<!-- banner-title: Your Project Name -->
<!-- banner-tagline: Your awesome project description -->
```

### 2. Generate Banners

```bash
# Generate with default settings (light theme, center alignment)
banner-gen ./my-project

# With custom theme and alignment
banner-gen ./my-project dark left

# All theme options
banner-gen ./my-project light   # Default: soft pastels
banner-gen ./my-project muted   # Subtle colors
banner-gen ./my-project dark    # Dark mode
```

### 3. Output

Generates two files in your project directory:
- `banner.svg` - Scalable vector graphic
- `banner.png` - 1600x600 PNG raster image

---

## Features

### 🚀 Core Features

- **Dual Rendering**: Automatically uses `rsvg-convert` if available, falls back to pure Go WASM renderer
- **Zero CGO**: Works as a single static binary with no C dependencies
- **Three Themes**: light, muted, dark - all with trans-pride inspired color palettes
- **Flexible Alignment**: center, left, right text alignment options
- **Nerd Font Support**: Proper rendering of Nerd Font icons and Unicode characters
- **Embedded Templates**: SVG templates bundled in binary via `go:embed`

### 🎨 Theme Examples

| Theme | Use Case | Colors |
|-------|----------|--------|
| **light** | Bright, welcoming projects | Soft pastels |
| **muted** | Professional, subtle look | Muted tones |
| **dark** | Dark mode UIs, modern projects | Deep colors |

### 🚀 Performance

- **Startup**: ~10ms (vs. ~100ms for Node.js version)
- **Binary Size**: ~3MB (includes embedded templates)
- **Memory**: <10MB peak usage
- **Cross-compilation**: Build for any platform from anywhere

---

## Usage

### Basic Command

```bash
banner-gen <project-dir> [theme] [align]
```

**Arguments**:
- `project-dir`: Path to project directory containing README.md (required)
- `theme`: `light|muted|dark` (default: `light`)
- `align`: `center|left|right` (default: `center`)

### Examples

#### Example 1: Default Settings

```bash
# Generate with light theme and center alignment
banner-gen ./my-awesome-project

# Output:
# Generated: ./my-awesome-project/banner.svg
# Generated: ./my-awesome-project/banner.png
```

#### Example 2: Dark Theme, Left Aligned

```bash
banner-gen ./my-project dark left
```

#### Example 3: All Theme Variations

```bash
# Light theme (default)
banner-gen ./project light center

# Muted theme for professional look
banner-gen ./project muted center

# Dark theme for modern aesthetic
banner-gen ./project dark center
```

### README.md Format

Your `README.md` must contain these HTML comment markers:

```markdown
<!-- banner-title: Project Name -->
<!-- banner-tagline: Optional one-line description -->
```

**Example**:
```markdown
<!-- banner-title: 🚀 My Amazing Tool -->
<!-- banner-tagline: Supercharge your workflow with zero configuration -->

# My Amazing Tool

Your regular README content goes here...
```

---

## Development

### Prerequisites

- **Go 1.21+** (required for development)
- **librsvg** (optional, for PNG rendering tests)
- **Nerd Fonts** (optional, for icon rendering tests)

### Project Structure

```
banner-kit-go/
├── main.go              # CLI entry point and argument parsing
├── generator.go         # SVG generation and PNG conversion logic
├── metadata.go          # README.md parsing for banner metadata
├── template.go          # Theme system and SVG template manipulation
├── templates/           # Embedded SVG templates
│   ├── banner.center.svg
│   ├── banner.left.svg
│   └── banner.right.svg
├── go.mod               # Go module definition
└── README.md            # This file
```

### Build & Test

```bash
# Install dependencies (pure Go, no CGO)
go mod download

# Build project
go build -o banner-gen

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o banner-gen-linux
GOOS=windows GOARCH=amd64 go build -o banner-gen.exe
GOOS=darwin GOARCH=arm64 go build -o banner-gen-mac

# Test with sample project
mkdir test-project
echo '<!-- banner-title: Test Project -->' > test-project/README.md
echo '<!-- banner-tagline: Testing banners -->' >> test-project/README.md
./banner-gen test-project
```

### Running Development Build

```bash
# Standard run
go run . ./test-project

# With different themes
go run . ./test-project dark left
go run . ./test-project muted right

# Format code before commit
go fmt ./...

# Check for issues
go vet ./...
```

### Key Dependencies

- **[github.com/kanrichan/resvg-go](https://github.com/kanrichan/resvg-go)**: Pure Go WASM-based SVG renderer (fallback)
- **[github.com/tetratelabs/wazero](https://github.com/tetratelabs/wazero)**: WebAssembly runtime (used by resvg-go)

### Coding Guidelines

- **`go fmt`** before commit
- Keep imports organized: stdlib → external → internal
- Wrap errors with context; never ignore errors
- Avoid global state; prefer explicit config/structs
- Follow existing code style and patterns
- No CGO dependencies (pure Go only)

---

## Architecture

### High-Level Design

```
┌─────────────────────────────────────────────────────────┐
│                   Banner Kit Go                         │
│                                                         │
│  ┌──────────────┐          ┌──────────────┐           │
│  │   CLI Args   │          │  README.md   │           │
│  │   Parser     │          │   Parser     │           │
│  └──────┬───────┘          └──────┬───────┘           │
│         │                         │                    │
│         └──────────┬──────────────┘                    │
│                    │                                   │
│         ┌──────────▼──────────┐                        │
│         │   Theme Engine      │                        │
│         │  (Color Palettes)   │                        │
│         └──────────┬──────────┘                        │
│                    │                                   │
│         ┌──────────▼──────────┐                        │
│         │  Template Engine    │                        │
│         │  (Embedded SVGs)    │                        │
│         └──────────┬──────────┘                        │
│                    │                                   │
│         ┌──────────▼──────────┐                        │
│         │   PNG Converter     │                        │
│         │  (rsvg or WASM)     │                        │
│         └──────────┬──────────┘                        │
│                    │                                   │
└────────────────────┼───────────────────────────────────┘
                     │
                     │ Output Files
                     │
          ┌──────────▼──────────┐
          │   banner.svg        │
          │   banner.png        │
          └─────────────────────┘
```

### Key Design Decisions

1. **Embedded Templates**: Used `go:embed` instead of filesystem reads for portability - single binary contains all templates
2. **Dual PNG Rendering**: Primary method uses system `rsvg-convert` for speed, automatically falls back to pure Go WASM renderer
3. **Zero CGO**: Avoided CGO to keep cross-compilation simple and deployment easy
4. **Standard Library First**: Minimized external dependencies - only 2 external packages (WASM renderer + runtime)
5. **Clean Error Propagation**: All errors bubble up to main() for consistent CLI error handling

### Data Flow

```
README.md → Parse Metadata → Apply Theme → Substitute Template → Generate SVG
                                                                        ↓
                                                                   PNG Converter
                                                                   (rsvg or WASM)
                                                                        ↓
                                                                  Write Files
```

---

## Acknowledgments

Built using:
- **[resvg-go](https://github.com/kanrichan/resvg-go)** - Pure Go WASM SVG renderer
- **[wazero](https://github.com/tetratelabs/wazero)** - Zero-dependency WebAssembly runtime
- **[librsvg](https://gitlab.gnome.org/GNOME/librsvg)** - Optional fast SVG renderer

Inspired by the original **[banner-kit](https://github.com/chhlga/banner-kit)** (Node.js version).

---

## License

This project is licensed under the [MIT License](LICENSE).

---

<p align="center">
  <strong>Create beautiful banners in seconds 🚀</strong>
</p>
