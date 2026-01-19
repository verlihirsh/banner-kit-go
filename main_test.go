package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateBanner(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("full banner generation with light theme", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-1")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: Test Project -->
<!-- banner-tagline: A comprehensive test -->
# Test Project`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "light", "center")
		require.NoError(t, err)

		svgPath := filepath.Join(projectDir, "banner.svg")
		_, err = os.Stat(svgPath)
		assert.NoError(t, err)

		svgContent, err := os.ReadFile(svgPath)
		require.NoError(t, err)
		assert.Contains(t, string(svgContent), "<svg")
		assert.Contains(t, string(svgContent), "Test Project")
		assert.Contains(t, string(svgContent), "A comprehensive test")
	})

	t.Run("banner generation with dark theme and left align", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-2")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: 🚀 Dark Project -->
<!-- banner-tagline: Dark mode enabled -->
# Dark Project`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "dark", "left")
		require.NoError(t, err)

		svgPath := filepath.Join(projectDir, "banner.svg")
		svgContent, err := os.ReadFile(svgPath)
		require.NoError(t, err)

		darkTheme, _ := getTheme("dark")
		assert.Contains(t, string(svgContent), darkTheme.BG0)
		assert.Contains(t, string(svgContent), "&#128640;")
	})

	t.Run("banner generation with muted theme and right align", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-3")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: Muted Project -->
<!-- banner-tagline: Subtle colors -->
# Muted Project`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "muted", "right")
		require.NoError(t, err)

		svgPath := filepath.Join(projectDir, "banner.svg")
		svgContent, err := os.ReadFile(svgPath)
		require.NoError(t, err)

		mutedTheme, _ := getTheme("muted")
		assert.Contains(t, string(svgContent), mutedTheme.BG0)
	})

	t.Run("banner generation without tagline", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-4")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: Solo Title -->
# Solo Title Project`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "light", "center")
		require.NoError(t, err)

		svgPath := filepath.Join(projectDir, "banner.svg")
		svgContent, err := os.ReadFile(svgPath)
		require.NoError(t, err)
		assert.Contains(t, string(svgContent), "Solo Title")
	})

	t.Run("error on invalid theme", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-5")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: Test -->
# Test`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "invalid-theme", "center")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown theme")
	})

	t.Run("error on missing README", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-6")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		err = generateBanner(projectDir, "light", "center")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read README.md")
	})

	t.Run("error on missing banner-title", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-7")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `# Project without banner metadata`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "light", "center")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no banner-title found")
	})

	t.Run("error on invalid alignment", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-8")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: Test -->
# Test`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "light", "invalid-align")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to load template")
	})

	t.Run("PNG warning but SVG still created", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-9")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: PNG Test -->
<!-- banner-tagline: Testing PNG fallback -->
# PNG Test`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "light", "center")
		require.NoError(t, err)

		svgPath := filepath.Join(projectDir, "banner.svg")
		_, err = os.Stat(svgPath)
		assert.NoError(t, err)
	})

	t.Run("banner generation with special characters", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "test-project-10")
		err := os.MkdirAll(projectDir, 0755)
		require.NoError(t, err)

		readmeContent := `<!-- banner-title: <Test> & "Project" -->
<!-- banner-tagline: Special 'chars' -->
# Test`
		readmePath := filepath.Join(projectDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		err = generateBanner(projectDir, "light", "center")
		require.NoError(t, err)

		svgPath := filepath.Join(projectDir, "banner.svg")
		svgContent, err := os.ReadFile(svgPath)
		require.NoError(t, err)
		assert.Contains(t, string(svgContent), "&lt;Test&gt;")
		assert.Contains(t, string(svgContent), "&amp;")
		assert.Contains(t, string(svgContent), "&quot;")
		assert.Contains(t, string(svgContent), "&apos;")
	})
}

func TestGenerateBannerPNGOutput(t *testing.T) {
	tempDir := t.TempDir()
	projectDir := filepath.Join(tempDir, "png-test")
	err := os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	readmeContent := `<!-- banner-title: PNG Test -->
<!-- banner-tagline: Testing PNG output -->
# PNG Test`
	readmePath := filepath.Join(projectDir, "README.md")
	err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
	require.NoError(t, err)

	err = generateBanner(projectDir, "light", "center")
	require.NoError(t, err)

	pngPath := filepath.Join(projectDir, "banner.png")
	pngInfo, err := os.Stat(pngPath)

	if err == nil {
		assert.Greater(t, pngInfo.Size(), int64(0))

		pngContent, err := os.ReadFile(pngPath)
		if err == nil {
			assert.Equal(t, byte(0x89), pngContent[0])
			assert.Equal(t, byte('P'), pngContent[1])
			assert.Equal(t, byte('N'), pngContent[2])
			assert.Equal(t, byte('G'), pngContent[3])
		}
	} else {
		t.Log("PNG file not created (renderer may not be available)")
	}
}

func TestGenerateBannerIntegration(t *testing.T) {
	tempDir := t.TempDir()

	themes := []string{"light", "muted", "dark"}
	alignments := []string{"center", "left", "right"}

	for _, theme := range themes {
		for _, align := range alignments {
			t.Run(theme+"_"+align, func(t *testing.T) {
				projectDir := filepath.Join(tempDir, theme+"_"+align)
				err := os.MkdirAll(projectDir, 0755)
				require.NoError(t, err)

				readmeContent := `<!-- banner-title: Integration Test -->
<!-- banner-tagline: Testing all combinations -->
# Integration Test`
				readmePath := filepath.Join(projectDir, "README.md")
				err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
				require.NoError(t, err)

				err = generateBanner(projectDir, theme, align)
				require.NoError(t, err)

				svgPath := filepath.Join(projectDir, "banner.svg")
				_, err = os.Stat(svgPath)
				assert.NoError(t, err)

				svgContent, err := os.ReadFile(svgPath)
				require.NoError(t, err)
				svgStr := string(svgContent)
				assert.True(t, strings.Contains(svgStr, "<svg") && (strings.HasPrefix(svgStr, "<?xml") || strings.HasPrefix(svgStr, "<svg")))
			})
		}
	}
}
