package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	// ANSI escape codes
	clearScreenCode     = "\033[2J"
	clearScrollbackCode = "\033[3J"
	cursorHomeCode      = "\033[H"
	hideCursorCode      = "\033[?25l"
	showCursorCode      = "\033[?25h"
	clearLineCode       = "\033[K"
	boldCode            = "\033[1m"
	cyanCode            = "\033[36m"
	yellowCode          = "\033[33m"
	blueCode            = "\033[34m"
	resetCode           = "\033[0m"
	bgCyanCode          = "\033[46m"
	bgBlackCode         = "\033[40m"
	bgLightGrayCode     = "\033[47m"
	bgDarkGrayCode      = "\033[100m"
	// Mouse tracking
	enableMouseCode  = "\033[?1000h\033[?1002h\033[?1006h"
	disableMouseCode = "\033[?1000l\033[?1002l\033[?1006l"
	// Use \r\n for proper line endings in raw mode
	newline = "\r\n"
)

type PromptUI struct {
	scripts         []Script
	filteredScripts []Script
	selectedIndex   int
	searchQuery     string
	maxHeight       int
	// Mouse tracking
	scriptStartLine int // Line where scripts start being displayed
	lastClickTime   int64
	lastClickLine   int
	// Scrolling state
	viewStartIdx int // Current view window start index
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
		nameMatch := fuzzyMatch(script.Name, ui.searchQuery)
		cmdMatch := fuzzyMatch(script.Command, ui.searchQuery)

		if nameMatch || cmdMatch {
			nameScore := fuzzyScore(script.Name, ui.searchQuery)
			cmdScore := fuzzyScore(script.Command, ui.searchQuery)
			bestScore := min(nameScore, cmdScore)

			scored = append(scored, scoredScript{
				script: script,
				score:  bestScore,
			})
		}
	}

	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score < scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	ui.filteredScripts = make([]Script, len(scored))
	for i, s := range scored {
		ui.filteredScripts[i] = s.script
	}

	if ui.selectedIndex >= len(ui.filteredScripts) {
		ui.selectedIndex = len(ui.filteredScripts) - 1
		if ui.selectedIndex < 0 {
			ui.selectedIndex = 0
		}
	}

	ui.viewStartIdx = 0
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

func (ui *PromptUI) handleMouseEvent(event *MouseEvent, startIdx, endIdx int) bool {
	switch event.Type {
	case "click":
		clickedLine := event.Y - ui.scriptStartLine
		clickedIndex := startIdx + clickedLine

		if clickedIndex >= startIdx && clickedIndex < endIdx && clickedIndex < len(ui.filteredScripts) {
			now := time.Now().UnixMilli()
			if ui.lastClickLine == event.Y && now-ui.lastClickTime < 500 {
				ui.selectedIndex = clickedIndex
				return true
			}

			ui.selectedIndex = clickedIndex
			ui.lastClickTime = now
			ui.lastClickLine = event.Y
		}

	case "scroll_up":
		if ui.selectedIndex > 0 {
			ui.selectedIndex--
		}

	case "scroll_down":
		if ui.selectedIndex < len(ui.filteredScripts)-1 {
			ui.selectedIndex++
		}
	}

	return false
}

func truncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return "..."
	}
	return text[:maxLen-3] + "..."
}

