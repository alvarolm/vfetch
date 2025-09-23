package main

import "strings"

func JSONRemoveTrailingCommas(jsonStr string) string {
	var result strings.Builder
	inString := false
	escaped := false

	runes := []rune(jsonStr)
	for i := 0; i < len(runes); i++ {
		char := runes[i]

		// Handle escape sequences
		if escaped {
			result.WriteRune(char)
			escaped = false
			continue
		}

		if char == '\\' && inString {
			escaped = true
			result.WriteRune(char)
			continue
		}

		// Handle string boundaries
		if char == '"' {
			inString = !inString
			result.WriteRune(char)
			continue
		}

		// If we're inside a string, just add the character
		if inString {
			result.WriteRune(char)
			continue
		}

		// Check for trailing comma outside of strings
		if char == ',' {
			// Look ahead to see if this is a trailing comma
			j := i + 1
			for j < len(runes) && (runes[j] == ' ' || runes[j] == '\t' || runes[j] == '\n' || runes[j] == '\r') {
				j++
			}

			// If the next non-whitespace character is } or ], skip the comma
			if j < len(runes) && (runes[j] == '}' || runes[j] == ']') {
				continue // Skip the trailing comma
			}
		}

		result.WriteRune(char)
	}

	return result.String()
}
