package utils

import (
	"strings"
	"unicode"
)

// Converts any string format to camelCase
func toCamelCase(str string) string {
	runes := []rune(str)
	var result []rune
	isFirst := true

	for _, r := range runes {
		if r == '_' || unicode.IsSpace(r) {
			// Set flag to capitalize the next non-underscore, non-space character
			isFirst = false
			continue
		}

		if isFirst {
			// Only the first character of the entire string is made lowercase
			result = append(result, unicode.ToLower(r))
			isFirst = false
		} else {
			// Preserve the original casing of all subsequent characters
			result = append(result, r)
		}
	}
	return string(result)
}

// Converts any string format to TitleCase (every word starts with an uppercase letter)
func toTitleCase(str string) string {
	runes := []rune(str)
	var result []rune
	nextUpper := true

	for i, r := range runes {
		if r == '_' || unicode.IsSpace(r) {
			nextUpper = true
			continue
		}

		if nextUpper {
			result = append(result, unicode.ToUpper(r))
			nextUpper = false
		} else {
			// Preserve the original casing of the character
			result = append(result, runes[i])
		}
	}
	return string(result)
}

// Converts any string format to snake_case
func toSnakeCase(str string) string {
	runes := []rune(str)
	var result []rune

	for i, r := range runes {
		if i > 0 && (unicode.IsUpper(r) || unicode.IsSpace(runes[i-1])) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return strings.Trim(string(result), "_")
}

// ProcessString processes the string based on given rules
func ProcessString(input string) (string, string, string) {
	camelCase := toCamelCase(input)
	titleCase := toTitleCase(input)
	snakeCase := toSnakeCase(input)

	return titleCase, snakeCase, camelCase
}
