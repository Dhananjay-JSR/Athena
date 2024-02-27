//go:build windows

package logManager

import "syscall"

func EnableVirtualTerminalProcessing(stream syscall.Handle, enable bool) error {
	kernel32Dll := syscall.NewLazyDLL("Kernel32.dll")
	setConsoleMode := kernel32Dll.NewProc("SetConsoleMode")
	//https://learn.microsoft.com/en-us/windows/console/setconsolemode
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x4

	var mode uint32
	err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	if err != nil {
		return err
	}

	if enable {
		mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	} else {
		mode &^= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	}

	ret, _, err := setConsoleMode.Call(uintptr(stream), uintptr(mode))
	if ret == 0 {
		return err
	}

	return nil
}