func highlightMatch(text, query string, isSelected bool, maxLen int) string {
	truncated := truncateText(text, maxLen)

	if query == "" {
		return truncated
	}

	lowerText := strings.ToLower(truncated)
	lowerQuery := strings.ToLower(query)

	index := strings.Index(lowerText, lowerQuery)
	if index != -1 {
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

	if textIdx < len(truncated) {
		result.WriteString(truncated[textIdx:])
	}

	return result.String()
}

func (ui *PromptUI) render() (startIdx, endIdx int) {
	var output strings.Builder

	output.WriteString(clearScreenCode)
	output.WriteString(clearScrollbackCode)
	output.WriteString(cursorHomeCode)
	output.WriteString(hideCursorCode)

	termWidth, termHeight := getTerminalSize()
	ui.maxHeight = termHeight - 2
	if ui.maxHeight < 5 {
		ui.maxHeight = 5
	}

	availableWidth := termWidth - 4
	if availableWidth < 20 {
		availableWidth = 20 // Minimum reasonable width
	}
	maxNameDisplayWidth := (availableWidth * 4) / 10
	maxCommandDisplayWidth := availableWidth - maxNameDisplayWidth - 2 // -2 for gap

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

	startIdx = ui.viewStartIdx

	if len(ui.filteredScripts) > ui.maxHeight {
		scrollThreshold := 3
		if ui.selectedIndex < startIdx+scrollThreshold {
			startIdx = ui.selectedIndex - scrollThreshold
			if startIdx < 0 {
				startIdx = 0
			}
		}
		if ui.selectedIndex >= startIdx+ui.maxHeight-scrollThreshold {
			startIdx = ui.selectedIndex - ui.maxHeight + scrollThreshold + 1
			if startIdx < 0 {
				startIdx = 0
			}
		}
		if startIdx+ui.maxHeight > len(ui.filteredScripts) {
			startIdx = len(ui.filteredScripts) - ui.maxHeight
		}
		endIdx = startIdx + ui.maxHeight
		if endIdx > len(ui.filteredScripts) {
			endIdx = len(ui.filteredScripts)
		}
		ui.viewStartIdx = startIdx
	} else {
		ui.viewStartIdx = 0
		startIdx = 0
		endIdx = len(ui.filteredScripts)
	}

	displayedCount := 0
	if len(ui.filteredScripts) == 0 {
		displayedCount = 0
	} else {
		displayedCount = endIdx - startIdx
	}

	listLines := termHeight - 2
	paddingLines := listLines - displayedCount
	if paddingLines < 0 {
		paddingLines = 0
	}
	for i := 0; i < paddingLines; i++ {
		output.WriteString(clearLineCode)
		output.WriteString(newline)
	}

	ui.scriptStartLine = paddingLines + 2

	for i := startIdx; i < endIdx; i++ {
		script := ui.filteredScripts[i]
		isSelected := i == ui.selectedIndex
		if isSelected {
			output.WriteString(boldCode)
			output.WriteString(cyanCode)
			output.WriteString(bgDarkGrayCode)
			output.WriteString("▌")
			output.WriteString(resetCode)
			output.WriteString(boldCode)
			output.WriteString(bgDarkGrayCode)
			output.WriteString(" ")
			output.WriteString(cyanCode)
		} else {
			output.WriteString("  ")
		}
		highlightedName := highlightMatch(script.Name, ui.searchQuery, isSelected, maxNameDisplayWidth)
		highlightedCommand := highlightMatch(script.Command, ui.searchQuery, isSelected, maxCommandDisplayWidth)
		truncatedNameLen := len(truncateText(script.Name, maxNameDisplayWidth))
		output.WriteString(highlightedName)
		if truncatedNameLen < maxNameWidth {
			padding := strings.Repeat(" ", maxNameWidth-truncatedNameLen)
			output.WriteString(padding)
		}
		output.WriteString("  ")
		output.WriteString(highlightedCommand)
		if isSelected {
			currentLen := 2 + maxNameWidth + 2 + len(truncateText(script.Command, maxCommandDisplayWidth))
			if currentLen < termWidth {
				output.WriteString(strings.Repeat(" ", termWidth-currentLen))
			}
			output.WriteString(resetCode)
		}
		output.WriteString(newline)
	}

	counter := ""
	if len(ui.filteredScripts) > 0 {
		counter = fmt.Sprintf("(%d/%d)", ui.selectedIndex+1, len(ui.filteredScripts))
	} else {
		counter = "(0/0)"
	}
	dividerWidth := termWidth - len(counter)
	if dividerWidth < 0 {
		dividerWidth = 0
	}
	output.WriteString(counter)
	output.WriteString(strings.Repeat("─", dividerWidth))
	output.WriteString(newline)

	output.WriteString(boldCode)
	output.WriteString(blueCode)
	output.WriteString("> ")
	output.WriteString(resetCode)
	output.WriteString(boldCode)
	output.WriteString(ui.searchQuery)
	output.WriteString(showCursorCode)

	fmt.Print(output.String())

	return startIdx, endIdx
}

func showScriptPrompt() (*Script, error) {
	scripts, err := getScriptsOrdered()
	if err != nil {
		return nil, err
	}

	if len(scripts) == 0 {
		return nil, fmt.Errorf("no scripts found in package.json")
	}

	oldState, err := enableRawMode()
	if err != nil {
		return nil, err
	}

	fmt.Print(enableMouseCode)

	defer func() {
		fmt.Print(disableMouseCode)
		restoreTerminal(oldState)
		fmt.Print(showCursorCode)
		fmt.Print(cursorHomeCode)
		fmt.Print(clearScreenCode)
	}()

	ui := NewPromptUI(scripts)

	// Set up terminal resize handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)
	defer signal.Stop(sigChan)

	// Handle resize signals in a goroutine
	resizeChan := make(chan bool, 1)
	go func() {
		for range sigChan {
			// Non-blocking send
			select {
			case resizeChan <- true:
			default:
			}
		}
	}()

	// Channel for key input
	keyChan := make(chan []byte)
	errorChan := make(chan error)

	// Read keys in a goroutine
	go func() {
		for {
			key, err := readKey()
			if err != nil {
				errorChan <- err
				return
			}
			keyChan <- key
		}
	}()

	// Initial render
	startIdx, endIdx := ui.render()

	for {
		select {
		case <-resizeChan:
			// Re-render on resize
			startIdx, endIdx = ui.render()
		case err := <-errorChan:
			return nil, err
		case key := <-keyChan:
			// Check if it's a mouse event
			if mouseEvent, ok := parseMouseEvent(key); ok {
				if ui.handleMouseEvent(mouseEvent, startIdx, endIdx) {
					// Double-click detected, run the script
					selected := ui.getSelectedScript()
					if selected != nil {
						return selected, nil
					}
				}
				startIdx, endIdx = ui.render()
			} else if isCtrlC(key) || isEscape(key) {
				return nil, fmt.Errorf("cancelled")
			} else if isEnter(key) {
				selected := ui.getSelectedScript()
				if selected != nil {
					return selected, nil
				}
			} else if isArrowUp(key) {
				ui.moveUp()
				startIdx, endIdx = ui.render()
			} else if isArrowDown(key) {
				ui.moveDown()
				startIdx, endIdx = ui.render()
			} else {
				ui.handleSearchInput(key)
				startIdx, endIdx = ui.render()
			}
		}
	}
}
