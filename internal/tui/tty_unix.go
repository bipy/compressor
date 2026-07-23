//go:build !windows

package tui

import (
	"os"

	"golang.org/x/term"
)

func isTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
