package internal

import (
	"runtime"
	"syscall"
)

func InitWindowsEscape() {
	if runtime.GOOS == "windows" {
		//this enables ANSI Escape on Windows Console
		EnableVirtualTerminalProcessing(syscall.Stdout, true)
	}
}
