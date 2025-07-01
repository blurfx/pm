package main

import (
	"fmt"
	"strings"
)

const (
	// ANSI escape codes
	clearScreenCode = "\033[2J"
	cursorHomeCode  = "\033[H"
	hideCursorCode  = "\033[?25l"
	showCursorCode  = "\033[?25h"
	clearLineCode   = "\033[K"
	cyanCode        = "\033[36m"
	yellowCode      = "\033[33m"
	resetCode       = "\033[0m"
	// Use \r\n for proper line endings in raw mode
	newline = "\r\n"
)

type PromptUI struct {
	scripts         []Script
	filteredScripts []Script
	selectedIndex   int
	searchQuery     string
	maxHeight       int
}

func NewPromptUI(scripts []Script) *PromptUI {
	ui := &PromptUI{
		scripts:         scripts,
		filteredScripts: scripts,
		selectedIndex:   0,
		searchQuery:     "",
	}
	return ui
}

// fuzzyMatch performs a simple fuzzy matching algorithm
// Returns true if all characters in query appear in text in order
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

// fuzzyScore calculates a score for fuzzy matching
// Lower score is better match
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
			// Add gap penalty
			if lastMatchIdx != -1 {
				score += textIdx - lastMatchIdx
			}
			lastMatchIdx = textIdx
			queryIdx++
		}
		textIdx++
	}
	
	// Penalize if not all characters matched
	if queryIdx < len(queryLower) {
		return 999999
	}
	
	// Bonus for exact match
	if strings.Contains(textLower, queryLower) {
		score -= 100
	}
	
	// Bonus for matching at start
	if strings.HasPrefix(textLower, queryLower) {
		score -= 200
	}
	
	return score
}

type scoredScript struct {
	script Script
	score  int
}

func (ui *PromptUI) filterScripts() {
	if ui.searchQuery == "" {
		ui.filteredScripts = ui.scripts
		return
	}

	scored := make([]scoredScript, 0)

	for _, script := range ui.scripts {
		// Check both name and command
		nameMatch := fuzzyMatch(script.Name, ui.searchQuery)
		cmdMatch := fuzzyMatch(script.Command, ui.searchQuery)
		
		if nameMatch || cmdMatch {
			// Calculate best score between name and command
			nameScore := fuzzyScore(script.Name, ui.searchQuery)
			cmdScore := fuzzyScore(script.Command, ui.searchQuery)
			
			bestScore := nameScore
			if cmdScore < bestScore {
				bestScore = cmdScore
			}
			
			scored = append(scored, scoredScript{
				script: script,
				score:  bestScore,
			})
		}
	}

	// Sort by score (better matches first)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score < scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	// Extract sorted scripts
	ui.filteredScripts = make([]Script, len(scored))
	for i, s := range scored {
		ui.filteredScripts[i] = s.script
	}

	// Reset selected index if it's out of bounds
	if ui.selectedIndex >= len(ui.filteredScripts) {
		ui.selectedIndex = len(ui.filteredScripts) - 1
		if ui.selectedIndex < 0 {
			ui.selectedIndex = 0
		}
	}
}

func (ui *PromptUI) moveUp() {
	if ui.selectedIndex > 0 {
		ui.selectedIndex--
	}
}

func (ui *PromptUI) moveDown() {
	if ui.selectedIndex < len(ui.filteredScripts)-1 {
		ui.selectedIndex++
	}
}

func (ui *PromptUI) handleSearchInput(key []byte) {
	if isBackspace(key) {
		if len(ui.searchQuery) > 0 {
			ui.searchQuery = ui.searchQuery[:len(ui.searchQuery)-1]
			ui.filterScripts()
		}
	} else if len(key) == 1 && key[0] >= 32 && key[0] <= 126 {
		// Printable character
		ui.searchQuery += string(key)
		ui.filterScripts()
	}
}

func (ui *PromptUI) getSelectedScript() *Script {
	if len(ui.filteredScripts) == 0 || ui.selectedIndex >= len(ui.filteredScripts) {
		return nil
	}
	return &ui.filteredScripts[ui.selectedIndex]
}

// truncateText truncates text to maxLen and adds ellipsis if needed
func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return "..."
	}
	return text[:maxLen-3] + "..."
}

// highlightMatch highlights the matching part of text based on the query
// and truncates the text if needed
func highlightMatch(text, query string, isSelected bool, maxLen int) string {
	// First truncate the text
	truncated := truncateText(text, maxLen)

	if query == "" {
		return truncated
	}

	// For fuzzy matching, we need to highlight individual characters
	lowerText := strings.ToLower(truncated)
	lowerQuery := strings.ToLower(query)
	
	// First check if it's a substring match (prioritize continuous matches)
	index := strings.Index(lowerText, lowerQuery)
	if index != -1 {
		// Build the highlighted string for continuous match
		highlightColor := yellowCode
		if isSelected {
			highlightColor = "\033[1;33m"
		}
		
		resetColor := resetCode
		if isSelected {
			resetColor = "\033[0m" + cyanCode
		}
		
		result := truncated[:index] + highlightColor + truncated[index:index+len(query)] + resetColor + truncated[index+len(query):]
		return result
	}
	
	// Fuzzy highlight - highlight matching characters
	var result strings.Builder
	textIdx := 0
	queryIdx := 0
	
	highlightColor := yellowCode
	if isSelected {
		highlightColor = "\033[1;33m"
	}
	
	resetColor := resetCode
	if isSelected {
		resetColor = "\033[0m" + cyanCode
	}
	
	for textIdx < len(truncated) && queryIdx < len(lowerQuery) {
		if strings.ToLower(string(truncated[textIdx])) == string(lowerQuery[queryIdx]) {
			result.WriteString(highlightColor)
			result.WriteByte(truncated[textIdx])
			result.WriteString(resetColor)
			queryIdx++
		} else {
			result.WriteByte(truncated[textIdx])
		}
		textIdx++
	}
	
	// Add remaining characters
	if textIdx < len(truncated) {
		result.WriteString(truncated[textIdx:])
	}
	
	return result.String()
}

