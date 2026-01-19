package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTheme(t *testing.T) {
	tests := []struct {
		name        string
		themeName   string
		expectError bool
	}{
		{
			name:        "light theme",
			themeName:   "light",
			expectError: false,
		},
		{
			name:        "muted theme",
			themeName:   "muted",
			expectError: false,
		},
		{
			name:        "dark theme",
			themeName:   "dark",
			expectError: false,
		},
		{
			name:        "invalid theme",
			themeName:   "invalid",
			expectError: true,
		},
		{
			name:        "empty theme name",
			themeName:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme, err := getTheme(tt.themeName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, theme)
				assert.Contains(t, err.Error(), "unknown theme")
			} else {
				require.NoError(t, err)
				require.NotNil(t, theme)
				assert.NotEmpty(t, theme.BG0)
				assert.NotEmpty(t, theme.BG1)
				assert.NotEmpty(t, theme.BG2)
				assert.NotEmpty(t, theme.WAVE0)
				assert.NotEmpty(t, theme.WAVE1)
			}
		})
	}
}

func TestGetThemeValidation(t *testing.T) {
	t.Run("light theme has correct colors", func(t *testing.T) {
		theme, err := getTheme("light")
		require.NoError(t, err)
		assert.Equal(t, "#8BCFE6", theme.BG0)
		assert.Equal(t, "#F2B5C8", theme.BG1)
		assert.Equal(t, "#F8F9FB", theme.BG2)
		assert.Equal(t, "#9DD7EC", theme.WAVE0)
		assert.Equal(t, "#F6AFC3", theme.WAVE1)
	})

	t.Run("muted theme has correct colors", func(t *testing.T) {
		theme, err := getTheme("muted")
		require.NoError(t, err)
		assert.Equal(t, "#7FC3DD", theme.BG0)
		assert.Equal(t, "#EFAEC2", theme.BG1)
		assert.Equal(t, "#F3F5F7", theme.BG2)
		assert.Equal(t, "#8FCFE3", theme.WAVE0)
		assert.Equal(t, "#F2A7BE", theme.WAVE1)
	})

	t.Run("dark theme has correct colors", func(t *testing.T) {
		theme, err := getTheme("dark")
		require.NoError(t, err)
		assert.Equal(t, "#245A74", theme.BG0)
		assert.Equal(t, "#7A3651", theme.BG1)
		assert.Equal(t, "#0F1720", theme.BG2)
		assert.Equal(t, "#3A7C96", theme.WAVE0)
		assert.Equal(t, "#A35A74", theme.WAVE1)
	})
}

