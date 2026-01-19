package main

import (
	"fmt"
	"strings"
)

type ThemePalette struct {
	BG0   string
	BG1   string
	BG2   string
	WAVE0 string
	WAVE1 string
}

var themes = map[string]ThemePalette{
	"light": {
		BG0:   "#8BCFE6",
		BG1:   "#F2B5C8",
		BG2:   "#F8F9FB",
		WAVE0: "#9DD7EC",
		WAVE1: "#F6AFC3",
	},
	"muted": {
		BG0:   "#7FC3DD",
		BG1:   "#EFAEC2",
		BG2:   "#F3F5F7",
		WAVE0: "#8FCFE3",
		WAVE1: "#F2A7BE",
	},
	"dark": {
		BG0:   "#245A74",
		BG1:   "#7A3651",
		BG2:   "#0F1720",
		WAVE0: "#3A7C96",
		WAVE1: "#A35A74",
	},
}

func getTheme(name string) (*ThemePalette, error) {
	theme, ok := themes[name]
	if !ok {
		availableThemes := make([]string, 0, len(themes))
		for k := range themes {
			availableThemes = append(availableThemes, k)
		}
		return nil, fmt.Errorf("unknown theme %q. Use: %s", name, strings.Join(availableThemes, ", "))
	}
	return &theme, nil
}

func escapeXML(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if r > 127 {
			result.WriteString(fmt.Sprintf("&#%d;", r))
		} else {
			switch r {
			case '&':
				result.WriteString("&amp;")
			case '<':
				result.WriteString("&lt;")
			case '>':
				result.WriteString("&gt;")
			case '"':
				result.WriteString("&quot;")
			case '\'':
				result.WriteString("&apos;")
			default:
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}

func replaceVariables(svg string, vars map[string]string) string {
	result := svg
	for k, v := range vars {
		placeholder := fmt.Sprintf("{{%s}}", k)
		result = strings.ReplaceAll(result, placeholder, escapeXML(v))
	}
	return result
}

func stripBadge(svg string, n int) string {
	startMarker := fmt.Sprintf("<!--BADGE%d_START-->", n)
	endMarker := fmt.Sprintf("<!--BADGE%d_END-->", n)

	startIdx := strings.Index(svg, startMarker)
	if startIdx == -1 {
		return svg
	}

	endIdx := strings.Index(svg, endMarker)
	if endIdx == -1 {
		return svg
	}

	endIdx += len(endMarker)
	for endIdx < len(svg) && (svg[endIdx] == ' ' || svg[endIdx] == '\n' || svg[endIdx] == '\r' || svg[endIdx] == '\t') {
		endIdx++
	}

	return svg[:startIdx] + svg[endIdx:]
}
