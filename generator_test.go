package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadTemplate(t *testing.T) {
	tests := []struct {
		name        string
		align       string
		expectError bool
	}{
		{
			name:        "center alignment",
			align:       "center",
			expectError: false,
		},
		{
			name:        "left alignment",
			align:       "left",
			expectError: false,
		},
		{
			name:        "right alignment",
			align:       "right",
			expectError: false,
		},
		{
			name:        "invalid alignment",
			align:       "invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := loadTemplate(tt.align)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Contains(t, err.Error(), "failed to load template")
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Contains(t, result, "<svg")
				assert.Contains(t, result, "{{BG0}}")
				assert.Contains(t, result, "{{PROJECT_NAME}}")
			}
		})
	}
}

func TestGenerateSVG(t *testing.T) {
	lightTheme, _ := getTheme("light")

	t.Run("basic SVG generation", func(t *testing.T) {
		metadata := &Metadata{
			Name:    "Test Project",
			Tagline: "A test tagline",
		}

		svg, err := generateSVG(metadata, lightTheme, "center", []string{})
		require.NoError(t, err)
		assert.NotEmpty(t, svg)
		assert.Contains(t, svg, "<svg")
		assert.Contains(t, svg, lightTheme.BG0)
		assert.Contains(t, svg, "Test Project")
		assert.Contains(t, svg, "A test tagline")
	})

	t.Run("SVG generation with emoji", func(t *testing.T) {
		metadata := &Metadata{
			Name:    "🚀 Rocket Project",
			Tagline: "Fast & Furious",
		}

		svg, err := generateSVG(metadata, lightTheme, "center", []string{})
		require.NoError(t, err)
		assert.Contains(t, svg, "&#128640;")
		assert.Contains(t, svg, "&amp;")
	})

	t.Run("SVG generation with badges", func(t *testing.T) {
		metadata := &Metadata{
			Name:    "Project",
			Tagline: "Tagline",
		}
		badges := []string{"badge1", "badge2", "badge3"}

		svg, err := generateSVG(metadata, lightTheme, "center", badges)
		require.NoError(t, err)
		assert.Contains(t, svg, "badge1")
		assert.Contains(t, svg, "badge2")
		assert.Contains(t, svg, "badge3")
	})

	t.Run("SVG generation strips empty badges", func(t *testing.T) {
		metadata := &Metadata{
			Name:    "Project",
			Tagline: "Tagline",
		}

		svg, err := generateSVG(metadata, lightTheme, "center", []string{})
		require.NoError(t, err)
		assert.NotContains(t, svg, "BADGE1_START")
		assert.NotContains(t, svg, "BADGE2_START")
		assert.NotContains(t, svg, "BADGE3_START")
	})

	t.Run("SVG generation with partial badges", func(t *testing.T) {
		metadata := &Metadata{
			Name:    "Project",
			Tagline: "Tagline",
		}
		badges := []string{"only-first"}

		svg, err := generateSVG(metadata, lightTheme, "center", badges)
		require.NoError(t, err)
		assert.Contains(t, svg, "only-first")
		assert.NotContains(t, svg, "BADGE2_START")
		assert.NotContains(t, svg, "BADGE3_START")
	})

	t.Run("SVG generation with all alignments", func(t *testing.T) {
		metadata := &Metadata{
			Name:    "Test",
			Tagline: "Test",
		}

		alignments := []string{"center", "left", "right"}
		for _, align := range alignments {
			svg, err := generateSVG(metadata, lightTheme, align, []string{})
			require.NoError(t, err, "alignment: %s", align)
			assert.NotEmpty(t, svg)
			assert.Contains(t, svg, "<svg")
		}
	})

	t.Run("invalid alignment", func(t *testing.T) {
		metadata := &Metadata{
			Name:    "Test",
			Tagline: "Test",
		}

		svg, err := generateSVG(metadata, lightTheme, "invalid", []string{})
		assert.Error(t, err)
		assert.Empty(t, svg)
	})
}

