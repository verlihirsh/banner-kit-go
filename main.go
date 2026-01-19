package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <project-dir> [theme] [align]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  project-dir   Path to project directory containing README.md\n")
		fmt.Fprintf(os.Stderr, "  theme         Theme name: light|muted|dark (default: light)\n")
		fmt.Fprintf(os.Stderr, "  align         Alignment: center|left|right (default: center)\n\n")
		fmt.Fprintf(os.Stderr, "Example:\n")
		fmt.Fprintf(os.Stderr, "  %s ./my-project dark left\n", os.Args[0])
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	projectDir := args[0]
	theme := "light"
	align := "center"

	if len(args) > 1 {
		theme = args[1]
	}
	if len(args) > 2 {
		align = args[2]
	}

	if err := generateBanner(projectDir, theme, align); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func generateBanner(projectDir, themeStr, align string) error {
	theme, err := getTheme(themeStr)
	if err != nil {
		return err
	}

	metadata, err := readProjectMetadata(projectDir)
	if err != nil {
		return err
	}

	svg, err := generateSVG(metadata, theme, align, []string{})
	if err != nil {
		return err
	}

	png, err := convertSVGToPNG(svg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		png = nil
	}

	return writeBannerFiles(projectDir, svg, png)
}
