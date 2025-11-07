package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"

	"pm/internal/detector"
	"pm/internal/project"
)

const (
	clearScreenCode      = "\033[2J"
	clearScrollbackCode  = "\033[3J"
	cursorHomeCode       = "\033[H"
	hideCursorCode       = "\033[?25l"
	showCursorCode       = "\033[?25h"
	enableMouseCode      = "\033[?1000h\033[?1002h\033[?1006h"
	disableMouseCode     = "\033[?1000l\033[?1002l\033[?1006l"
	enableAltScreenCode  = "\033[?1049h"
	disableAltScreenCode = "\033[?1049l"
	clearLineCode        = "\033[K"
	boldCode             = "\033[1m"
	resetCode            = "\033[0m"
	newline              = "\r\n"

	yellowCode  = "\033[33m"
	fgBlue      = "\033[34m"
	magentaCode = "\033[35m"
	fgMagenta   = "\033[38;5;198m"
	bgGrayCode  = "\033[48;5;237m"

	searchMarkerColor = fgBlue
	markerColor       = fgMagenta
	selectedBgColor   = bgGrayCode
)

// PromptUI manages the interactive script selection UI
type PromptUI struct {
	scripts         []project.Script
	filteredScripts []project.Script
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

// NewPromptUI creates a new prompt UI with the given scripts
func NewPromptUI(scripts []project.Script) *PromptUI {
	ui := &PromptUI{
		scripts:         scripts,
		filteredScripts: scripts,
		selectedIndex:   0,
		searchQuery:     "",
	}
	return ui
}

type scoredScript struct {
	script project.Script
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

	ui.filteredScripts = make([]project.Script, len(scored))
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (ui *PromptUI) moveUp() {
	if ui.selectedIndex > 0 {
		ui.selectedIndex--
	} else {
		ui.selectedIndex = len(ui.filteredScripts) - 1
	}
}

func (ui *PromptUI) moveDown() {
	if ui.selectedIndex < len(ui.filteredScripts)-1 {
		ui.selectedIndex++
	} else {
		ui.selectedIndex = 0
	}
}

func (ui *PromptUI) handleSearchInput(key []byte) {
	if isBackspace(key) {
		if len(ui.searchQuery) > 0 {
			runes := []rune(ui.searchQuery)
			if len(runes) > 0 {
				ui.searchQuery = string(runes[:len(runes)-1])
				ui.filterScripts()
			}
		}
		return
	}

	if len(key) == 0 || !utf8.Valid(key) {
		return
	}

	appended := false

	for _, r := range string(key) {
		if r == '\r' || r == '\n' {
			continue
		}
		if unicode.IsControl(r) && r != ' ' {
			continue
		}
		ui.searchQuery += string(r)
		appended = true
	}

	if appended {
		ui.filterScripts()
	}
}

func (ui *PromptUI) getSelectedScript() *project.Script {
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

func truncateText(text string, maxWidth int) string {
	if maxWidth <= 0 || len(text) == 0 {
		return ""
	}

	if runewidth.StringWidth(text) <= maxWidth {
		return text
	}

	const ellipsis = "..."
	ellipsisWidth := runewidth.StringWidth(ellipsis)
	if maxWidth <= ellipsisWidth {
		return ellipsis[:maxWidth]
	}

	var builder strings.Builder
	currentWidth := 0
	limit := maxWidth - ellipsisWidth

	for _, r := range text {
		rw := runewidth.RuneWidth(r)
		if currentWidth+rw > limit {
			break
		}
		builder.WriteRune(r)
		currentWidth += rw
	}

	builder.WriteString(ellipsis)
	return builder.String()
}

func highlightMatch(text, query string, isSelected bool) string {
	if query == "" {
		return text
	}

	textRunes := []rune(text)
	if len(textRunes) == 0 {
		return text
	}

	queryRunes := []rune(query)
	if len(queryRunes) == 0 {
		return text
	}

	lowerQuery := make([]rune, len(queryRunes))
	for i, r := range queryRunes {
		lowerQuery[i] = unicode.ToLower(r)
	}

	highlightColor := yellowCode
	if isSelected {
		highlightColor = "\033[1;33m"
	}

	resetColor := resetCode
	if isSelected {
		resetColor = "\033[0m" + selectedBgColor
	}

	// Try to highlight the longest contiguous match first.
	startIdx := -1
	for i := 0; i <= len(textRunes)-len(lowerQuery); i++ {
		match := true
		for j := 0; j < len(lowerQuery); j++ {
			if unicode.ToLower(textRunes[i+j]) != lowerQuery[j] {
				match = false
				break
			}
		}
		if match {
			startIdx = i
			break
		}
	}

	var builder strings.Builder

	if startIdx != -1 {
		builder.WriteString(string(textRunes[:startIdx]))
		builder.WriteString(highlightColor)
		builder.WriteString(string(textRunes[startIdx : startIdx+len(lowerQuery)]))
		builder.WriteString(resetColor)
		builder.WriteString(string(textRunes[startIdx+len(lowerQuery):]))
		return builder.String()
	}

	queryIdx := 0
	for _, r := range textRunes {
		if queryIdx < len(lowerQuery) && unicode.ToLower(r) == lowerQuery[queryIdx] {
			builder.WriteString(highlightColor)
			builder.WriteRune(r)
			builder.WriteString(resetColor)
			queryIdx++
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
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
		truncatedName := truncateText(script.Name, maxNameDisplayWidth)
		nameWidth := runewidth.StringWidth(truncatedName)
		if nameWidth > maxNameWidth {
			maxNameWidth = nameWidth
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

	ui.scriptStartLine = paddingLines + 1

	for i := startIdx; i < endIdx; i++ {
		script := ui.filteredScripts[i]
		isSelected := i == ui.selectedIndex
		if isSelected {
			output.WriteString(boldCode)
			output.WriteString(magentaCode)
			output.WriteString(selectedBgColor)
			output.WriteString(markerColor)
			output.WriteString("▌")
			output.WriteString(resetCode)
			output.WriteString(boldCode)
			output.WriteString(selectedBgColor)
			output.WriteString(" ")
		} else {
			output.WriteString("  ")
		}

		truncatedName := truncateText(script.Name, maxNameDisplayWidth)
		truncatedCommand := truncateText(script.Command, maxCommandDisplayWidth)
		nameWidth := runewidth.StringWidth(truncatedName)
		commandWidth := runewidth.StringWidth(truncatedCommand)

		highlightedName := highlightMatch(truncatedName, ui.searchQuery, isSelected)
		highlightedCommand := highlightMatch(truncatedCommand, ui.searchQuery, isSelected)

		output.WriteString(highlightedName)
		if nameWidth < maxNameWidth {
			output.WriteString(strings.Repeat(" ", maxNameWidth-nameWidth))
		}
		output.WriteString("  ")
		output.WriteString(highlightedCommand)
		if isSelected {
			currentWidth := 2 + maxNameWidth + 2 + commandWidth
			if currentWidth < termWidth {
				output.WriteString(strings.Repeat(" ", termWidth-currentWidth))
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
	output.WriteString(searchMarkerColor)
	output.WriteString("> ")
	output.WriteString(resetCode)
	output.WriteString(boldCode)
	output.WriteString(ui.searchQuery)
	output.WriteString(showCursorCode)

	fmt.Print(output.String())

	return startIdx, endIdx
}

// ShowScriptPrompt displays an interactive prompt for selecting a script
func ShowScriptPrompt() (*project.Script, error) {
	packageJSONPath, err := detector.FindPackageJSON()
	if err != nil {
		return nil, fmt.Errorf("cannot find package.json: %v", err)
	}

	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read package.json: %v", err)
	}

	var pkg project.PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, fmt.Errorf("cannot parse package.json: %v", err)
	}

	if len(pkg.OrderedScripts) == 0 {
		return nil, fmt.Errorf("no scripts found in package.json")
	}

	oldState, err := enableRawMode()
	if err != nil {
		return nil, err
	}

	fmt.Print(enableAltScreenCode)
	fmt.Print(enableMouseCode)

	defer func() {
		fmt.Print(disableMouseCode)
		fmt.Print(disableAltScreenCode)
		restoreTerminal(oldState)
		fmt.Print(showCursorCode)
	}()

	ui := NewPromptUI(pkg.OrderedScripts)

	// Set up terminal resize handling
	sigChan := make(chan os.Signal, 1)
	setupResizeSignal(sigChan)
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
