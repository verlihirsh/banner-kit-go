package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Metadata struct {
	Name    string
	Tagline string
}

func parseReadmeMetadata(content string) (*Metadata, error) {
	titleRe := regexp.MustCompile(`<!--\s*banner-title:\s*(.+?)\s*-->`)
	taglineRe := regexp.MustCompile(`<!--\s*banner-tagline:\s*(.+?)\s*-->`)

	titleMatch := titleRe.FindStringSubmatch(content)
	taglineMatch := taglineRe.FindStringSubmatch(content)

	if titleMatch == nil {
		return nil, fmt.Errorf("no banner-title found in README.md")
	}

	metadata := &Metadata{
		Name:    strings.TrimSpace(titleMatch[1]),
		Tagline: "",
	}

	if taglineMatch != nil {
		metadata.Tagline = strings.TrimSpace(taglineMatch[1])
	}

	return metadata, nil
}

func readProjectMetadata(projectDir string) (*Metadata, error) {
	readmePath := filepath.Join(projectDir, "README.md")

	content, err := os.ReadFile(readmePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read README.md from %s: %w", projectDir, err)
	}

	return parseReadmeMetadata(string(content))
}
