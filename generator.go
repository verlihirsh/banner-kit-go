package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kanrichan/resvg-go"
)

//go:embed templates/*.svg
var templateFS embed.FS

func loadTemplate(align string) (string, error) {
	templatePath := fmt.Sprintf("templates/banner.%s.svg", align)

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to load template %s: %w", templatePath, err)
	}

	return string(data), nil
}

func generateSVG(metadata *Metadata, theme *ThemePalette, align string, badges []string) (string, error) {
	template, err := loadTemplate(align)
	if err != nil {
		return "", err
	}

	vars := map[string]string{
		"BG0":          theme.BG0,
		"BG1":          theme.BG1,
		"BG2":          theme.BG2,
		"WAVE0":        theme.WAVE0,
		"WAVE1":        theme.WAVE1,
		"PROJECT_NAME": metadata.Name,
		"TAGLINE":      metadata.Tagline,
		"BADGE_1":      "",
		"BADGE_2":      "",
		"BADGE_3":      "",
	}

	if len(badges) > 0 {
		vars["BADGE_1"] = badges[0]
	}
	if len(badges) > 1 {
		vars["BADGE_2"] = badges[1]
	}
	if len(badges) > 2 {
		vars["BADGE_3"] = badges[2]
	}

	svg := replaceVariables(template, vars)

	if strings.TrimSpace(vars["BADGE_1"]) == "" {
		svg = stripBadge(svg, 1)
	}
	if strings.TrimSpace(vars["BADGE_2"]) == "" {
		svg = stripBadge(svg, 2)
	}
	if strings.TrimSpace(vars["BADGE_3"]) == "" {
		svg = stripBadge(svg, 3)
	}

	return svg, nil
}

func convertSVGToPNG(svg string) ([]byte, error) {
	// Try rsvg-convert first (fastest and most reliable)
	if _, err := exec.LookPath("rsvg-convert"); err == nil {
		return convertWithRsvgConvert([]byte(svg))
	}

	// Fallback to resvg-go (pure Go, no external dependencies)
	return convertWithResvg([]byte(svg))
}

func convertWithRsvgConvert(svgData []byte) ([]byte, error) {
	cmd := exec.Command("rsvg-convert", "-f", "png")
	cmd.Stdin = bytes.NewReader(svgData)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("rsvg-convert failed: %w (stderr: %s)", err, stderr.String())
	}

	return out.Bytes(), nil
}

func convertWithResvg(svgData []byte) ([]byte, error) {
	ctx := context.Background()
	resvgCtx, err := resvg.NewContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create resvg context: %w", err)
	}
	defer resvgCtx.Close()

	renderer, err := resvgCtx.NewRenderer()
	if err != nil {
		return nil, fmt.Errorf("failed to create renderer: %w", err)
	}
	defer renderer.Close()

	if err := renderer.LoadSystemFonts(); err != nil {
		return nil, fmt.Errorf("failed to load system fonts: %w", err)
	}

	pngData, err := renderer.Render(svgData)
	if err != nil {
		return nil, fmt.Errorf("resvg rendering failed: %w", err)
	}

	return pngData, nil
}

func writeBannerFiles(projectDir, svg string, png []byte) error {
	svgPath := filepath.Join(projectDir, "banner.svg")
	pngPath := filepath.Join(projectDir, "banner.png")

	if err := os.WriteFile(svgPath, []byte(svg), 0644); err != nil {
		return fmt.Errorf("failed to write SVG: %w", err)
	}
	fmt.Printf("Generated: %s\n", svgPath)

	if png != nil {
		if err := os.WriteFile(pngPath, png, 0644); err != nil {
			return fmt.Errorf("failed to write PNG: %w", err)
		}
		fmt.Printf("Generated: %s\n", pngPath)
	}

	return nil
}
