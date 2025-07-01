package main

import (
	"fmt"
	"os"

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

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}

func getTerminalSize() (width, height int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80, 24 // default size
	}
	return width, height
}

func readKey() ([]byte, error) {
	buf := make([]byte, 4)
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