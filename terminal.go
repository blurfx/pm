package main

import (
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"
)

func enableRawMode() (*term.State, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	return oldState, nil
}

func restoreTerminal(oldState *term.State) {
	term.Restore(int(os.Stdin.Fd()), oldState)
}

func getTerminalSize() (width, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24
	}
	return width, height
}

func readKey() ([]byte, error) {
	// Increase buffer size to handle mouse events (SGR format can be up to ~12 bytes)
	buf := make([]byte, 20)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func isArrowUp(key []byte) bool {
	return len(key) == 3 && key[0] == 27 && key[1] == 91 && key[2] == 65
}

func isArrowDown(key []byte) bool {
	return len(key) == 3 && key[0] == 27 && key[1] == 91 && key[2] == 66
}

func isEnter(key []byte) bool {
	return len(key) == 1 && (key[0] == 13 || key[0] == 10)
}

func isEscape(key []byte) bool {
	return len(key) == 1 && key[0] == 27
}

func isBackspace(key []byte) bool {
	return len(key) == 1 && (key[0] == 127 || key[0] == 8)
}

func isCtrlC(key []byte) bool {
	return len(key) == 1 && key[0] == 3
}

type MouseEvent struct {
	Button int
	X      int
	Y      int
	Type   string // "click", "scroll_up", "scroll_down"
}

// Parse SGR mouse event (ESC[<button;x;yM/m)
func parseMouseEvent(data []byte) (*MouseEvent, bool) {
	str := string(data)

	// Check for SGR mouse format
	if !strings.HasPrefix(str, "\033[<") {
		return nil, false
	}

	// Find the ending M or m
	endIdx := strings.IndexAny(str, "Mm")
	if endIdx == -1 {
		return nil, false
	}

	// Extract the parameters
	params := str[3:endIdx]
	parts := strings.Split(params, ";")
	if len(parts) != 3 {
		return nil, false
	}

	button, err1 := strconv.Atoi(parts[0])
	x, err2 := strconv.Atoi(parts[1])
	y, err3 := strconv.Atoi(parts[2])

	if err1 != nil || err2 != nil || err3 != nil {
		return nil, false
	}

	event := &MouseEvent{
		Button: button,
		X:      x,
		Y:      y,
	}

	// Determine event type based on button code
	isRelease := str[endIdx] == 'm'

	switch button {
	case 0, 1, 2: // Left, middle, right click
		if !isRelease {
			event.Type = "click"
		} else {
			return nil, false // Ignore release events
		}
	case 64: // Scroll up
		event.Type = "scroll_up"
	case 65: // Scroll down
		event.Type = "scroll_down"
	default:
		return nil, false
	}

	return event, event.Type != ""
}
