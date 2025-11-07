package ui

import (
	"strings"
	"unicode"
)

func toLowerRunes(s string) []rune {
	runes := []rune(s)
	for i, r := range runes {
		runes[i] = unicode.ToLower(r)
	}
	return runes
}

// fuzzyMatch checks if query matches text in a fuzzy manner
func fuzzyMatch(text, query string) bool {
	if query == "" {
		return true
	}

	textRunes := toLowerRunes(text)
	queryRunes := toLowerRunes(query)

	queryIdx := 0

	for _, r := range textRunes {
		if r == queryRunes[queryIdx] {
			queryIdx++
			if queryIdx == len(queryRunes) {
				return true
			}
		}
	}

	return false
}

// fuzzyScore calculates a fuzzy match score (lower is better)
func fuzzyScore(text, query string) int {
	if query == "" {
		return 0
	}

	textRunes := toLowerRunes(text)
	queryRunes := toLowerRunes(query)

	score := 0
	textIdx := 0
	queryIdx := 0
	lastMatchIdx := -1

	for textIdx < len(textRunes) && queryIdx < len(queryRunes) {
		if textRunes[textIdx] == queryRunes[queryIdx] {
			if lastMatchIdx != -1 {
				score += textIdx - lastMatchIdx
			}
			lastMatchIdx = textIdx
			queryIdx++
		}
		textIdx++
	}

	if queryIdx < len(queryRunes) {
		return 999999
	}

	textLower := string(textRunes)
	queryLower := string(queryRunes)

	if strings.Contains(textLower, queryLower) {
		score -= 100
	}

	if strings.HasPrefix(textLower, queryLower) {
		score -= 200
	}

	return score
}
