# Banner Kit Go

Go port of banner-kit - SVG/PNG banner generator for OSS projects.

## Status

**Work in Progress**: SVG generation is functional. PNG conversion requires implementation.

## Requirements

**Nerd Fonts** must be installed for proper text rendering in PNG output:

```bash
# macOS with Homebrew
brew install --cask font-hack-nerd-font

# Or download from https://www.nerdfonts.com/
```

## Build

```bash
go build -o banner-gen
```

## Usage

```bash
./banner-gen <project-dir> [theme] [align]
```

**Arguments**:
- `project-dir`: Path to project directory containing README.md
- `theme`: `light|muted|dark` (default: `light`)
- `align`: `center|left|right` (default: `center`)

**Example**:
```bash
./banner-gen ./my-project dark left
```

Outputs `banner.svg` and `banner.png` in the project directory.

## README.md Format

Add HTML comment markers to your project's README.md:

```markdown
<!-- banner-title: Your Project Name -->
<!-- banner-tagline: Your project description -->
```

See the original [banner-kit](https://github.com/yourusername/banner-kit) for examples.

## TODO

- [ ] Implement PNG conversion (SVG to PNG rendering)
- [ ] Add Nerd Font detection
- [ ] Add badge support via CLI arguments
- [ ] Add tests

## Original

This is a Go port of the Node.js [banner-kit](../banner-kit).
