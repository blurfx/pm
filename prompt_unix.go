//go:build darwin || linux
// +build darwin linux

package main

import (
	"os"
	"os/signal"
	"syscall"
)

func setupResizeSignal(sigChan chan os.Signal) {
	signal.Notify(sigChan, syscall.SIGWINCH)
}