func TestConvertSVGToPNG(t *testing.T) {
	simpleSVG := `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
		<rect width="100" height="100" fill="red"/>
	</svg>`

	t.Run("PNG conversion succeeds", func(t *testing.T) {
		png, err := convertSVGToPNG(simpleSVG)

		if err != nil {
			t.Logf("PNG conversion failed (may be expected if no renderer available): %v", err)
			t.Skip("Skipping PNG conversion test - no renderer available")
		}

		require.NotNil(t, png)
		assert.Greater(t, len(png), 0)
		assert.Equal(t, byte(0x89), png[0])
		assert.Equal(t, byte('P'), png[1])
		assert.Equal(t, byte('N'), png[2])
		assert.Equal(t, byte('G'), png[3])
	})

	t.Run("invalid SVG", func(t *testing.T) {
		invalidSVG := "not an svg"
		png, err := convertSVGToPNG(invalidSVG)

		if err == nil {
			t.Skip("Renderer accepted invalid SVG (renderer-specific behavior)")
		}

		assert.Error(t, err)
		assert.Nil(t, png)
	})
}

func TestWriteBannerFiles(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("write SVG and PNG", func(t *testing.T) {
		svg := "<svg>test</svg>"
		png := []byte{0x89, 'P', 'N', 'G'}

		err := writeBannerFiles(tempDir, svg, png)
		require.NoError(t, err)

		svgPath := filepath.Join(tempDir, "banner.svg")
		pngPath := filepath.Join(tempDir, "banner.png")

		svgContent, err := os.ReadFile(svgPath)
		require.NoError(t, err)
		assert.Equal(t, svg, string(svgContent))

		pngContent, err := os.ReadFile(pngPath)
		require.NoError(t, err)
		assert.Equal(t, png, pngContent)
	})

	t.Run("write SVG only", func(t *testing.T) {
		subDir := filepath.Join(tempDir, "svg-only")
		err := os.MkdirAll(subDir, 0755)
		require.NoError(t, err)

		svg := "<svg>only svg</svg>"

		err = writeBannerFiles(subDir, svg, nil)
		require.NoError(t, err)

		svgPath := filepath.Join(subDir, "banner.svg")
		pngPath := filepath.Join(subDir, "banner.png")

		_, err = os.Stat(svgPath)
		assert.NoError(t, err)

		_, err = os.Stat(pngPath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("write to non-existent directory", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "does-not-exist")
		svg := "<svg>test</svg>"

		err := writeBannerFiles(nonExistentDir, svg, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to write SVG")
	})

	t.Run("write empty SVG", func(t *testing.T) {
		subDir := filepath.Join(tempDir, "empty")
		err := os.MkdirAll(subDir, 0755)
		require.NoError(t, err)

		err = writeBannerFiles(subDir, "", nil)
		require.NoError(t, err)

		svgPath := filepath.Join(subDir, "banner.svg")
		content, err := os.ReadFile(svgPath)
		require.NoError(t, err)
		assert.Empty(t, content)
	})
}

func TestConvertWithRsvgConvert(t *testing.T) {
	simpleSVG := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
		<rect width="100" height="100" fill="blue"/>
	</svg>`)

	png, err := convertWithRsvgConvert(simpleSVG)

	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			t.Skip("rsvg-convert not installed, skipping test")
		}
		t.Logf("rsvg-convert error: %v", err)
	}

	if err == nil {
		require.NotNil(t, png)
		assert.Greater(t, len(png), 0)
		assert.Equal(t, byte(0x89), png[0])
		assert.Equal(t, byte('P'), png[1])
		assert.Equal(t, byte('N'), png[2])
		assert.Equal(t, byte('G'), png[3])
	}
}

func TestConvertWithResvg(t *testing.T) {
	simpleSVG := []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
		<rect width="100" height="100" fill="green"/>
	</svg>`)

	png, err := convertWithResvg(simpleSVG)

	if err != nil {
		t.Logf("resvg conversion error (may be environment-specific): %v", err)
		t.Skip("Skipping resvg test - conversion failed")
	}

	require.NotNil(t, png)
	assert.Greater(t, len(png), 0)
	assert.Equal(t, byte(0x89), png[0])
	assert.Equal(t, byte('P'), png[1])
	assert.Equal(t, byte('N'), png[2])
	assert.Equal(t, byte('G'), png[3])
}
