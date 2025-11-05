//go:build windows
// +build windows

package ui

import (
	"os"
)

func setupResizeSignal(sigChan chan os.Signal) {
	// No-op on Windows
}
