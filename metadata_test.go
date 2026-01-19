package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseReadmeMetadata(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    *Metadata
		expectError bool
	}{
		{
			name: "valid metadata with title and tagline",
			content: `# My Project
<!-- banner-title: Test Project -->
<!-- banner-tagline: A test project for testing -->
Some content`,
			expected: &Metadata{
				Name:    "Test Project",
				Tagline: "A test project for testing",
			},
			expectError: false,
		},
		{
			name: "valid metadata with only title",
			content: `<!-- banner-title: Solo Project -->
More content here`,
			expected: &Metadata{
				Name:    "Solo Project",
				Tagline: "",
			},
			expectError: false,
		},
		{
			name: "title with emoji and special characters",
			content: `<!-- banner-title: 🚀 Banner Kit Go -->
<!-- banner-tagline: SVG/PNG banner generator -->`,
			expected: &Metadata{
				Name:    "🚀 Banner Kit Go",
				Tagline: "SVG/PNG banner generator",
			},
			expectError: false,
		},
		{
			name: "title with extra whitespace",
			content: `<!--   banner-title:   Extra Spaces   -->
<!--   banner-tagline:   Spaced Out   -->`,
			expected: &Metadata{
				Name:    "Extra Spaces",
				Tagline: "Spaced Out",
			},
			expectError: false,
		},
		{
			name:        "missing banner-title",
			content:     `<!-- banner-tagline: Only Tagline -->`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty content",
			content:     "",
			expected:    nil,
			expectError: true,
		},
		{
			name: "title and tagline in reverse order",
			content: `<!-- banner-tagline: Tagline First -->
<!-- banner-title: Title Second -->`,
			expected: &Metadata{
				Name:    "Title Second",
				Tagline: "Tagline First",
			},
			expectError: false,
		},
		{
			name: "multiple titles uses first one",
			content: `<!-- banner-title: First Title -->
<!-- banner-title: Second Title -->
<!-- banner-tagline: The Tagline -->`,
			expected: &Metadata{
				Name:    "First Title",
				Tagline: "The Tagline",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseReadmeMetadata(tt.content)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.Name, result.Name)
				assert.Equal(t, tt.expected.Tagline, result.Tagline)
			}
		})
	}
}

func TestReadProjectMetadata(t *testing.T) {
	// Create temporary directory for tests
	tempDir := t.TempDir()

	t.Run("valid README.md", func(t *testing.T) {
		readmeContent := `<!-- banner-title: Test Project -->
<!-- banner-tagline: Testing metadata reading -->
# Test Project`

		readmePath := filepath.Join(tempDir, "README.md")
		err := os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		metadata, err := readProjectMetadata(tempDir)
		require.NoError(t, err)
		assert.Equal(t, "Test Project", metadata.Name)
		assert.Equal(t, "Testing metadata reading", metadata.Tagline)
	})

	t.Run("missing README.md", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "nonexistent")
		metadata, err := readProjectMetadata(nonExistentDir)
		assert.Error(t, err)
		assert.Nil(t, metadata)
		assert.Contains(t, err.Error(), "failed to read README.md")
	})

	t.Run("README.md without banner-title", func(t *testing.T) {
		invalidDir := filepath.Join(tempDir, "invalid")
		err := os.MkdirAll(invalidDir, 0755)
		require.NoError(t, err)

		readmeContent := `# Project without banner metadata`
		readmePath := filepath.Join(invalidDir, "README.md")
		err = os.WriteFile(readmePath, []byte(readmeContent), 0644)
		require.NoError(t, err)

		metadata, err := readProjectMetadata(invalidDir)
		assert.Error(t, err)
		assert.Nil(t, metadata)
		assert.Contains(t, err.Error(), "no banner-title found")
	})
}
