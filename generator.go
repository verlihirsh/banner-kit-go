package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	return nil, fmt.Errorf("PNG conversion not yet implemented")
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
