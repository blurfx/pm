//go:build windows
// +build windows

package main

import (
	"os"
)

func setupResizeSignal(sigChan chan os.Signal) {
	// No-op on Windows
}