func (ui *PromptUI) render() {
	// Build entire output as a string
	var output strings.Builder

	// Clear screen and move cursor to home
	output.WriteString(cursorHomeCode)
	output.WriteString(clearScreenCode)
	output.WriteString(hideCursorCode)

	// Title
	output.WriteString("Select a script to run (↑/↓ to navigate, Enter to select, Esc to cancel)")
	output.WriteString(newline)
	output.WriteString(newline)

	// Search box
	output.WriteString(fmt.Sprintf("Search: %s_", ui.searchQuery))
	output.WriteString(newline)
	output.WriteString(newline)

	// Get terminal dimensions
	termWidth, termHeight := getTerminalSize()
	ui.maxHeight = termHeight - 6 // Reserve lines for header and search
	if ui.maxHeight < 5 {
		ui.maxHeight = 5
	}

	// Calculate max widths for columns
	// Reserve space for: prefix (2) + gap (2) + some space for command
	availableWidth := termWidth - 4
	if availableWidth < 20 {
		availableWidth = 20 // Minimum reasonable width
	}

	// Split available width between name and command (40% name, 60% command)
	maxNameDisplayWidth := (availableWidth * 4) / 10
	maxCommandDisplayWidth := availableWidth - maxNameDisplayWidth - 2 // -2 for gap

	// Find the longest name (up to maxNameDisplayWidth)
	maxNameWidth := 0
	for _, script := range ui.filteredScripts {
		nameLen := len(script.Name)
		if nameLen > maxNameDisplayWidth {
			nameLen = maxNameDisplayWidth
		}
		if nameLen > maxNameWidth {
			maxNameWidth = nameLen
		}
	}

	startIdx := 0
	endIdx := len(ui.filteredScripts)

	// Implement scrolling if list is too long
	if len(ui.filteredScripts) > ui.maxHeight {
		if ui.selectedIndex >= ui.maxHeight/2 {
			startIdx = ui.selectedIndex - ui.maxHeight/2
			if startIdx+ui.maxHeight > len(ui.filteredScripts) {
				startIdx = len(ui.filteredScripts) - ui.maxHeight
			}
		}
		endIdx = startIdx + ui.maxHeight
		if endIdx > len(ui.filteredScripts) {
			endIdx = len(ui.filteredScripts)
		}
	}

	if len(ui.filteredScripts) == 0 {
		output.WriteString("  No scripts found")
		output.WriteString(newline)
	} else {
		for i := startIdx; i < endIdx; i++ {
			script := ui.filteredScripts[i]

			isSelected := i == ui.selectedIndex

			if isSelected {
				output.WriteString(cyanCode)
				output.WriteString("> ")
			} else {
				output.WriteString("  ")
			}

			// Highlight matching parts with truncation
			highlightedName := highlightMatch(script.Name, ui.searchQuery, isSelected, maxNameDisplayWidth)
			highlightedCommand := highlightMatch(script.Command, ui.searchQuery, isSelected, maxCommandDisplayWidth)

			// Calculate padding for name (considering truncated length)
			truncatedNameLen := len(truncateText(script.Name, maxNameDisplayWidth))
			output.WriteString(highlightedName)
			if truncatedNameLen < maxNameWidth {
				output.WriteString(strings.Repeat(" ", maxNameWidth-truncatedNameLen))
			}
			output.WriteString("  ")
			output.WriteString(highlightedCommand)

			if isSelected {
				output.WriteString(resetCode)
			}

			output.WriteString(newline)
		}
	}

	// Print everything at once
	fmt.Print(output.String())
}

func showScriptPrompt() (*Script, error) {
	scripts, err := getScriptsOrdered()
	if err != nil {
		return nil, err
	}

	if len(scripts) == 0 {
		return nil, fmt.Errorf("no scripts found in package.json")
	}

	// Enable raw mode
	oldState, err := enableRawMode()
	if err != nil {
		return nil, err
	}

	// Ensure cleanup
	defer func() {
		restoreTerminal(oldState)
		fmt.Print(showCursorCode)
		fmt.Print(cursorHomeCode)
		fmt.Print(clearScreenCode)
	}()

	ui := NewPromptUI(scripts)

	for {
		ui.render()

		key, err := readKey()
		if err != nil {
			return nil, err
		}

		if isCtrlC(key) || isEscape(key) {
			return nil, fmt.Errorf("cancelled")
		} else if isEnter(key) {
			selected := ui.getSelectedScript()
			if selected != nil {
				return selected, nil
			}
		} else if isArrowUp(key) {
			ui.moveUp()
		} else if isArrowDown(key) {
			ui.moveDown()
		} else {
			ui.handleSearchInput(key)
		}
	}
}