func TestEscapeXML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special characters",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "ampersand",
			input:    "Tom & Jerry",
			expected: "Tom &amp; Jerry",
		},
		{
			name:     "less than and greater than",
			input:    "<div>content</div>",
			expected: "&lt;div&gt;content&lt;/div&gt;",
		},
		{
			name:     "quotes",
			input:    `"double" and 'single'`,
			expected: "&quot;double&quot; and &apos;single&apos;",
		},
		{
			name:     "all special characters",
			input:    `<tag attr="value" other='value2'>&amp;</tag>`,
			expected: "&lt;tag attr=&quot;value&quot; other=&apos;value2&apos;&gt;&amp;amp;&lt;/tag&gt;",
		},
		{
			name:     "unicode emoji",
			input:    "🚀 Rocket",
			expected: "&#128640; Rocket",
		},
		{
			name:     "unicode characters",
			input:    "こんにちは",
			expected: "&#12371;&#12435;&#12395;&#12385;&#12399;",
		},
		{
			name:     "mixed unicode and special chars",
			input:    "🎨 <Art & Design>",
			expected: "&#127912; &lt;Art &amp; Design&gt;",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeXML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReplaceVariables(t *testing.T) {
	tests := []struct {
		name     string
		template string
		vars     map[string]string
		expected string
	}{
		{
			name:     "single variable",
			template: "Hello {{NAME}}",
			vars:     map[string]string{"NAME": "World"},
			expected: "Hello World",
		},
		{
			name:     "multiple variables",
			template: "{{GREETING}} {{NAME}}, welcome to {{PLACE}}",
			vars: map[string]string{
				"GREETING": "Hello",
				"NAME":     "Alice",
				"PLACE":    "Wonderland",
			},
			expected: "Hello Alice, welcome to Wonderland",
		},
		{
			name:     "variable with special characters gets escaped",
			template: "<text>{{CONTENT}}</text>",
			vars:     map[string]string{"CONTENT": "Tom & Jerry"},
			expected: "<text>Tom &amp; Jerry</text>",
		},
		{
			name:     "variable with unicode",
			template: "{{TITLE}}",
			vars:     map[string]string{"TITLE": "🚀 Project"},
			expected: "&#128640; Project",
		},
		{
			name:     "no variables",
			template: "Static content",
			vars:     map[string]string{},
			expected: "Static content",
		},
		{
			name:     "unused variables",
			template: "Hello {{NAME}}",
			vars: map[string]string{
				"NAME":  "World",
				"EXTRA": "Ignored",
			},
			expected: "Hello World",
		},
		{
			name:     "missing variables remain unchanged",
			template: "Hello {{NAME}}, from {{PLACE}}",
			vars:     map[string]string{"NAME": "World"},
			expected: "Hello World, from {{PLACE}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := replaceVariables(tt.template, tt.vars)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStripBadge(t *testing.T) {
	tests := []struct {
		name     string
		svg      string
		badgeNum int
		expected string
	}{
		{
			name: "strip badge 1",
			svg: `<svg>
<!--BADGE1_START-->
<rect x="10" y="10" width="100" height="20"/>
<!--BADGE1_END-->
<circle cx="50" cy="50" r="40"/>
</svg>`,
			badgeNum: 1,
			expected: `<svg>
<circle cx="50" cy="50" r="40"/>
</svg>`,
		},
		{
			name: "strip badge 2 with trailing whitespace",
			svg: `<svg>
<!--BADGE2_START-->
<text>Badge 2</text>
<!--BADGE2_END-->   
<circle/>
</svg>`,
			badgeNum: 2,
			expected: `<svg>
<circle/>
</svg>`,
		},
		{
			name: "strip badge 3 with newline after",
			svg: `<svg>
<!--BADGE3_START-->
<g>Badge 3</g>
<!--BADGE3_END-->

<rect/>
</svg>`,
			badgeNum: 3,
			expected: `<svg>
<rect/>
</svg>`,
		},
		{
			name: "badge markers not present",
			svg: `<svg>
<circle cx="50" cy="50" r="40"/>
</svg>`,
			badgeNum: 1,
			expected: `<svg>
<circle cx="50" cy="50" r="40"/>
</svg>`,
		},
		{
			name: "only start marker present",
			svg: `<svg>
<!--BADGE1_START-->
<rect/>
</svg>`,
			badgeNum: 1,
			expected: `<svg>
<!--BADGE1_START-->
<rect/>
</svg>`,
		},
		{
			name: "only end marker present",
			svg: `<svg>
<rect/>
<!--BADGE1_END-->
</svg>`,
			badgeNum: 1,
			expected: `<svg>
<rect/>
<!--BADGE1_END-->
</svg>`,
		},
		{
			name: "multiple badges strip specific one",
			svg: `<svg>
<!--BADGE1_START-->
<text>Badge 1</text>
<!--BADGE1_END-->
<!--BADGE2_START-->
<text>Badge 2</text>
<!--BADGE2_END-->
<circle/>
</svg>`,
			badgeNum: 1,
			expected: `<svg>
<!--BADGE2_START-->
<text>Badge 2</text>
<!--BADGE2_END-->
<circle/>
</svg>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripBadge(tt.svg, tt.badgeNum)
			assert.Equal(t, tt.expected, result)
		})
	}
}
