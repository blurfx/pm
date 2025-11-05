package ui

import "strings"

// fuzzyMatch checks if query matches text in a fuzzy manner
func fuzzyMatch(text, query string) bool {
	if query == "" {
		return true
	}

	textLower := strings.ToLower(text)
	queryLower := strings.ToLower(query)

	textIdx := 0
	queryIdx := 0

	for textIdx < len(textLower) && queryIdx < len(queryLower) {
		if textLower[textIdx] == queryLower[queryIdx] {
			queryIdx++
		}
		textIdx++
	}

	return queryIdx == len(queryLower)
}

// fuzzyScore calculates a fuzzy match score (lower is better)
func fuzzyScore(text, query string) int {
	if query == "" {
		return 0
	}

	textLower := strings.ToLower(text)
	queryLower := strings.ToLower(query)

	score := 0
	textIdx := 0
	queryIdx := 0
	lastMatchIdx := -1

	for textIdx < len(textLower) && queryIdx < len(queryLower) {
		if textLower[textIdx] == queryLower[queryIdx] {
			if lastMatchIdx != -1 {
				score += textIdx - lastMatchIdx
			}
			lastMatchIdx = textIdx
			queryIdx++
		}
		textIdx++
	}

	if queryIdx < len(queryLower) {
		return 999999
	}

	if strings.Contains(textLower, queryLower) {
		score -= 100
	}

	if strings.HasPrefix(textLower, queryLower) {
		score -= 200
	}

	return score
}
