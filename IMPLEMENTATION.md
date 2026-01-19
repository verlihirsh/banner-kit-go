# Implementation Notes

## Current Status

### ✅ Completed
- Core architecture and project structure
- README.md metadata parsing
- Theme system (light, muted, dark)
- SVG template loading and variable substitution
- XML entity escaping for Nerd Font icons
- Badge visibility management
- Template embedding in binary
- Command-line argument parsing

### ⏳ Pending: PNG Conversion

The PNG conversion functionality (`convertSVGToPNG` in `generator.go`) needs implementation.

## PNG Conversion Options

Based on the original JavaScript implementation using `@resvg/resvg-js`, here are Go alternatives:

### Option 1: CGO with librsvg (Recommended)
Use `librsvg` via CGO bindings.

**Pros**:
- Native SVG support
- Font rendering support
- Mature and stable

**Cons**:
- Requires CGO (complicates cross-compilation)
- System dependencies (librsvg, Cairo)

**Libraries**:
- Raw CGO bindings to librsvg
- System installation: `brew install librsvg` (macOS)

### Option 2: Pure Go with gg + oksvg
Use `github.com/fogleman/gg` for rendering and `github.com/srwiley/oksvg` for SVG parsing.

**Pros**:
- Pure Go (no CGO)
- Cross-compilation friendly

**Cons**:
- Limited SVG feature support
- May not handle complex gradients well
- Font handling can be tricky

**Libraries**:
- `github.com/fogleman/gg` - 2D rendering
- `github.com/srwiley/oksvg` - SVG parsing
- `github.com/srwiley/rasterx` - Rasterization

### Option 3: External Command (Simplest)
Shell out to system tools like `resvg` (Rust binary) or `inkscape`.

**Pros**:
- No implementation needed
- Best SVG compatibility
- Reliable font rendering

**Cons**:
- External dependency
- Requires resvg/inkscape installed
- Platform-specific installation

**Example**:
```go
import "os/exec"

func convertSVGToPNG(svgPath, pngPath string) error {
    cmd := exec.Command("resvg", svgPath, pngPath)
    return cmd.Run()
}
```

### Option 4: Use tdewolff/canvas
Pure Go library with comprehensive SVG and rendering support.

**Pros**:
- Pure Go
- Good SVG support
- Built-in PNG rendering

**Cons**:
- Relatively new
- Font handling needs verification

**Library**: `github.com/tdewolff/canvas`

## Recommended Approach

**Phase 1** (Immediate): Use external `resvg` command
- Fastest to implement
- Best compatibility with original
- Document installation requirements

**Phase 2** (Future): Evaluate pure Go solutions
- Try `tdewolff/canvas` first
- Fall back to `gg + oksvg` if needed
- Consider CGO option if font rendering is problematic

## Font Handling Notes

The original uses Nerd Fonts loaded via:
```javascript
font: {
  loadSystemFonts: true,
  fontFiles: [...nerdFonts, emojiFont],
  defaultFontFamily: "'Hack Nerd Font'",
}
```

For Go implementation:
- Font discovery: scan `~/Library/Fonts` (macOS), `/usr/share/fonts` (Linux)
- Font loading: depends on rendering library choice
- `golang.org/x/image/font` for TrueType parsing

## Implementation Skeleton for resvg

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

func convertSVGToPNG(svg string) ([]byte, error) {
    tmpDir, err := os.MkdirTemp("", "banner-kit-*")
    if err != nil {
        return nil, fmt.Errorf("failed to create temp dir: %w", err)
    }
    defer os.RemoveAll(tmpDir)

    svgPath := filepath.Join(tmpDir, "banner.svg")
    pngPath := filepath.Join(tmpDir, "banner.png")

    if err := os.WriteFile(svgPath, []byte(svg), 0644); err != nil {
        return nil, fmt.Errorf("failed to write temp SVG: %w", err)
    }

    // Check if resvg is available
    if _, err := exec.LookPath("resvg"); err != nil {
        return nil, fmt.Errorf("resvg not found. Install with: cargo install resvg")
    }

    // Run resvg
    cmd := exec.Command("resvg", svgPath, pngPath)
    if output, err := cmd.CombinedOutput(); err != nil {
        return nil, fmt.Errorf("resvg failed: %w\nOutput: %s", err, output)
    }

    // Read PNG
    png, err := os.ReadFile(pngPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read PNG: %w", err)
    }

    return png, nil
}
```

## Testing

Test with the existing test projects:
```bash
./banner-gen ~/work/banner-kit/test-projects/nerd-icon-example light center
./banner-gen ~/work/banner-kit/test-projects/emoji-test dark left
./banner-gen ~/work/banner-kit/test-projects/vibetunnel muted right
```

Compare output with Node.js version:
```bash
node ~/work/banner-kit/scripts/generate-banners.mjs ~/work/banner-kit/test-projects/nerd-icon-example light center
```
